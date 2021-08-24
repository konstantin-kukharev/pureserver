package http

import (
	"net/url"
	"strings"
)

type PatternServeMux struct {
	// NotFound, if set, is used whenever the request doesn't match any
	// pattern for its method. NotFound should be set before serving any
	// requests.
	NotFound func(w ResponseWriter, r HttpRequestInterface)
	Handlers map[string][]*PatHandler
}

func (p *PatternServeMux) ServeHTTP(w ResponseWriter, r HttpRequestInterface) {
	for _, ph := range p.Handlers[r.GetMethod()] {
		if params, ok := ph.try(r.GetPath().EscapedPath()); ok {
			if len(params) > 0 && !ph.redirect {
				r.GetPath().RawQuery = url.Values(params).Encode() + "&" + r.GetPath().RawQuery
			}
			ph.Handler.Handle(w, r)
			return
		}
	}

	if p.NotFound != nil {
		p.NotFound(w, r)
		return
	}

	allowed := make([]string, 0, len(p.Handlers))
	for meth, handlers := range p.Handlers {
		if meth == r.GetMethod() {
			continue
		}

		for _, ph := range handlers {
			if _, ok := ph.try(r.GetPath().EscapedPath()); ok {
				allowed = append(allowed, meth)
			}
		}
	}

	if len(allowed) == 0 {
		NotFound(w, r)
		return
	}

	w.Header().Add("Allow", strings.Join(allowed, ", "))
	Error(w, "Method Not Allowed", StatusMethodNotAllowed)
	return
}

// Head will register a pattern with a handler for HEAD requests.
func (p *PatternServeMux) Head(pat string, h CallableHandler) {
	p.Add("HEAD", pat, h)
}

// Get will register a pattern with a handler for GET requests.
// It also registers pat for HEAD requests. If this needs to be overridden, use
// Head before Get with pat.
func (p *PatternServeMux) Get(pat string, h CallableHandler) {
	p.Add("HEAD", pat, h)
	p.Add("GET", pat, h)
}

// Post will register a pattern with a handler for POST requests.
func (p *PatternServeMux) Post(pat string, h CallableHandler) {
	p.Add("POST", pat, h)
}

// Put will register a pattern with a handler for PUT requests.
func (p *PatternServeMux) Put(pat string, h CallableHandler) {
	p.Add("PUT", pat, h)
}

// Del will register a pattern with a handler for DELETE requests.
func (p *PatternServeMux) Del(pat string, h CallableHandler) {
	p.Add("DELETE", pat, h)
}

// Options will register a pattern with a handler for OPTIONS requests.
func (p *PatternServeMux) Options(pat string, h CallableHandler) {
	p.Add("OPTIONS", pat, h)
}

// Patch will register a pattern with a handler for PATCH requests.
func (p *PatternServeMux) Patch(pat string, h CallableHandler) {
	p.Add("PATCH", pat, h)
}

// Add will register a pattern with a handler for meth requests.
func (p *PatternServeMux) Add(meth, pat string, h CallableHandler) {
	p.add(meth, pat, h, false)
}

func (p *PatternServeMux) add(meth, pat string, h CallableHandler, redirect bool) {
	handlers := p.Handlers[meth]
	for _, p1 := range handlers {
		if p1.pat == pat {
			return // found existing pattern; do nothing
		}
	}
	handler := &PatHandler{
		pat:      pat,
		Handler:  h,
		redirect: redirect,
	}
	p.Handlers[meth] = append(handlers, handler)
}

type PatHandler struct {
	pat      string
	Handler  CallableHandler
	redirect bool
}

func (ph *PatHandler) try(path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(ph.pat):
			if ph.pat != "/" && len(ph.pat) > 0 && ph.pat[len(ph.pat)-1] == '/' {
				return p, true
			}
			return nil, false
		case ph.pat[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(ph.pat, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			escval, err := url.QueryUnescape(val)
			if err != nil {
				return nil, false
			}
			p.Add(":"+name, escval)
		case path[i] == ph.pat[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(ph.pat) {
		return nil, false
	}
	return p, true
}

func matchPart(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
