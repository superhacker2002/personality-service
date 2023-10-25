package httphandler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	s "github.com/superhacker2002/personality-service/internal/enricher/service"
	"github.com/superhacker2002/personality-service/internal/entity"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrInvalidPersonId = errors.New("invalid person ID")
	ErrReadRequestFail = errors.New("failed to read request body")
	ErrInvalidOffset   = errors.New("invalid offset parameter")
	ErrInvalidLimit    = errors.New("invalid limit parameter")
)

type service interface {
	DeletePerson(id int) error
	Person(id int) (entity.Person, error)
	AddPerson(name, surname, patronymic string) (int, error)
	PeopleByName(name string, offset int, limit int) ([]entity.Person, error)
	AllPeople(offset, limit int) ([]entity.Person, error)
	UpdatePerson(id int, name, surname, patronymic string) error
}

type person struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         string `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type page struct {
	offset int
	limit  int
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
	router.HandleFunc("/{personId}", h.getByIdHandler).Methods(http.MethodGet)
	router.HandleFunc("/", h.findPeopleHandler).Methods(http.MethodGet)
	router.HandleFunc("/{personId}", h.updatePersonHandler).Methods(http.MethodPut)
}

func (h HttpHandler) findPeopleHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	p, err := getPage(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var people []entity.Person

	if name != "" {
		people, err = h.s.PeopleByName(name, p.offset, p.limit)
	} else {
		people, err = h.s.AllPeople(p.offset, p.limit)
	}

	if errors.Is(err, s.ErrPeopleNotFound) {
		log.Println(s.ErrPeopleNotFound)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, entitiesToDTO(people), http.StatusOK)

}

func (h HttpHandler) addPersonHandler(w http.ResponseWriter, r *http.Request) {
	var p = struct {
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

	if err = json.Unmarshal(body, &p); err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.s.AddPerson(p.Name, p.Surname, p.Patronymic)

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

func (h HttpHandler) getByIdHandler(w http.ResponseWriter, r *http.Request) {
	personId, err := intPathParam(r, "personId")
	if err != nil {
		http.Error(w, ErrInvalidPersonId.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.s.Person(personId)
	if errors.Is(err, s.ErrPersonNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, entityToDTO(p), http.StatusOK)
}

func (h HttpHandler) updatePersonHandler(w http.ResponseWriter, r *http.Request) {
	personId, err := intPathParam(r, "personId")
	if err != nil {
		log.Println(err)
		http.Error(w, ErrInvalidPersonId.Error(), http.StatusBadRequest)
		return
	}

	var p = struct {
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

	if err = json.Unmarshal(body, &p); err != nil {
		log.Println(err)
		http.Error(w, ErrReadRequestFail.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.UpdatePerson(personId, p.Name, p.Surname, p.Patronymic)

	if errors.Is(err, s.ErrPersonNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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

func getPage(r *http.Request) (page, error) {
	const (
		defaultOffset = 0
		defaultLimit  = 10
	)
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	var p page
	if offsetStr == "" || limitStr == "" {
		p.offset = defaultOffset
		p.limit = defaultLimit
		log.Println("missing offset or limit, default values are used")
		return p, nil
	}

	var err error
	if p.offset, err = strconv.Atoi(offsetStr); err != nil {
		return p, ErrInvalidOffset
	}
	if p.limit, err = strconv.Atoi(limitStr); err != nil {
		return p, ErrInvalidLimit
	}

	if p.offset < 0 {
		return p, ErrInvalidOffset
	}
	if p.limit < 0 {
		return p, ErrInvalidLimit
	}

	return p, nil
}

func entitiesToDTO(people []entity.Person) []person {
	var DTOSessions []person
	for _, p := range people {
		DTOSessions = append(DTOSessions, entityToDTO(p))
	}
	return DTOSessions
}

func entityToDTO(p entity.Person) person {
	return person{
		Name:        p.Name,
		Surname:     p.Surname,
		Patronymic:  p.Patronymic,
		Age:         strconv.Itoa(p.Age),
		Gender:      p.Gender,
		Nationality: p.Nationality,
	}
}
