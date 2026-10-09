# ASCII Art Web — Stylize

A Go web server that turns text into ASCII art, with a simple CSS design.

## Run

```sh
go run .
```

Then open <http://localhost:8000>.

## Files

```
main.go              starts the server
root/handler.go      handles the pages
root/validation.go   checks the user's input
asciia/              makes the ASCII art
templates/           index.html and error.html
static/style.css     the design
style/               the banner files
```

Only Go standard packages are used.
