package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// APIError represents an error response with an HTTP status code.
type APIError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"error"`
}

func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates an APIError with the given HTTP status code and message.
func NewAPIError(statusCode int, msg string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    msg,
	}
}

// writeJSON serializes v as JSON and writes it to the response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// apiHandlerFunc is an HTTP handler that returns an error for centralized error handling.
type apiHandlerFunc func(w http.ResponseWriter, r *http.Request) error

// makeHTTPHandler wraps an apiHandlerFunc into a standard http.HandlerFunc,
// catching returned errors and writing appropriate JSON error responses.
func makeHTTPHandler(f apiHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) {
				writeJSON(w, apiErr.StatusCode, apiErr)
			} else {
				writeJSON(w, http.StatusInternalServerError, &APIError{
					StatusCode: http.StatusInternalServerError,
					Message:    "internal server error",
				})
			}
		}
	}
}

type APIServer struct {
	listenAddr string
	store      AccountStore
}

func NewAPIServer(listenAddr string, store AccountStore) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandler(s.handleGetAllAccounts)).Methods("GET")
	router.HandleFunc("/account", makeHTTPHandler(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleGetAccountByID)).Methods("GET")
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleUpdateAccount)).Methods("PATCH")
	router.HandleFunc("/account/{id}/transfer", makeHTTPHandler(s.handleTransfer)).Methods("PATCH")
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleDeleteAccount)).Methods("DELETE")

	log.Println("Server running at port:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "failed to retrieve accounts")
	}

	return writeJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, fmt.Sprintf("invalid account id: %s", idStr))
	}

	acc, err := s.store.GetAccountByID(id)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return NewAPIError(http.StatusNotFound, fmt.Sprintf("account %d not found", id))
		}
		return NewAPIError(http.StatusInternalServerError, "failed to retrieve account")
	}

	return writeJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return NewAPIError(http.StatusBadRequest, "invalid request body")
	}

	if req.FirstName == "" || req.LastName == "" {
		return NewAPIError(http.StatusBadRequest, "first_name and last_name are required")
	}

	account := NewAccount(req.FirstName, req.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return NewAPIError(http.StatusInternalServerError, "failed to create account")
	}

	return writeJSON(w, http.StatusCreated, account)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return NewAPIError(http.StatusNotImplemented, "update account not implemented")
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, fmt.Sprintf("invalid account id: %s", idStr))
	}

	if err := s.store.DeleteAccount(id); err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return NewAPIError(http.StatusNotFound, fmt.Sprintf("account %d not found", id))
		}
		return NewAPIError(http.StatusInternalServerError, "failed to delete account")
	}

	return writeJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return NewAPIError(http.StatusNotImplemented, "transfer not implemented")
}
