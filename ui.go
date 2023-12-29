package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle()

type FetchFunctionType func() []InstanceDetail

type InstanceDetail struct {
	InstanceId       string
	ImageId          string
	AZ               string
	State            string
	PrivateIpAddress string
	PublicIpAddress  string
	Name             string
}

type model struct {
	table         table.Model
	instances     []InstanceDetail
	fetchFunction FetchFunctionType
}

func (m model) Init() tea.Cmd {
	return getTableData(m.fetchFunction)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FetchFunctionType:
		m.instances = msg()
		m.table = getTableLayout(m.instances)
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

func getTableData(f FetchFunctionType) tea.Cmd {
	return func() tea.Msg {
		return f
	}
}

func getTableLayout(instances []InstanceDetail) table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "AZ", Width: 10},
		{Title: "Instance ID", Width: 25},
		{Title: "Public IP", Width: 20},
		{Title: "State", Width: 10},
	}

	rows := []table.Row{}
	for _, instance := range instances {
		rows = append(rows,
			table.Row{
				instance.Name,
				instance.AZ,
				instance.InstanceId,
				instance.PublicIpAddress,
				instance.State,
			})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
	)

	return t
}

func startUI(options AppOptions) {

	m := model{table: getTableLayout([]InstanceDetail{})}
	m.fetchFunction = func() []InstanceDetail {
		return []InstanceDetail{}
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Unable to start the UI:", err.Error())
		os.Exit(1)
	}
}
