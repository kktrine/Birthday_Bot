package tgBot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (b *Bot) signIn(username, password string) (string, error) {
	credentials := Credentials{Username: username, Password: password}
	body, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(b.apiBaseURL+"/sign_in", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return "", err
	}

	return authResponse.Token, nil
}

func (b *Bot) signUp(username, password string) error {
	credentials := Credentials{Username: username, Password: password}
	body, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	resp, err := http.Post(b.apiBaseURL+"/sign_up", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to sign up, status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *Bot) getEmployees(token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.apiBaseURL+"/api/employees", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
