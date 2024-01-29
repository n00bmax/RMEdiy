package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"k8s.io/klog/v2"
)

var colorMap = map[int]tcell.Color{
	4:  tcell.ColorGreen,
	7:  tcell.ColorBlue,
	10: tcell.ColorOrange,
}

func PrepareEQLayout(app *tview.Application, channel int) *tview.Flex {

	eq1 := tvxwidgets.NewSparkline()

	eq1.SetBorder(true)
	eq1.SetTitle(channelNameMap[channel-1] + " EQ")
	eq1.SetDataTitleColor(tcell.ColorDarkOrange)
	eq1.SetLineColor(colorMap[channel])

	update := func() {

		msgCh := StatusBroker.Subscribe()
		for status := range msgCh {

			val, ok := status[channel][4]
			if ok {
				app.QueueUpdate(func() {
					givenNumbers := []float64{10, 20, 30, 40, 50}

					scaleFactors := []float64{float64(mapRange(float64(val), -12, 12, 0, 10)), mapRange(float64(status[channel][7]), -12, 12, 0, 10), mapRange(float64(status[channel][10]), -12, 12, 0, 10), mapRange(float64(status[channel][13]), -12, 12, 0, 10), mapRange(float64(status[channel][17]), -12, 12, 0, 10)} // Adjust these scale values as needed

					_, yValues := generateContinuousCurve(givenNumbers, scaleFactors)
					klog.Info(yValues)

					eq1.SetData(yValues)

				})
			}

		}
	}

	go update()
	centerForm := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(eq1, 0, 1, false),
			0, 1, false)
	centerForm.SetRect(0, 0, 1000, 30)
	return centerForm

}
func generateContinuousCurve(givenNumbers []float64, scales []float64) ([]float64, []float64) {
	maxGivenNumber := givenNumbers[0]
	minGivenNumber := givenNumbers[0]

	for _, num := range givenNumbers {
		if num > maxGivenNumber {
			maxGivenNumber = num
		}
		if num < minGivenNumber {
			minGivenNumber = num
		}
	}

	xValues := make([]float64, int(maxGivenNumber)+20)
	for i := 0; i < len(xValues); i++ {
		xValues[i] = float64(i) + minGivenNumber
	}

	yValues := make([]float64, len(xValues))
	for i, x := range xValues {
		for j, givenNumber := range givenNumbers {
			distance := x - givenNumber

			yValues[i] += scales[j] * math.Exp(-0.5*(distance*distance)/20)
		}
	}

	return xValues, yValues

}
