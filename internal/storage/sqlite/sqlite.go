package sqlite

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
    db *sql.DB
}

func New(storagePath string) (*Storage, error) {
    const op = "storage.sqlite.New"

    db, err := sql.Open("sqlite3", storagePath)
    if err != nil {
        // возврвщаем fmt.Errorf() а не err, чтобы было наглядно видно
        // где происходит ошибка, видно по переменной 'op')
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    stmt, err := db.Prepare(`
        CREATE TABLE IF NOT EXISTS url(
            id INTEGER PRIMARY KEY,
            alias TEXT NOT NULL UNIQUE,
            url TEXT NOT NULL);
        CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
        )
    `)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err) 
    }

    _, err = stmt.Exec()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err) 
    }

    return &Storage{db: db}, nil
}

// add save url in databast function 
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
    const op = "storege.sqlite.SaveURL"

    // prepare request
    stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    // execute prepared request
    res, err := stmt.Exec(urlToSave, alias)
    if err != nil {
        // make error code check and return our error for SQL exception
        if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
            return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
        }
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    // get last isert item id
    id, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("%s: failed to last inseert id: %w", op, err)
    }
    
    return id, nil
}
