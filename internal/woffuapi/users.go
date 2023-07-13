package woffuapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type User struct {
	UserId        int     `json:"UserId"`
	Email         string  `json:"Email"`
	FirstName     string  `json:"FirstName"`
	CompanyId     string  `json:"CompanyId"`
	CompanyName   string  `json:"CompanyName"`
	AllocatedDays float64 `json:"AllocatedDays"`
	ApprovedDays  float64 `json:"ApprovedDays"`
	UsedDays      float64 `json:"UsedDays"`
}

func (w WoffuAPI) GetCurrentUser() (*User, error) {
	if w.auth.Type() != "TokenAuth" {
		return nil, errors.New("token authentication is required")
	}

	// Build API Request
	apiRequest := APIRequest{
		method:   "GET",
		endpoint: "/api/users",
	}

	// Get user
	resp, err := w.makeRequest(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting current user: %v", resp.Status)
	}

	// Parse response
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error parsing current user: %v", err)
	}

	return &user, nil
}

type WorkDay struct {
	StartTime         string  `json:"StartTime"`
	TrueStartTime     string  `json:"TrueStartTime"`
	EndTime           string  `json:"EndTime"`
	TrueEndTime       string  `json:"TrueEndTime"`
	ScheduleHours     float64 `json:"ScheduleHours"`
	TrueScheduleHours float64 `json:"TrueScheduleHours"`
	IsWeekend         bool    `json:"IsWeekend"`
	IsHoliday         bool    `json:"IsHoliday"`
	IsEvent           bool    `json:"IsEvent"`
	IsFlexible        bool    `json:"IsFlexible"`
}

func (w WoffuAPI) GetUserWorkDay(userId int) (*WorkDay, error) {
	if w.auth.Type() != "TokenAuth" {
		return nil, errors.New("token authentication is required")
	}

	// Build API Request
	apiRequest := APIRequest{
		method:   "GET",
		endpoint: fmt.Sprintf("/api/users/%d/workdaylite", userId),
	}

	// Get user
	resp, err := w.makeRequest(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting current user: %v", resp.Status)
	}

	// Parse response
	var workday WorkDay
	if err := json.NewDecoder(resp.Body).Decode(&workday); err != nil {
		return nil, fmt.Errorf("error parsing workday: %v", err)
	}

	return &workday, nil
}
