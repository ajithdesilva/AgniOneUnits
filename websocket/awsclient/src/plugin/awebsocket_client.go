/*
Author        :   D. Ajith Nilanta de Silva ajithdesilva@gmail.com | 02/02/2024

Class/module  :   awebsocket_client.go

Objective     :   Implement the web socket client using gorilla/websocket library based on the IAWSClient interface

	This package will be exported as plugin/library to the Application units.
	That will helps application units to implemnts it's features by using the framework suport

#########################################################################################

	Author			Date        	Action      	Descriptionrr

#########################################################################################

	Ajith de Silva		06/11/2004	Created 	Created the initial version

	Ajith de Silva		06/11/2004	Updated 	Implemented main functions

	Ajith de Silva		12/12/2004	Updated		optimized the memory usage

#########################################################################################
*/
package main

import (
	iawsclient "agnione/v1/src/afplugins/websocket/iawsclient"
	atypes "agnione/v1/src/appfm/types"
	build "agnione/v1/src/lib"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	buildinfo "awebsocket.client/src/build"

	"github.com/gorilla/websocket"
)

// AgniWSClient main structure to hold the AWSClient instance
type AgniWSClient struct {
	wsCon       *websocket.Conn
	wsURL       string
	id          int
	isConnected bool
}

// New creats a new insance of AWSClient and converts it into IZWSClient
//
// Returns the coverted IAWSClient interface
func (awsc *AgniWSClient) New() interface{} {

	return atypes.ConvertStoI[iawsclient.IAWSClient](new(AgniWSClient))
}

func (awsc *AgniWSClient) Initialize(pInstance_ID int) bool {
	awsc.id = pInstance_ID
	awsc.isConnected = false
	awsc.wsCon = nil
	return true
}

func (awsc *AgniWSClient) GetID() (pInstance_ID int) {
	return awsc.id
}

// DeInitialize clear the instace values
func (awsc *AgniWSClient) DeInitialize() {

}

// IsConnected checks the fetched web socket connection is connected by writing PING message.
//
// Returns true if the connection is live. unless false with error message
func (awsc *AgniWSClient) IsConnected() (bool, error) {

	if awsc.wsCon == nil {
		return false, fmt.Errorf("%d websocket client is not initialized", awsc.id)
	}
	_err := awsc.wsCon.WriteMessage(websocket.PingMessage, []byte("keepalive"))
	if _err != nil {
		awsc.isConnected = false

	} else {
		awsc.isConnected = true
	}

	return awsc.isConnected, nil
}

// Connect establishes the fetched connection from the pool.
//
// Given wsurl will be used with the request haders for connection.
//
// If success then returns the true,HTTP status code and nil.
//
// If failed then returns false,-1 and the error message
func (awsc *AgniWSClient)Connect(pWS_URL string, pRequest_Headers *map[string][]string, pSub_Pprotocols *[]string,pCompression bool) (bool, int, error) {

	if len(pWS_URL) == 0 {
		return false, -1, errors.New(strconv.Itoa(awsc.id)  + " invalid websocket url. " + pWS_URL)
	}

	awsc.wsURL = pWS_URL

	_httpHeaders := http.Header{}
	
	defer func(){
		_httpHeaders=nil
	}()
	
	if pRequest_Headers!=nil{
		if len(*pRequest_Headers) > 0 {
			for _key, _val := range *pRequest_Headers {
				_httpHeaders[_key] = _val
			}
		}
	}

	_tempwsCon, _httpResp, _err := websocket.DefaultDialer.Dial(awsc.wsURL, _httpHeaders)

	defer func(){
		_httpResp=nil
		_err=nil
	}()
	
	if _err != nil {
		return false, -1, errors.New( strconv.Itoa(awsc.id)  + " failed to connect to the websocket " +  awsc.wsURL + " ." +  _err.Error())
	} else {
		_status_code:=_httpResp.StatusCode
		awsc.wsCon = _tempwsCon
		awsc.isConnected = true
		awsc.wsCon.EnableWriteCompression(pCompression)	/// set the compression on/off
		return true, _status_code, nil
	}
}

// Disconnect disconnects the fetched connection from the pool.
//
// Returns true and nil if disconnection is success.
//
// Unless returns false and error message
func (awsc *AgniWSClient) Disconnect() (bool, error) {

	defer recover()
	
	if !awsc.isConnected {
		return false, errors.New("websocket connection " + strconv.Itoa(awsc.id) +  " already closed")
	}
	
	defer func ()  {
		awsc.wsCon.Close()
	}()
	_err := awsc.wsCon.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	
	defer func(){
		_err=nil
		awsc.wsCon = nil
		awsc.isConnected = false
	}()
	
	if _err != nil {
		return false, errors.New("failed to close connection " + strconv.Itoa(awsc.id) + " | "  + awsc.wsURL + ". " + _err.Error())
	} else {
		return true, nil
	}
}

// Read reads the data from the connection.
//
// If success returns message 1|0 (1=Text Message ,2=Binary Message), message bytes and nil for error
//
// Unless returns 0 as message type, nil and error
func (awsc *AgniWSClient) Read() (messageType int, data *[]byte, err error) {

	if !awsc.isConnected {
		return -1, nil, errors.New("websocket connection " + strconv.Itoa(awsc.id) + " is disconnected")
	}

	if awsc.wsCon==nil{
		return -1, nil, errors.New("websocket connection " + strconv.Itoa(awsc.id) + " is not initialized")
	}	
	
	_msgType, _msgData, _err := awsc.wsCon.ReadMessage()
	//fmt.Printf("reading a message from %s - DONE\n", awsc.wsURL)
	if _err != nil {
		awsc.isConnected = false
		return 0, nil, errors.New("failed to read from the connection " + _err.Error())
	} else {
		awsc.isConnected = true
		return _msgType, &_msgData, nil
	}
}



// Write writes the binary message to the fetched web socket connection.
//
// Returns true and nil if write is success.
//
// Unless returns false and error message
func (awsc *AgniWSClient) Write(pMessage_Type int, pMessage *[]byte) (bool, error) {

	if len(*pMessage) == 0 {
		return false, errors.New("invalid data provided to write")
	}

	if awsc.wsCon==nil{
		return false, errors.New("websocket connection " + strconv.Itoa(awsc.id) + " is not initialized")
	}
	
	if !awsc.isConnected {
		return false, errors.New("websocket connection " + strconv.Itoa(awsc.id) + " is disconnected")
	}
	
	if _err := awsc.wsCon.WriteMessage(pMessage_Type, *pMessage); _err != nil {
		awsc.isConnected = false
		return false, errors.New("websocket connection " + strconv.Itoa(awsc.id) + " error while writing to " +  awsc.wsURL  + ". " +_err.Error())
	} else {
		awsc.isConnected = true
		return true, nil
	}

}

// Info returns the build information of the library
func (awsc *AgniWSClient) Info() build.BuildInfo {
	return build.BuildInfo{Time: buildinfo.Time, User: buildinfo.User, Version: buildinfo.Version, BuildGoVersion: buildinfo.BuildGoVersion}
}

var IAWSClient AgniWSClient
