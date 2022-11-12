package request

import stdHTTP "net/http"

type Doer interface {
	Do(req *stdHTTP.Request) (*stdHTTP.Response, error)
}
