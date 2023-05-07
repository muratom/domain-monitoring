// Package inspector provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	. "github.com/muratom/domain-monitoring/api/rpc/v1/inspector/models"
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
	// AddDomain request
	AddDomain(ctx context.Context, params *AddDomainParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAllDomains request
	GetAllDomains(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetChangelog request
	GetChangelog(ctx context.Context, params *GetChangelogParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteDomain request
	DeleteDomain(ctx context.Context, params *DeleteDomainParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetDomain request
	GetDomain(ctx context.Context, params *GetDomainParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Ping request
	Ping(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateDomain request
	UpdateDomain(ctx context.Context, params *UpdateDomainParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) AddDomain(ctx context.Context, params *AddDomainParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddDomainRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAllDomains(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAllDomainsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetChangelog(ctx context.Context, params *GetChangelogParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetChangelogRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteDomain(ctx context.Context, params *DeleteDomainParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteDomainRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetDomain(ctx context.Context, params *GetDomainParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetDomainRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Ping(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPingRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateDomain(ctx context.Context, params *UpdateDomainParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateDomainRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewAddDomainRequest generates requests for AddDomain
func NewAddDomainRequest(server string, params *AddDomainParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/add-domain")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "fqdn", runtime.ParamLocationQuery, params.Fqdn); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetAllDomainsRequest generates requests for GetAllDomains
func NewGetAllDomainsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/all-domains")
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

// NewGetChangelogRequest generates requests for GetChangelog
func NewGetChangelogRequest(server string, params *GetChangelogParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/changelog")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "fqdn", runtime.ParamLocationQuery, params.Fqdn); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewDeleteDomainRequest generates requests for DeleteDomain
func NewDeleteDomainRequest(server string, params *DeleteDomainParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/delete-domain")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "fqdn", runtime.ParamLocationQuery, params.Fqdn); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetDomainRequest generates requests for GetDomain
func NewGetDomainRequest(server string, params *GetDomainParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/domain")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "fqdn", runtime.ParamLocationQuery, params.Fqdn); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPingRequest generates requests for Ping
func NewPingRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/ping")
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

// NewUpdateDomainRequest generates requests for UpdateDomain
func NewUpdateDomainRequest(server string, params *UpdateDomainParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/update-domain")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "fqdn", runtime.ParamLocationQuery, params.Fqdn); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

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
	// AddDomain request
	AddDomainWithResponse(ctx context.Context, params *AddDomainParams, reqEditors ...RequestEditorFn) (*AddDomainResponse, error)

	// GetAllDomains request
	GetAllDomainsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAllDomainsResponse, error)

	// GetChangelog request
	GetChangelogWithResponse(ctx context.Context, params *GetChangelogParams, reqEditors ...RequestEditorFn) (*GetChangelogResponse, error)

	// DeleteDomain request
	DeleteDomainWithResponse(ctx context.Context, params *DeleteDomainParams, reqEditors ...RequestEditorFn) (*DeleteDomainResponse, error)

	// GetDomain request
	GetDomainWithResponse(ctx context.Context, params *GetDomainParams, reqEditors ...RequestEditorFn) (*GetDomainResponse, error)

	// Ping request
	PingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingResponse, error)

	// UpdateDomain request
	UpdateDomainWithResponse(ctx context.Context, params *UpdateDomainParams, reqEditors ...RequestEditorFn) (*UpdateDomainResponse, error)
}

type AddDomainResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Domain
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r AddDomainResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddDomainResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAllDomainsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]string
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetAllDomainsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAllDomainsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetChangelogResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Changelog
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetChangelogResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetChangelogResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteDomainResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r DeleteDomainResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteDomainResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetDomainResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Domain
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetDomainResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetDomainResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r PingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateDomainResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Domain
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r UpdateDomainResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateDomainResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// AddDomainWithResponse request returning *AddDomainResponse
func (c *ClientWithResponses) AddDomainWithResponse(ctx context.Context, params *AddDomainParams, reqEditors ...RequestEditorFn) (*AddDomainResponse, error) {
	rsp, err := c.AddDomain(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddDomainResponse(rsp)
}

// GetAllDomainsWithResponse request returning *GetAllDomainsResponse
func (c *ClientWithResponses) GetAllDomainsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAllDomainsResponse, error) {
	rsp, err := c.GetAllDomains(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAllDomainsResponse(rsp)
}

// GetChangelogWithResponse request returning *GetChangelogResponse
func (c *ClientWithResponses) GetChangelogWithResponse(ctx context.Context, params *GetChangelogParams, reqEditors ...RequestEditorFn) (*GetChangelogResponse, error) {
	rsp, err := c.GetChangelog(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetChangelogResponse(rsp)
}

// DeleteDomainWithResponse request returning *DeleteDomainResponse
func (c *ClientWithResponses) DeleteDomainWithResponse(ctx context.Context, params *DeleteDomainParams, reqEditors ...RequestEditorFn) (*DeleteDomainResponse, error) {
	rsp, err := c.DeleteDomain(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteDomainResponse(rsp)
}

// GetDomainWithResponse request returning *GetDomainResponse
func (c *ClientWithResponses) GetDomainWithResponse(ctx context.Context, params *GetDomainParams, reqEditors ...RequestEditorFn) (*GetDomainResponse, error) {
	rsp, err := c.GetDomain(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetDomainResponse(rsp)
}

// PingWithResponse request returning *PingResponse
func (c *ClientWithResponses) PingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingResponse, error) {
	rsp, err := c.Ping(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePingResponse(rsp)
}

// UpdateDomainWithResponse request returning *UpdateDomainResponse
func (c *ClientWithResponses) UpdateDomainWithResponse(ctx context.Context, params *UpdateDomainParams, reqEditors ...RequestEditorFn) (*UpdateDomainResponse, error) {
	rsp, err := c.UpdateDomain(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateDomainResponse(rsp)
}

// ParseAddDomainResponse parses an HTTP response from a AddDomainWithResponse call
func ParseAddDomainResponse(rsp *http.Response) (*AddDomainResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AddDomainResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Domain
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

// ParseGetAllDomainsResponse parses an HTTP response from a GetAllDomainsWithResponse call
func ParseGetAllDomainsResponse(rsp *http.Response) (*GetAllDomainsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAllDomainsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []string
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

// ParseGetChangelogResponse parses an HTTP response from a GetChangelogWithResponse call
func ParseGetChangelogResponse(rsp *http.Response) (*GetChangelogResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetChangelogResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Changelog
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

// ParseDeleteDomainResponse parses an HTTP response from a DeleteDomainWithResponse call
func ParseDeleteDomainResponse(rsp *http.Response) (*DeleteDomainResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteDomainResponse{
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

// ParseGetDomainResponse parses an HTTP response from a GetDomainWithResponse call
func ParseGetDomainResponse(rsp *http.Response) (*GetDomainResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetDomainResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Domain
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

// ParsePingResponse parses an HTTP response from a PingWithResponse call
func ParsePingResponse(rsp *http.Response) (*PingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PingResponse{
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

// ParseUpdateDomainResponse parses an HTTP response from a UpdateDomainWithResponse call
func ParseUpdateDomainResponse(rsp *http.Response) (*UpdateDomainResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateDomainResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Domain
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
