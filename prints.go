package main

import (
	"bufio"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
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

func exportTimeEntries(start time.Time, end time.Time, filters []Param, timeEntries TimeEntriesResponse) {
	var filename string

	extention := "csv"

	if ParamContainsNested(filters, "Name", "project_id") {
		safeClientName := SafeFileName(timeEntries.TimeEntries[0].Client.Name, "-")

		filename = fmt.Sprintf("harvest_%s_%s_%s.%s", start.Format("2006-01-02"), end.Format("2006-01-02"), safeClientName, extention)
	} else {
		filename = fmt.Sprintf("harvest_%s_%s.%s", start.Format("2006-01-02"), end.Format("2006-01-02"), extention)
	}

	f, err := os.Create(filename)
	Check(err)
	defer f.Close()

	w := bufio.NewWriter(f)

	fmt.Fprintf(w, "%s,%s,%s,%s,%s\n", "Client", "Project", "Task", "Date", "Hours")
	for _, e := range timeEntries.TimeEntries {
		fmt.Fprintf(w, "\"%s\",\"%s\",\"%s\",\"%s\",%.2f\n", e.Client.Name, e.Project.Name, e.Task.Name, e.SpentDate, e.Hours)
	}

	defer w.Flush()

	fmt.Printf("File %s created.\n\n", filename)
}
