package cors

import (
	"net/http"
	"strconv"
	"strings"
)

var defaultMethods = []string{"GET", "DELETE", "HEAD", "PATCH", "POST", "PUT"}

type Client struct {
	// the default supports 'GET', 'DELETE', 'HEAD', 'OPTIONS', 'PATCH', 'POST', 'PUT' methods.
	AllowMethods []string

	// the default are 'Content-Type', 'Content-Length', 'Accept', 'Authorization'
	AllowHeaders []string

	// the default is 14400s(4 hours), the -1 disableds cache.
	MaxAge int

	// set Cookies, the default is true.
	WithoutCookies bool
}

func New() *Client {
	return &Client{}
}

func (c *Client) WrapH(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, ok := c.Handler(w, r)
		if ok {
			w.WriteHeader(status)
			return
		}
		h(w, r)
	}
}

func (c *Client) Handler(w http.ResponseWriter, r *http.Request) (status int, ok bool) {
	origin := r.URL.Query().Get("Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", c.allowMethods())
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Expose-Headers"))
		w.Header().Set("Access-Control-Allow-Headers", c.maxAge())
		if !c.WithoutCookies {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		return 204, false
	}
	return 0, true
}

func (c *Client) allowMethods() string {
	if len(c.AllowMethods) == 0 {
		c.AllowMethods = defaultMethods
	}
	return strings.Join(c.AllowMethods, ", ")
}

func (c *Client) maxAge() string {
	if c.MaxAge == 0 {
		c.MaxAge = 14400
	}
	if c.MaxAge < 0 {
		c.MaxAge = -1
	}
	return strconv.FormatInt(int64(c.MaxAge), 10)
}
