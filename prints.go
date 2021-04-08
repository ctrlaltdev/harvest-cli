package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

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

func exportTimeEntries(timeEntries TimeEntriesResponse) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', 0)
	fmt.Fprintf(w, "\n%s,%s,%s,%s,%s", "Client", "Project", "Task", "Date", "Hours")
	for _, e := range timeEntries.TimeEntries {
		fmt.Fprintf(w, "\n\"%s\",\"%s\",\"%s\",\"%s\",%.2f", e.Client.Name, e.Project.Name, e.Task.Name, e.SpentDate, e.Hours)
	}
	defer w.Flush()
}
