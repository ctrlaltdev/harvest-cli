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

func printTimeEntries(timeEntries TimeEntriesResponse) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 4, ' ', 0)
	for _, e := range timeEntries.TimeEntries {
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%.2fhrs\t", e.Client.Name, e.Project.Name, e.Task.Name, e.HoursRounded)
	}
	defer w.Flush()
}

func main() {
	printHeader()
	flag.Parse()

	if *_Token == "" && settings.Token == "" {
		var newToken string
		fmt.Print("Harvest Personal Token: ")
		fmt.Scanln(&newToken)
		settings.Token = newToken
		SaveSettings()
	}

	if *_Account == 0 && settings.Account == 0 {
		var newAccount int
		fmt.Print("Account ID: ")
		fmt.Scanln(&newAccount)
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
		fmt.Scanln(&action)
		fmt.Printf("\033[2J")

		if action == "1" {
			timeEntries := GetTimeEntries()
			printTimeEntries(timeEntries)
		} else if action == "2" {
			assignments := GetProjectAssignments()
			printAssignments(assignments)
		} else if action == "q" || action == "Q" {
			break
		} else {
			fmt.Println("Unrecognized Action")
		}

		fmt.Printf("\n\n")
		var next bool
		fmt.Scanln(&next)
		fmt.Printf("\033[2J")
	}
}
