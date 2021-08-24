package middleware

type Responder interface {
}

func NotImplemented(info string) Responder {
	type resp struct {
		Info string
	}
	return resp{info}
}