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

	mockInstances := []InstanceDetail{
		{
			InstanceId:       "a-instance-id-123",
			ImageId:          "a-image-id-456",
			AZ:               "us-east-1a",
			State:            "running",
			PrivateIpAddress: "10.0.0.23",
			PublicIpAddress:  "76.139.15.25",
			Name:             "Instance A",
		},
		{
			InstanceId:       "b-instance-id-123",
			ImageId:          "b-image-id-456",
			AZ:               "us-east-1b",
			State:            "terminated",
			PrivateIpAddress: "10.0.0.24",
			PublicIpAddress:  "76.139.15.26",
			Name:             "Instance B",
		},
	}

	expectedColumns := []string{
		"Name",
		"AZ",
		"Instance ID",
		"Public IP",
		"State",
	}

	m := model{
		table: getTableLayout([]InstanceDetail{}),
		fetchFunction: func() []InstanceDetail {
			return mockInstances
		},
	}

	p := tea.NewProgram(m, tea.WithInput(&in), tea.WithOutput(&buf))
	if _, err := p.Run(); err != nil {
		t.Fatal(err)
	}

	if buf.Len() == 0 {
		t.Fatalf("no output (we should at least see newlines)")
	}

	for _, column := range expectedColumns {
		if !bytes.Contains(buf.Bytes(), []byte(column)) {
			t.Errorf("expected column %q not found", column)
		}
	}

	for _, instance := range mockInstances {
		if !bytes.Contains(buf.Bytes(), []byte(instance.Name)) {
			t.Errorf("missing instance name: %q", instance.Name)
		}
	}
}
