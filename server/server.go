package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"wslp/internal/config"
	"wslp/internal/wsl"
)

type Server struct {
	port string
	// HTTP server instance (nil until Start() is called)
	httpServer *http.Server

	// Optional dependency injection fields for testing. When nil, defaults
	// are used (matching the same Real* implementations used in production).
	// Tests can inject mocks via New*Server constructors or direct field assignment.
	lister             wsl.Lister
	defaultGetter      wsl.DefaultGetter
	defaultSetter      wsl.DefaultSetter
	unregisterer       wsl.Unregisterer
	backuper           wsl.Backuper
	terminator         wsl.Terminator
	renamer            wsl.Renamer
	copier             wsl.Copier
	workshopRunner     wsl.WorkshopRunner
	workshopController wsl.WorkshopController
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
		// DI defaults - same Real* implementations used today
		lister:             wsl.RealLister{},
		defaultGetter:      wsl.RealDefaultGetter{},
		defaultSetter:      wsl.RealDefaultSetter{},
		unregisterer:       wsl.RealUnregisterer{},
		backuper:           wsl.RealBackuper{},
		terminator:         wsl.RealTerminator{},
		renamer:            wsl.RealRenamer{},
		copier:             wsl.RealCopier{},
		workshopRunner:     wsl.RealWorkshopRunner{},
		workshopController: wsl.RealWorkshopController{},
	}
}

// Start runs the HTTP server, blocking until it is shut down via Shutdown
// (triggered by an OS signal / Ctrl+C in cmd/serve.go, or via the
// /api/shutdown endpoint used by the GUI when its window closes so the
// server it depends on doesn't linger). It returns nil on a graceful
// shutdown, or an error if the server failed to start/run for any other
// reason.
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/distros", s.handleListDistros)
	mux.HandleFunc("/api/default", s.handleGetDefault)
	mux.HandleFunc("/api/available", s.handleListAvailable)
	mux.HandleFunc("/api/install", s.handleInstall)
	mux.HandleFunc("/api/unregister", s.handleUnregister)
	mux.HandleFunc("/api/set-default", s.handleSetDefault)
	mux.HandleFunc("/api/backup", s.handleBackup)
	mux.HandleFunc("/api/terminate", s.handleTerminate)
	mux.HandleFunc("/api/launch", s.handleLaunch)
	mux.HandleFunc("/api/rename", s.handleRename)
	mux.HandleFunc("/api/copy", s.handleCopy)
	mux.HandleFunc("/api/ubuntu-telemetry", s.handleUbuntuTelemetry)
	mux.HandleFunc("/api/wsl-info", s.handleWSLInfo)
	mux.HandleFunc("/api/distro-info", s.handleDistroInfo)
	mux.HandleFunc("/api/workshops", s.handleWorkshops)
	mux.HandleFunc("/api/workshop-action", s.handleWorkshopAction)
	mux.HandleFunc("/api/workshop-shell", s.handleWorkshopShell)
	mux.HandleFunc("/api/shutdown", s.handleShutdown)

	// Add CORS middleware for Flutter
	handler := corsMiddleware(mux)

	addr := fmt.Sprintf(":%s", s.port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	fmt.Printf("Starting server on http://localhost%s\n", addr)
	err := s.httpServer.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Shutdown gracefully stops the running server, letting in-flight requests
// finish (bounded by ctx). Safe to call even if the server hasn't finished
// starting yet.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

// handleShutdown gracefully stops the server. Intended to be called
// automatically by the GUI when its window closes (not exposed as a
// user-facing button — the GUI has no manual "stop server" control), so
// the server doesn't keep running after the app that depends on it has
// exited. The HTTP response is written before shutdown begins so the
// caller reliably sees success.
func (s *Server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Server shutting down",
	})
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(ctx)
	}()
}

func (s *Server) handleListDistros(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	distros, err := wsl.ListDistros(context.Background(), s.lister)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"distros": distros,
		"count":   len(distros),
	})
}

func (s *Server) handleGetDefault(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defaultDistro, err := wsl.GetDefaultDistro(context.Background(), s.defaultGetter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"default": defaultDistro,
	})
}

func (s *Server) handleListAvailable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	distros, err := wsl.GetAvailableDistros(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"available": distros,
		"count":     len(distros),
	})
}

func (s *Server) handleInstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Distros []string `json:"distros"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Distros) == 0 {
		http.Error(w, "No distros specified", http.StatusBadRequest)
		return
	}

	results := wsl.InstallDistros(context.Background(), request.Distros, false)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
}

func (s *Server) handleUnregister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Distros []string `json:"distros"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Distros) == 0 {
		http.Error(w, "No distros specified", http.StatusBadRequest)
		return
	}

	results := wsl.UnregisterDistros(context.Background(), s.unregisterer, request.Distros)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
}

func (s *Server) handleSetDefault(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "No distro name specified", http.StatusBadRequest)
		return
	}

	if err := wsl.SetDefaultDistro(context.Background(), request.Name, s.defaultSetter); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Successfully set %s as default", request.Name),
	})
}

func (s *Server) handleBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Distros    []string `json:"distros"`
		CustomName string   `json:"customName,omitempty"`
		BackupDir  string   `json:"backupDir,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Distros) == 0 {
		http.Error(w, "No distros specified", http.StatusBadRequest)
		return
	}

	// Validate custom name usage
	if request.CustomName != "" && len(request.Distros) > 1 {
		http.Error(w, "Custom name can only be used when backing up a single distribution", http.StatusBadRequest)
		return
	}

	// Determine backup directory
	backupDir := request.BackupDir
	if backupDir == "" {
		backupDir = config.GetBackupDir()
	}

	// Ensure backup directory exists
	if err := config.EnsureBackupDir(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create backup directory: %v", err), http.StatusInternalServerError)
		return
	}

	opts := wsl.BackupOptions{
		CustomName: request.CustomName,
	}

	results := wsl.BackupDistros(context.Background(), s.backuper, request.Distros, backupDir, opts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
}

func (s *Server) handleTerminate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Distros []string `json:"distros"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Distros) == 0 {
		http.Error(w, "No distros specified", http.StatusBadRequest)
		return
	}

	results := wsl.TerminateDistros(context.Background(), s.terminator, request.Distros)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "No distro name specified", http.StatusBadRequest)
		return
	}

	// Launch in terminal (non-blocking)
	if err := wsl.LaunchInTerminal(context.Background(), request.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Launched %s in new terminal", request.Name),
	})
}

func (s *Server) handleRename(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.OldName == "" || request.NewName == "" {
		http.Error(w, "Both old and new names are required", http.StatusBadRequest)
		return
	}

	result := wsl.RenameDistro(context.Background(), s.renamer, request.OldName, request.NewName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleWSLInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info, err := wsl.GetWSLSystemInfo(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func (s *Server) handleDistroInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "No distro name specified", http.StatusBadRequest)
		return
	}

	info, err := wsl.GetDistroDetailInfo(context.Background(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// handleWorkshops reports the Canonical Workshop (canonical/workshop)
// environments running inside a distro, if any. Workshop is optional
// third-party tooling, so this endpoint never errors when it's absent —
// it just reports zero workshops.
func (s *Server) handleWorkshops(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("distro")
	if name == "" {
		http.Error(w, "No distro name specified", http.StatusBadRequest)
		return
	}

	workshops := wsl.GetWorkshops(context.Background(), name, s.workshopRunner)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"workshops": workshops,
	})
}

// handleWorkshopAction starts or stops a single Workshop environment inside
// a distro (blocking; workshop start/stop typically complete in a few
// seconds).
func (s *Server) handleWorkshopAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Distro  string `json:"distro"`
		Project string `json:"project"`
		Name    string `json:"name"`
		Action  string `json:"action"` // "start" or "stop"
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Distro == "" || request.Project == "" || request.Name == "" {
		http.Error(w, "distro, project and name are required", http.StatusBadRequest)
		return
	}

	var err error
	switch request.Action {
	case "start":
		err = wsl.StartWorkshop(context.Background(), request.Distro, request.Project, request.Name, s.workshopController)
	case "stop":
		err = wsl.StopWorkshop(context.Background(), request.Distro, request.Project, request.Name, s.workshopController)
	default:
		http.Error(w, "action must be 'start' or 'stop'", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	verb := "started"
	if request.Action == "stop" {
		verb = "stopped"
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Workshop %s %s", request.Name, verb),
	})
}

// handleWorkshopShell opens an interactive `workshop shell` session for a
// workshop in a new terminal window (non-blocking), mirroring /api/launch.
func (s *Server) handleWorkshopShell(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Distro  string `json:"distro"`
		Project string `json:"project"`
		Name    string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Distro == "" || request.Project == "" || request.Name == "" {
		http.Error(w, "distro, project and name are required", http.StatusBadRequest)
		return
	}

	if err := wsl.LaunchWorkshopShell(context.Background(), request.Distro, request.Project, request.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Launched shell for workshop %s", request.Name),
	})
}

func (s *Server) handleCopy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Source     string `json:"source"`
		NewName    string `json:"newName"`
		InstallDir string `json:"installDir,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Source == "" || request.NewName == "" {
		http.Error(w, "Both source and newName are required", http.StatusBadRequest)
		return
	}

	result := wsl.CopyDistro(context.Background(), s.copier, request.Source, request.NewName, request.InstallDir)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleUbuntuTelemetry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		enabled := wsl.GetUbuntuTelemetryStatus()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled": enabled,
		})

	case http.MethodPost:
		var request struct {
			Enabled bool `json:"enabled"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := wsl.SetUbuntuTelemetryStatus(request.Enabled); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled": request.Enabled,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
