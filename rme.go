package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"k8s.io/klog/v2"
)

var SysExHeader = []byte{0x00, 0x20, 0x0d}
var lastUpdate time.Time
var lastChange time.Time

type (
	DeviceID        uint8
	CommandID       uint8
	ParamTransfer   uint8
	RequestCommand  uint8
	ParamCommand    CommandID
	TransferCommand CommandID
)

var ADIDevices = map[DeviceID]string{
	ADI2DAC: "ADI-2 DAC",
	ADI2Pro: "ADI-2 PRO",
	ADI24:   "ADI-2 2/4",
}

type RMERequest struct {
	CommandID
	Option uint8
}

// Device IDs
const (
	ADI2DAC DeviceID = iota + 0x71
	ADI2Pro
	ADI24
)

// CommandIDs
const (
	RemoteParam ParamCommand = iota + 1
	DeviceParam
	DeviceRequest RequestCommand = iota + 1
	DeviceRequestPreset
	NameTransferToRemote TransferCommand = iota + 1
	NameTransferToDevice
	Status
)

// Parameter-Transfers
const (
	InputChannel ParamTransfer = iota
	InputEQLeftAndBassTreble
	InputEQRight
	LineOutEQChannel
	LineOutEQLeftAndBassTreble
	LineOutEQRight
	Phones12EQChannel
	Phones12EQLeftAndBassTreble
	Phones12EQRight
	Phones34EQChannel
	Phones34EQLeftAndBassTreble
	Phones34EQRight
	Device
	PresetEQChannel
	PresetEQLeftAndBassTreble
	PresetEQRight
)

type RMESysExMessage struct {
	Header [3]byte
	DeviceID
	CommandID
}

var CurrentDeviceStatusMap = map[int]map[int]int{}

func GetRMEStatus() error {
	send, err := midi.SendTo(out)
	if err != nil {
		return fmt.Errorf("ERROR: %s\n", err)
	}
	send(generateSysexStatusRequest())
	return nil
}
func SetCurrentStatusMap(vals map[int]map[int]int) {
	klog.V(3).Info(vals)
	change := false
	for i, p := range vals {
		for k, v := range p {
			if k == 0 {
				continue
			}
			_, exists := CurrentDeviceStatusMap[i]
			if CurrentDeviceStatusMap[i][k] != v {
				change = true
			}
			if !exists {
				CurrentDeviceStatusMap[i] = make(map[int]int)
			}
			CurrentDeviceStatusMap[i][k] = v

		}
	}
	if change {
		// for k, p := range vals {
		// 	if vals[k] == nil {
		// 		continue
		// 	}
		// 	if CurrentDeviceStatusMap[p] == nil {
		// 		CurrentDeviceStatusMap[k] = make(map[int]int)
		// 	}
		// 	CurrentDeviceStatusMap[k]
		// 	maps.Copy(CurrentDeviceStatusMap[k], vals[k])
		// }
		// maps.Copy(CurrentDeviceStatusMap, vals)
		klog.Infoln(CurrentDeviceStatusMap)
		lastChange = time.Now()
		StatusBroker.Publish(CurrentDeviceStatusMap)
	}

}
func SendCommand(channel, parameterIndex, parameterValue int) {
	var mess = []byte{}
	if channel == 12 {
		klog.V(3).Infoln(channel, parameterIndex)

		mess = append(rmeSysExCommandBase, generateDeviceCommand(parameterIndex, parameterValue)...)

	} else if channel == 3 || channel == 6 || channel == 9 {
		klog.V(3).Infoln(channel, parameterIndex)
		mess = append(rmeSysExCommandBase, generateChannelCommand(channel, parameterIndex, parameterValue)...)
		// log.Fatal(mess)
	}
	send, err := midi.SendTo(out)
	mess = append(mess, SysExClose...)

	if err != nil {
		klog.Infof("ERROR: %s\n", err)
		return
	}

	klog.V(3).Infof("%#v", mess)

	send(mess)
}
