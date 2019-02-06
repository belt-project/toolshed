package toolshed

import "log"

// Run starts the server.
func Run(listen string, logger *log.Logger) error {
	fetcher := &githubFetcher{
		logger: logger,
		repo:   "belt-sh/belt.sh",
		cache:  make(map[string]string),
	}

	s := &server{logger, listen, fetcher}
	s.Routes()

	return s.Run()
}
