package main

import (
	"ascii-art-web/root"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", root.HandleRoot)
	http.HandleFunc("/ascii-art", root.HandleAsciiArt)
	http.Handle("/static/", root.HandleStatic("static"))

	log.Println("Server running at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
