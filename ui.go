package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	table_selected = iota
	viewport_selected
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder())

var focusedStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder())

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
	viewport      viewport.Model
	selected      uint
	instances     []InstanceDetail
	fetchFunction FetchFunctionType
}

func (m model) Init() tea.Cmd {
	return getTableData(m.fetchFunction)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case FetchFunctionType:
		m.instances = msg()
		m.table = getTableLayout(m.instances)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.selected == table_selected {
				m.selected = viewport_selected
			} else {
				m.selected = table_selected
			}
		}
	}

	if m.selected == table_selected {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.selected == table_selected {
		return fmt.Sprintf("%s\n\n%s", focusedStyle.Render(m.table.View()), baseStyle.Render(m.viewport.View()))
	}
	return fmt.Sprintf("%s\n\n%s", baseStyle.Render(m.table.View()), focusedStyle.Render(m.viewport.View()))
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
		table.WithHeight(10),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return t
}

func getJsonViewport(content string) viewport.Model {
	const width = 78

	vp := viewport.New(width, 20)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		fmt.Println("unable to get the viewport renderer: ", err.Error())
		os.Exit(1)
	}

	str, err := renderer.Render(content)
	if err != nil {
		fmt.Println("unable to render the JSON viewport: ", err.Error())
		os.Exit(1)
	}

	vp.SetContent(str)

	return vp
}

func startUI(options AppOptions) {

	m := model{table: getTableLayout([]InstanceDetail{}), selected: table_selected}
	m.viewport = getJsonViewport("*Initializing...*")
	m.fetchFunction = fetchEc2Data

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Unable to start the UI:", err.Error())
		os.Exit(1)
	}
}
