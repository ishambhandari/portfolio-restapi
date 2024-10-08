package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func (s *APIServer) Run() {
	log.Println("JSON API server running on %S port", s.listenAddr)
	router := mux.NewRouter()
	router.HandleFunc("/works", makeHTTPHandleFunc(s.handleWork))
	router.HandleFunc("/works/{id}", makeHTTPHandleFunc(s.handleWorkById))
	router.HandleFunc("/mail", makeHTTPHandleFunc(s.handleEmail))
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleWork(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	fmt.Println("asdf", r.Method)
	if r.Method == "GET" {
		return s.handleGetWork(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateWork(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetWork(w http.ResponseWriter, r *http.Request) error {
	works, err := s.store.getWorks()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, works)
}

func (s *APIServer) handleCreateWork(w http.ResponseWriter, r *http.Request) error {
	createWorkRequest := &Work{}
	if err := json.NewDecoder(r.Body).Decode(createWorkRequest); err != nil {
		return err
	}
	work := NewWork(createWorkRequest.Title, createWorkRequest.Description, createWorkRequest.ImageUrl, createWorkRequest.Code_link, createWorkRequest.Live_link)
	if err := s.store.createWork(work); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, work)
}

func (s *APIServer) handleWorkById(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing work ID", http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid work id", http.StatusBadRequest)
	}
	if r.Method == "GET" {
		work, err := s.store.getWorkById(id)
		if err != nil {
			return err
		}
		return WriteJson(w, http.StatusOK, work)

	}
	if r.Method == "DELETE" {
		err = s.store.deleteWork(id)
		if err != nil {
			return err
		}
	}
	return WriteJson(w, http.StatusOK, "Done")

}

func (s *APIServer) handleEmail(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	if r.Method == "GET" {
		http.Error(w, "Get request not supported", http.StatusBadRequest)
	}
	if r.Method == "POST" {
		return s.handlePostEmail(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handlePostEmail(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	createEmail := &PostContactDetails{}
	if err := json.NewDecoder(r.Body).Decode(createEmail); err != nil {
		return err
	}
	if err := sendMail(getEnv("EMAIL_SENDER"), getEnv("EMAIL_RECEIVER"), createEmail.Name, createEmail.Description); err != nil {
		return err
	}
	fmt.Println("herer!!!!")
	return WriteJson(w, http.StatusOK, "Done")
}

func sendMail(sender string, receiver string, receiver_name string, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// The email message headers and body.
	message := []byte("Subject: " + "Personal Website Message by: " + receiver_name + "\r\n" +
		"\r\n" + body)

	// Authentication for the email sender.
	auth := smtp.PlainAuth("", sender, getEnv("EMAIL_PASSWORD"), smtpHost)

	// Send the email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{receiver}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
