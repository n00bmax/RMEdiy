package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"k8s.io/klog/v2"
)

func RmediyColor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	col, err := strconv.Atoi(ps.ByName("color"))
	klog.V(3).Infof("got color sync  %v request\n", col)
	if err != nil {
		klog.Error(err)
		return
	}
	meterColorChangeRequest(col)
}

func RmediyStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jData, err := json.Marshal(config.ADI)
	if err != nil {
		// handle error
	}
	w.Write(jData)
}
func RmediySet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := strconv.Atoi(ps.ByName("channel"))
	if err != nil {
		klog.Error(err)
		return
	}
	p, err := strconv.Atoi(ps.ByName("param"))
	if err != nil {
		klog.Error(err)
		return
	}
	v, err := strconv.Atoi(ps.ByName("value"))
	if err != nil {
		klog.Error(err)
		return
	}
	SendCommand(c, p, v)

}

func RmediySettings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setting := ps.ByName("setting")
	jData := []byte{}
	var err error
	if setting == "all" {
		jData, err = json.Marshal(CurrentDeviceStatusMap)
		if err != nil {
			// handle error
		}
	} else if setting, err := strconv.Atoi(setting); err == nil {
		jData, err = json.Marshal(CurrentDeviceStatusMap[setting])
		if err != nil {
			// handle error
		}
	}

	w.Write(jData)
}
func initRmediator() {
	router := httprouter.New()
	router.GET("/status", RmediyStatus)
	router.GET("/metercolor/:color", RmediyColor)
	router.GET("/settings/:setting", RmediySettings)
	router.POST("/set/:channel/:param/:value", RmediySet)

	klog.Info("RMEdiator serving requests at ", config.Rmediator.Port)
	go http.ListenAndServe(":"+config.Rmediator.Port, router)
}

type RmediatorSettings struct {
	Parameter
	Current     int
	CurrentText string
}
