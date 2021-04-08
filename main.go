package main

import (
	"flag"
	"fmt"
	"strconv"
)

var (
	version = "v1.0.0"

	settings = GetSettings()

	_Token   = flag.String("token", "", "Override Harvest Personal Token")
	_Account = flag.Int("account", 0, "Override Account ID")
)

var Actions = []Action{
	{label: "Start a new Time Entry", code: "new"},
	{label: "Restart a Time Entry", code: "restart"},
	{label: "Stop a Time Entry", code: "stop"},
	{label: "See Time Entries", code: "list-time"},
	{label: "Export Time Entries", code: "export-time"},
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

		} else if action == "export-time" {

			fmt.Printf("\n\nStart Date\n")
			start, err := askDate()
			Check(err)

			fmt.Printf("\n\nEnd Date\n")
			end, err := askDate()
			Check(err)
			fmt.Printf("\033[2J")

			params := []Param{}

			for {
				action, err := askExportAction()
				Check(err)
				fmt.Printf("\033[2J")

				if action == "filter-proj" {
					project, err := askProject(assignments)
					Check(err)
					fmt.Printf("\033[2J")

					projectParam := Param{Name: "project_id", Value: fmt.Sprint(project.Project.ID)}
					clientParam := Param{Name: "client_id", Value: fmt.Sprint(project.Client.ID)}
					params = append(params, projectParam)
					params = append(params, clientParam)

				} else if action == "export" {
					HandleExportTimeEntries(start, end, params)
					break
				} else {
					fmt.Println("Unrecognized Action")
				}
			}

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
