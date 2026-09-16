package upstream

import (
	"net/http"
	"time"
)

// NewClient returns the single *http.Client shared by every upstream call.
// It deliberately does not set Client.Timeout — the only per-call timeout
// mechanism is context.WithTimeout in the tool handler, so there is exactly
// one timeout source to reason about.
//
// Redirects are never followed. Go's default client will copy custom auth
// headers (anything other than Authorization/Cookie) to any host, put
// query-auth secrets into the Referer of the next hop, and fetch whatever
// Location points at — including link-local/metadata URLs. A 301/302 on
// POST is also rewritten to GET, so a failed create can look like a
// successful list. The 3xx is surfaced as a tool error instead.
func NewClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}
	return &http.Client{
		Transport:     transport,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
}
