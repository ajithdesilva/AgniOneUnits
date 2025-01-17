/*
#########################################################################################

	Author        :   D. Ajith Nilanta de Silva ajithdesilva@gmail.com
	Class/module  :   demohttp AgniOne Unit
	Objective     :   Demonstrate AgniOne Unit implementation with the help of AgniOne Application
						Framework utilizing the build-in AgniOne HTTP plugin

#########################################################################################

	Author                 	Date        	Action      	Description

------------------------------------------------------------------------------------------------------

	Ajith de Silva		05/01/2025	Created 	Created the initial version

	Ajith de Silva		08/01/2025	Updated 	Defined functions with parameters & return values

	Ajith de Silva		10/01/2025	Updated 	Implements methods.

#########################################################################################
*/
package main

import (
	httptypes "agnione/v1/src/afplugins/http/types"
	fmtypes "agnione/v1/src/appfm/types"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DoGet a sample go routine that performs HTTP GET and print results via framework
func (appc *AUDemoHttp) DoGet() {
	appc.base.AppFramework.Write2Log("starting the DoGet routine",fmtypes.LOG_INFO)

	defer func() {
		appc.base.Remove_Routine()
	}()
	
	_httpClient, _err := appc.base.AppFramework.Get_RESTClient(&appc.config.Get.Plugin_Type)
	if _err != nil {
		appc.base.AppFramework.Write2Log(appc.base.App_UID + "  - Failed to get the HTTP client plugin " +  _err.Error(),fmtypes.LOG_ERROR)
		
		//// send the monitoring message to web-socket monitor
		go appc.base.Send_Monitor_Message(appc.base.Generate_Monitoring_Message(
			appc.base.App_UID,
			"DoGet",
			"ERROR",
			map[string]string{
				"action":"HTTP Plugin Load",
				"entry":"",
				"info":appc.base.App_UID + "  -Failed to get the HTTP client plugin " +  _err.Error(),
				"time":strconv.FormatInt(time.Now().Unix(), 10) }))
				
		return
	}

	_httpReq := new(httptypes.AHTTPRequest)
	_httpReq.URL = appc.config.Get.Url
	 _headers:= new(strings.Builder)
	_ticker := time.NewTicker(time.Second * time.Duration(appc.config.Get.Fequancy_Secs)) /// ticker to read in intervals
	 
	defer func ()  {
		_ticker.Stop()
		_ticker=nil
		_headers=nil
		_httpReq=nil
	}()
	
	_count := 0
	for {
		select {
		case <-appc.base.Stopper:
			appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - demo http stopped",fmtypes.LOG_INFO)
			return
		case <-appc.base.AppFramework.Is_Interrupted():
			appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - appliction forced to stop",fmtypes.LOG_INFO)
			return
		case <- _ticker.C:
			appc.base.Increse_Active_Count()
			_httResp, _err := _httpClient.Get(_httpReq)
			
			if _err != nil {
				
				go appc.base.Add_Request_Failed_Count()
				
				appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - Failed to load fetch the request - " +  _err.Error(),fmtypes.LOG_ERROR)
				
				//// send the monitoring message to web-socket monitor
				go appc.base.Send_Monitor_Message(appc.base.Generate_Monitoring_Message(
					appc.base.App_UID,
					"DoGet",
					"ERROR",
					map[string]string{
						"action":"Get",
						"entry":_httpReq.URL,
						"info":appc.base.App_UID  +  " - DoGet failed to load fetch the request -" +_err.Error(),
						"time":strconv.FormatInt(time.Now().Unix(), 10) }))
				
				appc.base.Decrease_Active_Count()
				
			} else {
				
				
				go appc.base.Add_Request_Handled_Count()
				
				appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - [" + strconv.Itoa(_count) + "] - GET Result :: Status code: " + strconv.Itoa(_httResp.StatusCode),fmtypes.LOG_INFO)
				
				var _key string
				var _val string
				for _key, _val = range _httResp.Headers {
					_headers.WriteString(_key + ":" + _val + "\n")
				}
				_key=""
				_val=""
				
				appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - Headers :\r\n" + _headers.String(),fmtypes.LOG_INFO)
				_headers.Reset()
				
				//// send the monitoring message to web-socket monitor
				go appc.base.Send_Monitor_Message(appc.base.Generate_Monitoring_Message(
					appc.base.App_UID,
					"GET",
					"SUCCESS",
					map[string]string{
						"action":"POST",
						"entry":_httpReq.URL,
						"info":appc.base.App_UID  +  " - [" + strconv.Itoa(_count) + "] - GET Result :: Status code: " + strconv.Itoa(_httResp.StatusCode),
						"time":strconv.FormatInt(time.Now().Unix(), 10) }))

			}
			
			appc.base.Decrease_Active_Count()
			_count++
		}

	}
}

// DoPost a sample go routine that performs the HTTP POST and prints the results via framework
func (appc *AUDemoHttp) DoPost() {

	appc.base.AppFramework.Write2Console("starting the DoPost routine")

	defer func() {
		appc.base.Remove_Routine()
	}()
	
	
	_httpClient, _err := appc.base.AppFramework.Get_RESTClient(&appc.config.Get.Plugin_Type)
	if _err != nil {
		appc.base.AppFramework.Write2Log(appc.base.App_UID +  " - DoPost Failed to get the HTTP client plugin " + _err.Error(),fmtypes.LOG_ERROR)
		
		//// send the monitoring message to web-socket monitor
		go appc.base.Send_Monitor_Message(appc.base.Generate_Monitoring_Message(
			appc.base.App_UID,
			"Post",
			"ERROR",
			map[string]string{
				"action":"HTTP Plugin Load",
				"entry":"",
				"info":appc.base.App_UID + "  - Failed to get the HTTP client plugin " +  _err.Error(),
				"time":strconv.FormatInt(time.Now().Unix(), 10) }))
				
				
		return
	}

	_httpReq := new(httptypes.AHTTPRequest)
	_httResp := new(httptypes.AHTTPResponse)
	_headers := new(strings.Builder)
	_ticker := time.NewTicker(time.Second * time.Duration(appc.config.Post.Fequancy_Secs)) /// ticker to read in intervals
	
	defer func ()  {
		_ticker.Stop()
		_ticker=nil
		_httpReq=nil
		_httResp=nil
		_httpClient=nil
		_headers=nil
	}()
	
	_httpReq.URL = appc.config.Post.Url
	
	_count := 0

	for {
		select {
		case <-appc.base.Stopper:
			appc.base.AppFramework.Write2Log(appc.base.App_UID +  " - demo http stopped. stopper channel closed ",fmtypes.LOG_INFO)
			return
		case <-appc.base.AppFramework.Is_Interrupted():
			appc.base.AppFramework.Write2Log(appc.base.App_UID +  " - demo http stopped. application forced to stop",fmtypes.LOG_INFO)
			return
		case <-_ticker.C :
			
			appc.base.Increse_Active_Count()
			
			if _reqBody, _err := json.Marshal(
				AuthReq{
					ApiKey:         fmt.Sprintf("b3-45a4-826d-23453%v", time.Now().UnixMilli()),
					AppName:        "3f93916b-e9b3-45a4-53@TestAPP",
					ConversationId: fmt.Sprintf("%s -DEMO-POST Body -%d",appc.base.App_UID , _count),
				}); _err != nil {
				
					appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - DoPost error occurred while constructing body " + _err.Error(),fmtypes.LOG_ERROR)
					appc.base.Decrease_Active_Count()
					go appc.base.Add_Request_Failed_Count()
					
				return
			} else {
				_httpReq.Body = _reqBody
			}

			_httResp, _err = _httpClient.Post(_httpReq)
			if _err != nil {
				appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - DoPost Failed to load fetch the request -" +_err.Error(),fmtypes.LOG_ERROR)
				appc.base.Add_Request_Failed_Count()
				
				//// send the monitoring message to web-socker monitor
				go appc.base.Send_Monitor_Message(appc.base.Generate_Monitoring_Message(
					appc.base.App_UID,
					"Post",
					"ERROR",
					map[string]string{
						"action":"POST",
						"entry":_httpReq.URL,
						"info":appc.base.App_UID  +  " - DoPost Failed to load fetch the request -" +_err.Error(),
						"time":strconv.FormatInt(time.Now().Unix(), 10) }))

			}else{
				
				go appc.base.Add_Request_Handled_Count()
				appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - DoPost Result :: Status code: " +strconv.Itoa(_httResp.StatusCode),fmtypes.LOG_INFO)
				
				var _key string
				var _val string
				
				for _key, _val = range _httResp.Headers {
					_headers.WriteString(_key + ":" + _val + "\n")
				}
				_key=""
				_val=""
				
				appc.base.AppFramework.Write2Log(appc.base.App_UID  +  " - Headers :\r\n" + _headers.String(),fmtypes.LOG_INFO)
				_headers.Reset()
				
				//// send the monitoring message to web-socker monitor
				go appc.base.Send_Monitor_Message(appc.base.Generate_Monitoring_Message(
					appc.base.App_UID,
					"Post",
					"SUCCESS",
					map[string]string{
						"action":"POST",
						"entry":_httpReq.URL,
						"info":appc.base.App_UID  +  " - DoPost Result :: Status code: " +strconv.Itoa(_httResp.StatusCode),
						"time":strconv.FormatInt(time.Now().Unix(), 10) }))

			}

			appc.base.Decrease_Active_Count()
		}
	}

}

