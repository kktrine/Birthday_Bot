package tgBot

import (
	"birthday_bot/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (b *Bot) signIn(username, password string, chatId int64) (string, int, error) {
	credentials := Credentials{Username: username, Password: password, ChatId: chatId}
	body, err := json.Marshal(credentials)
	if err != nil {
		return "", 0, err
	}

	resp, err := http.Post(b.apiBaseURL+"/sign_in", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var signInResp signInResponse
	if err := json.NewDecoder(resp.Body).Decode(&signInResp); err != nil {
		return "", 0, err
	}

	return signInResp.Token, signInResp.Id, nil
}

func (b *Bot) signUp(username, password string, chatId int64) error {
	credentials := Credentials{Username: username, Password: password, ChatId: chatId}
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
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.New("ошибка")
		}
		return errors.New(string(respBody))
	}

	return nil
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
	if err != nil {
		return err
	}
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

func (b *Bot) subscribe(sub model.Subscribe, token string) error {
	bodyReq, err := json.Marshal(sub)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", b.apiBaseURL+"/api/subscribe", bytes.NewBuffer(bodyReq))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

func (b *Bot) unsubscribe(unsubscribe model.Subscribe, token string) error {
	bodyReq, err := json.Marshal(unsubscribe)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", b.apiBaseURL+"/api/unsubscribe", bytes.NewBuffer(bodyReq))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}
