// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/pytorchtw/go-jwt-server/gen/models"
)

// GetGreetingOKCode is the HTTP code returned for type GetGreetingOK
const GetGreetingOKCode int = 200

/*GetGreetingOK returns a greeting

swagger:response getGreetingOK
*/
type GetGreetingOK struct {

	/*
	  In: Body
	*/
	Payload *models.User `json:"body,omitempty"`
}

// NewGetGreetingOK creates GetGreetingOK with default headers values
func NewGetGreetingOK() *GetGreetingOK {

	return &GetGreetingOK{}
}

// WithPayload adds the payload to the get greeting o k response
func (o *GetGreetingOK) WithPayload(payload *models.User) *GetGreetingOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get greeting o k response
func (o *GetGreetingOK) SetPayload(payload *models.User) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetGreetingOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
