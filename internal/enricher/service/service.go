package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrNoAge         = errors.New("failed to get person's age")
	ErrNoGender      = errors.New("failed to get person's gender")
	ErrNoNationality = errors.New("failed to get person's nationality")
	ErrInternal      = errors.New("internal server error")
)

type repository interface {
	AddPerson(name, surname, patronymic, gender, nationality string, age int) (int, error)
}

type Service struct {
	r repository
}

func New(r repository) Service {
	return Service{r: r}
}

func (s Service) DeletePerson(id int) error {
	return nil
}

func (s Service) AddPerson(name, surname, patronymic string) (int, error) {
	age, err := callExternalAPI(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		log.Println(ErrNoAge, ":", err)
		return 0, ErrNoAge
	}
	ageInt, _ := strconv.Atoi(age)

	gender, err := callExternalAPI(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		log.Println(ErrNoGender, ":", err)
		return 0, ErrNoGender
	}

	nationality, err := callExternalAPI(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		log.Println(ErrNoNationality, ":", err)
		return 0, ErrNoNationality
	}

	id, err := s.r.AddPerson(name, surname, patronymic, gender, nationality, ageInt)
	if err != nil {
		log.Println(err)
		return 0, ErrInternal
	}

	return id, nil
}

func callExternalAPI(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get response: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get information from response: %w", err)
	}

	return string(b), nil

}
