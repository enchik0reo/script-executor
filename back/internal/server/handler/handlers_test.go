package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/enchik0reo/commandApi/internal/logs"
	"github.com/enchik0reo/commandApi/internal/models"
	"github.com/enchik0reo/commandApi/internal/server/handler/mocks"
	"github.com/enchik0reo/commandApi/internal/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCustomRouter_create(t *testing.T) {
	type want struct {
		resBody string
		status  int
	}
	type fields struct {
		Commander *mocks.Commander
	}
	tests := []struct {
		name            string
		want            want
		script          string
		fields          fields
		calledCommander bool
		prepare         func(fields fields, reqBody createRequest) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "test_1, OK",
			want: want{
				resBody: `{"status":201,"body":{"command_id":1}}`,
				status:  http.StatusOK,
			},
			calledCommander: true,
			script:          "whoami",
			prepare: func(fields fields, reqBody createRequest) (*httptest.ResponseRecorder, *http.Request) {
				body, _ := json.Marshal(reqBody)

				req := httptest.NewRequest("POST", "/create", strings.NewReader(string(body)))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("CreateNewCommand", mock.Anything, reqBody.Script).Return(int64(1), nil)

				return rr, req
			},
		},
		{
			name: "test_2, BadRequest",
			want: want{
				resBody: `{"status":400,"body":{"error":"Bad Request"}}`,
				status:  http.StatusOK,
			},
			prepare: func(fields fields, reqBody createRequest) (*httptest.ResponseRecorder, *http.Request) {
				badBody := fmt.Sprintf("{%s{}", reqBody.Script)

				req := httptest.NewRequest("POST", "/create", strings.NewReader(badBody))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				return rr, req
			},
		},
		{
			name: "test_3, InternalServerError",
			want: want{
				resBody: `{"status":500,"body":{"error":"Internal Server Error"}}`,
				status:  http.StatusOK,
			},
			calledCommander: true,
			script:          "whoami",
			prepare: func(fields fields, reqBody createRequest) (*httptest.ResponseRecorder, *http.Request) {
				body, _ := json.Marshal(reqBody)

				req := httptest.NewRequest("POST", "/create", strings.NewReader(string(body)))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("CreateNewCommand", mock.Anything, reqBody.Script).
					Return(int64(0), errors.New("some error in commander"))

				return rr, req
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Commander := new(mocks.Commander)
			tt.fields.Commander = Commander

			dlog := logs.NewDiscardLogger()

			router := &CustomRouter{
				cmdr:    Commander,
				timeout: 10 * time.Second,
				log:     dlog,
			}

			handler := router.create()

			reqBody := createRequest{
				Script: tt.script,
			}

			rr, req := tt.prepare(tt.fields, reqBody)

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.want.status, rr.Code)
			require.Equal(t, tt.want.resBody, rr.Body.String())

			if tt.calledCommander {
				if !Commander.AssertCalled(t, "CreateNewCommand", mock.Anything, reqBody.Script) {
					t.Errorf("Expected call Commander")
				}
			}
		})
	}
}

func TestCustomRouter_createUpload(t *testing.T) {
	type want struct {
		resBody string
		status  int
	}
	type fields struct {
		Commander *mocks.Commander
	}
	tests := []struct {
		name            string
		want            want
		fields          fields
		calledCommander bool
		prepare         func(fields fields) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "test_1, OK",
			want: want{
				resBody: `{"status":201,"body":{"command_id":1}}`,
				status:  http.StatusOK,
			},
			calledCommander: true,
			prepare: func(fields fields) (*httptest.ResponseRecorder, *http.Request) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.txt")
				part.Write([]byte("whoami"))
				writer.Close()

				req := httptest.NewRequest("POST", "/create/upload", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				rr := httptest.NewRecorder()

				fields.Commander.On("CreateNewCommand", mock.Anything, mock.Anything).Return(int64(1), nil)

				return rr, req
			},
		},
		{
			name: "test_2, InternalServerError",
			want: want{
				resBody: `{"status":500,"body":{"error":"Internal Server Error"}}`,
				status:  http.StatusOK,
			},
			prepare: func(fields fields) (*httptest.ResponseRecorder, *http.Request) {
				badBody := bytes.NewReader([]byte("test data"))

				req := httptest.NewRequest("POST", "/create/upload", badBody)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				return rr, req
			},
		},
		{
			name: "test_3, InternalServerError",
			want: want{
				resBody: `{"status":500,"body":{"error":"Internal Server Error"}}`,
				status:  http.StatusOK,
			},
			calledCommander: true,
			prepare: func(fields fields) (*httptest.ResponseRecorder, *http.Request) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.txt")
				part.Write([]byte("whoami"))
				writer.Close()

				req := httptest.NewRequest("POST", "/create/upload", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				rr := httptest.NewRecorder()

				fields.Commander.On("CreateNewCommand", mock.Anything, mock.Anything).
					Return(int64(0), errors.New("some error"))

				return rr, req
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Commander := new(mocks.Commander)
			tt.fields.Commander = Commander

			dlog := logs.NewDiscardLogger()

			router := &CustomRouter{
				cmdr:    Commander,
				timeout: 10 * time.Second,
				log:     dlog,
			}

			handler := router.createUpload()

			rr, req := tt.prepare(tt.fields)

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.want.status, rr.Code)
			require.Equal(t, tt.want.resBody, rr.Body.String())

			if tt.calledCommander {
				if !Commander.AssertCalled(t, "CreateNewCommand", mock.Anything, mock.Anything) {
					t.Errorf("Expected call Commander")
				}
			}
		})
	}
}

/* func TestCreateUploadHandler(t *testing.T) {
	// Setup
	mockCommander := new(mocks.Commander)
	mockCommander.On("CreateNewCommand", mock.Anything, mock.Anything).Return("123", nil)

	r := chi.NewRouter()
	h := &CustomRouter{cmdr: mockCommander, timeout: 5 * time.Second}
	r.Post("/upload", h.createUpload())

	t.Run("successful upload", func(t *testing.T) {
		// Prepare request
		body := bytes.NewReader([]byte("test data"))
		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", "multipart/form-data")

		// Perform request
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Check response status code
		require.Equal(t, http.StatusCreated, w.Code)

		// Check response body
		expectedBody := `{"command_id":"123"}`
		require.Equal(t, expectedBody, w.Body.String())

		// Check interactions with mock
		mockCommander.AssertCalled(t, "CreateNewCommand", mock.Anything, "test data")
	})

	t.Run("failed upload", func(t *testing.T) {
		// Prepare request
		req := httptest.NewRequest("POST", "/upload", nil)

		// Perform request
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Check response status code
		require.Equal(t, http.StatusInternalServerError, w.Code)

		// Check interactions with mock
		mockCommander.AssertNotCalled(t, "CreateNewCommand")
	})
} */

func TestCustomRouter_commands(t *testing.T) {
	type want struct {
		resBody string
		status  int
	}
	type fields struct {
		Commander *mocks.Commander
	}
	tests := []struct {
		name            string
		want            want
		limit           int64
		fields          fields
		calledCommander bool
		prepare         func(fields fields, limit int64) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "test_1, OK",
			want: want{
				resBody: "{\"status\":200,\"body\":{\"commands\":" +
					"[{\"id\":1,\"command_name\":\"whoami\",\"created_at\":\"8:40PM\",\"output\":[\"root\"],\"is_working\":false}," +
					"{\"id\":2,\"command_name\":\"sleep 1000\",\"created_at\":\"8:40PM\",\"is_working\":true}]}}",
				status: http.StatusOK,
			},
			calledCommander: true,
			limit:           2,
			prepare: func(fields fields, limit int64) (*httptest.ResponseRecorder, *http.Request) {

				req := httptest.NewRequest("GET", fmt.Sprintf("/list?limit=%d", limit), nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("GetCommandList", mock.Anything, limit).Return([]models.Command{
					{
						ID:        1,
						Name:      "whoami",
						StartedAt: "8:40PM",
						Output:    []string{"root"},
						IsWorking: false,
					},
					{
						ID:        2,
						Name:      "sleep 1000",
						StartedAt: "8:40PM",
						IsWorking: true,
					},
				}, nil)

				return rr, req
			},
		},
		{
			name: "test_2, BadRequest",
			want: want{
				resBody: `{"status":400,"body":{"error":"Bad Request"}}`,
				status:  http.StatusOK,
			},
			prepare: func(fields fields, limit int64) (*httptest.ResponseRecorder, *http.Request) {

				req := httptest.NewRequest("GET", "/list", nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				return rr, req
			},
		},
		{
			name: "test_3, InternalServerError",
			want: want{
				resBody: `{"status":500,"body":{"error":"Internal Server Error"}}`,
				status:  http.StatusOK,
			},
			limit:           2,
			calledCommander: true,
			prepare: func(fields fields, limit int64) (*httptest.ResponseRecorder, *http.Request) {

				req := httptest.NewRequest("GET", fmt.Sprintf("/list?limit=%d", limit), nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("GetCommandList", mock.Anything, limit).Return(nil, errors.New("some error"))

				return rr, req
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Commander := new(mocks.Commander)
			tt.fields.Commander = Commander

			dlog := logs.NewDiscardLogger()

			router := &CustomRouter{
				cmdr:    Commander,
				timeout: 10 * time.Second,
				log:     dlog,
			}

			handler := router.commands()

			rr, req := tt.prepare(tt.fields, tt.limit)

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.want.status, rr.Code)
			require.Equal(t, tt.want.resBody, rr.Body.String())

			if tt.calledCommander {
				if !Commander.AssertCalled(t, "GetCommandList", mock.Anything, tt.limit) {
					t.Errorf("Expected call Commander")
				}
			}
		})
	}
}

func TestCustomRouter_command(t *testing.T) {
	type want struct {
		resBody string
		status  int
	}
	type fields struct {
		Commander *mocks.Commander
	}
	tests := []struct {
		name            string
		want            want
		id              int64
		fields          fields
		calledCommander bool
		prepare         func(fields fields, id int64) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "test_1, OK",
			want: want{
				resBody: "{\"status\":200,\"body\":{\"command\":" +
					"{\"id\":1,\"command_name\":\"whoami\",\"created_at\":\"8:40PM\",\"output\":[\"root\"],\"is_working\":false}}}",
				status: http.StatusOK,
			},
			calledCommander: true,
			id:              1,
			prepare: func(fields fields, id int64) (*httptest.ResponseRecorder, *http.Request) {

				req := httptest.NewRequest("GET", fmt.Sprintf("/cmd?id=%d", id), nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("GetOneCommandDescription", mock.Anything, id).Return(&models.Command{
					ID:        1,
					Name:      "whoami",
					StartedAt: "8:40PM",
					Output:    []string{"root"},
					IsWorking: false,
				}, nil)

				return rr, req
			},
		},
		{
			name: "test_2, BadRequest",
			want: want{
				resBody: `{"status":400,"body":{"error":"Bad Request"}}`,
				status:  http.StatusOK,
			},
			prepare: func(fields fields, id int64) (*httptest.ResponseRecorder, *http.Request) {

				req := httptest.NewRequest("GET", "/cmd", nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				return rr, req
			},
		},
		{
			name: "test_3, InternalServerError",
			want: want{
				resBody: `{"status":500,"body":{"error":"Internal Server Error"}}`,
				status:  http.StatusOK,
			},
			id:              1,
			calledCommander: true,
			prepare: func(fields fields, id int64) (*httptest.ResponseRecorder, *http.Request) {

				req := httptest.NewRequest("GET", fmt.Sprintf("/cmd?id=%d", id), nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("GetOneCommandDescription", mock.Anything, id).
					Return(nil, errors.New("some error"))

				return rr, req
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Commander := new(mocks.Commander)
			tt.fields.Commander = Commander

			dlog := logs.NewDiscardLogger()

			router := &CustomRouter{
				cmdr:    Commander,
				timeout: 10 * time.Second,
				log:     dlog,
			}

			handler := router.command()

			rr, req := tt.prepare(tt.fields, tt.id)

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.want.status, rr.Code)
			require.Equal(t, tt.want.resBody, rr.Body.String())

			if tt.calledCommander {
				if !Commander.AssertCalled(t, "GetOneCommandDescription", mock.Anything, tt.id) {
					t.Errorf("Expected call Commander")
				}
			}
		})
	}
}

func TestCustomRouter_stopCommand(t *testing.T) {
	type want struct {
		resBody string
		status  int
	}
	type fields struct {
		Commander *mocks.Commander
	}
	tests := []struct {
		name            string
		want            want
		id              string
		fields          fields
		calledCommander bool
		prepare         func(fields fields, reqBody stopCommandRequest) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "test_1, OK",
			want: want{
				resBody: `{"status":202,"body":{"command_id":1}}`,
				status:  http.StatusOK,
			},
			calledCommander: true,
			id:              "1",
			prepare: func(fields fields, reqBody stopCommandRequest) (*httptest.ResponseRecorder, *http.Request) {
				body, _ := json.Marshal(reqBody)

				req := httptest.NewRequest("PUT", "/stop", strings.NewReader(string(body)))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("StopCommand", mock.Anything, mock.Anything).Return(int64(1), nil)

				return rr, req
			},
		},
		{
			name: "test_2, BadRequest",
			want: want{
				resBody: `{"status":400,"body":{"error":"Bad Request"}}`,
				status:  http.StatusOK,
			},
			id: "1",
			prepare: func(fields fields, reqBody stopCommandRequest) (*httptest.ResponseRecorder, *http.Request) {
				badBody := fmt.Sprintf("{%s{}", reqBody.ID)

				req := httptest.NewRequest("PUT", "/stop", strings.NewReader(badBody))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				return rr, req
			},
		},
		{
			name: "test_3, BadRequest",
			want: want{
				resBody: `{"status":400,"body":{"error":"Bad Request"}}`,
				status:  http.StatusOK,
			},
			id: "invalid",
			prepare: func(fields fields, reqBody stopCommandRequest) (*httptest.ResponseRecorder, *http.Request) {
				body, _ := json.Marshal(reqBody)

				req := httptest.NewRequest("PUT", "/stop", strings.NewReader(string(body)))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				return rr, req
			},
		},
		{
			name: "test_4, StatusNotModified",
			want: want{
				resBody: `{"status":304,"body":{"error":"Not Modified"}}`,
				status:  http.StatusOK,
			},
			calledCommander: true,
			id:              "1",
			prepare: func(fields fields, reqBody stopCommandRequest) (*httptest.ResponseRecorder, *http.Request) {
				body, _ := json.Marshal(reqBody)

				req := httptest.NewRequest("PUT", "/stop", strings.NewReader(string(body)))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("StopCommand", mock.Anything, mock.Anything).
					Return(int64(0), services.ErrNoExecutingCommand)

				return rr, req
			},
		},
		{
			name: "test_5, InternalServerError",
			want: want{
				resBody: `{"status":500,"body":{"error":"Internal Server Error"}}`,
				status:  http.StatusOK,
			},
			calledCommander: true,
			id:              "1",
			prepare: func(fields fields, reqBody stopCommandRequest) (*httptest.ResponseRecorder, *http.Request) {
				body, _ := json.Marshal(reqBody)

				req := httptest.NewRequest("PUT", "/stop", strings.NewReader(string(body)))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				fields.Commander.On("StopCommand", mock.Anything, mock.Anything).
					Return(int64(0), errors.New("some error"))

				return rr, req
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Commander := new(mocks.Commander)
			tt.fields.Commander = Commander

			dlog := logs.NewDiscardLogger()

			router := &CustomRouter{
				cmdr:    Commander,
				timeout: 10 * time.Second,
				log:     dlog,
			}

			handler := router.stopCommand()

			reqBody := stopCommandRequest{
				ID: tt.id,
			}

			rr, req := tt.prepare(tt.fields, reqBody)

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.want.status, rr.Code)
			require.Equal(t, tt.want.resBody, rr.Body.String())

			if tt.calledCommander {
				if !Commander.AssertCalled(t, "StopCommand", mock.Anything, mock.Anything) {
					t.Errorf("Expected call Commander")
				}
			}
		})
	}
}
