// Package main provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// SignUpJSONBody defines parameters for SignUp.
type SignUpJSONBody struct {
	ConfirmPassword string `json:"confirm_password"`

	// Email valid email
	Email string `json:"email"`

	// Name full user name
	Name string `json:"name"`

	// Password strong password
	Password string `json:"password"`
}

// SignUpJSONRequestBody defines body for SignUp for application/json ContentType.
type SignUpJSONRequestBody SignUpJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// SignUpWithBody request with any body
	SignUpWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SignUp(ctx context.Context, body SignUpJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) SignUpWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSignUpRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SignUp(ctx context.Context, body SignUpJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSignUpRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewSignUpRequest calls the generic SignUp builder with application/json body
func NewSignUpRequest(server string, body SignUpJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSignUpRequestWithBody(server, "application/json", bodyReader)
}

// NewSignUpRequestWithBody generates requests for SignUp with any type of body
func NewSignUpRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/signup")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// SignUpWithBodyWithResponse request with any body
	SignUpWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignUpResponse, error)

	SignUpWithResponse(ctx context.Context, body SignUpJSONRequestBody, reqEditors ...RequestEditorFn) (*SignUpResponse, error)
}

type SignUpResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Email *string `json:"email,omitempty"`
		Name  *string `json:"name,omitempty"`
	}
	JSONDefault *struct {
		Error   *string `json:"error,omitempty"`
		Message *string `json:"message,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r SignUpResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SignUpResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// SignUpWithBodyWithResponse request with arbitrary body returning *SignUpResponse
func (c *ClientWithResponses) SignUpWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignUpResponse, error) {
	rsp, err := c.SignUpWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignUpResponse(rsp)
}

func (c *ClientWithResponses) SignUpWithResponse(ctx context.Context, body SignUpJSONRequestBody, reqEditors ...RequestEditorFn) (*SignUpResponse, error) {
	rsp, err := c.SignUp(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignUpResponse(rsp)
}

// ParseSignUpResponse parses an HTTP response from a SignUpWithResponse call
func ParseSignUpResponse(rsp *http.Response) (*SignUpResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SignUpResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Email *string `json:"email,omitempty"`
			Name  *string `json:"name,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest struct {
			Error   *string `json:"error,omitempty"`
			Message *string `json:"message,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
