package root

import (
	"html/template"
	"log"
	"net/http"
	"unicode/utf8"

	"ascii-art-web/asciia"
)

const maxTextLength = 200

// pageData is what the HTML templates show.
type pageData struct {
	Text   string
	Style  string
	Result string
	Error  string
	Status int
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		showError(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		showError(w, http.StatusMethodNotAllowed)
		return
	}

	showPage(w, http.StatusOK, pageData{Style: "standard"})
}

func HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		showError(w, http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	style := r.FormValue("style")

	if !validStyle(style) {
		showError(w, http.StatusBadRequest)
		return
	}

	data := pageData{Text: text, Style: style}

	// Wrong input: show the form again with a message.
	if text == "" {
		data.Error = "Please type some text."
	} else if utf8.RuneCountInString(text) > maxTextLength {
		data.Error = "Your text is too long (maximum 200 characters)."
	} else if !validText(text) {
		data.Error = "Only English letters, numbers, spaces and symbols are allowed."
	}
	if data.Error != "" {
		showPage(w, http.StatusBadRequest, data)
		return
	}

	result, err := asciiArt.AsciiArt(text, style)
	if err != nil {
		log.Println("ascii art:", err)
		showError(w, http.StatusInternalServerError)
		return
	}

	data.Result = result
	showPage(w, http.StatusOK, data)
}

func showPage(w http.ResponseWriter, status int, data pageData) {
	render(w, "templates/index.html", status, data)
}

func showError(w http.ResponseWriter, status int) {
	render(w, "templates/error.html", status, pageData{Status: status, Error: http.StatusText(status)})
}

func render(w http.ResponseWriter, file string, status int, data pageData) {
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		log.Println("template:", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template:", err)
	}
}
