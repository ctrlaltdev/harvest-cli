package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	te "github.com/muesli/termenv"
)

var (
	color = te.ColorProfile().Color
)

type errMsg error

type model struct {
	loading  bool
	quitting bool
	colors   bool
	spinner  spinner.Model
	err      error
}

func header() {
	header := []string{
		"  _   _                           _   ",
		" | | | | __ _ _ ____   _____  ___| |_ ",
		" | |_| |/ _\\ | '__\\ \\ / / _ \\/ __| __|",
		" |  _  | (_| | |   \\ V /  __/\\__ \\ |_ ",
		" |_| |_|\\__,_|_|    \\_/ \\___||___/\\__|",
		"                                      ",
	}

	fmt.Printf("%s\n\n", te.String(strings.Join(header, "\n")).Foreground(color("202")).String())
}

func main() {
	header()

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	s := spinner.NewModel()
	s.Spinner = spinner.MiniDot

	return model{
		loading: true,
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.Type {

		case tea.KeyHome:
			m.colors = true
			m.loading = false
			return m, nil

		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var str string

	if m.loading {
		s := te.String(m.spinner.View()).Foreground(color("202")).String()
		str = fmt.Sprintf("\n\n\t%s Loading...\n\n\tpress ESC to quit\n\n", s)
	} else if m.colors {
		for i := 1; i < 256; i++ {
			str += te.String(fmt.Sprintf("%03d\n", i)).Foreground(color(fmt.Sprintf("%03d", i))).String()
		}
	}

	if m.quitting {
		return str + "\n"
	}

	return str
}
