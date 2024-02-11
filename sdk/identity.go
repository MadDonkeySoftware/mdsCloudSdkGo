package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// IdentityClient Client to interact with MDS Cloud identity service
type IdentityClient struct {
	identityURL       string
	allowSelfSignCert bool
	authManager       *AuthManager
}

// RegisterAccountArgs Data needed to register a new account
type RegisterAccountArgs struct {
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	FriendlyName string `json:"friendlyName"`
	AccountName  string `json:"accountName"`
}

// RegisterResult Result of the account registration operation
type RegisterResult struct {
	Status    string `json:"status"`
	AccountID string `json:"accountId"`
}

func (ic *IdentityClient) getHTTPClient() *http.Client {
	var client *http.Client

	if ic.allowSelfSignCert {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Timeout: API_TIMEOUT, Transport: tr}
	} else {
		client = &http.Client{Timeout: API_TIMEOUT}
	}

	return client
}

// Register Attempts to register a new account with the MDS Cloud deployment
func (ic *IdentityClient) Register(data *RegisterAccountArgs) (*RegisterResult, error) {
	client := ic.getHTTPClient()

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/register", ic.identityURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("could not build request to register user")
	}

	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 200:
		payload := RegisterResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return nil, errors.New("could not decode response from API of resource")
		}
		return &payload, nil
	default:
		body, _ = io.ReadAll(r.Body)
		return nil, fmt.Errorf("did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// AuthenticateArgs Data needed to authenticate a user
type AuthenticateArgs struct {
	AccountID string `json:"accountId"`
	UserID    string `json:"userId"`
	Password  string `json:"password"`
}

// AuthenticateResult Authentication results
type AuthenticateResult struct {
	Token string `json:"token"`
}

// Authenticate Attempts to authenticate a user against the MDS Cloud deployment
func (ic *IdentityClient) Authenticate(data *AuthenticateArgs) (*AuthenticateResult, error) {
	overrides := map[string]string{
		"accountId": data.AccountID,
		"userId":    data.UserID,
		"password":  data.Password,
	}
	token, err := ic.authManager.GetAuthenticationToken(overrides)

	if err != nil {
		return nil, err
	}

	return &AuthenticateResult{
		Token: token,
	}, nil
}

// UpdateUserArgs Data needed to update user details
type UpdateUserArgs struct {
	Email        string `json:"email"`
	OldPassword  string `json:"oldPassword"`
	NewPassword  string `json:"newPassword"`
	FriendlyName string `json:"friendlyName"`
}

// UpdateUser Attempts to update various aspects of the user
func (ic *IdentityClient) UpdateUser(data *UpdateUserArgs) error {
	client := ic.getHTTPClient()

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/updateUser", ic.identityURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if err != nil {
		return errors.New("could not build request to register user")
	}

	token, err := ic.authManager.GetAuthenticationToken(nil)
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
		body, _ = io.ReadAll(r.Body)
		return fmt.Errorf("did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// ImpersonateUserArgs Data needed to authenticate a user
type ImpersonateUserArgs struct {
	AccountID string `json:"accountId"`
}

// ImpersonateUserResult Impersonation result
type ImpersonateUserResult struct {
	Token string `json:"token"`
}

// ImpersonateUser Get impersonation token for a user on a given account
func (ic *IdentityClient) ImpersonateUser(data *ImpersonateUserArgs) (*ImpersonateUserResult, error) {
	client := ic.getHTTPClient()

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/impersonate", ic.identityURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("could not build request to register user")
	}

	token, err := ic.authManager.GetAuthenticationToken(nil)
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
		payload := ImpersonateUserResult{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return nil, errors.New("could not decode response from API of resource")
		}
		return &payload, nil
	default:
		body, _ = io.ReadAll(r.Body)
		return nil, fmt.Errorf("did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}

// PublicSignatureResponse Response object holding the public signature
type PublicSignatureResponse struct {
	Signature string `json:"signature"`
}

// GetPublicSignature Gets the active public signature from the MDS Cloud deployment
func (ic *IdentityClient) GetPublicSignature() (*PublicSignatureResponse, error) {
	client := ic.getHTTPClient()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/publicSignature", ic.identityURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 200:
		payload := PublicSignatureResponse{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return nil, errors.New("could not decode response from API of resource")
		}
		return &payload, nil
	default:
		body, _ := io.ReadAll(r.Body)
		return nil, fmt.Errorf("did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}
