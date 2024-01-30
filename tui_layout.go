package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

func PrepareTUI(app *tview.Application) {
	prepFormItemsMap()

	addFormItemsToMap(app)

}

func PrepareDeviceLayout() *tview.Flex {
	deviceLeft := tview.NewForm()
	AddFormItemsFromMap(deviceLeft, 12, 6, 2, 3)

	deviceRight := tview.NewForm()
	AddFormItemsFromMap(deviceRight, 12, 7, 4, 5)

	displayForm := tview.NewForm()
	AddFormItemsFromMap(displayForm, 12, 32, 33, 34, 35, 36, 37)

	displayForm.SetBorder(true).SetTitle("Display")
	displayForm.SetHorizontal(true)
	displayForm.AddFormItem(prepSyncRGBButton())
	leftForm := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(deviceLeft, 0, 1, false).
				AddItem(deviceRight, 0, 1, false),
			0, 1, false).
		SetDirection(tview.FlexRow).AddItem(displayForm, 0, 1, false)

	leftForm.SetBorder(true).SetTitle("Device Mode")
	return leftForm
}

func prepSyncRGBButton() *tview.Checkbox {
	rgbBt := tview.NewCheckbox().SetLabel("Meter Color Sync")
	SyncRGB()
	rgbBt.SetChangedFunc(func(checked bool) {
		if checked {
			if RGBticker == nil {
				RGBticker = time.NewTicker(time.Duration(config.Sync.RGB.RefreshRate))
			} else {
				RGBticker.Reset(time.Duration(config.Sync.RGB.RefreshRate))
			}
			metr := GetFormItemsFromMap(12, 33)
			metr[0].SetDisabled(true)
		} else {
			RGBticker.Stop()
			metr := GetFormItemsFromMap(12, 33)
			metr[0].SetDisabled(false)
		}

	})
	return rgbBt
}
func PrepareDeviceLayoutRight() *tview.Flex {
	centerFormLeft := tview.NewForm()
	AddFormItemsFromMap(centerFormLeft, 12, 8, 10)
	centerFormRight := tview.NewForm()
	AddFormItemsFromMap(centerFormRight, 12, 9, 11, 12)

	clockForm := tview.NewForm()
	clockForm.SetBorder(true).SetTitle("Clock Settings")
	clockForm.SetHorizontal(true)

	AddFormItemsFromMap(clockForm, 12, 15, 16)
	centerForm := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(centerFormLeft, 0, 1, false).
				AddItem(centerFormRight, 0, 1, false),
			0, 1, false).
		SetDirection(tview.FlexRow).AddItem(clockForm, 0, 1, false)
	centerForm.SetBorder(true).SetTitle("Phones")
	return centerForm
}

func prepResizeButton() *tview.Button {
	button := tview.NewButton(">")
	button.
		SetActivatedStyle(tcell.Style{}).
		SetDisabledStyle(tcell.Style{}).
		SetStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
		SetSelectedFunc(func() {
			if button.GetLabel() == "<" {
				MainFlex.ResizeItem(tView, 0, 2)
				button.SetLabel(">")
			} else {
				MainFlex.ResizeItem(tView, 0, 0)
				button.SetLabel("<")
			}
		})
	return button
}

func PrepareChannelLayout(channel int) *tview.Flex {
	deviceLeft := tview.NewForm()
	AddFormItemsFromMap(deviceLeft, channel, 1, 12, 2, 3, 9)
	deviceRight := tview.NewForm()
	AddFormItemsFromMap(deviceRight, channel, 4, 7, 6, 5)

	channelLeft2 := tview.NewForm()
	AddFormItemsFromMap(channelLeft2, channel, 18, 20, 10, 17)
	channelRight2 := tview.NewForm()
	AddFormItemsFromMap(channelRight2, channel, 11, 13, 14, 16)

	channelForm2 := tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(channelLeft2, 0, 1, false).
				AddItem(channelRight2, 0, 1, false),
			0, 1, false)

	channelForm2.SetBorder(true).SetTitle(channelNameMap[channel] + " Advanced")
	leftForm := tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(deviceLeft, 0, 1, false).
				AddItem(deviceRight, 0, 1, false),
			0, 1, false).
		SetDirection(tview.FlexRow).
		AddItem(channelForm2,
			0, 1, false)

	leftForm.SetBorder(true).SetTitle(channelNameMap[channel])
	return leftForm
}

func PrepareEQLayoutVolume(app *tview.Application) *tview.Flex {

	barGraph := tvxwidgets.NewBarChart()
	barGraph.SetBorder(true)
	barGraph.SetTitle("Volume")
	barGraph.AddBar(channelNameMap[3], 0, colorMap[4])
	barGraph.AddBar(channelNameMap[6], 0, colorMap[7])
	barGraph.AddBar(channelNameMap[9], 0, colorMap[10])
	barGraph.SetMaxValue(100)

	barGraph.SetAxesColor(tcell.ColorAntiqueWhite)
	barGraph.SetAxesLabelColor(tcell.ColorAntiqueWhite)

	update := func() {

		msgCh := StatusBroker.Subscribe()
		for status := range msgCh {
			for i := 3; i <= 9; i += 3 {
				val, ok := status[i][12]
				if ok {
					app.QueueUpdate(func() {

						barGraph.SetBarValue(channelNameMap[i], int(mapRange(float64(val), -1145, 60, 0, 100)))
					})
				}
			}

		}
	}

	go update()

	centerForm := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(barGraph, 0, 1, false).
				AddItem(PrepareEQLayout(app, 4), 0, 1, false).
				AddItem(PrepareEQLayout(app, 7), 0, 1, false).
				AddItem(PrepareEQLayout(app, 10), 0, 1, false),

			0, 1, false)
	return centerForm
}

var channelNameMap = map[int]string{
	1: "Line In",
	3: "Line Out",
	6: "Phones 1/2",
	9: "Phones 3/4",
}

func mapRange(value, fromMin, fromMax, toMin, toMax float64) float64 {
	if value < fromMin {
		value = fromMin
	} else if value > fromMax {
		value = fromMax
	}

	normalized := (value - fromMin) / (fromMax - fromMin)
	return normalized*(toMax-toMin) + toMin
}
