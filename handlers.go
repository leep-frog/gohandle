package gohandle

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type TemplateFile string

type TemplateFiles []TemplateFile

func Join(ss ...string) TemplateFile {
	return TemplateFile(filepath.Join(ss...))
}

type Handler interface {
	GetPattern() string
	GetTemplate() TemplateFile
	GetTemplates() TemplateFiles
	GetFunctions() []Function
	GetData() any
}

type SimpleHandler struct {
	Pattern   string
	Template  TemplateFile
	Templates TemplateFiles
	Functions []Function
	Data      any
}

func (sh *SimpleHandler) GetPattern() string {
	return sh.Pattern
}

func (sh *SimpleHandler) GetTemplate() TemplateFile {
	return sh.Template
}

func (sh *SimpleHandler) GetTemplates() TemplateFiles {
	return sh.Templates
}

func (sh *SimpleHandler) GetFunctions() []Function {
	return sh.Functions
}

func (sh *SimpleHandler) GetData() any {
	return sh.Data
}

type Request struct{}

func convertTemplates(tfs TemplateFiles) []string {
	var r []string
	for _, tf := range tfs {
		r = append(r, string(tf))
	}
	return r
}

func Handle(h Handler) {
	http.HandleFunc(h.GetPattern(), func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement our own pattern handling since the pattern handling
		// by net/http includes all sub paths
		// https://pkg.go.dev/net/http#ServeMux
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		funcs := h.GetFunctions()
		funcMap := map[string]any{}
		for _, f := range funcs {
			funcMap[f.Name()] = f.Func()
		}
		baseName := filepath.Base(string(h.GetTemplate()))
		tmpl, err := template.New(baseName).Funcs(funcMap).ParseFiles(convertTemplates(h.GetTemplates())...)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Server error: failed to render template: %v", err)))
			return
		}

		if err := tmpl.Execute(w, h.GetData()); err != nil {
			w.Write([]byte(fmt.Sprintf("Template error: %v", err)))
		}
	})
}

func Run(port string, handlers []Handler) {
	for _, h := range handlers {
		Handle(h)
	}

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failure: %v", err)
	}
}
