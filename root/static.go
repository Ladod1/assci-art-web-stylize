package root

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// HandleStatic serves the files inside dir under the /static/ prefix.
// Directory listings and missing files get the styled 404 page.
func HandleStatic(dir string) http.Handler {
	files := http.StripPrefix("/static/", http.FileServer(http.Dir(dir)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			writeHTTPError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/static")
		info, err := os.Stat(filepath.Join(dir, filepath.FromSlash(name)))
		if err != nil || info.IsDir() {
			writeHTTPError(w, http.StatusNotFound, "Page not found")
			return
		}

		files.ServeHTTP(w, r)
	})
}
