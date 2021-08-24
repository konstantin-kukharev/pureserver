package http

import (
	"errors"
	"fmt"
	ps "github.com/konstantin-kukharev/pureserver/internal"
	"log"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	A          PatternServeMuxInterface
	loops      int
	port       []int
	unixSocket []string
	router     PatternServeMuxInterface
}

func (server *Server) SetLoops(loops int) {
	server.loops = loops
}

func (server *Server) SetPort(port ...int) {
	server.port = port
}

func (server *Server) SetUnixSocket(socket ...string) {
	server.unixSocket = socket
}

func (server *Server) Serve() error {
	var events ps.Events
	var addresses []string

	if server.loops == 0 {
		server.loops = -1
	}

	events.NumLoops = server.loops
	events.Serving = func(srv ps.Server) (action ps.Action) {
		log.Printf("http server started on port %d (loops: %d)", server.port, srv.NumLoops)
		if len(server.unixSocket) != 0 {
			log.Printf("http server started at %v", server.unixSocket)
		}
		return
	}

	events.Opened = func(c ps.Conn) (out []byte, opts ps.Options, action ps.Action) {
		c.SetContext(&ps.InputStream{})
		return
	}

	events.Closed = func(c ps.Conn, err error) (action ps.Action) {
		return
	}

	events.Data = func(c ps.Conn, in []byte) (out []byte, action ps.Action) {
		if in == nil {
			return
		}
		is := c.Context().(*ps.InputStream)
		data := is.Begin(in)
		// process the pipeline
		var req Request
		for {
			leftover, err := server.parseRequest(data, &req)
			if err != nil {
				out = server.appendResponse(out, "500 Error", "", err.Error()+"\n")
				action = ps.Close
				break
			} else if len(leftover) == len(data) {
				// request not ready, yet
				break
			}
			// handle the request
			req.RemoteAddr = c.RemoteAddr().String()
			out = server.appendHandle(out, &req)
			data = leftover
		}
		is.End(data)
		return
	}

	if len(server.port) != 0 {
		for _, curPort := range server.port {
			addresses = append(addresses, fmt.Sprintf("tcp://:%d", curPort))
		}
	}

	if len(server.unixSocket) != 0 {
		for _, curAddr := range server.unixSocket {
			addresses = append(addresses, fmt.Sprintf("unix://%s", curAddr))
		}
	}

	if len(addresses) == 0 {
		return errors.New("no address specified")
	}

	return ps.Serve(events, addresses...)
}

func (server *Server) SetHandler(handler PatternServeMuxInterface) {
	server.router = handler
}

// appendHandle handles the incoming request and appends the response to
// the provided bytes, which is then returned to the caller.
func (server *Server) appendHandle(_ []byte, req HttpRequestInterface) []byte {
	writer := Writer{
		head: Header{},
	}
	server.router.ServeHTTP(&writer, req)
	writer.Write()
	return writer.response
}

// A Handler responds to an HTTP request.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
//
// Depending on the HTTP client software, HTTP protocol version, and
// any intermediaries between the client and the Go server, it may not
// be possible to read from the Request.Body after writing to the
// ResponseWriter. Cautious handlers should read the Request.Body
// first, and then reply.
//
// Except for reading the body, handlers should not modify the
// provided Request.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
// that the effect of the panic was isolated to the active request.
// It recovers the panic, logs a stack trace to the server error log,
// and either closes the network connection or sends an HTTP/2
// RST_STREAM, depending on the HTTP protocol. To abort a handler so
// the client sees an interrupted response but the server doesn't log
// an error, panic with the value ErrAbortHandler.
//
// appendResponse will append a valid http response to provide bytes.
// The status param should be the code plus text such as "200 OK".
// The head parameter should be a series of lines ending with "\r\n" or empty.
func (server *Server) appendResponse(b []byte, status, head, body string) []byte {
	b = append(b, "HTTP/1.1"...)
	b = append(b, ' ')
	b = append(b, status...)
	b = append(b, '\r', '\n')
	b = append(b, "Server: pure server\r\n"...)
	b = append(b, "Date: "...)
	b = time.Now().AppendFormat(b, "Mon, 02 Jan 2006 15:04:05 GMT")
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, "Content-Length: "...)
		b = strconv.AppendInt(b, int64(len(body)), 10)
		b = append(b, '\r', '\n')
	}
	b = append(b, head...)
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, body...)
	}
	return b
}

// parseRequest is a very simple http request parser. This operation
// waits for the entire payload to be buffered before returning a
// valid request.
func (server *Server) parseRequest(data []byte, req HttpRequestInterface) (leftover []byte, err error) {
	sData := string(data)
	var i, s int
	var top string
	var clen int
	var q = -1
	// method, path, proto line
	for ; i < len(sData); i++ {
		if sData[i] == ' ' {
			req.SetMethod(sData[s:i])
			for i, s = i+1, i+1; i < len(sData); i++ {
				if sData[i] == '?' && q == -1 {
					q = i - s
				} else if sData[i] == ' ' {
					if q != -1 {
						req.SetPath(sData[s:q])
						req.SetQuery(req.GetPath().EscapedPath()[q+1 : i])
					} else {
						req.SetPath(sData[s:i])
					}
					for i, s = i+1, i+1; i < len(sData); i++ {
						if sData[i] == '\n' && sData[i-1] == '\r' {
							req.SetProto(sData[s : i-1])
							i, s = i+1, i+1
							break
						}
					}
					break
				}
			}
			break
		}
	}
	if req.GetProto() == "" {
		return data, fmt.Errorf("malformed request")
	}
	top = sData[:s]
	for ; i < len(sData); i++ {
		if i > 1 && sData[i] == '\n' && sData[i-1] == '\r' {
			line := sData[s : i-1]
			s = i + 1
			if line == "" {
				req.SetHead(sData[len(top) : i+1])
				i++
				if clen > 0 {
					if len(sData[i:]) < clen {
						break
					}
					aa := sData[i : i+clen]
					req.SetBody(aa)
					i += clen
				}
				return data[i:], nil
			}
			if strings.HasPrefix(line, "Content-Length:") {
				n, err := strconv.ParseInt(strings.TrimSpace(line[len("Content-Length:"):]), 10, 64)
				if err == nil {
					clen = int(n)
				}
			}
		}
	}
	// not enough data
	return data, nil
}
