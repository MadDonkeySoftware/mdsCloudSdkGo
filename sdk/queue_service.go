package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// QueueServiceClient Client to interact with the MDS Cloud queue service
type QueueServiceClient struct {
	queueServiceURL string
	authManager     *AuthManager
}

// CreateQueueArgs Data needed to create a new queue
type CreateQueueArgs struct {
	Name     string `json:"name"`
	Resource string `json:"resource,omitempty"`
	Dlq      string `json:"dlq,omitempty"`
}

// CreateQueueResult Create queue results
type CreateQueueResult struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	Orid   string `json:"orid"`
}

// CreateQueue Attempts to create a new queue with the MDS Cloud deployment
func (qs *QueueServiceClient) CreateQueue(data *CreateQueueArgs) (*CreateQueueResult, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/queue", qs.queueServiceURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("Could not build request to create queue")
	}

	token, err := qs.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 200:
		fallthrough
	case 201:
		payload := CreateQueueResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			// return nil, errors.New("Could not decode response from API of resource")
			return nil, err
		}

		if r.StatusCode == 200 {
			payload.Status = "exists"
		} else {
			payload.Status = "created"
		}

		return &payload, nil
	default:
		body, _ = ioutil.ReadAll(r.Body)
		return nil, fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// func DeleteMessage() {}

// DeleteQueueArgs Data needed to delete a queue
type DeleteQueueArgs struct {
	Orid string `json:"orid"`
}

// DeleteQueue Attempts to delete a queue from the MDS Cloud deployment
func (qs *QueueServiceClient) DeleteQueue(data *DeleteQueueArgs) error {
	client := &http.Client{Timeout: API_TIMEOUT}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/queue/%s", qs.queueServiceURL, data.Orid), nil)
	if err != nil {
		return errors.New("Could not build request to delete queue")
	}

	token, err := qs.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 204:
		return nil
	default:
		body, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// func EnqueueMessage() {}

// func FetchMessage() {}

// GetQueueDetailsArgs Data needed to fetch queue details
type GetQueueDetailsArgs struct {
	Orid string `json:"orid"`
}

// GetQueueDetailsResult The details of the given queue
type GetQueueDetailsResult struct {
	Orid     string `json:"orid"`
	Resource string `json:"resource,omitempty"`
	Dlq      string `json:"dlq,omitempty"`
}

// GetQueueDetails Gets details for the specified queue
func (qs *QueueServiceClient) GetQueueDetails(data *GetQueueDetailsArgs) (*GetQueueDetailsResult, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/queue/%s/details", qs.queueServiceURL, data.Orid), nil)
	if err != nil {
		return nil, errors.New("Could not build request to delete queue")
	}

	token, err := qs.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 200:
		payload := GetQueueDetailsResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			// return nil, errors.New("Could not decode response from API of resource")
			return nil, err
		}

		payload.Orid = data.Orid
		return &payload, nil
	default:
		body, _ := ioutil.ReadAll(r.Body)
		return nil, fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// func GetQueueLength() {}

// func ListQueues() {}

// func UpdateQueue() {}

// UpdateQueueArgs Data needed to create a new queue
type UpdateQueueArgs struct {
	Orid     string `json:"orid"`
	Resource string `json:"resource,omitempty"`
	Dlq      string `json:"dlq,omitempty"`
}

// UpdateQueue Attempts to create a new queue with the MDS Cloud deployment
func (qs *QueueServiceClient) UpdateQueue(data *UpdateQueueArgs) error {
	client := &http.Client{Timeout: API_TIMEOUT}

	type updateQueuePayload struct {
		Resource interface{} `json:"resource,omitempty"`
		Dlq      interface{} `json:"dlq,omitempty"`
	}

	postPayload := &updateQueuePayload{}
	if data.Resource == "NULL" {
		postPayload.Resource = ""
	} else {
		postPayload.Resource = data.Resource
	}

	if data.Dlq == "NULL" {
		postPayload.Dlq = ""
	} else {
		postPayload.Dlq = data.Dlq
	}

	body, err := json.Marshal(postPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/queue/%s", qs.queueServiceURL, data.Orid), bytes.NewBuffer(body))
	if err != nil {
		return errors.New("Could not build request to create queue")
	}

	token, err := qs.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 200:
		return nil
	default:
		body, _ = ioutil.ReadAll(r.Body)
		return fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}
