# ASCII Art Web — Stylize

A Go web server that turns text into ASCII art banners, with a styled,
responsive and interactive interface.

## Usage

```sh
go run .
```

Then open <http://localhost:8000>.

## Features

- **Banner picker**: visual cards preview the `standard`, `shadow` and `thinkertoy` banners.
- **Live feedback**: a character counter (amber near the 200-character limit, red at the limit)
  and an instant warning when a non-ASCII character is typed.
- **Inline errors**: invalid input re-renders the form with a clear message and keeps what was typed.
- **Result tools**: copy to clipboard, download as `.txt`, and a size slider for wide output.
- **Shortcut**: <kbd>Ctrl</kbd>/<kbd>⌘</kbd> + <kbd>Enter</kbd> generates.
- **Light and dark themes**: follows the system setting, with a toggle in the header.
- **Styled error pages** for 400, 404, 405 and 500 responses.
- **Accessible and readable**: AA-contrast colour pairs in both themes, visible focus rings,
  a skip link, labelled controls and reduced-motion support.

All JavaScript is a progressive enhancement: the site works fully without it.

## HTTP endpoints

| Method | Path          | Description                         |
| ------ | ------------- | ----------------------------------- |
| GET    | `/`           | Main page                           |
| POST   | `/ascii-art`  | Generates the art (`text`, `style`) |
| GET    | `/static/...` | CSS and JavaScript assets           |

Status codes: `200` OK, `400` bad request, `404` not found, `405` method not
allowed, `500` internal server error.

## Project structure

```
main.go            server setup and routes
root/              HTTP handlers, validation and static file serving
asciia/            ASCII art rendering
templates/         HTML templates (index, error)
static/            style.css, theme.js, app.js
style/             banner files
```

Only Go standard library packages are used.
