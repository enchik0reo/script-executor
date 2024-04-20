package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enchik0reo/commandApi/internal/models"
)

type idRespOK struct {
	Status int          `json:"status"`
	Body   idRespBodyOK `json:"body"`
}

type idRespBodyOK struct {
	CommandID int64 `json:"command_id,omitempty"`
}

func idRespJSONOk(w http.ResponseWriter, status int, body idRespBodyOK) error {
	resp := idRespOK{
		Status: status,
		Body:   body,
	}

	w.Header().Add("Content-Type", "application/json")

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(respJSON)
	if err != nil {
		return err
	}

	return nil
}

type commandsRespOK struct {
	Status int                `json:"status"`
	Body   commandsRespBodyOK `json:"body"`
}

type commandsRespBodyOK struct {
	Commands []models.Command `json:"commands,omitempty"`
}

func commandsRespJSONOk(w http.ResponseWriter, status int, body commandsRespBodyOK) error {
	resp := commandsRespOK{
		Status: status,
		Body:   body,
	}

	w.Header().Add("Content-Type", "application/json")

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(respJSON)
	if err != nil {
		return err
	}

	return nil
}

type commandRespOK struct {
	Status int               `json:"status"`
	Body   commandRespBodyOK `json:"body"`
}

type commandRespBodyOK struct {
	CommandDescription *models.Command `json:"command,omitempty"`
}

func commandRespJSONOk(w http.ResponseWriter, status int, body commandRespBodyOK) error {
	resp := commandRespOK{
		Status: status,
		Body:   body,
	}

	w.Header().Add("Content-Type", "application/json")

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(respJSON)
	if err != nil {
		return err
	}

	return nil
}

type responseErr struct {
	Status int         `json:"status"`
	Body   respBodyErr `json:"body"`
}

type respBodyErr struct {
	Error string `json:"error,omitempty"`
}

func responseJSONError(w http.ResponseWriter, status int, error string) error {
	resp := responseErr{
		Status: status,
	}

	if error != "" {
		resp.Body.Error = error
	}

	w.Header().Add("Content-Type", "application/json")

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(respJSON)
	if err != nil {
		return err
	}

	return nil
}
