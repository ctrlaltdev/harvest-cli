package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

var (
	version = "v1.0.0"

	settings = GetSettings()

	_Token   = flag.String("token", "", "Override Harvest Personal Token")
	_Account = flag.Int("account", 0, "Override Account ID")
)

func printActions() {
	choices := []string{
		"Start a new Time Entry",
		"Restart a Time Entry",
		"Stop a Time Entry",
		"See Time Entries",
		"See Projects",
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 4, 2, ' ', 0)

	fmt.Fprintf(w, "What do you want to do?\n\n")
	for i, e := range choices {
		fmt.Fprintf(w, "\n %d.\t%s\t", i+1, e)
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

func printActionTimeEntries(timeEntries TimeEntriesResponse, isRunning bool) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 4, ' ', 0)
	for i, e := range timeEntries.TimeEntries {
		fmt.Fprintf(w, "\n %d.\t%s\t%s\t%s\t%.2fhrs\t", i+1, e.Client.Name, e.Project.Name, e.Task.Name, e.HoursRounded)
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

func main() {
	printHeader()
	flag.Parse()

	if *_Token == "" && settings.Token == "" {
		var newToken string
		fmt.Print("Harvest Personal Token: ")
		_, err := fmt.Scanln(&newToken)
		Check(err)
		settings.Token = newToken
		SaveSettings()
	}

	if *_Account == 0 && settings.Account == 0 {
		var newAccount int
		fmt.Print("Account ID: ")
		_, err := fmt.Scanln(&newAccount)
		Check(err)
		settings.Account = newAccount
		SaveSettings()
	}

	if settings.User.ID == 0 {
		userinfo := GetUserInfo()
		settings.User.ID = userinfo.ID
		SaveSettings()
	}

	for {
		printActions()

		var action string
		_, err := fmt.Scanln(&action)
		Check(err)
		fmt.Printf("\033[2J")

		if action == "1" {

		} else if action == "2" {
			timeEntries := GetTimeEntriesToggled(false)
			printActionTimeEntries(timeEntries, false)
		} else if action == "3" {
			timeEntries := GetTimeEntriesToggled(true)
			printActionTimeEntries(timeEntries, true)
		} else if action == "4" {
			timeEntries := GetTimeEntries()
			printTimeEntries(timeEntries)
		} else if action == "5" {
			assignments := GetProjectAssignments()
			printAssignments(assignments)
		} else if action == "q" || action == "Q" || action == "" {
			break
		} else {
			fmt.Println("Unrecognized Action")
		}

		fmt.Printf("\n\n")
		var next bool
		_, err = fmt.Scanln(&next)
		Check(err)
		fmt.Printf("\033[2J")
	}
}
