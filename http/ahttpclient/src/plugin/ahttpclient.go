//*****************************************************************************************************
//# Author        :   D. Ajith Nilanta de Silva  | 12/12/2024 ajithdesilva@gmail.com
//# Class/module  :   ahttpclient.go
//# Objective     :   Implement the http client library using github.com/valyala/fasthttp
//#					Uses the common data types defined in the src/httpclient/httpclient.go as parameters
//#######################################################################################################
//# Author                        Date        Action      Description
//#------------------------------------------------------------------------------------------------------
//# Ajith de Silva				01/12/2004	Created 	Created the initial version
//# Ajith de Silva				03/12/2004	Updated 	implemented main functions
//# Ajith de Silva				03/12/2004	Updated		Added content decoding
//# Ajith de Silva				05/12/2004	Updated		Added error handling
//#######################################################################################################

package main

import (
	"agnione/v1/src/afplugins/http/iahttpclient"
	httptypes "agnione/v1/src/afplugins/http/types"
	atypes "agnione/v1/src/appfm/types"
	build "agnione/v1/src/lib"
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"time"

	libbuild "ahttp.client/src/build"

	"github.com/valyala/fasthttp"
)

// AHTTPClient implements all the defined methods in the iahttpclient
type AHTTPClient struct {
	id int
}

func (ahttp *AHTTPClient) content_decode(pResponse *fasthttp.Response) []byte {

	/// Do we need to decompress the response?
	/// if we have the Header then do it
	_contentEncoding := pResponse.Header.Peek("Content-Encoding")
	
	defer func ()  {
		_contentEncoding=nil
	}()
	
	var _body []byte
	
	if bytes.EqualFold(_contentEncoding, []byte("gzip")) {
		_body,_= pResponse.BodyGunzip()
	} else {
		_body=pResponse.Body()
	}
	
	return _body
	
}

func (ahttp *AHTTPClient) add_request_headers(http_request httptypes.AHTTPRequest, request *fasthttp.Request) *fasthttp.Request {
	
	if http_request.Headers==nil{
		return request
	}
	
	var _key string 
	var _val string
	
	if len(http_request.Headers) > 0 {
		for _key, _val= range http_request.Headers {
			request.Header.Add(_key, _val)
		}
	}
	
	return request
}

func (ahttp *AHTTPClient) httpConnError(err error) (string, bool) {
	var (
		errName string
		known   = true
	)

	switch {
	case errors.Is(err, fasthttp.ErrTimeout):
		errName = "timeout"
	case errors.Is(err, fasthttp.ErrNoFreeConns):
		errName = "conn_limit"
	case errors.Is(err, fasthttp.ErrConnectionClosed):
		errName = "connection close"
	case reflect.TypeOf(err).String() == "*net.OpError":
		errName = "timeout"
	default:
		known = false
	}

	return errName, known
}

func (ahttp *AHTTPClient) do_request(pHTTP_Request *httptypes.AHTTPRequest, http_method string) (*httptypes.AHTTPResponse, error) {
	
	if len(pHTTP_Request.URL) == 0 {
		return nil, fmt.Errorf("invalid url. please provide the valid http url")
	}

	_hrequest := fasthttp.AcquireRequest()   /// get request instance
	_hresponse := fasthttp.AcquireResponse() // Acquire a response instance
	var _err error

	defer func() {
		fasthttp.ReleaseRequest(_hrequest)
		fasthttp.ReleaseResponse(_hresponse)
		_err = nil
		_hrequest = nil
		_hresponse = nil
	}()

	_hrequest.SetRequestURI(pHTTP_Request.URL) /// sets the request URL

	_hrequest.Header.SetMethod(http_method) /// sets the http method

	/// check what HTTP request have been requested
	switch http_method {

	/// try to set body if HTTP methods are POST,PUT and DELETE
	case fasthttp.MethodPost, fasthttp.MethodPut, fasthttp.MethodDelete:
		{
			/// if we have body then set it
			if len(pHTTP_Request.Body) > 0 {
				_hrequest.SetBodyRaw(pHTTP_Request.Body)
			}
		}
	}

	/// insert HTTP headers to the request. (if provided)
	_hrequest = ahttp.add_request_headers(*pHTTP_Request, _hrequest)

	/// perform the request based on the timeout
	if pHTTP_Request.Timeout > 0 {
		_err = fasthttp.DoTimeout(_hrequest, _hresponse,
			time.Duration(pHTTP_Request.Timeout*int(time.Millisecond)))
	} else {
		_err = fasthttp.Do(_hrequest, _hresponse)
	}
	
	if _err!=nil{
		return nil, errors.New("request failed: " +  _err.Error())
	}
	
	/// First.. fetch the response status code and body to the response.
	_response := httptypes.AHTTPResponse{StatusCode: _hresponse.StatusCode(),
		Body: ahttp.content_decode(_hresponse),
	}

	/// add the response headers
	_bHKeys := _hresponse.Header.PeekKeys()
	
	defer func(){
		_bHKeys = nil
	}()
	
	_response.Headers = make(map[string]string,len(_bHKeys))
	
	var _key int
	for _key = range _bHKeys {
		_response.Headers[string(_bHKeys[_key])] = string(_hresponse.Header.Peek(string(_bHKeys[_key])))
	}

	///Second. If we have error then work on it
	if _err != nil {
		_errName, _ := ahttp.httpConnError(_err)
		return nil, errors.New("request failed: " + _errName  + ". " +  _err.Error()) 
	} else {
		return &_response, nil /// return response with nil error
	}
}

// Cretes a new isntance of IAHTTPClient
func (ahttp *AHTTPClient) New() interface{} {

	return atypes.ConvertStoI[iahttpclient.IAHTTPClient](new(AHTTPClient))
}

// Initialize initializes the given id to the instance. Used to identify the intance
// Returns true if success. Unless false
func (ahttp *AHTTPClient) Initialize(pInstance_ID int) bool {
	ahttp.id = pInstance_ID
	return true
}

// GetID retuns the pre-set id of the current instance
func (ahttp *AHTTPClient) GetID() (instance_id int) {
	return ahttp.id
}

// Get perfoms a HTTP GET request based on the given ZHTTPRequest.
// If success returns HTTPResponse with StatusCode 200 and Body with result data in []bytes
// If failed then returns HTTPResponse with valid status code and error with menaningful error
func (ahttp *AHTTPClient) Get (pHTTP_Request *httptypes.AHTTPRequest) (*httptypes.AHTTPResponse, error){

	return ahttp.do_request(pHTTP_Request, fasthttp.MethodGet)

}

// Post perfoms a HTTP POST request based on the given ZHTTPRequest.
// If success, returns HTTPResponse with StatusCode 200 and Body with result data in []bytes
// If failed then returns HTTPResponse with valid status code and error with menaningful error
func (ahttp *AHTTPClient) Post(pHTTP_Request *httptypes.AHTTPRequest) (*httptypes.AHTTPResponse, error) {

	return ahttp.do_request(pHTTP_Request, fasthttp.MethodPost)

}

// Put perfoms a HTTP PUT request based on the given ZHTTPRequest.
// If success, returns HTTPResponse with StatusCode 200 and Body with result data in []bytes
// If failed then returns HTTPResponse with valid status code and error with menaningful error
func (ahttp *AHTTPClient) Put(pHTTP_Request *httptypes.AHTTPRequest) (*httptypes.AHTTPResponse, error) {

	return ahttp.do_request(pHTTP_Request, fasthttp.MethodPut)

}

// Delete perfoms a HTTP DELETE request based on the given ZHTTPRequest.
// If success, returns HTTPResponse with StatusCode 200 and Body with result data in []bytes
// If failed then returns HTTPResponse with valid status code and error with menaningful error
func (ahttp *AHTTPClient) Delete(pHTTP_Request *httptypes.AHTTPRequest) (*httptypes.AHTTPResponse, error) {

	return ahttp.do_request(pHTTP_Request, fasthttp.MethodDelete)

}

func (ahttp *AHTTPClient) Info() build.BuildInfo {
	return build.BuildInfo{Time: libbuild.Time, User: libbuild.User, Version: libbuild.Version, BuildGoVersion: libbuild.BuildGoVersion}
}

var IAHTTPClient AHTTPClient
