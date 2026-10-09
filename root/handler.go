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
	Text   string
	Style  string
	Result string
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

	renderPage(w, pageData{Style: "standard"})
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

	if text == "" {
		writeHTTPError(w, http.StatusBadRequest, "Text is required")
		return
	}
	if !validStyle(style) {
		writeHTTPError(w, http.StatusBadRequest, "Invalid banner style")
		return
	}
	if utf8.RuneCountInString(text) > maxTextLength {
		writeHTTPError(w, http.StatusBadRequest, "Text is too long (maximum 200 characters)")
		return
	}
	if !validText(text) {
		writeHTTPError(w, http.StatusBadRequest, "Invalid text: only printable ASCII characters are supported")
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

	renderPage(w, pageData{Text: text, Style: style, Result: result})
}

func renderPage(w http.ResponseWriter, data pageData) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("parse template: %v", err)
		if errors.Is(err, os.ErrNotExist) {
			writeHTTPError(w, http.StatusNotFound, "Page not found")
			return
		}
		writeHTTPError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var page bytes.Buffer
	if err := tmpl.Execute(&page, data); err != nil {
		log.Printf("render template: %v", err)
		writeHTTPError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := page.WriteTo(w); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeHTTPError(w http.ResponseWriter, status int, message string) {
	http.Error(w, fmt.Sprintf("%d %s: %s", status, http.StatusText(status), message), status)
}