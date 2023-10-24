package controller

import (
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrInvalidPersonId = errors.New("invalid person ID")
)

type service interface {
	DeletePerson(id int) error
}

type HttpHandler struct {
	s service
}

func New(router *mux.Router, s service) HttpHandler {
	handler := HttpHandler{s: s}
	handler.setRoutes(router)

	return handler
}

func (h HttpHandler) setRoutes(router *mux.Router) {
	router.HandleFunc("/{personId}", h.deleteByIdHandler).Methods(http.MethodDelete)
}

func (h HttpHandler) addPersonHandler(w http.ResponseWriter, r *http.Request) {}

func (h HttpHandler) deleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	personId, err := intPathParam(r, "personId")
	if err != nil {
		http.Error(w, ErrInvalidPersonId.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.DeletePerson(personId)
	//if errors.Is(err, service.ErrPersonNotFound) {
	//	http.Error(w, err.Error(), http.StatusNotFound)
	//	return
	//}

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
