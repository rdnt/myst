// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
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
	// Keystore request
	Keystore(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UnlockKeystore request  with any body
	UnlockKeystoreWithBody(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UnlockKeystore(ctx context.Context, keystoreId string, body UnlockKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetKeystore request
	GetKeystore(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateEntry request  with any body
	CreateEntryWithBody(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateEntry(ctx context.Context, keystoreId string, body CreateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Keystores request
	Keystores(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateKeystore request  with any body
	CreateKeystoreWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateKeystore(ctx context.Context, body CreateKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
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

func (c *Client) UnlockKeystoreWithBody(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUnlockKeystoreRequestWithBody(c.Server, keystoreId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UnlockKeystore(ctx context.Context, keystoreId string, body UnlockKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUnlockKeystoreRequest(c.Server, keystoreId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetKeystore(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetKeystoreRequest(c.Server, keystoreId)
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
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUnlockKeystoreRequest calls the generic UnlockKeystore builder with application/json body
func NewUnlockKeystoreRequest(server string, keystoreId string, body UnlockKeystoreJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUnlockKeystoreRequestWithBody(server, keystoreId, "application/json", bodyReader)
}

// NewUnlockKeystoreRequestWithBody generates requests for UnlockKeystore with any type of body
func NewUnlockKeystoreRequestWithBody(server string, keystoreId string, contentType string, body io.Reader) (*http.Request, error) {
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
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetKeystoreRequest generates requests for GetKeystore
func NewGetKeystoreRequest(server string, keystoreId string) (*http.Request, error) {
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
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

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
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("POST", queryURL.String(), body)
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
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

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
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

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
	// Keystore request
	KeystoreWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*KeystoreResponse, error)

	// UnlockKeystore request  with any body
	UnlockKeystoreWithBodyWithResponse(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UnlockKeystoreResponse, error)

	UnlockKeystoreWithResponse(ctx context.Context, keystoreId string, body UnlockKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*UnlockKeystoreResponse, error)

	// GetKeystore request
	GetKeystoreWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*GetKeystoreResponse, error)

	// CreateEntry request  with any body
	CreateEntryWithBodyWithResponse(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateEntryResponse, error)

	CreateEntryWithResponse(ctx context.Context, keystoreId string, body CreateEntryJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateEntryResponse, error)

	// Keystores request
	KeystoresWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*KeystoresResponse, error)

	// CreateKeystore request  with any body
	CreateKeystoreWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateKeystoreResponse, error)

	CreateKeystoreWithResponse(ctx context.Context, body CreateKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateKeystoreResponse, error)
}

type KeystoreResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Keystore
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

type UnlockKeystoreResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Keystore
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r UnlockKeystoreResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UnlockKeystoreResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetKeystoreResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Keystore
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetKeystoreResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetKeystoreResponse) StatusCode() int {
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

type KeystoresResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Keystore
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
	JSON200      *[]Keystore
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

// KeystoreWithResponse request returning *KeystoreResponse
func (c *ClientWithResponses) KeystoreWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*KeystoreResponse, error) {
	rsp, err := c.Keystore(ctx, keystoreId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseKeystoreResponse(rsp)
}

// UnlockKeystoreWithBodyWithResponse request with arbitrary body returning *UnlockKeystoreResponse
func (c *ClientWithResponses) UnlockKeystoreWithBodyWithResponse(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UnlockKeystoreResponse, error) {
	rsp, err := c.UnlockKeystoreWithBody(ctx, keystoreId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUnlockKeystoreResponse(rsp)
}

func (c *ClientWithResponses) UnlockKeystoreWithResponse(ctx context.Context, keystoreId string, body UnlockKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*UnlockKeystoreResponse, error) {
	rsp, err := c.UnlockKeystore(ctx, keystoreId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUnlockKeystoreResponse(rsp)
}

// GetKeystoreWithResponse request returning *GetKeystoreResponse
func (c *ClientWithResponses) GetKeystoreWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*GetKeystoreResponse, error) {
	rsp, err := c.GetKeystore(ctx, keystoreId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetKeystoreResponse(rsp)
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

// ParseKeystoreResponse parses an HTTP response from a KeystoreWithResponse call
func ParseKeystoreResponse(rsp *http.Response) (*KeystoreResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &KeystoreResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Keystore
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

// ParseUnlockKeystoreResponse parses an HTTP response from a UnlockKeystoreWithResponse call
func ParseUnlockKeystoreResponse(rsp *http.Response) (*UnlockKeystoreResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &UnlockKeystoreResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Keystore
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

// ParseGetKeystoreResponse parses an HTTP response from a GetKeystoreWithResponse call
func ParseGetKeystoreResponse(rsp *http.Response) (*GetKeystoreResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetKeystoreResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Keystore
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
	defer rsp.Body.Close()
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

// ParseKeystoresResponse parses an HTTP response from a KeystoresWithResponse call
func ParseKeystoresResponse(rsp *http.Response) (*KeystoresResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &KeystoresResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Keystore
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
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &CreateKeystoreResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Keystore
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
