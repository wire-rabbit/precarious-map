package main

import (
	"bytes"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func TestStartUI(t *testing.T) {

	var buf bytes.Buffer
	var in bytes.Buffer
	in.Write([]byte("q"))

	m := model{getTableLayout()}
	p := tea.NewProgram(m, tea.WithInput(&in), tea.WithOutput(&buf))
	if _, err := p.Run(); err != nil {
		t.Fatal(err)
	}

	if buf.Len() == 0 {
		t.Fatalf("no output (we should at least see newlines)")
	}

	expectedColumns := []string{
		"Name",
		"AZ",
		"Instance ID",
		"Public IP",
		"State",
	}
	for _, column := range expectedColumns {
		if !bytes.Contains(buf.Bytes(), []byte(column)) {
			t.Errorf("expected column %q not found", column)
		}
	}
}
