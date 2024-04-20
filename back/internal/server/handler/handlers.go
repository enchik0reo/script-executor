package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/enchik0reo/commandApi/internal/services"
)

type createRequest struct {
	Script string `json:"script"`
}

// create godoc
// @Summary Create new command
// @Description Run new command and add it to DB
// @Tags  commands
// @Accept  json
// @Produce  json
// @Param command body createRequest true "Script for execution"
// @Success 201 {object} idRespOK "Sucess"
// @Failure 400 {object} responseErr "Bad request"
// @Failure 500 {object} responseErr "Internal server error"
// @Router /create [post]
func (h *CustomRouter) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		req := createRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if !strings.Contains(err.Error(), "EOF") {
				h.log.Debug("Can't decode body from create command request", h.log.Attr("error", err))

				err = responseJSONError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				if err != nil {
					h.log.Error("Can't make response", h.log.Attr("error", err))
				}
				return
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()

		id, err := h.cmdr.CreateNewCommand(ctx, req.Script)
		if err != nil {
			h.log.Error("Can't create new command", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		respBody := idRespBodyOK{
			CommandID: id,
		}

		if err = idRespJSONOk(w, http.StatusCreated, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}

// createUpload godoc
// @Summary Create new command from file
// @Description Run new command from file and add it to DB
// @Tags  commands
// @Accept multipart/form-data
// @Produce  json
// @Param file formData file true "Upload file"
// @Success 201 {object} idRespOK "Sucess"
// @Failure 400 {object} responseErr "Bad request"
// @Failure 500 {object} responseErr "Internal server error"
// @Router /create/upload [post]
func (h *CustomRouter) createUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		file, header, err := r.FormFile("file")
		if err != nil {
			h.log.Error("Can't get file", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			h.log.Error("Can't read from file", h.log.Attr("error", err), h.log.Attr("fileName", header.Filename))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()

		id, err := h.cmdr.CreateNewCommand(ctx, string(data))
		if err != nil {
			h.log.Error("Can't create new command", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		respBody := idRespBodyOK{
			CommandID: id,
		}

		if err = idRespJSONOk(w, http.StatusCreated, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}

// commands godoc
// @Summary Show commands
// @Description Show last n commands
// @Tags  commands
// @Accept  json
// @Produce  json
// @Param limit query int true  "Limit for commands"
// @Success 200 {object} commandsRespOK "Sucess"
// @Failure 400 {object} responseErr "Bad request"
// @Failure 500 {object} responseErr "Internal server error"
// @Router /list [get]
func (h *CustomRouter) commands() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		l := r.URL.Query().Get("limit")

		limit, err := strconv.Atoi(l)
		if err != nil {
			h.log.Debug("Can't convert limit to int", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()

		cmds, err := h.cmdr.GetCommandList(ctx, int64(limit))
		if err != nil {
			h.log.Error("Can't get list of commands", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return

		}

		respBody := commandsRespBodyOK{
			Commands: cmds,
		}

		if err = commandsRespJSONOk(w, http.StatusOK, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}

// command godoc
// @Summary Show one command
// @Description Show command description by id
// @Tags  commands
// @Accept  json
// @Produce  json
// @Param id query int true  "Command id"
// @Success 200 {object} commandRespOK "Sucess"
// @Failure 400 {object} responseErr "Bad request"
// @Failure 500 {object} responseErr "Internal server error"
// @Router /cmd [get]
func (h *CustomRouter) command() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		i := r.URL.Query().Get("id")

		id, err := strconv.Atoi(i)
		if err != nil {
			h.log.Debug("Can't convert id to int", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()

		cmd, err := h.cmdr.GetOneCommandDescription(ctx, int64(id))
		if err != nil {
			h.log.Error("Can't get command's description", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return

		}

		respBody := commandRespBodyOK{
			CommandDescription: cmd,
		}

		if err = commandRespJSONOk(w, http.StatusOK, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}

type stopCommandRequest struct {
	ID string `json:"id"`
}

// stopCommand godoc
// @Summary Stop one command
// @Description Stop command's execution by id
// @Tags  commands
// @Accept  json
// @Produce  json
// @Param id body stopCommandRequest true "Command id"
// @Success 202 {object} idRespOK "Sucess"
// @Failure 304 {object} responseErr "Not Modified"
// @Failure 400 {object} responseErr "Bad request"
// @Failure 500 {object} responseErr "Internal server error"
// @Router /stop [put]
func (h *CustomRouter) stopCommand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		req := stopCommandRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if !strings.Contains(err.Error(), "EOF") {
				h.log.Debug("Can't decode body from stop command request", h.log.Attr("error", err))

				err = responseJSONError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				if err != nil {
					h.log.Error("Can't make response", h.log.Attr("error", err))
				}
				return
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()

		id, err := strconv.Atoi(req.ID)
		if err != nil {
			h.log.Debug("Can't convert id to int", h.log.Attr("error", err))

			err = responseJSONError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			if err != nil {
				h.log.Error("Can't make response", h.log.Attr("error", err))
			}
			return
		}

		delId, err := h.cmdr.StopCommand(ctx, int64(id))
		if err != nil {
			switch {
			case errors.Is(services.ErrNoExecutingCommand, err):
				h.log.Debug("Can't stop service", h.log.Attr("command_id", req.ID), h.log.Attr("error", err))

				err = responseJSONError(w, http.StatusNotModified, http.StatusText(http.StatusNotModified))
				if err != nil {
					h.log.Error("Can't make response", h.log.Attr("error", err))
				}
				return

			default:
				h.log.Error("Can't stop command", h.log.Attr("error", err))

				err = responseJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				if err != nil {
					h.log.Error("Can't make response", h.log.Attr("error", err))
				}
				return
			}
		}

		respBody := idRespBodyOK{
			CommandID: delId,
		}

		if err = idRespJSONOk(w, http.StatusAccepted, respBody); err != nil {
			h.log.Error("Can't make response", h.log.Attr("error", err))
		}
	}
}
