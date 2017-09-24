package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var (
	version string
	build string
	app     = struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Build    string `json:"build"`
		Hostname string `json:"hostname"`
		Port     string `json:"port"`
		Status   string `json:"status"`
		TimeNow  string `json:"time_now"`
	}{
		getName(),
		version,
		build,
		getHostname(),
		getPort(),
		checkHealth(),
		"",
	}
)

func getName() string {
	stdoutStderr, err := exec.Command("go", "list", "-f", "{{.ImportPath}}").CombinedOutput()
	if err != nil {
		log.Print(err)
	}
	return path.Base(strings.TrimSpace(string(stdoutStderr)))
}

func getHostname() string {
	h, err := os.Hostname()
	if err != nil {
		log.Print(err)
	}
	return h
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func checkHealth() string {
	if 2+2 == 4 {
		return "healthy"
	}
	return "unhealthy"
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	app.TimeNow = time.Now().UTC().Format(time.RFC3339)
	if r.URL.Path != "" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, app)
	if err != nil {
		log.Print(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	app.TimeNow = time.Now().UTC().Format(time.RFC3339)
	err := json.NewEncoder(w).Encode(app)
	if err != nil {
		log.Print(err)
	}
}

func logHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
	})
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/health/", handleHealth)
	log.Printf("%s %s now listening on port %s\n", app.Name, app.Version, app.Port)
	log.Fatal(http.ListenAndServe(app.Port, logHTTP(http.DefaultServeMux)))
}
