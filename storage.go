package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
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
	fmt.Println("DB_HOST:", getEnv("DB_HOST"))
	fmt.Println("DB_USER:", getEnv("DB_USER"))
	fmt.Println("DB_NAME:", getEnv("DB_NAME"))
	fmt.Println("DB_PASSWORD:", getEnv("DB_PASSWORD"))
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=require",
		getEnv("DB_HOST"),
		getEnv("DB_USER"),
		getEnv("DB_NAME"),
		getEnv("DB_PASSWORD"))
	var db *sql.DB
	var err error

	// Retry loop to attempt connection multiple times
	for i := 0; i < 5; i++ { // Adjust the retry count as needed
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				fmt.Println("Connected to the database successfully!")
				break
			}
		}

		fmt.Printf("Failed to connect to the database: %v. Retrying in 2 seconds...\n", err)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to the database after multiple attempts: %v", err)
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
