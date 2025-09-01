package client

import (
	"fmt"
	"net/url"
	"strings"
)

type RequestData struct {
	Host        string
	Port        int
	ContextPath string
	Path        string
	PathVars    map[string]string
	QueryParams map[string]string
	Headers     map[string]string
	Body        []byte
	UseHTTPS    bool
}

// BuildURL arma la URL final a partir de los componentes
func (r *RequestData) BuildURL() (string, error) {
	u := url.URL{}
	if r.UseHTTPS {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	u.Host = fmt.Sprintf("%s:%d", r.Host, r.Port)

	// Construir el path con context + path
	fullPath := strings.TrimSuffix(r.ContextPath, "/") + "/" + strings.TrimPrefix(r.Path, "/")

	// Reemplazar path variables
	for k, v := range r.PathVars {
		fullPath = strings.ReplaceAll(fullPath, fmt.Sprintf("{%s}", k), v)
	}

	u.Path = fullPath

	// Agregar query params
	if len(r.QueryParams) > 0 {
		q := u.Query()
		for k, v := range r.QueryParams {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}
