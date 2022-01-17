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
	// Login request  with any body
	LoginWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Login(ctx context.Context, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Keystore request
	Keystore(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Invitation request
	Invitation(ctx context.Context, keystoreId string, invitationId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AcceptInvitation request  with any body
	AcceptInvitationWithBody(ctx context.Context, keystoreId string, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AcceptInvitation(ctx context.Context, keystoreId string, invitationId string, body AcceptInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// FinalizeInvitation request  with any body
	FinalizeInvitationWithBody(ctx context.Context, keystoreId string, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	FinalizeInvitation(ctx context.Context, keystoreId string, invitationId string, body FinalizeInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateInvitation request  with any body
	CreateInvitationWithBody(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateInvitation(ctx context.Context, keystoreId string, body CreateInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Keystores request
	Keystores(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateKeystore request  with any body
	CreateKeystoreWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateKeystore(ctx context.Context, body CreateKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) LoginWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLoginRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Login(ctx context.Context, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLoginRequest(c.Server, body)
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

func (c *Client) Invitation(ctx context.Context, keystoreId string, invitationId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewInvitationRequest(c.Server, keystoreId, invitationId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AcceptInvitationWithBody(ctx context.Context, keystoreId string, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAcceptInvitationRequestWithBody(c.Server, keystoreId, invitationId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AcceptInvitation(ctx context.Context, keystoreId string, invitationId string, body AcceptInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAcceptInvitationRequest(c.Server, keystoreId, invitationId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) FinalizeInvitationWithBody(ctx context.Context, keystoreId string, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewFinalizeInvitationRequestWithBody(c.Server, keystoreId, invitationId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) FinalizeInvitation(ctx context.Context, keystoreId string, invitationId string, body FinalizeInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewFinalizeInvitationRequest(c.Server, keystoreId, invitationId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateInvitationWithBody(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateInvitationRequestWithBody(c.Server, keystoreId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateInvitation(ctx context.Context, keystoreId string, body CreateInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateInvitationRequest(c.Server, keystoreId, body)
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

// NewLoginRequest calls the generic Login builder with application/json body
func NewLoginRequest(server string, body LoginJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewLoginRequestWithBody(server, "application/json", bodyReader)
}

// NewLoginRequestWithBody generates requests for Login with any type of body
func NewLoginRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/auth/login")
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

	operationPath := fmt.Sprintf("/api/keystore/%s", pathParam0)
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

// NewInvitationRequest generates requests for Invitation
func NewInvitationRequest(server string, keystoreId string, invitationId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "keystoreId", runtime.ParamLocationPath, keystoreId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "invitationId", runtime.ParamLocationPath, invitationId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/keystore/%s/invitation/%s", pathParam0, pathParam1)
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

// NewAcceptInvitationRequest calls the generic AcceptInvitation builder with application/json body
func NewAcceptInvitationRequest(server string, keystoreId string, invitationId string, body AcceptInvitationJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAcceptInvitationRequestWithBody(server, keystoreId, invitationId, "application/json", bodyReader)
}

// NewAcceptInvitationRequestWithBody generates requests for AcceptInvitation with any type of body
func NewAcceptInvitationRequestWithBody(server string, keystoreId string, invitationId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "keystoreId", runtime.ParamLocationPath, keystoreId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "invitationId", runtime.ParamLocationPath, invitationId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/keystore/%s/invitation/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewFinalizeInvitationRequest calls the generic FinalizeInvitation builder with application/json body
func NewFinalizeInvitationRequest(server string, keystoreId string, invitationId string, body FinalizeInvitationJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewFinalizeInvitationRequestWithBody(server, keystoreId, invitationId, "application/json", bodyReader)
}

// NewFinalizeInvitationRequestWithBody generates requests for FinalizeInvitation with any type of body
func NewFinalizeInvitationRequestWithBody(server string, keystoreId string, invitationId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "keystoreId", runtime.ParamLocationPath, keystoreId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "invitationId", runtime.ParamLocationPath, invitationId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/keystore/%s/invitation/%s", pathParam0, pathParam1)
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

// NewCreateInvitationRequest calls the generic CreateInvitation builder with application/json body
func NewCreateInvitationRequest(server string, keystoreId string, body CreateInvitationJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateInvitationRequestWithBody(server, keystoreId, "application/json", bodyReader)
}

// NewCreateInvitationRequestWithBody generates requests for CreateInvitation with any type of body
func NewCreateInvitationRequestWithBody(server string, keystoreId string, contentType string, body io.Reader) (*http.Request, error) {
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

	operationPath := fmt.Sprintf("/api/keystore/%s/invitations", pathParam0)
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

	operationPath := fmt.Sprintf("/api/keystores")
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

	operationPath := fmt.Sprintf("/api/keystores")
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
	// Login request  with any body
	LoginWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LoginResponse, error)

	LoginWithResponse(ctx context.Context, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*LoginResponse, error)

	// Keystore request
	KeystoreWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*KeystoreResponse, error)

	// Invitation request
	InvitationWithResponse(ctx context.Context, keystoreId string, invitationId string, reqEditors ...RequestEditorFn) (*InvitationResponse, error)

	// AcceptInvitation request  with any body
	AcceptInvitationWithBodyWithResponse(ctx context.Context, keystoreId string, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AcceptInvitationResponse, error)

	AcceptInvitationWithResponse(ctx context.Context, keystoreId string, invitationId string, body AcceptInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*AcceptInvitationResponse, error)

	// FinalizeInvitation request  with any body
	FinalizeInvitationWithBodyWithResponse(ctx context.Context, keystoreId string, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*FinalizeInvitationResponse, error)

	FinalizeInvitationWithResponse(ctx context.Context, keystoreId string, invitationId string, body FinalizeInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*FinalizeInvitationResponse, error)

	// CreateInvitation request  with any body
	CreateInvitationWithBodyWithResponse(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateInvitationResponse, error)

	CreateInvitationWithResponse(ctx context.Context, keystoreId string, body CreateInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateInvitationResponse, error)

	// Keystores request
	KeystoresWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*KeystoresResponse, error)

	// CreateKeystore request  with any body
	CreateKeystoreWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateKeystoreResponse, error)

	CreateKeystoreWithResponse(ctx context.Context, body CreateKeystoreJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateKeystoreResponse, error)
}

type LoginResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *AuthToken
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r LoginResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r LoginResponse) StatusCode() int {
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

type InvitationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Invitation
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r InvitationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r InvitationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AcceptInvitationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Invitation
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r AcceptInvitationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AcceptInvitationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type FinalizeInvitationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Invitation
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r FinalizeInvitationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r FinalizeInvitationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateInvitationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Invitation
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r CreateInvitationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateInvitationResponse) StatusCode() int {
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

// LoginWithBodyWithResponse request with arbitrary body returning *LoginResponse
func (c *ClientWithResponses) LoginWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LoginResponse, error) {
	rsp, err := c.LoginWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLoginResponse(rsp)
}

func (c *ClientWithResponses) LoginWithResponse(ctx context.Context, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*LoginResponse, error) {
	rsp, err := c.Login(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLoginResponse(rsp)
}

// KeystoreWithResponse request returning *KeystoreResponse
func (c *ClientWithResponses) KeystoreWithResponse(ctx context.Context, keystoreId string, reqEditors ...RequestEditorFn) (*KeystoreResponse, error) {
	rsp, err := c.Keystore(ctx, keystoreId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseKeystoreResponse(rsp)
}

// InvitationWithResponse request returning *InvitationResponse
func (c *ClientWithResponses) InvitationWithResponse(ctx context.Context, keystoreId string, invitationId string, reqEditors ...RequestEditorFn) (*InvitationResponse, error) {
	rsp, err := c.Invitation(ctx, keystoreId, invitationId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseInvitationResponse(rsp)
}

// AcceptInvitationWithBodyWithResponse request with arbitrary body returning *AcceptInvitationResponse
func (c *ClientWithResponses) AcceptInvitationWithBodyWithResponse(ctx context.Context, keystoreId string, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AcceptInvitationResponse, error) {
	rsp, err := c.AcceptInvitationWithBody(ctx, keystoreId, invitationId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAcceptInvitationResponse(rsp)
}

func (c *ClientWithResponses) AcceptInvitationWithResponse(ctx context.Context, keystoreId string, invitationId string, body AcceptInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*AcceptInvitationResponse, error) {
	rsp, err := c.AcceptInvitation(ctx, keystoreId, invitationId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAcceptInvitationResponse(rsp)
}

// FinalizeInvitationWithBodyWithResponse request with arbitrary body returning *FinalizeInvitationResponse
func (c *ClientWithResponses) FinalizeInvitationWithBodyWithResponse(ctx context.Context, keystoreId string, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*FinalizeInvitationResponse, error) {
	rsp, err := c.FinalizeInvitationWithBody(ctx, keystoreId, invitationId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseFinalizeInvitationResponse(rsp)
}

func (c *ClientWithResponses) FinalizeInvitationWithResponse(ctx context.Context, keystoreId string, invitationId string, body FinalizeInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*FinalizeInvitationResponse, error) {
	rsp, err := c.FinalizeInvitation(ctx, keystoreId, invitationId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseFinalizeInvitationResponse(rsp)
}

// CreateInvitationWithBodyWithResponse request with arbitrary body returning *CreateInvitationResponse
func (c *ClientWithResponses) CreateInvitationWithBodyWithResponse(ctx context.Context, keystoreId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateInvitationResponse, error) {
	rsp, err := c.CreateInvitationWithBody(ctx, keystoreId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateInvitationResponse(rsp)
}

func (c *ClientWithResponses) CreateInvitationWithResponse(ctx context.Context, keystoreId string, body CreateInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateInvitationResponse, error) {
	rsp, err := c.CreateInvitation(ctx, keystoreId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateInvitationResponse(rsp)
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

// ParseLoginResponse parses an HTTP response from a LoginWithResponse call
func ParseLoginResponse(rsp *http.Response) (*LoginResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &LoginResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest AuthToken
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

// ParseInvitationResponse parses an HTTP response from a InvitationWithResponse call
func ParseInvitationResponse(rsp *http.Response) (*InvitationResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &InvitationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Invitation
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

// ParseAcceptInvitationResponse parses an HTTP response from a AcceptInvitationWithResponse call
func ParseAcceptInvitationResponse(rsp *http.Response) (*AcceptInvitationResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &AcceptInvitationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Invitation
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

// ParseFinalizeInvitationResponse parses an HTTP response from a FinalizeInvitationWithResponse call
func ParseFinalizeInvitationResponse(rsp *http.Response) (*FinalizeInvitationResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &FinalizeInvitationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Invitation
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

// ParseCreateInvitationResponse parses an HTTP response from a CreateInvitationWithResponse call
func ParseCreateInvitationResponse(rsp *http.Response) (*CreateInvitationResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &CreateInvitationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Invitation
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
