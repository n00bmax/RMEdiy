package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"k8s.io/klog/v2"
)

func initKlog() {
	flags := &flag.FlagSet{}
	klog.InitFlags(flags)
	flags.Set("logtostderr", "false")
	flags.Set("alsologtostderr", "false")
	flags.Set("stderrthreshold", "4")
}

func initADIDevice() {

}
func main() {
	klog.Info("out ports: \n" + midi.GetOutPorts().String())
	klog.Info("in ports: \n" + midi.GetInPorts().String())
	initRmediy()
	if !config.Rmediator.Disabled {
		initRmediator()
	}
	StartTUI()
}

func initRmediy() {
	//set base sysex address
	setSysexBase()
	// setCommandBase()

	klog.Info(sysexDeviceBase)
	klog.Info(rmeSysExCommandBase)
	go StatusBroker.Start()

	go SysExListener()

	if GetRMEStatus() != nil {
		fmt.Print("error starting rmediy")

		time.Sleep(10 * time.Second)
		klog.Exit()
	} else {
		initKlog()
	}
}

func init() {
	if err := getConfig(); err != nil {
		fmt.Print(err)
		time.Sleep(20 * time.Second)
		os.Exit(0)
	}
	out, _ = midi.OutPort(config.ADI.OutPort)
}
