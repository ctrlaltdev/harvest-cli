package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"gopkg.in/yaml.v2"
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
		fmt.Fprintf(w, "\n %d.\t%s\t%s\t%s\t%.2fhrs\t%s\t%s\t", i+1, e.Client.Name, e.Project.Name, e.Task.Name, e.HoursRounded, state, e.Notes)
	}
	defer w.Flush()
}

func exportTimeEntries(start time.Time, end time.Time, filters []Param, timeEntries TimeEntriesResponse, extension string) {
	var filename string

	if ParamContainsNested(filters, "Name", "project_id") {
		safeClientName := SafeFileName(timeEntries.TimeEntries[0].Client.Name, "-")

		filename = fmt.Sprintf("harvest_%s_%s_%s.%s", start.Format("2006-01-02"), end.Format("2006-01-02"), safeClientName, extension)
	} else {
		filename = fmt.Sprintf("harvest_%s_%s.%s", start.Format("2006-01-02"), end.Format("2006-01-02"), extension)
	}

	f, err := os.Create(filename)
	Check(err)

	w := bufio.NewWriter(f)

	exportEntries := []TimeEntriesExport{}

	for _, e := range timeEntries.TimeEntries {
		exportEntries = append(exportEntries, TimeEntriesExport{
			ClientName:  e.Client.Name,
			ProjectName: e.Project.Name,
			TaskName:    e.Task.Name,
			SpentDate:   e.SpentDate,
			Hours:       e.Hours,
			Notes:       e.Notes,
		})
	}

	if extension == "csv" {
		fmt.Fprintf(w, "%s,%s,%s,%s,%s,%s\n", "Client", "Project", "Task", "Date", "Hours", "Notes")
		for _, e := range exportEntries {
			fmt.Fprintf(w, "\"%s\",\"%s\",\"%s\",\"%s\",%.2f,\"%s\"\n", e.ClientName, e.ProjectName, e.TaskName, e.SpentDate, e.Hours, e.Notes)
		}
	} else if extension == "json" {
		jsonPayload, err := json.Marshal(exportEntries)
		Check(err)

		fmt.Fprintf(w, "%s\n", jsonPayload)
	} else if extension == "yml" {
		yamlPayload, err := yaml.Marshal(exportEntries)
		Check(err)

		fmt.Fprintf(w, "%s\n", yamlPayload)
	} else {
		panic(fmt.Sprintf("extension unrecognized: %s", extension))
	}

	err = w.Flush()
	Check(err)

	err = f.Close()
	Check(err)

	fmt.Printf("File %s created.\n\n", filename)
}
