package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func Test_getName(t *testing.T) {
	p := getName()
	if p != "bluewhales" {
		t.Error()
	}
}

func Test_getHostname(t *testing.T) {
	h := getHostname()
	if len(h) < 1 {
		t.Error()
	}
}

func Test_getPort(t *testing.T) {
	p := getPort()
	r := regexp.MustCompile(`^:\d+$`)
	if !r.MatchString(p) {
		t.Error()
	}
}

func Test_checkHealth(t *testing.T) {
	s := checkHealth()
	if s != "healthy" && s != "unhealthy" {
		t.Error()
	}
}

func Test_handleRoot(t *testing.T) {
	testTable := []struct {
		path       string
		statusCode int
	}{
		{"", 200},
		{"/", 200},
		{"//", 200},
		{"/foo", 404},
		{"/foo/", 404},
	}
	for _, tt := range testTable {
		r, _ := http.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		handleRoot(w, r)
		if w.Code != tt.statusCode {
			t.Errorf("GET %s returned %d, expected %d", tt.path, w.Code, tt.statusCode)
		}
	}
}

func Test_handleHealth(t *testing.T) {
	r, _ := http.NewRequest("GET", "/health/", nil)
	w := httptest.NewRecorder()
	handleHealth(w, r)
	if w.Code != 200 {
		t.Error()
	}
}
