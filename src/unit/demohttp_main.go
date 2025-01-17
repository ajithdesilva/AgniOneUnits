package main

import (
	izappunit "agnione/v1/src/aau/iappunit"
	fmtypes "agnione/v1/src/appfm/types"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	libbuild "agniunit.demo.http/src/build"

	build "agnione/v1/src/lib"

	unitbase "agnione/v1/src/aau/base"
	izappfm "agnione/v1/src/appfm/iappfw"
)

type AuthReq struct {
	ApiKey         string `json:"apiKey"`
	AppName        string `json:"appName"`
	ConversationId string `json:"conversationId"`
}

/* defin d the config type */
type section struct{
	Plugin_Type string `json:"plugin_type"`
	Url string `json:"url"`
	Fequancy_Secs int `json:"frequancy_secs"`
}

type Config struct{
	Get section `json:"get"`
	Post section `json:"post"`
}


type AUDemoHttp struct {
	base unitbase.AUBase
	config Config
}


func (appc *AUDemoHttp) New() interface{} {
	return fmtypes.ConvertStoI[izappunit.IAppUnit](new(AUDemoHttp))
}


func (appc *AUDemoHttp) Initialize(pFM_Instance izappfm.IAgniApp, pInstance_ID int,
	pAppunit_Name string, pAppunit_Path string,pConfig_File string) (bool, error) {
		
	/// 1. Try to nitialize the base
	_, _err := appc.base.Initialize(pFM_Instance, pInstance_ID, pAppunit_Name, pAppunit_Path,pConfig_File)
	if _err != nil {
		return false, _err
	}
	
	/// loads the configuration
	var _fileData *[]byte
	
	defer func(){
		_fileData=nil
	}()
	
	_fileData,_err=appc.base.AppFramework.Get_File_Content(&pConfig_File)
	if _err != nil {
		
		_info:=appc.base.App_UID + " - Application Unit " + pAppunit_Name + " failed to read the config file " + pConfig_File + _err.Error()
		appc.base.Write2Log(_info, fmtypes.LOG_ERROR)
		
		//// send monitor message to websocker monitor clients
		go appc.base.Send_Monitor_Message(appc.base.Generate_Monitoring_Message(
			appc.base.App_UID,
			"NONE",
			"ERROR",
			map[string]string{
				"action":"Read Configuration",
				"entry":"",
				"info":_info,
				"time":strconv.FormatInt(time.Now().Unix(), 10) }))
				
		return false,errors.New(_info)
	}
	
	_err = json.Unmarshal(*_fileData, &appc.config)
	if _err != nil {
		_info:=appc.base.App_UID + " - ZAU failed to parse the config file " + pConfig_File + _err.Error()
		appc.base.Write2Log(_info, fmtypes.LOG_ERROR)
		
		//// send monitor message to websocker monitor clients
		go appc.base.Send_Monitor_Message(appc.base.Generate_Monitoring_Message(
			appc.base.App_UID,
			"NONE",
			"ERROR",
			map[string]string{
				"action":"Parse configuration",
				"entry":"",
				"info":_info,
				"time":strconv.FormatInt(time.Now().Unix(), 10) }))
				
		return false,errors.New(_info)
	}
	
	appc.base.Unit_Info.Info.Version = libbuild.Version
		
	return true,nil

}

// Deinitialize clear all the related objects
func (appc *AUDemoHttp) Deinitialize() {
	/// clear objects here

	if !appc.base.Is_Initialized {
		fmt.Println("Framework  NOT Initialize")
		return
	}
	
	appc.config=Config{}
	
	appc.base.Write2Log( appc.base.App_UID + " - Deinitializing....", fmtypes.LOG_INFO)
	/// 1. clear Unit specific objects here
	
	
	// 2. base.Deinitialize()
	appc.base.Deinitialize()
	appc.base.Write2Log(appc.base.App_UID + " - Deinitialize..... DONE", fmtypes.LOG_INFO)
}




func (appc *AUDemoHttp) Start() (bool, error)  {

	appc.base.AppFramework.Write2Log(appc.base.App_UID + " - Starting .....",fmtypes.LOG_INFO)
	
	appc.base.Stopper = make(chan bool)	//// main channle to control units's routines

	//// start the rotines that does perform tasks
	appc.base.Add_Routine()
	go appc.DoGet()
	
	appc.base.Add_Routine()
	go appc.DoPost()
	
	appc.base.AppFramework.Write2Log(appc.base.App_UID + " - Starting the ........... DONE",fmtypes.LOG_INFO)
	
	appc.base.Is_Started=true	/// flag that Unit started
	return true,nil
}

func (appc *AUDemoHttp) Stop() (bool, error)  {

	defer func ()  {
		recover()
	}()
	
	appc.base.AppFramework.Write2Log(appc.base.App_UID + " - Stopping HTTPZAU demo......", fmtypes.LOG_INFO)
	if appc.base.Stopper != nil {
		close(appc.base.Stopper)
	}

	appc.base.AppFramework.Write2Log(appc.base.App_UID + " - Stopping HTTPZAU demo........... DONE", fmtypes.LOG_INFO)
	return true,nil
}


func (appc *AUDemoHttp) IsInitialized() bool{
	return appc.base.Is_Started
}


func (appc *AUDemoHttp) IsStarted() bool{
	return appc.base.Is_Started
}

func (appc *AUDemoHttp)GetID() (instance_id int) {
	return appc.base.Get_ID()
}

func (appc *AUDemoHttp) Status() *fmtypes.AppUnitInfo {
	return appc.base.Status()
}

// Info returns the information of the library
func (appc *AUDemoHttp) Info() build.BuildInfo {
	
	return build.BuildInfo{Time: libbuild.Time, User: libbuild.User, Version: libbuild.Version, BuildGoVersion: libbuild.BuildGoVersion}
}

var IAppUnit AUDemoHttp
