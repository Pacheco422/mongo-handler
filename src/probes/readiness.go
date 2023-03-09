// Implements a Readiness probe with ReadinessHandler
package probes

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ReadinessHandler struct {
	Ready bool
}


func NewReadinessHandler() *ReadinessHandler {
	return &ReadinessHandler{false}
}

// Readiness probe .
func (p *ReadinessHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Debug("Readiness Probe was called .",)
	if p.Ready == true{
		rw.WriteHeader(http.StatusOK)
	} else {
		http.Error(rw, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}
}
