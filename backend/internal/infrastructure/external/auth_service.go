package external

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type AuthClient struct {
	baseURL string
	client  *http.Client
}

func NewAuthClient(url string) *AuthClient {
	return &AuthClient{
		baseURL: url,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *AuthClient) Login(email, password string) (bool, error) {
	if c.baseURL == "" {
		return false, nil
	}

	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	resp, err := c.client.Post(c.baseURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}
