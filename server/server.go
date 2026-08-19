package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"wslp/internal/config"
	"wslp/internal/wsl"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

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

	// Add CORS middleware for Flutter
	handler := corsMiddleware(mux)

	addr := fmt.Sprintf(":%s", s.port)
	fmt.Printf("Starting server on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, handler)
}

func (s *Server) handleListDistros(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	distros, err := wsl.ListDistros(context.Background(), wsl.RealLister{})
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

	defaultDistro, err := wsl.GetDefaultDistro(context.Background(), wsl.RealDefaultGetter{})
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

	unregisterer := wsl.RealUnregisterer{}
	results := wsl.UnregisterDistros(context.Background(), unregisterer, request.Distros)

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

	if err := wsl.SetDefaultDistro(context.Background(), request.Name, wsl.RealDefaultSetter{}); err != nil {
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

	backuper := wsl.RealBackuper{}
	results := wsl.BackupDistros(context.Background(), backuper, request.Distros, backupDir, opts)

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

	terminator := wsl.RealTerminator{}
	results := wsl.TerminateDistros(context.Background(), terminator, request.Distros)

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

	renamer := wsl.RealRenamer{}
	result := wsl.RenameDistro(context.Background(), renamer, request.OldName, request.NewName)

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

	workshops := wsl.GetWorkshops(context.Background(), name, wsl.RealWorkshopRunner{})

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

	controller := wsl.RealWorkshopController{}

	var err error
	switch request.Action {
	case "start":
		err = wsl.StartWorkshop(context.Background(), request.Distro, request.Project, request.Name, controller)
	case "stop":
		err = wsl.StopWorkshop(context.Background(), request.Distro, request.Project, request.Name, controller)
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

	copier := wsl.RealCopier{}
	result := wsl.CopyDistro(context.Background(), copier, request.Source, request.NewName, request.InstallDir)

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
