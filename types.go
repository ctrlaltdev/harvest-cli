package main

type Action struct {
	label string
	code  string
}

type Settings struct {
	Token   string `yaml:"token"`
	Account int    `yaml:"account"`
	User    struct {
		ID int `yaml:"id"`
	}
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

type TaskAssignment struct {
	ID   int `json:"id"`
	Task struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"task"`
}

type ProjectAssignment struct {
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
	TaskAssignments []TaskAssignment `json:"task_assignments"`
}

type ProjectAssignmentsResponse struct {
	ProjectAssignments []ProjectAssignment `json:"project_assignments"`
}

type TimeEntry struct {
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
}

type TimeEntriesResponse struct {
	TimeEntries  []TimeEntry `json:"time_entries"`
	Page         int         `json:"page"`
	PerPage      int         `json:"per_page"`
	TotalPages   int         `json:"total_pages"`
	TotalEntries int         `json:"total_entries"`
	NextPage     *int        `json:"next_page"`
	PreviousPage *int        `json:"previous_page"`
}

type CreateTimeEntryRequest struct {
	UserID    int    `json:"user_id"`
	ProjectID int    `json:"project_id"`
	TaskID    int    `json:"task_id"`
	SpentDate string `json:"spent_date"`
}

type TimeEntriesExport struct {
	ClientName  string  `json:"clientName" yaml:"clientName"`
	ProjectName string  `json:"projectName" yaml:"projectName"`
	TaskName    string  `json:"taskName" yaml:"taskName"`
	SpentDate   string  `json:"date" yaml:"date"`
	Hours       float64 `json:"hours" yaml:"hours"`
}
