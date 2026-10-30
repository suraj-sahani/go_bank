package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// This is a helper function that will convert the response
// to JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) // Set the status for the response
	return json.NewEncoder(w).Encode(v)
}

type APIError struct {
	Error string
}

// Our custom handler function that allow us to handle
// errors within the functions we want
type apiFunc func(w http.ResponseWriter, r *http.Request) error

// This acts as a higher order function that satisfies the
// expected type by mux but also fulfills our functionality
// of handling error the way we want
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error here
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

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

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// We will get an error if we use our custom handle functions because
	// of the type difference between was mux wants and what we are providing
	// The change is simple, mux does not expect an error to be returned.
	// Thus, we will make our own custom functions that will allow us
	// to handle error where we want
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleGetAllAccounts)).Methods("GET")
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleUpdateAccount)).Methods("PATCH")
	router.HandleFunc("/account/{id}/transfer", makeHTTPHandleFunc(s.handleTransfer)).Methods("PATCH")
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")

	log.Println("Server running at port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	res, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, res)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("suraj", "sahani")
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccPayload := new(CreateAccountDTO)
	if err := json.NewDecoder(r.Body).Decode(&createAccPayload); err != nil {
		return err
	}

	account := NewAccount(createAccPayload.FirstName, createAccPayload.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIServer) handleUpdateAccount(W http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
