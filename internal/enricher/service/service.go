package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/superhacker2002/personality-service/internal/entity"
	"io"
	"log"
	"net/http"
)

var (
	ErrNoAge          = errors.New("failed to get person's age")
	ErrNoGender       = errors.New("failed to get person's gender")
	ErrNoNationality  = errors.New("failed to get person's nationality")
	ErrInternal       = errors.New("internal server error")
	ErrPersonNotFound = errors.New("person was not found")
	ErrPeopleNotFound = errors.New("people were not found")
)

type repository interface {
	AddPerson(name, surname, patronymic, gender, nationality string, age int) (int, error)
	DeletePerson(id int) (found bool, err error)
	PeopleByName(name string, offset int, limit int) ([]entity.Person, error)
	AllPeople(offset, limit int) ([]entity.Person, error)
	Person(id int) (entity.Person, error)
}

type Service struct {
	r repository
}

func New(r repository) Service {
	return Service{r: r}
}

func (s Service) Person(id int) (entity.Person, error) {
	people, err := s.r.Person(id)
	if errors.Is(err, ErrPersonNotFound) {
		return entity.Person{}, err
	}

	if err != nil {
		return entity.Person{}, ErrInternal
	}
	return people, nil
}

func (s Service) AllPeople(offset, limit int) ([]entity.Person, error) {
	people, err := s.r.AllPeople(offset, limit)
	if errors.Is(err, ErrPeopleNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, ErrInternal
	}
	return people, nil
}

func (s Service) PeopleByName(name string, offset, limit int) ([]entity.Person, error) {
	people, err := s.r.PeopleByName(name, offset, limit)
	if errors.Is(err, ErrPeopleNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, ErrInternal
	}
	return people, nil
}

func (s Service) DeletePerson(id int) error {
	found, err := s.r.DeletePerson(id)
	if err != nil {
		return ErrInternal
	}
	if !found {
		log.Println("failed to delete person from database:", ErrPersonNotFound)
		return ErrPersonNotFound
	}

	return nil
}

func (s Service) AddPerson(name, surname, patronymic string) (int, error) {
	averageAge, err := age(name)
	if err != nil {
		return 0, err
	}

	averageGender, err := gender(name)
	if err != nil {
		return 0, err
	}

	averageNationality, err := nationality(name)
	if err != nil {
		return 0, err
	}

	id, err := s.r.AddPerson(name, surname, patronymic, averageGender, averageNationality, averageAge)
	if err != nil {
		return 0, ErrInternal
	}

	return id, nil
}

func age(name string) (int, error) {
	ageStruct := struct {
		Age int `json:"age"`
	}{}

	b, err := callExternalAPI(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		log.Println(ErrNoAge, ":", err)
		return 0, ErrNoAge
	}

	err = json.Unmarshal(b, &ageStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal response to struct: %w", err)
	}

	return ageStruct.Age, nil
}

func gender(name string) (string, error) {
	genderStruct := struct {
		Gender string `json:"gender"`
	}{}

	b, err := callExternalAPI(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		log.Println(ErrNoGender, ":", err)
		return "", ErrNoGender
	}

	err = json.Unmarshal(b, &genderStruct)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response to struct: %w", err)
	}

	return genderStruct.Gender, nil
}

func nationality(name string) (string, error) {
	type Country struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	}

	nationalityStruct := struct {
		Country []Country `json:"country"`
	}{}

	b, err := callExternalAPI(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		log.Println(ErrNoNationality, ":", err)
		return "", ErrNoNationality
	}

	err = json.Unmarshal(b, &nationalityStruct)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response to struct: %w", err)
	}

	return nationalityStruct.Country[0].CountryID, nil

}

func callExternalAPI(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			err = fmt.Errorf("failed to close response body: %w", err)
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to get information from response: %w", err)
	}

	return b, err

}
