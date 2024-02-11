package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// StateMachineServiceClient Client to interact with the MDS Cloud state machine service
type StateMachineServiceClient struct {
	stateMachineServiceURL string
	authManager            *AuthManager
}

// CreateStateMachineArgs Data needed to create a new state machine
type CreateStateMachineArgs struct {
	Definition string
}

// CreateStateMachineResult Create state machine results
type CreateStateMachineResult struct {
	Orid string `json:"orid"`
}

// CreateStateMachine Attempts to create a new state machine within the MDS Cloud deployment
func (cs *StateMachineServiceClient) CreateStateMachine(data *CreateStateMachineArgs) (*CreateStateMachineResult, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	// body, err := json.Marshal(data)
	// if err != nil {
	// 	return err
	// }

	body := bytes.NewBuffer([]byte(data.Definition))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/machine", cs.stateMachineServiceURL), body)
	if err != nil {
		return nil, errors.New("could not build request to create state machine")
	}

	token, err := cs.authManager.GetAuthenticationToken(nil)
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
		payload := CreateStateMachineResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			// return nil, errors.New("Could not decode response from API of resource")
			return nil, err
		}

		return &payload, nil
	default:
		body, _ := io.ReadAll(r.Body)
		return nil, fmt.Errorf("did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// GetStateMachineDetailsArgs Data needed to fetch state machine details
type GetStateMachineDetailsArgs struct {
	Orid string
}

// GetStateMachineDetailsResult Get state machine details result
type GetStateMachineDetailsResult struct {
	Orid       string      `json:"orid"`
	Name       string      `json:"name"`
	Definition interface{} `json:"definition"`
}

// GetStateMachineDetails Attempts to fetch the details of a state machine within the MDS Cloud deployment
func (cs *StateMachineServiceClient) GetStateMachineDetails(data *GetStateMachineDetailsArgs) (*GetStateMachineDetailsResult, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/machine/%s", cs.stateMachineServiceURL, data.Orid), nil)
	if err != nil {
		return nil, errors.New("could not build request to get state machine details")
	}

	token, err := cs.authManager.GetAuthenticationToken(nil)
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
		payload := GetStateMachineDetailsResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			// return nil, errors.New("Could not decode response from API of resource")
			return nil, err
		}

		return &payload, nil
	default:
		body, _ := io.ReadAll(r.Body)
		return nil, fmt.Errorf("did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// UpdateStateMachineArgs Data needed to create a new state machine
type UpdateStateMachineArgs struct {
	Orid       string
	Definition string
}

// UpdateStateMachineResult Create state machine results
type UpdateStateMachineResult struct {
	Orid string `json:"orid"`
}

// UpdateStateMachine Attempts to create a new state machine within the MDS Cloud deployment
func (cs *StateMachineServiceClient) UpdateStateMachine(data *UpdateStateMachineArgs) (*UpdateStateMachineResult, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	body := bytes.NewBuffer([]byte(data.Definition))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/machine/%s", cs.stateMachineServiceURL, data.Orid), body)
	if err != nil {
		return nil, errors.New("could not build request to create state machine")
	}

	token, err := cs.authManager.GetAuthenticationToken(nil)
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
		payload := UpdateStateMachineResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			// return nil, errors.New("Could not decode response from API of resource")
			return nil, err
		}

		return &payload, nil
	default:
		body, _ := io.ReadAll(r.Body)
		return nil, fmt.Errorf("did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// DeleteStateMachineArgs Data needed to delete a state machine
type DeleteStateMachineArgs struct {
	Orid string
}

// DeleteStateMachineResult Delete state machine result
type DeleteStateMachineResult struct {
	Orid string `json:"orid"`
}

// DeleteStateMachine Attempts to delete a state machine within the MDS Cloud deployment
func (cs *StateMachineServiceClient) DeleteStateMachine(data *DeleteStateMachineArgs) (*DeleteStateMachineResult, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/machine/%s", cs.stateMachineServiceURL, data.Orid), nil)
	if err != nil {
		return nil, errors.New("could not build request to create state machine")
	}

	token, err := cs.authManager.GetAuthenticationToken(nil)
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
		payload := DeleteStateMachineResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			// return nil, errors.New("Could not decode response from API of resource")
			return nil, err
		}

		return &payload, nil
	default:
		body, _ := io.ReadAll(r.Body)
		return nil, fmt.Errorf("did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}
