package services

import "errors"

var (
	ErrNoExecutingCommand = errors.New("there's no executing script")
	ErrStoppedManually    = errors.New("script was stopped manually")
)
