package toolshed

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const errorScript = `echo "failed to fetch upstream script" && exit 1`

type server struct {
	logger  *log.Logger
	listen  string
	fetcher fetcher
}

func (s *server) handleFavicon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version := parseVersion(r.URL.Path)

		script, err := s.fetcher.Fetch(version)
		if err != nil {
			s.logger.Printf("request failed for %s: %s", version, err)
			http.Error(w, errorScript, http.StatusInternalServerError)
			return
		}

		v := fmt.Sprintf(`BELT_VERSION="%s"`, version)

		script = strings.Replace(script, `BELT_VERSION="master"`, v, 1)

		s.logger.Printf("request succeeded for %s", version)
		w.Write([]byte(script))
	}
}

func (s *server) Routes() {
	http.HandleFunc("/favicon.ico", s.handleFavicon())
	http.HandleFunc("/", s.handleIndex())
}

func (s *server) Run() error {
	s.logger.Printf("server running at %s", s.listen)
	return http.ListenAndServe(s.listen, nil)
}

func parseVersion(path string) string {
	base := filepath.Base(path)

	if base == "/" {
		return "master"
	}

	return base
}
