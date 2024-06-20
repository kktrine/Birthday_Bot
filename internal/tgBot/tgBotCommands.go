package tgBot

import (
	"birthday_bot/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (b *Bot) signIn(username, password string, chatId int64) (string, error) {
	credentials := Credentials{Username: username, Password: password, ChatId: chatId}
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

func (b *Bot) signUp(username, password string, chatId int64) (int, error) {
	credentials := Credentials{Username: username, Password: password, ChatId: chatId}
	body, err := json.Marshal(credentials)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(b.apiBaseURL+"/sign_up", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, errors.New("ошибка")
		}
		return 0, errors.New(string(respBody))
	}
	var response SignInResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	return response.Id, nil
}

func (b *Bot) getEmployees(token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.apiBaseURL+"/api/employees", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (b *Bot) info(info *model.Employee, token string) error {
	bodyReq, err := json.Marshal(info)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", b.apiBaseURL+"/api/info", bytes.NewBuffer(bodyReq))
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		return errors.New(resp.Status)
	}
	return nil

}
