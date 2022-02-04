// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package generated

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

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
	// Authenticate request with any body
	AuthenticateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Authenticate(ctx context.Context, body AuthenticateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// HealthCheck request
	HealthCheck(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Keystore request
	Keystore(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// KeystoreEntries request
	KeystoreEntries(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateEntry request with any body
	CreateEntryWithBody(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateEntry(ctx context.Context, keystoreId string, body CreateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateEntry request with any body
	UpdateEntryWithBody(ctx context.Context, keystoreId string, entryId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateEntry(ctx context.Context, keystoreId string, entryId string, body UpdateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Keystores request
	Keystores(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateKeystore request with any body
	CreateKeystoreWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateKeystore(ctx context.Context, body CreateKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) AuthenticateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAuthenticateRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Authenticate(ctx context.Context, body AuthenticateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAuthenticateRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) HealthCheck(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewHealthCheckRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Keystore(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewKeystoreRequest(c.Server, keystoreId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) KeystoreEntries(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewKeystoreEntriesRequest(c.Server, keystoreId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateEntryWithBody(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateEntryRequestWithBody(c.Server, keystoreId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateEntry(ctx context.Context, keystoreId string, body CreateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateEntryRequest(c.Server, keystoreId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateEntryWithBody(ctx context.Context, keystoreId string, entryId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateEntryRequestWithBody(c.Server, keystoreId, entryId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateEntry(ctx context.Context, keystoreId string, entryId string, body UpdateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateEntryRequest(c.Server, keystoreId, entryId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Keystores(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewKeystoresRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateKeystoreWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateKeystoreRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateKeystore(ctx context.Context, body CreateKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateKeystoreRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewAuthenticateRequest calls the generic Authenticate builder with application/json body
func NewAuthenticateRequest(server string, body AuthenticateJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAuthenticateRequestWithBody(server, "application/json", bodyReader)
}

// NewAuthenticateRequestWithBody generates requests for Authenticate with any type of body
func NewAuthenticateRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/authenticate")
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

// NewHealthCheckRequest generates requests for HealthCheck
func NewHealthCheckRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/health")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewKeystoreRequest generates requests for Keystore
func NewKeystoreRequest(server string, keystoreId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "keystoreId", runtime.ParamLocationPath, keystoreId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/keystore/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewKeystoreEntriesRequest generates requests for KeystoreEntries
func NewKeystoreEntriesRequest(server string, keystoreId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "keystoreId", runtime.ParamLocationPath, keystoreId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/keystore/%s/entries", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateEntryRequest calls the generic CreateEntry builder with application/json body
func NewCreateEntryRequest(server string, keystoreId string, body CreateEntryJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateEntryRequestWithBody(server, keystoreId, "application/json", bodyReader)
}

// NewCreateEntryRequestWithBody generates requests for CreateEntry with any type of body
func NewCreateEntryRequestWithBody(server string, keystoreId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "keystoreId", runtime.ParamLocationPath, keystoreId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/keystore/%s/entries", pathParam0)
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

// NewUpdateEntryRequest calls the generic UpdateEntry builder with application/json body
func NewUpdateEntryRequest(server string, keystoreId string, entryId string, body UpdateEntryJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateEntryRequestWithBody(server, keystoreId, entryId, "application/json", bodyReader)
}

// NewUpdateEntryRequestWithBody generates requests for UpdateEntry with any type of body
func NewUpdateEntryRequestWithBody(server string, keystoreId string, entryId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "keystoreId", runtime.ParamLocationPath, keystoreId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "entryId", runtime.ParamLocationPath, entryId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/keystore/%s/entry/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewKeystoresRequest generates requests for Keystores
func NewKeystoresRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/keystores")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateKeystoreRequest calls the generic CreateKeystore builder with application/json body
func NewCreateKeystoreRequest(server string, body CreateKeystoreJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateKeystoreRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateKeystoreRequestWithBody generates requests for CreateKeystore with any type of body
func NewCreateKeystoreRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/keystores")
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
	// Authenticate request with any body
	AuthenticateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AuthenticateResponse, error)

	AuthenticateWithResponse(ctx context.Context, body AuthenticateJSONRequestBody, reqEditors ...RequestEditorFn) (*AuthenticateResponse, error)

	// HealthCheck request
	HealthCheckWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthCheckResponse, error)

	// Keystore request
	KeystoreWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*KeystoreResponse, error)

	// KeystoreEntries request
	KeystoreEntriesWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*KeystoreEntriesResponse, error)

	// CreateEntry request with any body
	CreateEntryWithBodyWithResponse(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateEntryResponse, error)

	CreateEntryWithResponse(ctx context.Context, keystoreId string, body CreateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateEntryResponse, error)

	// UpdateEntry request with any body
	UpdateEntryWithBodyWithResponse(ctx context.Context, keystoreId string, entryId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateEntryResponse, error)

	UpdateEntryWithResponse(ctx context.Context, keystoreId string, entryId string, body UpdateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateEntryResponse, error)

	// Keystores request
	KeystoresWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*KeystoresResponse, error)

	// CreateKeystore request with any body
	CreateKeystoreWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateKeystoreResponse, error)

	CreateKeystoreWithResponse(ctx context.Context, body CreateKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateKeystoreResponse, error)
}

type AuthenticateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r AuthenticateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AuthenticateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type HealthCheckResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r HealthCheckResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r HealthCheckResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type KeystoreResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Keystore
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r KeystoreResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r KeystoreResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type KeystoreEntriesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Entry
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r KeystoreEntriesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r KeystoreEntriesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateEntryResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Entry
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r CreateEntryResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateEntryResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateEntryResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Entry
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r UpdateEntryResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateEntryResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type KeystoresResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Keystores
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r KeystoresResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r KeystoresResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateKeystoreResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Keystore
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r CreateKeystoreResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateKeystoreResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// AuthenticateWithBodyWithResponse request with arbitrary body returning *AuthenticateResponse
func (c *ClientWithResponses) AuthenticateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AuthenticateResponse, error) {
	rsp, err := c.AuthenticateWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAuthenticateResponse(rsp)
}

func (c *ClientWithResponses) AuthenticateWithResponse(ctx context.Context, body AuthenticateJSONRequestBody, reqEditors ...RequestEditorFn) (*AuthenticateResponse, error) {
	rsp, err := c.Authenticate(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAuthenticateResponse(rsp)
}

// HealthCheckWithResponse request returning *HealthCheckResponse
func (c *ClientWithResponses) HealthCheckWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthCheckResponse, error) {
	rsp, err := c.HealthCheck(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseHealthCheckResponse(rsp)
}

// KeystoreWithResponse request returning *KeystoreResponse
func (c *ClientWithResponses) KeystoreWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*KeystoreResponse, error) {
	rsp, err := c.Keystore(ctx, keystoreId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseKeystoreResponse(rsp)
}

// KeystoreEntriesWithResponse request returning *KeystoreEntriesResponse
func (c *ClientWithResponses) KeystoreEntriesWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*KeystoreEntriesResponse, error) {
	rsp, err := c.KeystoreEntries(ctx, keystoreId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseKeystoreEntriesResponse(rsp)
}

// CreateEntryWithBodyWithResponse request with arbitrary body returning *CreateEntryResponse
func (c *ClientWithResponses) CreateEntryWithBodyWithResponse(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateEntryResponse, error) {
	rsp, err := c.CreateEntryWithBody(ctx, keystoreId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateEntryResponse(rsp)
}

func (c *ClientWithResponses) CreateEntryWithResponse(ctx context.Context, keystoreId string, body CreateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateEntryResponse, error) {
	rsp, err := c.CreateEntry(ctx, keystoreId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateEntryResponse(rsp)
}

// UpdateEntryWithBodyWithResponse request with arbitrary body returning *UpdateEntryResponse
func (c *ClientWithResponses) UpdateEntryWithBodyWithResponse(ctx context.Context, keystoreId string, entryId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateEntryResponse, error) {
	rsp, err := c.UpdateEntryWithBody(ctx, keystoreId, entryId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateEntryResponse(rsp)
}

func (c *ClientWithResponses) UpdateEntryWithResponse(ctx context.Context, keystoreId string, entryId string, body UpdateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateEntryResponse, error) {
	rsp, err := c.UpdateEntry(ctx, keystoreId, entryId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateEntryResponse(rsp)
}

// KeystoresWithResponse request returning *KeystoresResponse
func (c *ClientWithResponses) KeystoresWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*KeystoresResponse, error) {
	rsp, err := c.Keystores(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseKeystoresResponse(rsp)
}

// CreateKeystoreWithBodyWithResponse request with arbitrary body returning *CreateKeystoreResponse
func (c *ClientWithResponses) CreateKeystoreWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateKeystoreResponse, error) {
	rsp, err := c.CreateKeystoreWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateKeystoreResponse(rsp)
}

func (c *ClientWithResponses) CreateKeystoreWithResponse(ctx context.Context, body CreateKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateKeystoreResponse, error) {
	rsp, err := c.CreateKeystore(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateKeystoreResponse(rsp)
}

// ParseAuthenticateResponse parses an HTTP response from a AuthenticateWithResponse call
func ParseAuthenticateResponse(rsp *http.Response) (*AuthenticateResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AuthenticateResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseHealthCheckResponse parses an HTTP response from a HealthCheckWithResponse call
func ParseHealthCheckResponse(rsp *http.Response) (*HealthCheckResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &HealthCheckResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseKeystoreResponse parses an HTTP response from a KeystoreWithResponse call
func ParseKeystoreResponse(rsp *http.Response) (*KeystoreResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &KeystoreResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Keystore
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseKeystoreEntriesResponse parses an HTTP response from a KeystoreEntriesWithResponse call
func ParseKeystoreEntriesResponse(rsp *http.Response) (*KeystoreEntriesResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &KeystoreEntriesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Entry
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseCreateEntryResponse parses an HTTP response from a CreateEntryWithResponse call
func ParseCreateEntryResponse(rsp *http.Response) (*CreateEntryResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateEntryResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Entry
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseUpdateEntryResponse parses an HTTP response from a UpdateEntryWithResponse call
func ParseUpdateEntryResponse(rsp *http.Response) (*UpdateEntryResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateEntryResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Entry
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseKeystoresResponse parses an HTTP response from a KeystoresWithResponse call
func ParseKeystoresResponse(rsp *http.Response) (*KeystoresResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &KeystoresResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Keystores
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseCreateKeystoreResponse parses an HTTP response from a CreateKeystoreWithResponse call
func ParseCreateKeystoreResponse(rsp *http.Response) (*CreateKeystoreResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateKeystoreResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Keystore
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
