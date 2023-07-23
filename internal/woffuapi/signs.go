package woffuapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Signs []struct {
	SignIn   bool   `json:"SignIn"`
	IP       string `json:"IP"`
	Date     string `json:"Date"`     // UTC
	TrueDate string `json:"TrueDate"` // Local Time
}

func (s Signs) HasAlreadyClockedIn() bool {
	if len(s) > 0 {
		return s[len(s)-1].SignIn
	}
	return false
}

func (w WoffuAPI) GetSigns() (Signs, error) {
	if w.auth.Type() != "TokenAuth" {
		return nil, errors.New("token authentication is required")
	}

	// Build API Request
	apiRequest := APIRequest{
		method:   "GET",
		endpoint: "/api/signs",
	}

	// Get signs
	resp, err := w.makeRequest(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting current signs: %v", resp.Status)
	}

	// Parse response
	var signs Signs
	if err := json.NewDecoder(resp.Body).Decode(&signs); err != nil {
		return nil, fmt.Errorf("error parsing current signs: %v", err)
	}

	return signs, nil
}

type SignSlots []struct {
	ID     string        `json:"$id"`
	In     SignSlotEvent `json:"In"`
	Out    SignSlotEvent `json:"Out"`
	Motive struct {
		ID   string `json:"$id"`
		Name string `json:"Name"`
	} `json:"Motive"`
}

type SignSlotEvent struct {
	ID            string `json:"$id"`
	ShortTrueTime string `json:"ShortTrueTime"`
	SignType      int    `json:"SignType"`
	SignEventId   string `json:"SignEventId"`
}

func (w WoffuAPI) GetSignSlots() (SignSlots, error) {
	if w.auth.Type() != "TokenAuth" {
		return nil, errors.New("token authentication is required")
	}

	// Build API Request
	apiRequest := APIRequest{
		method:   "GET",
		endpoint: "/api/signs/slots",
	}

	// Get signs
	resp, err := w.makeRequest(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting current sign slots: %v", resp.Status)
	}

	// Parse response
	var slots SignSlots
	if err := json.NewDecoder(resp.Body).Decode(&slots); err != nil {
		return nil, fmt.Errorf("error parsing current sign slots: %v", err)
	}

	return slots, nil
}

type SignRequest struct {
	UserId         int `json:"UserId"`
	TimezoneOffset int `json:"TimezoneOffset"`
}

func (w WoffuAPI) Sign(userId int) (err error) {
	if w.auth.Type() != "TokenAuth" {
		return errors.New("token authentication is required")
	}

	// Build request body
	signRequestBody := SignRequest{
		UserId:         userId,
		TimezoneOffset: getTimezoneOffsetInMinutes(),
	}

	body, err := json.Marshal(signRequestBody)
	if err != nil {
		return fmt.Errorf("error marshalling signRequest JSON: %v", err)
	}

	// Build API Request
	apiRequest := APIRequest{
		method:   "POST",
		endpoint: "/api/svc/signs/signs",
		headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json;charset=utf-8",
		},
		body: body,
	}

	// Sign
	resp, err := w.makeRequest(apiRequest)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error sign in/out: %v", resp.Status)
	}

	return
}
