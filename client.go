package cors

import (
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	allowMethods  []string
	exposeHeaders []string

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
		if ok := c.Handler(w, r); ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h(w, r)
	}
}

func (c *Client) Handler(w http.ResponseWriter, r *http.Request) (isOptions bool) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)

	if c.withCookies {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.allowMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.Header().Set("Access-Control-Max-Age", c.maxAge)
		return true
	}

	if len(c.exposeHeaders) != 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(c.exposeHeaders, ", "))
	}
	return false
}

func (c *Client) SetCookie(take bool) *Client {
	c.withCookies = take
	return c
}

func (c *Client) SetAllowMethods(methods []string) *Client {
	var mm []string
	copy(mm, methods)
	c.allowMethods = mm
	return c
}

func (c *Client) SetExposeHeaders(headers []string) *Client {
	var hh []string
	copy(hh, headers)
	c.exposeHeaders = headers
	return c
}

func (c *Client) SetMaxAge(maxAge int) *Client {
	if maxAge < 0 {
		maxAge = -1
	}
	c.maxAge = strconv.FormatInt(int64(maxAge), 10)
	return c
}
