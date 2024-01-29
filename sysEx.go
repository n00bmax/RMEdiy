package main

import (
	"bytes"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"k8s.io/klog/v2"
)

var sysexBase = []byte{0xF0, 0x00, 0x20, 0x0D}

var sysexDeviceBase = []byte{}

var rmeSysExCommandBase = []byte{}

var SysExClose = []byte{0xF7}

func setSysexBase() {
	klog.Info(byte(config.ADI.ID))
	sysexDeviceBase = append(sysexBase, byte(config.ADI.ID))
	rmeSysExCommandBase = append(append(sysexBase, byte(config.ADI.ID)), byte(0x02))
}

func SysExListener() {
	in, err := midi.InPort(config.InPort)
	if err != nil {
		klog.Exit(err)
	}
	_, err = midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		switch {

		case msg.GetSysEx(&bt):
			switch msgByte := bytes.Clone(bt); int(msgByte[4]) {
			case 1:
				lastUpdate = time.Now()

				klog.V(3).Infof("\ngot sysex setting device: % X\n", msgByte)

				parsedStatus := parseStatus(msgByte[5:])
				klog.V(3).Infof("\ngot sysex setting device: %#v\n", parsedStatus)

				if len(parsedStatus) != 0 {
					go SetCurrentStatusMap(parsedStatus)

				}
			case 7:
				// klog.Infof("\ngot sysex status: % X\n", bt)
			case 5:
				// klog.Infof("\ngot sysex eq: % X\n", bt)
			}
		default:
			klog.V(2).Infof("unidentified: % X\n", bt)
			// ignore
		}
	}, midi.UseSysEx())
	if err != nil {
		klog.Infof("ERROR: %s\n", err)
		return
	}
}
func generateSysexStatusRequest() []byte {
	mess := append(sysexDeviceBase, byte(0x03))
	mess = append(mess, byte(0x09))
	mess = append(mess, SysExClose...)
	return mess

}
