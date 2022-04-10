package transports

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/sohhamm/string-svc/services"
	s "github.com/sohhamm/string-svc/services"
)

var UppercaseHandler = httpTransport.NewServer(makeUppercaseEndpoint(services.Svc), decodeUppercaseRequest, encodeResponse)

var CountHandler = httpTransport.NewServer(makeCountEndpoint(services.Svc), decodeCountRequest, encodeResponse)

// * Request and Responses (Endpoint layer)
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't json unmarshal so we use strings
}

type countRequest struct {
	S string `json:"s"`
}
type countResponse struct {
	V int `json:"v"`
}

// ? go kit endpoint function signature
// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
// Adapter will take the StringService interface and returns and endpoint corresponding to that one method

// * Adapters

func makeUppercaseEndpoint(svc s.StringService) endpoint.Endpoint {

	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}

		return uppercaseResponse{v, ""}, nil

	}

}

func makeCountEndpoint(svc s.StringService) endpoint.Endpoint {

	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}

}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil

}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil

}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
