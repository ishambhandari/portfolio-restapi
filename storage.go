package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	getWorks() ([]*Work, error)
	createWork(work *Work) error
	deleteWork(id int) error
	getWorkById(id int) (*Work, error)
}

type PostGresStore struct {
	db *sql.DB
}

func NewPostGresStore() (*PostGresStore, error) {
	// postgres_password := getEnv("POSTGRES_PASSWORD")
	connStr := "user=postgres dbname=postgres password=Isham@123 sslmode=disable"
	// connStr := fmt.Sprintf("user=postgres dbname=Main-postgres password=%s sslmode=disable", postgres_password)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		fmt.Println("here")
		return nil, err
	}
	return &PostGresStore{
		db: db,
	}, nil

}

func (s *PostGresStore) createWork(work *Work) error {
	query := `INSERT INTO work (title, description, imageUrl, code_link, live_link) VALUES ($1, $2, $3,$4,$5)`
	_, err := s.db.Exec(query, work.Title, work.Description, work.ImageUrl, work.Code_link, work.Live_link)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostGresStore) deleteWork(id int) error {
	query := `delete from work where id = $1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostGresStore) getWorkById(id int) (*Work, error) {
	query := `select * from work where id = $1`
	resp, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	resp.Next()
	response, err := scanIntoWorks(resp)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *PostGresStore) getWorks() ([]*Work, error) {
	resp, err := s.db.Query(`select * from work`)
	if err != nil {
		return nil, err
	}
	works := []*Work{}
	for resp.Next() {
		work, err := scanIntoWorks(resp)
		if err != nil {
			return nil, err
		}
		works = append(works, work)
	}
	return works, nil
}

func scanIntoWorks(resp *sql.Rows) (*Work, error) {
	work := new(Work)
	err := resp.Scan(
		&work.ID,
		&work.Title,
		&work.Description,
		&work.ImageUrl,
		&work.Code_link,
		&work.Live_link,
	)
	if err != nil {
		return nil, err
	}
	return work, nil
}
