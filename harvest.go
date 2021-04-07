package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	client  = http.Client{Timeout: time.Duration(5 * time.Second)}
	baseURL = "https://api.harvestapp.com/v2/"
)

func newRequest(method string, path string, body *bytes.Buffer) *http.Request {
	req, err := http.NewRequest("GET", baseURL+path, body)
	Check(err)

	if *_Token != "" {
		req.Header.Set("Authorization", "Bearer "+*_Token)
	} else {
		req.Header.Set("Authorization", "Bearer "+settings.Token)
	}
	if *_Account != 0 {
		req.Header.Set("Harvest-Account-Id", fmt.Sprint(*_Account))
	} else {
		req.Header.Set("Harvest-Account-Id", fmt.Sprint(settings.Account))
	}
	req.Header.Set("User-Agent", "THI "+version)

	return req
}

type Settings struct {
	Token   string `yaml:"token"`
	Account int    `yaml:"account"`
	User    struct {
		ID int `yaml:"id"`
	}
}

func GetSettings() Settings {
	settingsPath := path.Join(GetUserPath(), ".harvest", "settings.yaml")

	data, err := ioutil.ReadFile(settingsPath) // #nosec G304
	if os.IsNotExist(err) {
		return Settings{}
	} else {
		var settings Settings

		err := yaml.Unmarshal(data, &settings)
		Check(err)

		return settings
	}
}

func SaveSettings() {
	yamlSettings, err := yaml.Marshal(settings)
	Check(err)
	SaveToFile(".harvest", "settings.yaml", yamlSettings)
}

type User struct {
	ID                int    `json:"id"`
	Firstname         string `json:"first_name"`
	Lastname          string `json:"last_name"`
	Email             string `json:"email"`
	Telephone         string `json:"telephone"`
	Timezone          string `json:"timezone"`
	Weekly_capacity   int    `json:"weekly_capacity"`
	IsContractor      bool   `json:"is_contractor"`
	IsAdmin           bool   `json:"is_admin"`
	IsProject_manager bool   `json:"is_project_manager"`
	IsActive          bool   `json:"is_active"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

func GetUserInfo() User {
	req := newRequest("GET", "users/me", bytes.NewBuffer(nil))

	res, err := client.Do(req)
	Check(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	Check(err)

	var user User
	err = json.Unmarshal(body, &user)
	Check(err)

	return user
}

type ProjectAssignmentsResponse struct {
	ProjectAssignments []struct {
		ID      int `json:"id"`
		Project struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"project"`
		Client struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"client"`
		TaskAssignments []struct {
			ID   int `json:"id"`
			Task struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"task"`
		} `json:"task_assignments"`
	} `json:"project_assignments"`
}

func GetProjectAssignments() ProjectAssignmentsResponse {
	req := newRequest("GET", "users/me/project_assignments", bytes.NewBuffer(nil))

	res, err := client.Do(req)
	Check(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	Check(err)

	var assignments ProjectAssignmentsResponse
	err = json.Unmarshal(body, &assignments)
	Check(err)

	return assignments
}

type TimeEntriesResponse struct {
	TimeEntries []struct {
		ID           int     `json:"id"`
		SpentDate    string  `json:"spent_date"`
		Hours        float64 `json:"hours"`
		HoursRounded float64 `json:"rounded_hours"`
		IsLocked     bool    `json:"is_locked"`
		IsClosed     bool    `json:"is_closed"`
		IsRunning    bool    `json:"is_running"`
		Client       struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"client"`
		Project struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Task struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"task"`
	} `json:"time_entries"`
}

func GetFilteredTimeEntries(params []Param) TimeEntriesResponse {
	req := newRequest("GET", "time_entries"+CreateGETParams(params), bytes.NewBuffer(nil))

	res, err := client.Do(req)
	Check(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	Check(err)

	var timeEntries TimeEntriesResponse
	err = json.Unmarshal(body, &timeEntries)
	Check(err)

	return timeEntries
}

func GetTimeEntriesToggled(isRunning bool) TimeEntriesResponse {
	t := time.Now()

	params := []Param{
		{Name: "user_id", Value: fmt.Sprint(settings.User.ID)},
		{Name: "from", Value: t.Format("2006-01-02")},
		{Name: "is_running", Value: fmt.Sprintf("%t", isRunning)},
	}

	return GetFilteredTimeEntries(params)
}

func GetTimeEntries() TimeEntriesResponse {
	t := time.Now()

	params := []Param{
		{Name: "user_id", Value: fmt.Sprint(settings.User.ID)},
		{Name: "from", Value: t.Format("2006-01-02")},
	}

	return GetFilteredTimeEntries(params)
}
