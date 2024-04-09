package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enchik0reo/commandApi/internal/models"
)

type response struct {
	Status int      `json:"status"`
	Body   respBody `json:"body"`
}

type respBody struct {
	CommandID          int64            `json:"command_id,omitempty"`
	CommandDescription *models.Command  `json:"command,omitempty"`
	Commands           []models.Command `json:"commands,omitempty"`
	Error              string           `json:"error,omitempty"`
}

// responseJSONOk writes to ResponseWriter status and body ...
func responseJSONOk(w http.ResponseWriter, status int, body respBody) error {
	resp := response{
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

// rresponseJSONError writes to ResponseWriter status and error ...
func responseJSONError(w http.ResponseWriter, status int, error string) error {
	resp := response{
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
