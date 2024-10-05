package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
}

func (s *APIServer) Run() {
	log.Println("JSON API server running on %S port", s.listenAddr)
	router := mux.NewRouter()
	router.HandleFunc("/works", makeHTTPHandleFunc(s.handleWork))
	router.HandleFunc("/works/{id}", makeHTTPHandleFunc(s.handleWorkById))
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
