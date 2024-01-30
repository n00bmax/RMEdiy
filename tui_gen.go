package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

func addFormItemsToMap(app *tview.Application) {

	gauge := tvxwidgets.NewActivityModeGauge()
	gauge.SetTitle("[blue]RME[purple]ote Control")
	gauge.SetPgBgColor(tcell.ColorDarkCyan)
	gauge.SetRect(10, 4, 50, 3)
	gauge.SetBorder(true)

	for channel, ChannelParameters := range Parameters {
		for _, parameter := range ChannelParameters {

			parameter := parameter
			channel := channel
			if len(parameter.Options) < 2 {
				if parameter.Max != 1 {
					cb := tview.NewInputField().SetLabel(parameter.Name)
					update := func() {
						parameter := parameter
						msgCh := StatusBroker.Subscribe()
						for status := range msgCh {
							val, ok := status[channel][parameter.Index]
							if ok && cb.GetText() != strconv.Itoa(val) {
								app.QueueUpdate(func() {
									IsStatusUpdate = true
									cb.SetText(strconv.Itoa(val))
									IsStatusUpdate = false
								})
							}
						}
					}
					cb.SetChangedFunc(func(text string) {
						if IsStatusUpdate {
							// app.QueueUpdate(func() { cb.SetText(text) })

							return
						}
						parameter := parameter
						app.RWMutex.Lock()
						val, _ := strconv.Atoi(text)
						CurrentDeviceStatusMap[channel][parameter.Index] = val
						app.RWMutex.Unlock()

						SendCommand(channel, parameter.Index, val)
					})
					go update()

					TUIFormMap[channel][parameter.ParameterKey.Index] = cb
				} else {
					cb := tview.NewCheckbox().SetLabel(parameter.Name)
					cb.SetChangedFunc(func(checked bool) {
						if IsStatusUpdate {
							return
						}
						parameter := parameter
						itemIndex := 0
						if checked {
							itemIndex = 1
						}
						app.RWMutex.Lock()
						CurrentDeviceStatusMap[channel][parameter.Index] = itemIndex
						app.RWMutex.Unlock()

						SendCommand(channel, parameter.Index, itemIndex)
					})
					update := func() {
						parameter := parameter
						msgCh := StatusBroker.Subscribe()

						for status := range msgCh {
							val, ok := status[channel][parameter.Index]
							if ok {
								app.QueueUpdate(func() {
									checked := false
									if val == 1 {
										checked = true
									}
									IsStatusUpdate = true
									cb.SetChecked(checked)
									IsStatusUpdate = false
								})
							}
						}
					}

					go update()
					TUIFormMap[channel][parameter.ParameterKey.Index] = cb
				}
			} else {
				dd :=
					tview.NewDropDown().SetLabel(parameter.Name)

				dd.SetOptions(parameter.Options, func(text string, index int) {
					if IsStatusUpdate || dd.IsOpen() {
						return
					}
					parameter := parameter
					app.RWMutex.Lock()
					CurrentDeviceStatusMap[channel][parameter.Index] = index
					app.RWMutex.Unlock()

					SendCommand(channel, parameter.Index, index)

				})
				update := func() {
					parameter := parameter
					msgCh := StatusBroker.Subscribe()
					for status := range msgCh {
						val, ok := status[channel][parameter.Index]
						if ok {
							app.QueueUpdate(func() {
								IsStatusUpdate = true
								dd.SetCurrentOption(val)
								IsStatusUpdate = false
							})
						}
					}
				}

				go update()
				TUIFormMap[channel][parameter.ParameterKey.Index] = dd
			}

		}
		actGauge = gauge
	}
}

func GetFormItemsFromMap(channel int, key int, keys ...int) (items []tview.FormItem) {
	for _, key := range append(keys, key) {
		item, ok := TUIFormMap[channel][key]
		if !ok {
			continue
		}
		items = append(items, item)
	}
	return
}

func AddFormItems(form *tview.Form, items []tview.FormItem) {
	for _, item := range items {
		form.AddFormItem(item)
	}
}

func AddFormItemsFromMap(form *tview.Form, channel int, key int, keys ...int) {
	AddFormItems(form, GetFormItemsFromMap(channel, key, keys...))
}

var TUIFormMap = map[int]map[int]tview.FormItem{}

func prepFormItemsMap() {
	for i := 1; i <= 12; i++ {
		TUIFormMap[i] = make(map[int]tview.FormItem)
	}
}
