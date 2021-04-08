package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

var (
	version = "v1.0.0"

	settings = GetSettings()

	_Token   = flag.String("token", "", "Override Harvest Personal Token")
	_Account = flag.Int("account", 0, "Override Account ID")
)

type Action struct {
	label string
	code  string
}

var Actions = []Action{
	{label: "Start a new Time Entry", code: "new"},
	{label: "Restart a Time Entry", code: "restart"},
	{label: "Stop a Time Entry", code: "stop"},
	{label: "See Time Entries", code: "list-time"},
	{label: "See Projects", code: "list-proj"},
}

func translateAction(input string) string {
	index, err := strconv.Atoi(input)
	if err != nil {
		return input
	}
	index = index - 1

	if index < len(Actions) && index >= 0 {
		return Actions[index].code
	} else {
		return "Incorrect Input"
	}
}

func printActions() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 4, 2, ' ', 0)

	fmt.Fprintf(w, "What do you want to do?\n\n")
	for i, e := range Actions {
		fmt.Fprintf(w, "\n %d.\t%s\t", i+1, e.label)
	}
	fmt.Fprintf(w, "\n\n Q.\t%s\t", "Exit")
	fmt.Fprintf(w, "\n\n")
	defer w.Flush()
}

func printAssignments(assignments ProjectAssignmentsResponse) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 2, ' ', 0)
	for _, e := range assignments.ProjectAssignments {
		fmt.Fprintf(w, "\n [%s]\t%s\t%s\t", e.Project.Code, e.Client.Name, e.Project.Name)
	}
	defer w.Flush()
}

func printTimeEntries(timeEntries TimeEntriesResponse) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 4, ' ', 0)
	for i, e := range timeEntries.TimeEntries {
		var state string
		if e.IsRunning {
			state = "running"
		} else {
			state = "stopped"
		}
		fmt.Fprintf(w, "\n %d.\t%s\t%s\t%s\t%.2fhrs\t%s\t", i+1, e.Client.Name, e.Project.Name, e.Task.Name, e.HoursRounded, state)
	}
	defer w.Flush()
}

func askProject(assignments ProjectAssignmentsResponse) (ProjectAssignment, error) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 2, ' ', 0)
	fmt.Fprintf(w, "For what project?\n\n")
	for i, e := range assignments.ProjectAssignments {
		fmt.Fprintf(w, "\n %d.\t[%s]\t%s\t%s\t", i+1, e.Project.Code, e.Client.Name, e.Project.Name)
	}
	w.Flush() // #nosec G104
	fmt.Printf("\n\n")

	var input string
	fmt.Scanln(&input) // #nosec G104

	index, err := strconv.Atoi(input)
	if err != nil {
		return ProjectAssignment{}, errors.New("You must enter a valid Project Index")
	}
	index = index - 1

	var project ProjectAssignment

	if index < len(assignments.ProjectAssignments) && index >= 0 {
		project = assignments.ProjectAssignments[index]
	} else {
		return ProjectAssignment{}, errors.New("You must enter a valid Project Index")
	}

	return project, nil
}

func askTask(project ProjectAssignment) (TaskAssignment, error) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 2, ' ', 0)
	fmt.Fprintf(w, "For what task?\n\n")
	for i, e := range project.TaskAssignments {
		fmt.Fprintf(w, "\n %d.\t%s\t", i+1, e.Task.Name)
	}
	w.Flush() // #nosec G104
	fmt.Printf("\n\n")

	var input string
	fmt.Scanln(&input) // #nosec G104

	index, err := strconv.Atoi(input)
	if err != nil {
		return TaskAssignment{}, errors.New("You must enter a valid Task Index")
	}
	index = index - 1

	var task TaskAssignment

	if index < len(project.TaskAssignments) && index >= 0 {
		task = project.TaskAssignments[index]
	} else {
		return TaskAssignment{}, errors.New("You must enter a valid Task Index")
	}

	return task, nil
}

func checkSettings() {
	if *_Token == "" && settings.Token == "" {
		var newToken string
		fmt.Print("Harvest Personal Token: ")
		fmt.Scanln(&newToken) // #nosec G104
		settings.Token = newToken
		SaveSettings()
	}

	if *_Account == 0 && settings.Account == 0 {
		var newAccount int
		fmt.Print("Account ID: ")
		fmt.Scanln(&newAccount) // #nosec G104
		settings.Account = newAccount
		SaveSettings()
	}

	if settings.User.ID == 0 {
		userinfo := GetUserInfo()
		settings.User.ID = userinfo.ID
		SaveSettings()
	}
}

func main() {
	printHeader()
	flag.Parse()

	checkSettings()
	assignments := GetProjectAssignments()

	for {
		printActions()

		var input string
		fmt.Scanln(&input) // #nosec G104
		fmt.Printf("\033[2J")

		action := translateAction(input)

		if action == "new" {

			project, err := askProject(assignments)
			Check(err)
			fmt.Printf("\033[2J")

			task, err := askTask(project)
			Check(err)
			fmt.Printf("\033[2J")

			_, err = CreateTimeEntry(project.Project.ID, task.Task.ID)
			Check(err)

			fmt.Println("Time Entry Started")

		} else if action == "restart" {

			timeEntries := GetTimeEntriesToggled(false)
			printTimeEntries(timeEntries)
			fmt.Printf("\n\n")

			var input string
			fmt.Scanln(&input) // #nosec G104

			HandleTimeEntryUpdate(timeEntries, input, false)

		} else if action == "stop" {

			timeEntries := GetTimeEntriesToggled(true)
			printTimeEntries(timeEntries)
			fmt.Printf("\n\n")

			var input string
			fmt.Scanln(&input) // #nosec G104

			HandleTimeEntryUpdate(timeEntries, input, true)

		} else if action == "list-time" {

			timeEntries := GetTimeEntries()
			printTimeEntries(timeEntries)

		} else if action == "list-proj" {

			printAssignments(assignments)

		} else if action == "q" || action == "Q" || action == "" {

			break

		} else {

			fmt.Println("Unrecognized Action")

		}

		fmt.Printf("\n\n")
		var next bool
		fmt.Scanln(&next) // #nosec G104
		fmt.Printf("\033[2J")
	}
}
