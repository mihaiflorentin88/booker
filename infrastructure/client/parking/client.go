package parking

import (
	"booking/port/config"
	"booking/port/dto"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const LoginURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key="
const BookURL = "https://us-central1-project-3687381701726997562.cloudfunctions.net/reserve3"

type ParkingClient struct {
	credentials *config.Credentials
	accessToken string
}

func NewParkingClient(credentials *config.Credentials) *ParkingClient {
	return &ParkingClient{credentials: credentials}
}

func (p *ParkingClient) Login() error {
	payload, err := NewLoginPayload(p.credentials.Username, p.credentials.Password).ToJson()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", LoginURL+p.credentials.GoogleAPIKey, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create login request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(p.credentials.Username, p.credentials.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status: %v, body: %v", resp.Status, string(body))
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}

	accessToken, ok := result["idToken"].(string)
	if !ok {
		return errors.New("login response did not contain access token")
	}

	p.accessToken = accessToken
	return nil
}

func (p *ParkingClient) Book(payload *dto.BookParkingPayload) error {
	payloadBytes, err := payload.ToJson()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", BookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create the booking request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", p.GetAccessToken())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute booking request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status: %v, body: %v", resp.Status, string(body))
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}
	return nil
}

func (p *ParkingClient) GetAccessToken() string {
	return "Bearer " + p.accessToken
}
