package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// структура нужна для парсинга json из реквеста
type Request struct {
    // validate это надо для пакета валидации(validator/v10)
    // показывае что это поле reqired и url
    URL string `json:"url" validate:"reqired,url"`
    // omitempty этот параметр можно указать в strucTag json
    // если он есть и параметр alias не заполнен
    // то в место пустой строки параметра этого не будет
    Alias string `json:"alias,omitempty"`
}

type Response struct {
    // вынесем Status и Error отдельно, так как они будуи использоваться
    // в нескольких местах
    resp.Response
    // возвращаем alias, так как если его не быдо в запросе
    // то рандомно его генерим и присылаем пользователю в ответе
    Alias string `json:"alias,omitempty"`
}

// TODO: лучше вынести в конфиг файл
const aliasLength = 6

// сигнатура этого метода должна быть такая же как у метода SaveURL в файле internal/storage/sqlite/sqlite.go 
// так как наследует этот интерфейс
type URLSaver interface {
    SaveURL(urlToSave string, alias string) (int64, error)
}

// функция конструктор для handler
// при подключении этого handler будем вызывать функцию New
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handler.url.save.New"

        // добавляем доп параметры для логгера
        log = log.With(
            slog.String("op", op),
            // трейсинг запросов, достаем из контекста реквеста
            slog.String("request_id", middleware.GetReqID(r.Context())),
        )

        // создаем объект запроса в который будем анмаршалить поступивший запрос
        var req Request

        // анмаршалим запрос тут(парсим запрос)
        err := render.DecodeJSON(r.Body, &req)  
        if err != nil {
            log.Error("failed to decode request body", sl.Err(err))
            render.JSON(w, r, resp.Error("failed to decode request"))
            return
        }

        // логгируем распарсенный запрос, для проверки корректности параметров
        log.Info("request body decoded", slog.Any("request", req))

        // валидируем поля в реквесте
        if err := validator.New().Struct(req); err != nil {
            // приводим ошибку к нужному типу
            validateErr := err.(validator.ValidationErrors) 

            // пишем ошибку в лог
            log.Error("invalid request", sl.Err(err))

            // функция resp.ValidationError() формирует человекочитаймую ошибку
            // для ответа на request
            render.JSON(w, r, resp.ValidationError(validateErr))
            
            // возвращаем ответ на запрос с ошибкой
            return
        }
        
        alias := req.Alias
        // если alias пустой, то генерим новый рандомный
        if alias == "" {
            // TODO: может сгенериться уже существующий alias, надо это разрулить
            alias = random.NewRandomString(aliasLength)
        }

        id, err := urlSaver.SaveURL(req.URL, alias)
        // обрабатываем ошибку ErrURLExists, так как функция SaveURL() может вернуть ошибку ErrURLExists
        if errors.Is(err, storage.ErrURLExists) {
            log.Info("url already exists", slog.String("url", req.URL))

            render.JSON(w, r, resp.Error("url already exists"))

            return
        }

        // обрабатываем все остальные ошибки
        if err != nil {
            log.Error("failed to add url", sl.Err(err))

            render.JSON(w, r, resp.Error("failed to add url"))
            
            return
        }

        // запрос успешно сохранен, запишем это в лог
        log.Info("url added", slog.Int64("id", id))

        // возвращаем успешный ответ
        responseOK(w, r, alias)
    }
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
    render.JSON(w, r, Response{
        Response: resp.OK(),
        Alias: alias,
    })
}
