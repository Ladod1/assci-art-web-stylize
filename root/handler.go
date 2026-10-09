package root

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"unicode/utf8"

	"ascii-art-web/asciia"
)

const maxTextLength = 200

type pageData struct {
	Text      string
	Style     string
	Result    string
	Error     string
	MaxLength int
}

type errorData struct {
	Status     int
	StatusText string
	Message    string
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches every path that has no more specific handler,
	// so reject anything other than the exact home route.
	if r.URL.Path != "/" {
		writeHTTPError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeHTTPError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	renderPage(w, http.StatusOK, pageData{Style: "standard"})
}

func HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		writeHTTPError(w, http.StatusNotFound, "Page not found")
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeHTTPError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := r.ParseForm(); err != nil {
		writeHTTPError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	text := r.PostForm.Get("text")
	style := r.PostForm.Get("style")

	if !validStyle(style) {
		writeHTTPError(w, http.StatusBadRequest, "Invalid banner style")
		return
	}

	// Input problems are shown next to the form so the user can fix them
	// without losing what they typed.
	data := pageData{Text: text, Style: style}
	switch {
	case text == "":
		data.Error = "Please enter some text to convert."
	case utf8.RuneCountInString(text) > maxTextLength:
		data.Error = fmt.Sprintf("Text is too long (maximum %d characters).", maxTextLength)
	case !validText(text):
		data.Error = "Only printable ASCII characters (letters, digits, spaces and symbols like !?#) are supported."
	}
	if data.Error != "" {
		renderPage(w, http.StatusBadRequest, data)
		return
	}

	result, err := asciiArt.AsciiArt(text, style)
	if err != nil {
		log.Printf("generate ASCII art: %v", err)
		if errors.Is(err, os.ErrNotExist) {
			writeHTTPError(w, http.StatusNotFound, "Banner not found")
			return
		}
		writeHTTPError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	data.Result = result
	renderPage(w, http.StatusOK, data)
}

func renderPage(w http.ResponseWriter, status int, data pageData) {
	data.MaxLength = maxTextLength

	page, err := executeTemplate("templates/index.html", data)
	if err != nil {
		log.Printf("render page: %v", err)
		if errors.Is(err, os.ErrNotExist) {
			writeHTTPError(w, http.StatusNotFound, "Page not found")
			return
		}
		writeHTTPError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writePage(w, status, page)
}

// writeHTTPError renders the styled error page, falling back to plain text
// if the error template itself cannot be rendered.
func writeHTTPError(w http.ResponseWriter, status int, message string) {
	page, err := executeTemplate("templates/error.html", errorData{
		Status:     status,
		StatusText: http.StatusText(status),
		Message:    message,
	})
	if err != nil {
		log.Printf("render error page: %v", err)
		http.Error(w, fmt.Sprintf("%d %s: %s", status, http.StatusText(status), message), status)
		return
	}

	writePage(w, status, page)
}

func executeTemplate(file string, data any) (*bytes.Buffer, error) {
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		return nil, err
	}

	var page bytes.Buffer
	if err := tmpl.Execute(&page, data); err != nil {
		return nil, err
	}
	return &page, nil
}

func writePage(w http.ResponseWriter, status int, page *bytes.Buffer) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := page.WriteTo(w); err != nil {
		log.Printf("write response: %v", err)
	}
}
