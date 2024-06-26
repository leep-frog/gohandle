package gohandle

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TemplateFile string

type TemplateFiles []TemplateFile

func Join(ss ...string) TemplateFile {
	return TemplateFile(filepath.Join(ss...))
}

type Handler interface {
	http.Handler
	GetPattern() string
}

type TemplateHandler struct {
	Pattern      string
	Template     TemplateFile
	Templates    TemplateFiles
	Functions    []Function
	GenerateData func(requestBodyBytes []byte) (any, error)
}

func (sh *TemplateHandler) GetPattern() string {
	return sh.Pattern
}

func (sh *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: implement our own pattern handling since the pattern handling
	// by net/http includes all sub paths
	// https://pkg.go.dev/net/http#ServeMux
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	funcs := sh.Functions
	funcMap := map[string]any{}
	for _, f := range funcs {
		funcMap[f.Name()] = f.Func()
	}
	baseName := filepath.Base(string(sh.Template))
	tmpl, err := template.New(baseName).Funcs(funcMap).ParseFiles(convertTemplates(sh.Templates)...)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Server error: failed to render template: %v", err)))
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Server error: failed to process request body: %v", err)))
		return
	}

	var data any
	if sh.GenerateData != nil {
		data, err = sh.GenerateData(b)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Server error: failed to generate template data: %v", err)))
			return
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		w.Write([]byte(fmt.Sprintf("Template error: %v", err)))
	}
}

type RedirectHandler struct {
	Pattern string
	Dest    string
}

func (rh *RedirectHandler) GetPattern() string {
	return rh.Pattern
}

func (rh *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, rh.Dest, http.StatusPermanentRedirect)
}

type Request struct{}

func convertTemplates(tfs TemplateFiles) []string {
	var r []string
	for _, tf := range tfs {
		r = append(r, string(tf))
	}
	return r
}

const LOCAL_ENV_VAR = "GOHANDLE_LOCAL"

func Run(handlers []Handler) {
	for _, h := range handlers {
		http.HandleFunc(h.GetPattern(), h.ServeHTTP)
	}

	port := ":8080"

	local := strings.TrimSpace(os.Getenv(LOCAL_ENV_VAR)) != ""
	if local {
		// Update port (so windows doesn't ask allow/block prompt)
		port = "127.0.0.1:8080"
		// Open Chrome
		if err := exec.Command("cmd", `/c`, "start", `C:\Program Files\Google\Chrome\Application\chrome.exe`, `http://localhost:8080`).Start(); err != nil {
			log.Fatalf("Failed to run chrome: %v", err)
		}
	}

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failure: %v", err)
	}
}
