package main

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"k8s.io/klog/v2"
)

var (
	actGauge       *tvxwidgets.ActivityModeGauge
	TermWriter     io.Writer
	tView          *tview.Flex
	IsStatusUpdate = false
	MainFlex       *tview.Flex
	logPane        *tview.TextView
	verbosity      int
)

func StartTUI() {

	app := tview.NewApplication()

	PrepareTUI(app)

	logPane = tview.NewTextView().SetDynamicColors(true)
	klog.SetOutput(tview.ANSIWriter(logPane))

	tView = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(logPane, 0, 1, false)

	MainFlex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(PrepareDeviceLayout(), 0, 3, false).
		AddItem(PrepareDeviceLayoutRight(), 0, 3, false).
		AddItem(GenerateChannelLayouts(), 0, 4, false).
		AddItem(PrepareEQLayoutVolume(app), 42, 1, false).
		AddItem(prepResizeButton(), 3, 1, false).
		AddItem(tView, 0, 2, false)

	footer := tview.NewTextView().SetText(config.Name + "\t\tLast Update: ")
	footer.SetDynamicColors(true)

	footer.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		footer.SetText(fmt.Sprintf("[green]%v \t\t [darkcyan]Last Update:[blue] %v\t\t [darkcyan]Last Change: [purple]%v", config.Name, lastUpdate.Format("Monday, 02-Jan-06 10:04:05PM"), lastChange.Format("10:04:05PM")))
		return x, y, width, height
	})
	TUIMain := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(actGauge, 3, 1, false).
		AddItem(MainFlex, 0, 1, false).
		AddItem(footer, 1, 1, false)

	update := func() {
		for range time.Tick(time.Duration(config.Sync.Interval) * time.Second) {
			GetRMEStatus()
			app.QueueUpdateDraw(func() {
				actGauge.Pulse()
			})
		}
	}
	go update()

	if err := app.SetRoot(TUIMain, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func footerUpdate() {

}

func init() {
	RGBticker = time.NewTicker(100 * time.Hour)
	RGBticker.Stop()
}

func GenerateChannelLayouts() *tview.Flex {
	pages := tview.NewPages()
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			if len(added) != 0 {
				pages.SwitchToPage(added[0])
			}
		})

	for i, channel := range []int{3, 6, 9} {
		pages.AddPage(strconv.Itoa(i), PrepareChannelLayout(channel), true, channel == 3)
		fmt.Fprintf(info, `%d ["%d"][%s]%s[white][""]  `, i+1, i, channelColorMap[channel], channelNameMap[channel])
	}
	info.Highlight("1")
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(info, 1, 1, false).
		AddItem(pages, 0, 1, true)

}

var channelColorMap = map[int]string{
	1: "red",
	3: "green",
	6: "blue",
	9: "purple",
}
