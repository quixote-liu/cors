package cors

import (
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	allowMethods []string

	// the default is 14400s(4 hours), the -1 disableds cache.
	maxAge string

	// set Cookies, the default is true.
	withCookies bool
}

func New() *Client {
	return &Client{
		allowMethods: []string{"GET", "DELETE", "HEAD", "PATCH", "POST", "PUT"},
		maxAge:       "14400",
		withCookies:  true,
	}
}

func (c *Client) WrapH(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, ok := c.Handler(w, r)
		if !ok {
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
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.allowMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Expose-Headers"))
		w.Header().Set("Access-Control-Allow-Headers", c.maxAge)
		if c.withCookies {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		return 204, false
	}
	return 0, true
}

func (c *Client) SetAllowMethods(methods []string) *Client {
	var mm []string
	copy(mm, methods)
	c.allowMethods = mm
	return c
}

func (c *Client) SetMaxAge(maxAge int) *Client {
	if maxAge < 0 {
		maxAge = -1
	}
	c.maxAge = strconv.FormatInt(int64(maxAge), 10)
	return c
}
