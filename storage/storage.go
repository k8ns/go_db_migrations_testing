package storage

import (
	"database/sql"
	"fmt"
)

type Article struct {
	Id      int
	Title   string
	Created string
}

func (a *Article) String() string {
	return fmt.Sprintf("%d:%s-%s", a.Id, a.Title, a.Created)
}

type Storage struct {
	db *sql.DB
}

func (s *Storage) Get(id int) (*Article, error) {
	row := s.db.QueryRow("SELECT * FROM articles WHERE id = ?", id)

	a := &Article{}
	err := row.Scan(&a.Id, &a.Title, &a.Created)
	if err != nil {
		return nil, err
	}

	return a, nil
}
