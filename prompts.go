package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

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
		return ProjectAssignment{}, errors.New("you must enter a valid Project Index")
	}
	index = index - 1

	var project ProjectAssignment

	if index < len(assignments.ProjectAssignments) && index >= 0 {
		project = assignments.ProjectAssignments[index]
	} else {
		return ProjectAssignment{}, errors.New("you must enter a valid Project Index")
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
		return TaskAssignment{}, errors.New("you must enter a valid Task Index")
	}
	index = index - 1

	var task TaskAssignment

	if index < len(project.TaskAssignments) && index >= 0 {
		task = project.TaskAssignments[index]
	} else {
		return TaskAssignment{}, errors.New("you must enter a valid Task Index")
	}

	return task, nil
}

func askDate() (time.Time, error) {
	fmt.Print("YYYY MM DD (type 'now' or press 'enter' for today): ")

	var args [3]string
	argsN, _ := fmt.Scanln(&args[0], &args[1], &args[2]) // #nosec G104

	if strings.ToLower(args[0]) == "now" || argsN == 0 {
		return time.Now(), nil
	}

	return time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", args[0], args[1], args[2]))
}

func askExportAction() (string, error) {
	actions := []Action{
		{label: "Export Time Entries", code: "export"},
		{label: "Filter by Project", code: "filter-proj"},
		{label: "Export Format (default: CSV)", code: "format"},
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 2, ' ', 0)
	fmt.Fprintf(w, "Please choose next action\n\n")

	for i, e := range actions {
		fmt.Fprintf(w, "\n %d.\t%s\t", i+1, e.label)
	}

	w.Flush() // #nosec G104
	fmt.Printf("\n\n")

	var input string
	fmt.Scanln(&input) // #nosec G104

	index, err := strconv.Atoi(input)
	if err != nil {
		return "", errors.New("you must enter a valid action Index")
	}
	index = index - 1

	var action string

	if index < len(actions) && index >= 0 {
		action = actions[index].code
	} else {
		return "", errors.New("you must enter a valid action Index")
	}

	return action, nil
}

func askExportFormat() (string, error) {
	formats := []Action{
		{label: "CSV", code: "csv"},
		{label: "JSON", code: "json"},
		{label: "YAML", code: "yml"},
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 2, ' ', 0)
	fmt.Fprintf(w, "Please choose export format\n\n")

	for i, e := range formats {
		fmt.Fprintf(w, "\n %d.\t%s\t", i+1, e.label)
	}

	w.Flush() // #nosec G104
	fmt.Printf("\n\n")

	var input string
	fmt.Scanln(&input) // #nosec G104

	index, err := strconv.Atoi(input)
	if err != nil {
		return "", errors.New("you must enter a valid format Index")
	}
	index = index - 1

	var format string

	if index < len(formats) && index >= 0 {
		format = formats[index].code
	} else {
		return "", errors.New("you must enter a valid format Index")
	}

	return format, nil
}
