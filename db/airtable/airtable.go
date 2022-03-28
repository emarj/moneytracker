package airtable

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"ronche.se/moneytracker/domain"
)

type Airtable struct {
	client *http.Client
	apiKey string
}

const url = "https://api.airtable.com/v0/appfbyqOSj86V5Wku/Transactions"

func NewAirtable(apiKey string) *Airtable {
	return &Airtable{client: &http.Client{
		Timeout: time.Second * 10,
	}, apiKey: apiKey}
}

func (at *Airtable) request(method, url string, body io.Reader) (string, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+at.apiKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := at.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()
	buf := new(strings.Builder)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		return "", fmt.Errorf("Got error %s", err.Error())
	}
	responseBody := buf.String()

	if response.StatusCode >= 299 {
		if err != nil {
			return "", fmt.Errorf("Response: %d (%s), %s", response.StatusCode, response.Status, responseBody)
		}
	}

	return responseBody, nil

	//return "", nil
}

//DELETE
//DeleteTransaction(tID uuid.UUID) error
//UPDATE
//UpdateTransaction(t *domain.Transaction) (uuid.UUID, error)
//GET
////GetTransaction(id uuid.UUID) (*domain.Transaction, error)
//GetTransactionsByAccount(aID string) ([]*domain.Transaction, error)
//GetTransactionsByUser(uID string) ([]*domain.Transaction, error)

func (at *Airtable) InsertTransaction(t *domain.Transaction) (uuid.UUID, error) {
	if t == nil {
		return uuid.Nil, errors.New("nil transaction")
	}

	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	t.ID = id

	s, err := json.Marshal(t)
	if err != nil {
		return uuid.Nil, err
	}

	payload := fmt.Sprintf("{\"records\" : [{\"fields\" : %s}]}", string(s))

	fmt.Println(payload)

	fmt.Println(IsJSON(payload))

	response, err := at.request(http.MethodPost, url, strings.NewReader(payload))

	fmt.Printf("Response: %s\n\r", response)

	return id, err
	//return id, err
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
