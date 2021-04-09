package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	client  = http.Client{Timeout: time.Duration(5 * time.Second)}
	baseURL = "https://api.harvestapp.com/v2/"
)

func newRequest(method string, path string, body *bytes.Buffer) *http.Request {
	req, err := http.NewRequest(method, baseURL+path, body)
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

func CreateTimeEntry(projectId int, taskId int) (TimeEntry, error) {
	t := time.Now()

	payload := CreateTimeEntryRequest{
		UserID:    settings.User.ID,
		ProjectID: projectId,
		TaskID:    taskId,
		SpentDate: t.Format("2006-01-02"),
	}

	jsonPayload, err := json.Marshal(payload)
	Check(err)

	req := newRequest("POST", "time_entries", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return TimeEntry{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return TimeEntry{}, err
	}

	var entry TimeEntry
	err = json.Unmarshal(body, &entry)

	if err != nil {
		return TimeEntry{}, err
	}

	return entry, nil
}

func RestartTimeEntry(id int) TimeEntry {
	req := newRequest("PATCH", "time_entries/"+fmt.Sprint(id)+"/restart", bytes.NewBuffer(nil))

	res, err := client.Do(req)
	Check(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	Check(err)

	var entry TimeEntry
	err = json.Unmarshal(body, &entry)

	Check(err)

	return entry
}

func StopTimeEntry(id int) TimeEntry {
	req := newRequest("PATCH", "time_entries/"+fmt.Sprint(id)+"/stop", bytes.NewBuffer(nil))

	res, err := client.Do(req)
	Check(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	Check(err)

	var entry TimeEntry
	err = json.Unmarshal(body, &entry)
	Check(err)

	return entry
}

func HandleTimeEntryUpdate(timeEntries TimeEntriesResponse, input string, isRunning bool) {
	index, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("You must enter a valid Task Index")
		return
	}
	index = index - 1

	var entry TimeEntry

	if index < len(Actions) && index >= 0 {
		entry = timeEntries.TimeEntries[index]
	} else {
		fmt.Println("You must enter a valid Task Index")
		return
	}

	if isRunning {
		StopTimeEntry(entry.ID)
	} else {
		RestartTimeEntry(entry.ID)
	}
}

func HandleExportTimeEntries(start time.Time, end time.Time, filters []Param) {
	filters = append(filters, Param{Name: "user_id", Value: fmt.Sprint(settings.User.ID)})
	filters = append(filters, Param{Name: "from", Value: start.Format("2006-01-02")})
	filters = append(filters, Param{Name: "to", Value: end.Format("2006-01-02")})

	timeEntries := TimeEntriesResponse{}

	res := GetFilteredTimeEntries(filters)
	timeEntries.TimeEntries = append(timeEntries.TimeEntries, res.TimeEntries...)

	for res.NextPage != nil {
		pagedFilters := filters
		pagedFilters = append(pagedFilters, Param{Name: "page", Value: fmt.Sprint(*res.NextPage)})

		res = GetFilteredTimeEntries(pagedFilters)

		timeEntries.TimeEntries = append(timeEntries.TimeEntries, res.TimeEntries...)
	}

	exportTimeEntries(start, end, filters, timeEntries)
}
