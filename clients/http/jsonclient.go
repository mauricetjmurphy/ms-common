package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

// Doer is alias by *(net/http).Client, which it is easy and safe to do JSON requests to REST services.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// JSONClient presents a way of easy and safe to do JSON requests to REST services.
type JSONClient struct {
	baseURL    string
	httpClient Doer
}

// JSONRequest presents the request details.
type JSONRequest struct {
	Method string
	Path   string
	Query  url.Values
	Body   interface{}
	Header http.Header
}

// NewJSONClient create the HTTP JSON client instance.
// Usage:
//
//	    httpClient := New(&http.ClientConf{})
//	    client := NewJSONClient("http://localhost:8000", httpClient)
//
//	    1) Call request with no options:
//	   	client.GET(ctx, path,query, output)
//
//	    2) Call request with header options.
//	    	ctx2 := http.WithReqOpts(ctx, http.ReqOptions{
//				Header: map[string]string{
//					"APIKEY": "GVS_NBCU_TOKEN",
//				},
//			})
//	    client.GET(ctx2, path,query, output)
func NewJSONClient(baseURL string, httpClient Doer) *JSONClient {
	return &JSONClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

// NewDefaultJSONClient creates the default HTTP JSON Client
func NewDefaultJSONClient(baseURL string) *JSONClient {
	return NewJSONClient(baseURL, NewDefaultClient())
}

func (c *JSONClient) Get(ctx context.Context, path string, query url.Values, output interface{}) (*http.Response, error) {
	return c.DoWithJSON(ctx, JSONRequest{
		Method: http.MethodGet,
		Path:   path,
		Query:  query,
	}, output)
}

func (c *JSONClient) Post(ctx context.Context, path string, body interface{}, output interface{}) (*http.Response, error) {
	return c.DoWithJSON(ctx, JSONRequest{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	}, output)
}

func (c *JSONClient) Put(ctx context.Context, path string, body interface{}, output interface{}) (*http.Response, error) {
	return c.DoWithJSON(ctx, JSONRequest{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	}, output)
}

func (c *JSONClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.DoWithJSON(ctx, JSONRequest{
		Method: http.MethodDelete,
		Path:   path,
	}, nil)
}

func (c *JSONClient) Do(ctx context.Context, jsonReq JSONRequest, reqBody io.Reader, output interface{}) (response *http.Response, err error) {
	urlStr := c.urlStr(jsonReq.Path, jsonReq.Query)

	request, err := http.NewRequest(jsonReq.Method, urlStr, reqBody)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	addRequestOpts(ctx, jsonReq, request)

	response, err = c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		e := response.Body.Close()
		if e != nil {
			err = e
		}
	}()

	if IsErrorStatus(response) {
		return response, ParseJSONErr(response)
	}

	if output != nil {
		err = json.NewDecoder(response.Body).Decode(output)
		return response, err
	}

	return response, nil
}

func (c *JSONClient) DoWithJSON(ctx context.Context, jsonReq JSONRequest, output interface{}) (*http.Response, error) {
	reqBody, err := serializeJSON(jsonReq.Body)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, jsonReq, reqBody, output)
}

func addRequestOpts(ctx context.Context, jsonReq JSONRequest, request *http.Request) {
	header := jsonReq.Header
	if header == nil {
		header = http.Header{ContentType: {ContentTypeJSON}}
	}
	if header.Get(ContentType) == "" {
		header.Set(ContentType, ContentTypeJSON)
	}

	reqOpts := GetReqOpts(ctx)
	for k, v := range reqOpts.Header {
		header.Set(k, v)
	}

	request.Header = header
}

func serializeJSON(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonBytes), nil
}

func (c *JSONClient) urlStr(path string, query url.Values) string {
	if path == "" {
		return ""
	}

	sep := ""
	if path[0:1] != "/" {
		sep = "/"
	}

	requestURL := c.baseURL + sep + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}
	return requestURL
}
