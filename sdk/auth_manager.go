package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"gopkg.in/square/go-jose.v2/jwt"
)

// AuthManager Client to manage authorization credentials for various MDS Cloud calls
type AuthManager struct {
	cache             BaseCache
	identityURL       string
	userID            string
	password          string
	account           string
	allowSelfSignCert bool
	enableSemaphore   bool
}

func defaultIfNilOrEmpty(value interface{}, def interface{}) interface{} {
	if value != nil && value.(string) != "" {
		return value
	}
	return def
}

// NewAuthManager Creates a new AuthManager client
func NewAuthManager(identityURL string, userID string, password string, account string, allowSelfSignCert bool, enableSemaphore bool) *AuthManager {
	manager := AuthManager{
		cache:             NewInMemoryCache(),
		identityURL:       identityURL,
		userID:            userID,
		password:          password,
		account:           account,
		allowSelfSignCert: allowSelfSignCert,
		enableSemaphore:   enableSemaphore,
	}

	return &manager
}

var semaphore = make(chan int, 1)

// GetAuthenticationToken Gets an authentication token to use against the MDS apis
func (am *AuthManager) GetAuthenticationToken(overrides map[string]string) (string, error) {

	if am.enableSemaphore {
		semaphore <- 1
	}

	data, err := am.getAuthenticationTokenWork(overrides)

	if am.enableSemaphore {
		<-semaphore
	}

	return data, err
}

func (am *AuthManager) getAuthenticationTokenWork(overrides map[string]string) (string, error) {
	account := defaultIfNilOrEmpty(overrides["accountId"], am.account).(string)
	user := defaultIfNilOrEmpty(overrides["userId"], am.userID).(string)
	password := defaultIfNilOrEmpty(overrides["password"], am.password).(string)

	cacheKey := fmt.Sprintf("%s|%s|%s", am.identityURL, account, user)
	token := am.cache.Get(cacheKey)
	if token != nil {
		// Parse old token for expiration before giving it back to the caller
		var claims map[string]interface{}
		payload, err := jwt.ParseSigned(token.(string))
		if err != nil {
			return "", err
		}
		err = payload.UnsafeClaimsWithoutVerification(&claims)
		if err != nil {
			return "", err
		}

		// NOTE: Add a 60 second buffer to ensure calls will succeed.
		nowSec := time.Now().Unix() + 60
		expSec := int64(math.Floor(claims["exp"].(float64)))
		if nowSec < expSec {
			return token.(string), nil
		}
		am.cache.Remove(cacheKey)
	}

	// Acquire new token
	token, err := am.getNewToken(account, user, password)
	if err != nil {
		return "", err
	}
	am.cache.Set(cacheKey, token.(string))
	return token.(string), nil
}

func (am *AuthManager) getNewToken(account string, userName string, password string) (string, error) {
	var client *http.Client

	if am.allowSelfSignCert {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Timeout: 10 * time.Second, Transport: tr}
	} else {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	body := []byte(fmt.Sprintf(`{"accountId":"%s","userId":"%s","password":"%s"}`, account, userName, password))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/authenticate", am.identityURL), bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("Could not build request to authenticate user")
	}

	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		// return "", errors.New("Could not execute request to authenticate user")
		return "", err
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 200:
		payload := make(map[string]interface{})
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return "", errors.New("Could not decode response from API of resource")
		}
		return payload["token"].(string), nil
	default:
		body, _ = ioutil.ReadAll(r.Body)
		return "", fmt.Errorf("Did not understand response from API: %d, %s", r.StatusCode, string(body))
	}
}
