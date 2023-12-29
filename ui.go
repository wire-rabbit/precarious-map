package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle()

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return baseStyle.Render(m.table.View())
}

func getTableLayout() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "AZ", Width: 10},
		{Title: "Instance ID", Width: 25},
		{Title: "Public IP", Width: 20},
		{Title: "State", Width: 10},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
	)

	return t
}

func startUI(options AppOptions) {

	m := model{getTableLayout()}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Unable to start the UI:", err.Error())
		os.Exit(1)
	}
}
