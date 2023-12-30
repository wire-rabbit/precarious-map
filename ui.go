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

const viewport_width = 78

const UiHelpText = `
tab selects | up/k and down/j to select or scroll | q to quit	
`

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder())

var focusedStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder())

var helpStyle = lipgloss.NewStyle().
	Width(viewport_width).
	Foreground(lipgloss.Color("#A8A8F8")).
	Align(lipgloss.Left).
	PaddingLeft(1)

type FetchFunctionType func() []InstanceDetail

type InstanceDetail struct {
	InstanceId       string
	ImageId          string
	AZ               string
	State            string
	PrivateIpAddress string
	PublicIpAddress  string
	Name             string
	JSON             string
}

type model struct {
	table            table.Model
	viewport         viewport.Model
	selectedWidget   uint
	selectedInstance uint
	instances        []InstanceDetail
	fetchFunction    FetchFunctionType
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
		if len(m.instances) > 0 {
			content := fmt.Sprintf("```json\n%s\n```", m.instances[m.table.Cursor()].JSON)
			m.viewport.SetContent(getMarkdown(content))
			m.viewport.GotoTop()
		} else {
			content := fmt.Sprintf("# AWS returned no data!")
			m.viewport.SetContent(getMarkdown(content))
			m.viewport.GotoTop()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "k":
			if m.selectedWidget == table_selected {
				if len(m.instances) > 0 {
					content := fmt.Sprintf("```json\n%s\n```", m.instances[m.table.Cursor()].JSON)
					m.viewport.SetContent(getMarkdown(content))
					m.viewport.GotoTop()
				}
			}
		case "tab":
			if m.selectedWidget == table_selected {
				m.selectedWidget = viewport_selected
			} else {
				m.selectedWidget = table_selected
			}
		}
	}

	if m.selectedWidget == table_selected {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	baseContent := ""
	if m.selectedWidget == table_selected {
		baseContent = fmt.Sprintf("%s\n\n%s", focusedStyle.Render(m.table.View()), baseStyle.Render(m.viewport.View()))
	} else {
		baseContent = fmt.Sprintf("%s\n\n%s", baseStyle.Render(m.table.View()), focusedStyle.Render(m.viewport.View()))
	}
	return baseContent + helpStyle.Render(UiHelpText)
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

func getMarkdown(content string) string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(viewport_width),
	)
	if err != nil {
		fmt.Printf("unable to get the viewport renderer: %s\n", err)
		os.Exit(1)
	}

	str, err := renderer.Render(content)
	if err != nil {
		fmt.Printf("unable to render the JSON viewport: %s\n", err)
		os.Exit(1)
	}

	return str
}

func getJsonViewport(content string) viewport.Model {
	vp := viewport.New(viewport_width, 20)
	vp.SetContent(getMarkdown(content))
	return vp
}

func startUI(options AppOptions) {

	m := model{table: getTableLayout([]InstanceDetail{}), selectedWidget: table_selected}
	m.viewport = getJsonViewport("*Waiting for data...*")
	m.fetchFunction = fetchEc2Data

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Unable to start the UI: %s\n", err)
		os.Exit(1)
	}
}
