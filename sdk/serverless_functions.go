package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ServerlessFunctionsClient Client to interact with MDS Cloud serverless functions
type ServerlessFunctionsClient struct {
	serviceURL  string
	authManager *AuthManager
}

// ServerlessFunctionSummary Function summary details
type ServerlessFunctionSummary struct {
	Orid string
	Name string
}

// ServerlessFunctionDetails Data describing the serverless function
type ServerlessFunctionDetails struct {
	Orid       string
	Name       string
	Version    string
	Runtime    string
	EntryPoint string
	Created    string
	LastUpdate string
	LastInvoke string
}

// CreateFunction Create a new serverless function
func (c *ServerlessFunctionsClient) CreateFunction(name string) (*ServerlessFunctionSummary, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	token, err := c.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return nil, errors.New("Could not acquire authentication token")
	}

	body := []byte(fmt.Sprintf(`{"name":"%s"}`, name))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/create", c.serviceURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("Could not build request to create new function")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Could not execute request to create new function")
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 201:
		function := make(map[string]interface{})
		err = json.NewDecoder(r.Body).Decode(&function)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Could not decode response from API of resource")
		}
		data := ServerlessFunctionSummary{
			Name: name,
			Orid: function["orid"].(string),
		}
		return &data, nil
	case 400:
		body, _ = ioutil.ReadAll(r.Body)
		return nil, errors.New(string(body))
	case 409:
		return nil, fmt.Errorf("Function with name \"%s\" appears to already exist", name)
	default:
		body, _ = ioutil.ReadAll(r.Body)
		return nil, fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// ListFunctions List the available functions
func (c *ServerlessFunctionsClient) ListFunctions() (*[]ServerlessFunctionSummary, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	token, err := c.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return nil, errors.New("Could not acquire authentication token")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/list", c.serviceURL), nil)
	if err != nil {
		return nil, errors.New("Could not build request to fetch list of functions from API")
	}

	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Could not execute request to fetch list of functions from serverless functions API")
	}
	defer r.Body.Close()

	apiFunctions := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&apiFunctions)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not decode response from API of resource")
	}

	functions := make([]ServerlessFunctionSummary, 0)
	for _, f := range apiFunctions {
		functions = append(functions, ServerlessFunctionSummary{
			Orid: f["orid"].(string),
			Name: f["name"].(string),
		})
	}

	return &functions, nil
}

// DeleteFunction .
func (c *ServerlessFunctionsClient) DeleteFunction(orid string) error {
	client := &http.Client{Timeout: API_TIMEOUT}

	token, err := c.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return errors.New("Could not acquire authentication token")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/%s", c.serviceURL, orid), nil)
	if err != nil {
		return errors.New("Could not build request to delete function")
	}

	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return errors.New("Could not execute request to delete function")
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

// InvokeFunction .
func (c *ServerlessFunctionsClient) InvokeFunction(orid string, body interface{}) (interface{}, error) {
	client := &http.Client{Timeout: 30 * time.Minute}

	token, err := c.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return nil, errors.New("Could not acquire authentication token")
	}

	bodyBytes, _ := json.Marshal(body)
	payload := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/invoke/%s", c.serviceURL, orid), payload)
	if err != nil {
		return nil, errors.New("Could not build request to invoke function")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Could not execute request to invoke function")
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 200:
		body, _ := ioutil.ReadAll(r.Body)
		return body, nil
	case 400:
		body, _ := ioutil.ReadAll(r.Body)
		return nil, errors.New(string(body))
	default:
		body, _ := ioutil.ReadAll(r.Body)
		return nil, fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// GetFunctionDetails Gets details for a function
func (c *ServerlessFunctionsClient) GetFunctionDetails(orid string) (*ServerlessFunctionDetails, error) {
	client := &http.Client{Timeout: API_TIMEOUT}

	token, err := c.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return nil, errors.New("Could not acquire authentication token")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/inspect/%s", c.serviceURL, orid), nil)
	if err != nil {
		return nil, errors.New("Could not build request to fetch function from API")
	}

	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Could not execute request to fetch function from serverless functions API")
	}
	defer r.Body.Close()

	function := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&function)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not decode response from API of resource")
	}

	safeToString := func(data interface{}) string {
		if data == nil {
			return ""
		}
		return data.(string)
	}

	data := &ServerlessFunctionDetails{
		Orid:       safeToString(function["orid"]),
		Name:       safeToString(function["name"]),
		Version:    safeToString(function["version"]),
		Runtime:    safeToString(function["runtime"]),
		EntryPoint: safeToString(function["entryPoint"]),
		Created:    safeToString(function["created"]),
		LastUpdate: safeToString(function["lastUpdate"]),
		LastInvoke: safeToString(function["lastInvoke"]),
	}
	return data, nil
}

// UpdateFunctionCode .
type UpdateFunctionCodeArgs struct {
	Orid             string
	Runtime          string
	EntryPoint       string
	SourcePathOrFile string
	Context          string
}

func (c *ServerlessFunctionsClient) UpdateFunctionCode(data *UpdateFunctionCodeArgs) error {
	client := &http.Client{Timeout: 30 * time.Minute}

	token, err := c.authManager.GetAuthenticationToken(nil)
	if err != nil {
		return errors.New("Could not acquire authentication token")
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(data.SourcePathOrFile)
	defer file.Close()

	fileWriter, err := writer.CreateFormFile("sourceArchive", filepath.Base(data.SourcePathOrFile))
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}
	_ = writer.WriteField("runtime", data.Runtime)
	_ = writer.WriteField("entryPoint", data.EntryPoint)

	if data.Context != "" {
		_ = writer.WriteField("context", data.Context)
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/uploadCode/%s", c.serviceURL, data.Orid), payload)
	if err != nil {
		return errors.New("Could not build request to create new function")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Token", token)
	r, err := client.Do(req)
	if err != nil {
		return errors.New("Could not execute request to create new function")
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 201:
		result := make(map[string]interface{})
		err = json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			fmt.Println(err)
			return errors.New("Could not decode response from API of resource")
		}

		buildStatus := result["status"].(string)
		if buildStatus == "buildFailed" {
			return errors.New("function failed to build")
		}

		return nil
	case 400:
		body, _ := ioutil.ReadAll(r.Body)
		return errors.New(string(body))
	default:
		body, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}
