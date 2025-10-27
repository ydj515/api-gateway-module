package http

import "errors"

type HttpMethod string

const (
	GET    = HttpMethod("GET")
	POST   = HttpMethod("POST")
	DELETE = HttpMethod("DELETE")
	PUT    = HttpMethod("PUT")
)

func (h HttpMethod) String() string {
	return string(h)
}

type GetType string

const (
	QUERY = GetType("query")
	URL   = GetType("url")
)

func (h GetType) String() string {
	return string(h)
}

func (h GetType) CheckType() error {
	switch h {
	case QUERY, URL:
		return nil
	default:
		return errors.New("Failed to check get type")
	}
}
