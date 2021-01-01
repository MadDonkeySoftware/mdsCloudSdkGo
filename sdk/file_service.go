package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// FileServiceClient Client to interact with the MDS Cloud container service
type FileServiceClient struct {
	fileServiceURL string
	authManager    *AuthManager
}

// CreateContainerArgs Data needed to create a new container
type CreateContainerArgs struct {
	Name string `json:"name"`
}

// CreateContainerResult Create container results
type CreateContainerResult struct {
	Orid string `json:"orid"`
}

// CreateContainer Attempts to create a new container with the MDS Cloud deployment
func (cs *FileServiceClient) CreateContainer(data *CreateContainerArgs) (*CreateContainerResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/createContainer/%s", cs.fileServiceURL, data.Name), nil)
	if err != nil {
		return nil, errors.New("Could not build request to create container")
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
	case 201:
		payload := CreateContainerResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			// return nil, errors.New("Could not decode response from API of resource")
			return nil, err
		}

		return &payload, nil
	case 409:
		return nil, errors.New("container already exists")
	default:
		body, _ := ioutil.ReadAll(r.Body)
		return nil, fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// ListContainerContentsArgs Data needed to create a new container
type ListContainerContentsArgs struct {
	Orid string `json:"orid"`
}

// ListContainerContentsResult Create container results
type ListContainerContentsResult struct {
	Directories []string `json:"directories"`
	Files       []string `json:"files"`
}

// ListContainerContents Attempts to create a new container with the MDS Cloud deployment
func (cs *FileServiceClient) ListContainerContents(data *ListContainerContentsArgs) (*ListContainerContentsResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/list/%s", cs.fileServiceURL, data.Orid), nil)
	if err != nil {
		return nil, errors.New("Could not build request to create container")
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
		payload := ListContainerContentsResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			// return nil, errors.New("Could not decode response from API of resource")
			return nil, err
		}

		return &payload, nil
	default:
		body, _ := ioutil.ReadAll(r.Body)
		return nil, fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// DeleteContainerArgs Data needed to delete a container
type DeleteContainerArgs struct {
	Orid string `json:"orid"`
}

// DeleteContainerOrPath Attempts to delete a container or path within a container in the MDS Cloud deployment
func (cs *FileServiceClient) DeleteContainerOrPath(data *DeleteContainerArgs) error {
	client := &http.Client{Timeout: 10 * time.Second}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/%s", cs.fileServiceURL, data.Orid), bytes.NewBuffer(body))
	if err != nil {
		return errors.New("Could not build request to create container")
	}

	token, err := cs.authManager.GetAuthenticationToken(nil)
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
		body, _ = ioutil.ReadAll(r.Body)
		return fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}
