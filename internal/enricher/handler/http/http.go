package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	s "github.com/superhacker2002/personality-service/internal/enricher/service"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrInvalidPersonId = errors.New("invalid person ID")
	ErrReadRequestFail = errors.New("failed to read request body")
)

type service interface {
	DeletePerson(id int) error
	AddPerson(name, surname, patronymic string) (int, error)
}

type HttpHandler struct {
	s service
}

func New(s service) HttpHandler {
	handler := HttpHandler{s: s}
	return handler
}

func (h HttpHandler) SetRoutes(router *mux.Router) {
	router.HandleFunc("/", h.addPersonHandler).Methods(http.MethodPost)
	router.HandleFunc("/{personId}", h.deleteByIdHandler).Methods(http.MethodDelete)
}

func (h HttpHandler) addPersonHandler(w http.ResponseWriter, r *http.Request) {
	var person = struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic,omitempty"`
	}{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &person); err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.s.AddPerson(person.Name, person.Surname, person.Patronymic)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeResponse(w, map[string]int{"personId": id}, http.StatusCreated)

}

func (h HttpHandler) deleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	personId, err := intPathParam(r, "personId")
	if err != nil {
		http.Error(w, ErrInvalidPersonId.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.DeletePerson(personId)
	if errors.Is(err, s.ErrPersonNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpHandler) changePersonHandler(w http.ResponseWriter, r *http.Request) {}

func intPathParam(r *http.Request, varName string) (int, error) {
	vars := mux.Vars(r)
	varStr := vars[varName]
	varInt, err := strconv.Atoi(varStr)
	if err != nil {
		log.Printf("%v: %s\n", err, varStr)
		return 0, err
	}
	if varInt <= 0 {
		return 0, errors.New("parameter is less than zero")
	}
	return varInt, nil
}

func writeResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
		http.Error(w, s.ErrInternal.Error(), http.StatusInternalServerError)
	}
}
