package main

import (
	"context"
	// "database/sql"
	"fmt"
	// "log"
	// "time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	_ "github.com/lib/pq"
)

type Storage interface {
	getWorks() ([]*Work, error)
	createWork(work *Work) error
	deleteWork(id int) error
	getWorkById(id int) (*Work, error)
}

// type PostGresStore struct {
// 	db *sql.DB
// }

type FirebaseStore struct {
	client *firestore.Client
}

const (
	projectID      string = "portfolio-4b08a"
	collectionName string = "works"
)

func NewFirebaseStore() (*FirebaseStore, error) {
	ctx := context.Background()
	// Initialize Firebase App with credentials
	conf := &firebase.Config{ProjectID: projectID}
	// opt := option.WithCredentialsFile("~/.config/firebase/portfolio-4b08a-firebase-adminsdk-najxl-533618384f.json")
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firestore client: %v", err)
	}

	return &FirebaseStore{
		client: client,
	}, nil
}

func (fs *FirebaseStore) createWork(work *Work) error {
	ctx := context.Background()
	docRef := fs.client.Collection(collectionName).NewDoc()
	_, err := docRef.Set(ctx, map[string]interface{}{
		"id":          docRef.ID,
		"title":       work.Title,
		"description": work.Description,
		"imageUrl":    work.ImageUrl,
		"code_link":   work.Code_link,
		"live_link":   work.Live_link,
	})
	if err != nil {
		return fmt.Errorf("Failed to create work %v", err)
	}
	work.ID = docRef.ID
	return nil
}

func (fs *FirebaseStore) deleteWork(id int) error {
	ctx := context.Background()

	// Get the document reference by ID
	query := fs.client.Collection(collectionName).Where("id", "==", id).Limit(1)
	iter := query.Documents(ctx)

	// Get the document and delete it
	for {
		doc, err := iter.Next()
		if err != nil {
			return fmt.Errorf("failed to find document to delete: %v", err)
		}
		// Delete the document
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete document: %v", err)
		}
		break
	}

	return nil
}

func (fs *FirebaseStore) getWorkById(id int) (*Work, error) {
	ctx := context.Background()

	// Query Firestore to find the work by ID
	query := fs.client.Collection(collectionName).Where("id", "==", id).Limit(1)
	iter := query.Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get work by ID: %v", err)
	}

	var work Work
	err = doc.DataTo(&work)
	if err != nil {
		return nil, fmt.Errorf("failed to map document to Work struct: %v", err)
	}

	return &work, nil
}

func (fs *FirebaseStore) getWorks() ([]*Work, error) {
	ctx := context.Background()
	// Query all documents in the "works" collection
	iter := fs.client.Collection(collectionName).Documents(ctx)

	var works []*Work
	for {
		doc, err := iter.Next()
		fmt.Println("This is doc", iter)
		if err != nil {
			break
		}

		var work Work
		err = doc.DataTo(&work)
		if err != nil {
			return nil, fmt.Errorf("failed to map document to Work struct: %v", err)
		}

		works = append(works, &work)
	}

	if len(works) == 0 {
		return nil, fmt.Errorf("no works found")
	}

	return works, nil
}

// func NewPostGresStore() (*PostGresStore, error) {
// 	fmt.Println("DB_HOST:", getEnv("DB_HOST"))
// 	fmt.Println("DB_USER:", getEnv("DB_USER"))
// 	fmt.Println("DB_NAME:", getEnv("DB_NAME"))
// 	fmt.Println("DB_PASSWORD:", getEnv("DB_PASSWORD"))
// 	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=require",
// 		getEnv("DB_HOST"),
// 		getEnv("DB_USER"),
// 		getEnv("DB_NAME"),
// 		getEnv("DB_PASSWORD"))
// 	var db *sql.DB
// 	var err error
//
// 	// Retry loop to attempt connection multiple times
// 	for i := 0; i < 5; i++ { // Adjust the retry count as needed
// 		db, err = sql.Open("postgres", connStr)
// 		if err == nil {
// 			err = db.Ping()
// 			if err == nil {
// 				fmt.Println("Connected to the database successfully!")
// 				break
// 			}
// 		}
//
// 		fmt.Printf("Failed to connect to the database: %v. Retrying in 2 seconds...\n", err)
// 		time.Sleep(2 * time.Second) // Wait before retrying
// 	}
//
// 	if err != nil {
// 		return nil, fmt.Errorf("could not connect to the database after multiple attempts: %v", err)
// 	}
//
// 	return &PostGresStore{
// 		db: db,
// 	}, nil
// }
//
// func (s *PostGresStore) createWork(work *Work) error {
// 	query := `INSERT INTO work (title, description, imageUrl, code_link, live_link) VALUES ($1, $2, $3,$4,$5)`
// 	_, err := s.db.Exec(query, work.Title, work.Description, work.ImageUrl, work.Code_link, work.Live_link)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (s *PostGresStore) deleteWork(id int) error {
// 	query := `delete from work where id = $1`
// 	_, err := s.db.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (s *PostGresStore) getWorkById(id int) (*Work, error) {
// 	query := `select * from work where id = $1`
// 	resp, err := s.db.Query(query, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp.Next()
// 	response, err := scanIntoWorks(resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil
// }
//
// func (s *PostGresStore) getWorks() ([]*Work, error) {
// 	resp, err := s.db.Query(`select * from work`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	works := []*Work{}
// 	for resp.Next() {
// 		work, err := scanIntoWorks(resp)
// 		if err != nil {
// 			return nil, err
// 		}
// 		works = append(works, work)
// 	}
// 	return works, nil
// }
//
// func scanIntoWorks(resp *sql.Rows) (*Work, error) {
// 	work := new(Work)
// 	err := resp.Scan(
// 		&work.ID,
// 		&work.Title,
// 		&work.Description,
// 		&work.ImageUrl,
// 		&work.Code_link,
// 		&work.Live_link,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return work, nil
// }
