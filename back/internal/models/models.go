package models

import (
	"os"
)

type Script struct {
	Name string
	Cmd  os.File
}

type Command struct {
	ID        int64    `json:"id"`
	Name      string   `json:"command_name"`
	StartedAt string   `json:"created_at"`
	Output    []string `json:"output,omitempty"`
	IsWorking bool     `json:"is_working"`
}
