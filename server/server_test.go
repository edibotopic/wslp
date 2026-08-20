package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Mock implementations for testing

type mockLister struct {
	distros []string
	err     error
}

func (m *mockLister) List(ctx context.Context) ([]string, error) {
	return m.distros, m.err
}

type mockDefaultGetter struct {
	defaultDistro string
	err           error
}

func (m *mockDefaultGetter) GetDefault(ctx context.Context) (string, error) {
	return m.defaultDistro, m.err
}

type mockDefaultSetter struct {
	registered bool
	err        error
	setErr     error
}

func (m *mockDefaultSetter) IsRegistered(ctx context.Context, name string) (bool, error) {
	return m.registered, m.err
}

func (m *mockDefaultSetter) SetAsDefault(ctx context.Context, name string) error {
	return m.setErr
}

type mockUnregisterer struct {
	registered   bool
	checkErr     error
	unregErr     error
	unregistered []string
}

func (m *mockUnregisterer) IsRegistered(ctx context.Context, name string) (bool, error) {
	return m.registered, m.checkErr
}

func (m *mockUnregisterer) Unregister(ctx context.Context, name string) error {
	if m.unregErr != nil {
		return m.unregErr
	}
	m.unregistered = append(m.unregistered, name)
	return nil
}

type mockBackuper struct {
	registered bool
	exportErr  error
}

func (m *mockBackuper) IsRegistered(ctx context.Context, name string) (bool, error) {
	return m.registered, nil
}

func (m *mockBackuper) Export(ctx context.Context, distroName, outputPath string) error {
	return m.exportErr
}

type mockTerminator struct {
	registered bool
	err        error
}

func (m *mockTerminator) IsRegistered(ctx context.Context, name string) (bool, error) {
	return m.registered, nil
}

func (m *mockTerminator) Terminate(ctx context.Context, name string) error {
	return m.err
}

type mockRenamer struct {
	registered bool
	err        error
	guid       string
}

func (m *mockRenamer) IsRegistered(ctx context.Context, name string) (bool, error) {
	return m.registered, nil
}

func (m *mockRenamer) GetDistroGUID(ctx context.Context, name string) (string, error) {
	return m.guid, nil
}

func (m *mockRenamer) RenameInRegistry(guid, newName string) error {
	return m.err
}

type mockCopier struct {
	registered bool
	exportErr  error
	importErr  error
}

func (m *mockCopier) IsRegistered(ctx context.Context, name string) (bool, error) {
	return m.registered, nil
}

func (m *mockCopier) Export(ctx context.Context, distroName, outputPath string) error {
	return m.exportErr
}

func (m *mockCopier) Import(ctx context.Context, newName, tarPath, installDir string) error {
	return m.importErr
}

type mockWorkshopRunner struct {
	output []byte
	err    error
}

func (m *mockWorkshopRunner) ListWorkshops(ctx context.Context, distro string) ([]byte, error) {
	return m.output, m.err
}

type mockWorkshopController struct {
	err error
}

func (m *mockWorkshopController) RunAction(ctx context.Context, distro, project, name, action string) ([]byte, error) {
	return nil, m.err
}

// Test helpers

func testRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func parseJSONResponse(t *testing.T, body []byte, v interface{}) {
	if err := json.Unmarshal(body, v); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
}

// CORS middleware tests

func TestCORSMiddleware(t *testing.T) {
	t.Run("adds CORS headers to responses", func(t *testing.T) {
		handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		handler.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("missing Access-Control-Allow-Origin header")
		}
		if rec.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
			t.Error("missing or incorrect Access-Control-Allow-Methods header")
		}
		if rec.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
			t.Error("missing Access-Control-Allow-Headers header")
		}
	})

	t.Run("handles OPTIONS requests", func(t *testing.T) {
		handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("OPTIONS", "/", nil)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

// handleListDistros tests

func TestHandleListDistros(t *testing.T) {
	t.Run("returns 405 for non-GET methods", func(t *testing.T) {
		srv := &Server{lister: &mockLister{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/distros", nil)

		srv.handleListDistros(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns 500 on lister error", func(t *testing.T) {
		srv := &Server{lister: &mockLister{err: errors.New("mock error")}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/distros", nil)

		srv.handleListDistros(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})

	t.Run("returns JSON with distros and count", func(t *testing.T) {
		distros := []string{"Ubuntu", "Debian"}
		srv := &Server{lister: &mockLister{distros: distros}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/distros", nil)

		srv.handleListDistros(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		parseJSONResponse(t, rec.Body.Bytes(), &response)

		if response["count"].(float64) != 2 {
			t.Errorf("expected count 2, got %v", response["count"])
		}
	})
}

// handleGetDefault tests

func TestHandleGetDefault(t *testing.T) {
	t.Run("returns 405 for non-GET methods", func(t *testing.T) {
		srv := &Server{defaultGetter: &mockDefaultGetter{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/default", nil)

		srv.handleGetDefault(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns 500 on getter error", func(t *testing.T) {
		srv := &Server{defaultGetter: &mockDefaultGetter{err: errors.New("mock error")}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/default", nil)

		srv.handleGetDefault(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})

	t.Run("returns JSON with default distro", func(t *testing.T) {
		srv := &Server{defaultGetter: &mockDefaultGetter{defaultDistro: "Ubuntu"}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/default", nil)

		srv.handleGetDefault(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		parseJSONResponse(t, rec.Body.Bytes(), &response)

		if response["default"] != "Ubuntu" {
			t.Errorf("expected default 'Ubuntu', got %v", response["default"])
		}
	})
}

// handleSetDefault tests

func TestHandleSetDefault(t *testing.T) {
	t.Run("returns 405 for non-POST methods", func(t *testing.T) {
		srv := &Server{defaultSetter: &mockDefaultSetter{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/set-default", nil)

		srv.handleSetDefault(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		srv := &Server{defaultSetter: &mockDefaultSetter{}}
		rec := httptest.NewRecorder()
		req := testRequest("POST", "/api/set-default", []byte("{invalid"))

		srv.handleSetDefault(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for empty name", func(t *testing.T) {
		srv := &Server{defaultSetter: &mockDefaultSetter{}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name": ""})
		req := testRequest("POST", "/api/set-default", body)

		srv.handleSetDefault(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 500 on setter error", func(t *testing.T) {
		srv := &Server{defaultSetter: &mockDefaultSetter{setErr: errors.New("mock error")}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name": "Ubuntu"})
		req := testRequest("POST", "/api/set-default", body)

		srv.handleSetDefault(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})

	t.Run("returns success on valid request", func(t *testing.T) {
		srv := &Server{defaultSetter: &mockDefaultSetter{registered: true}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name": "Ubuntu"})
		req := testRequest("POST", "/api/set-default", body)

		srv.handleSetDefault(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		parseJSONResponse(t, rec.Body.Bytes(), &response)

		if response["success"] != true {
			t.Errorf("expected success true, got %v", response["success"])
		}
	})
}

// handleUnregister tests

func TestHandleUnregister(t *testing.T) {
	t.Run("returns 405 for non-POST methods", func(t *testing.T) {
		srv := &Server{unregisterer: &mockUnregisterer{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/unregister", nil)

		srv.handleUnregister(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		srv := &Server{unregisterer: &mockUnregisterer{}}
		rec := httptest.NewRecorder()
		req := testRequest("POST", "/api/unregister", []byte("{invalid"))

		srv.handleUnregister(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for empty distros list", func(t *testing.T) {
		srv := &Server{unregisterer: &mockUnregisterer{}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string][]string{"distros": {}})
		req := testRequest("POST", "/api/unregister", body)

		srv.handleUnregister(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns results JSON on success", func(t *testing.T) {
		srv := &Server{unregisterer: &mockUnregisterer{registered: true}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string][]string{"distros": {"Ubuntu"}})
		req := testRequest("POST", "/api/unregister", body)

		srv.handleUnregister(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		parseJSONResponse(t, rec.Body.Bytes(), &response)

		if response["results"] == nil {
			t.Error("expected results in response")
		}
	})
}

// handleTerminate tests

func TestHandleTerminate(t *testing.T) {
	t.Run("returns 405 for non-POST methods", func(t *testing.T) {
		srv := &Server{terminator: &mockTerminator{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/terminate", nil)

		srv.handleTerminate(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		srv := &Server{terminator: &mockTerminator{}}
		rec := httptest.NewRecorder()
		req := testRequest("POST", "/api/terminate", []byte("{invalid"))

		srv.handleTerminate(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for empty distros list", func(t *testing.T) {
		srv := &Server{terminator: &mockTerminator{}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string][]string{"distros": {}})
		req := testRequest("POST", "/api/terminate", body)

		srv.handleTerminate(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns results JSON on success", func(t *testing.T) {
		srv := &Server{terminator: &mockTerminator{registered: true}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string][]string{"distros": {"Ubuntu"}})
		req := testRequest("POST", "/api/terminate", body)

		srv.handleTerminate(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		parseJSONResponse(t, rec.Body.Bytes(), &response)

		if response["results"] == nil {
			t.Error("expected results in response")
		}
	})
}

// handleRename tests

func TestHandleRename(t *testing.T) {
	t.Run("returns 405 for non-POST methods", func(t *testing.T) {
		srv := &Server{renamer: &mockRenamer{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/rename", nil)

		srv.handleRename(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		srv := &Server{renamer: &mockRenamer{}}
		rec := httptest.NewRecorder()
		req := testRequest("POST", "/api/rename", []byte("{invalid"))

		srv.handleRename(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for missing oldName", func(t *testing.T) {
		srv := &Server{renamer: &mockRenamer{}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"newName": "NewUbuntu"})
		req := testRequest("POST", "/api/rename", body)

		srv.handleRename(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for missing newName", func(t *testing.T) {
		srv := &Server{renamer: &mockRenamer{}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"oldName": "Ubuntu"})
		req := testRequest("POST", "/api/rename", body)

		srv.handleRename(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})
}

// handleWorkshops tests

func TestHandleWorkshops(t *testing.T) {
	t.Run("returns 405 for non-GET methods", func(t *testing.T) {
		srv := &Server{workshopRunner: &mockWorkshopRunner{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/workshops", nil)

		srv.handleWorkshops(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for missing distro query param", func(t *testing.T) {
		srv := &Server{workshopRunner: &mockWorkshopRunner{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/workshops", nil)

		srv.handleWorkshops(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns workshops JSON on success", func(t *testing.T) {
		srv := &Server{workshopRunner: &mockWorkshopRunner{output: []byte("")}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/workshops?distro=Ubuntu", nil)

		srv.handleWorkshops(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		parseJSONResponse(t, rec.Body.Bytes(), &response)

		if response["workshops"] == nil {
			t.Error("expected workshops in response")
		}
	})
}

// handleWorkshopAction tests

func TestHandleWorkshopAction(t *testing.T) {
	t.Run("returns 405 for non-POST methods", func(t *testing.T) {
		srv := &Server{workshopController: &mockWorkshopController{}}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/workshop-action", nil)

		srv.handleWorkshopAction(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		srv := &Server{workshopController: &mockWorkshopController{}}
		rec := httptest.NewRecorder()
		req := testRequest("POST", "/api/workshop-action", []byte("{invalid"))

		srv.handleWorkshopAction(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for missing required fields", func(t *testing.T) {
		srv := &Server{workshopController: &mockWorkshopController{}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"distro": "Ubuntu"})
		req := testRequest("POST", "/api/workshop-action", body)

		srv.handleWorkshopAction(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns 400 for invalid action", func(t *testing.T) {
		srv := &Server{workshopController: &mockWorkshopController{}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{
			"distro":  "Ubuntu",
			"project": "proj",
			"name":    "ws",
			"action":  "invalid",
		})
		req := testRequest("POST", "/api/workshop-action", body)

		srv.handleWorkshopAction(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("returns success on start", func(t *testing.T) {
		srv := &Server{workshopController: &mockWorkshopController{}}
		rec := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{
			"distro":  "Ubuntu",
			"project": "proj",
			"name":    "ws",
			"action":  "start",
		})
		req := testRequest("POST", "/api/workshop-action", body)

		srv.handleWorkshopAction(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		parseJSONResponse(t, rec.Body.Bytes(), &response)

		if response["success"] != true {
			t.Errorf("expected success true, got %v", response["success"])
		}
	})
}

// handleShutdown tests

func TestHandleShutdown(t *testing.T) {
	t.Run("returns 405 for non-POST methods", func(t *testing.T) {
		srv := &Server{}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/shutdown", nil)

		srv.handleShutdown(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rec.Code)
		}
	})

	t.Run("returns success JSON on POST", func(t *testing.T) {
		srv := &Server{}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/shutdown", nil)

		srv.handleShutdown(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response map[string]interface{}
		parseJSONResponse(t, rec.Body.Bytes(), &response)

		if response["success"] != true {
			t.Errorf("expected success true, got %v", response["success"])
		}
	})
}

// Shutdown method tests

func TestShutdown(t *testing.T) {
	t.Run("returns nil when httpServer is nil", func(t *testing.T) {
		srv := &Server{}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

// NewServer constructor tests

func TestNewServer(t *testing.T) {
	t.Run("creates server with port", func(t *testing.T) {
		srv := NewServer("8080")

		if srv.port != "8080" {
			t.Errorf("expected port 8080, got %s", srv.port)
		}
	})

	t.Run("initializes with default implementations", func(t *testing.T) {
		srv := NewServer("8080")

		if srv.lister == nil {
			t.Error("lister should not be nil")
		}
		if srv.defaultGetter == nil {
			t.Error("defaultGetter should not be nil")
		}
		if srv.defaultSetter == nil {
			t.Error("defaultSetter should not be nil")
		}
		if srv.unregisterer == nil {
			t.Error("unregisterer should not be nil")
		}
		if srv.backuper == nil {
			t.Error("backuper should not be nil")
		}
		if srv.terminator == nil {
			t.Error("terminator should not be nil")
		}
		if srv.renamer == nil {
			t.Error("renamer should not be nil")
		}
		if srv.copier == nil {
			t.Error("copier should not be nil")
		}
		if srv.workshopRunner == nil {
			t.Error("workshopRunner should not be nil")
		}
		if srv.workshopController == nil {
			t.Error("workshopController should not be nil")
		}
	})

	t.Run("allows DI field injection for testing", func(t *testing.T) {
		mock := &mockLister{distros: []string{"Test"}}
		srv := NewServer("8080")
		srv.lister = mock

		if srv.lister != mock {
			t.Error("DI injection failed")
		}
	})
}
