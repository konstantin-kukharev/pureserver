package http

import "net/url"

type Request struct {
	Proto, Method     string
	Path              *url.URL
	Query, Head, Body string
	RemoteAddr        string
}

func (r *Request) GetProto() string {
	return r.Proto
}

func (r *Request) GetMethod() string {
	return r.Method
}

func (r *Request) GetPath() *url.URL {
	return r.Path
}

func (r *Request) GetQuery() string {
	return r.Query
}

func (r *Request) GetHead() string {
	return r.Head
}

func (r *Request) GetBody() string {
	return r.Body
}

func (r *Request) GetRemoteAddr() string {
	return r.RemoteAddr
}

func (r *Request) SetProto(proto string) {
	r.Proto = proto
}

func (r *Request) SetMethod(method string) {
	r.Method = method
}

func (r *Request) SetPath(path string) {
	r.Path, _ = url.Parse(path)
}

func (r *Request) SetQuery(query string) {
	r.Query = query
}

func (r *Request) SetHead(head string) {
	r.Head = head
}

func (r *Request) SetBody(body string) {
	r.Body = body
}

func (r *Request) SetRemoteAddr(address string) {
	r.RemoteAddr = address
}

func (r *Request) GetParam(key string) (value string, isEmpty bool) {
	val := r.GetPath().Query().Get(`:` + key)

	return val, val == ""
}
