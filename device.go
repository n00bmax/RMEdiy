package main

import "maps"

type Parameter struct {
	ParameterKey
	Min       int
	Max       int
	Default   int
	Increment int
	Options   []string
}

var Parameters map[int]map[int]Parameter

func init() {
	DeviceParameters := map[int]map[int]Parameter{
		12: {
			1:  {ParameterKey{1, "Mute Line vs."}, 0, 3, 1, 1, []string{"OFF", "vs. Phones", "Toggle Ph/Line", "Toggle plugged"}},
			2:  {ParameterKey{2, "Auto Standby"}, 0, 4, 0, 1, []string{"OFF", "30min", "1h", "2h", "4h"}},
			3:  {ParameterKey{3, "DSD Detection"}, 0, 1, 1, 1, nil},
			4:  {ParameterKey{4, "DSD Filter"}, 0, 1, 0, 1, []string{"50 kHz", "150 kHz"}},
			5:  {ParameterKey{5, "DSD Direct (Line)"}, 0, 1, 0, 1, nil},
			6:  {ParameterKey{6, "Basic Mode"}, 0, 5, 0, 1, []string{"Auto", "AD/DA", "USB", "Preamp", "Dig Thru", "DAC"}},
			7:  {ParameterKey{7, "Digital Out Source"}, 0, 1, 0, 1, []string{"Default", "Main Out"}},
			8:  {ParameterKey{8, "Dual Phones"}, 0, 1, 0, 0, nil},
			9:  {ParameterKey{9, "Bal. TRS Phones Mode"}, 0, 0, 0, 0, []string{"OFF", "ON", "Auto"}},
			10: {ParameterKey{10, "Toggle Phones/Line"}, 0, 0, 0, 0, []string{"OFF", "Ph 1/2", "Ph 3/4", "1/2+3/4", "all plugged", "Line/Digit."}},
			11: {ParameterKey{11, "Mute Line vs. PH12"}, 0, 1, 1, 1, nil},
			12: {ParameterKey{12, "Mute LIne vs. PH 34"}, 0, 1, 1, 1, nil},
			13: {ParameterKey{13, "Unused Parameter"}, 0, 0, 0, 0, nil},
			14: {ParameterKey{14, "Unused Parameter"}, 0, 0, 0, 0, nil},
			15: {ParameterKey{15, "Clock Source"}, 0, 3, 0, 1, []string{"Auto", "INT", "AES", "SPDIF"}},
			16: {ParameterKey{16, "Sample Rate"}, 0, 9, 1, 1, []string{"44.1 kHz", "48 kHz", "88 kHz", "96 kHz", "176.4 kHz", "192 kHz", "352.8 kHz", "384 kHz", "705.6 kHz", "768 kHz"}},
			17: {ParameterKey{17, "Unused Parameter"}, 0, 0, 0, 0, nil},
			18: {ParameterKey{18, "Unused Parameter"}, 0, 0, 0, 0, nil},
			19: {ParameterKey{19, "Unused Parameter"}, 0, 0, 0, 0, nil},
			20: {ParameterKey{20, "Unused Parameter"}, 0, 0, 0, 0, nil},
			21: {ParameterKey{21, "IR 5"}, 0, 2, 2, 1, []string{"ADI-2/4", "ADI-2 Pro"}},
			22: {ParameterKey{22, "IR 6"}, 0, 2, 2, 1, []string{"ADI-2/4", "ADI-2 Pro"}},
			23: {ParameterKey{23, "IR 7"}, 0, 2, 2, 1, []string{"ADI-2/4", "ADI-2 Pro"}},
			24: {ParameterKey{24, "Unused Parameter"}, 0, 0, 0, 0, nil},
			25: {ParameterKey{25, "Remap Keys"}, 0, 2, 2, 1, []string{"OFF", "ON", "IR-Remote"}},
			26: {ParameterKey{26, "VOL (1)"}, 0, 0, 0, 0, nil},
			27: {ParameterKey{27, "I/O (2)"}, 0, 0, 0, 0, nil},
			28: {ParameterKey{28, "EQ (3)"}, 0, 0, 0, 0, nil},
			29: {ParameterKey{29, "SETUP (4)"}, 0, 0, 0, 0, nil},
			30: {ParameterKey{30, "Unused Parameter"}, 0, 0, 0, 0, nil},
			31: {ParameterKey{31, "Unused Parameter"}, 0, 0, 0, 0, nil},
			32: {ParameterKey{32, "Display Mode"}, 0, 1, 0, 1, []string{"Default", "Dark"}},
			33: {ParameterKey{33, "Meter Color"}, 0, 5, 0, 1, []string{"Green", "Cyan", "Amber", "Monochrome", "Red", "Orange"}},
			34: {ParameterKey{34, "Hor. Meter"}, 0, 3, 0, 1, []string{"Post-FX", "Pre-FX", "Dual", "Post-FX dBu"}},
			35: {ParameterKey{35, "AutoDark Mode"}, 0, 1, 0, 1, nil},
			36: {ParameterKey{36, "Show Vol. Screen"}, 0, 1, 1, 1, nil},
			37: {ParameterKey{37, "Lock UI"}, 0, 3, 0, 1, []string{"OFF", "Remote", "Keys", "Keys+Rem."}},
		}}

	ChannelParameters := map[int]map[int]Parameter{
		3: {
			1:  {ParameterKey{1, "Source"}, 0, 6, 0, 1, []string{"Auto", "AES", "SPDIF", "Analog", "USB 1/2", "USB 3/4"}},
			2:  {ParameterKey{2, "Ref Level"}, 0, 3, 2, 1, []string{"+4 dBu", "+13 dBu", "+19 dBu", "+24 dBu"}},
			3:  {ParameterKey{3, "Auto Ref Level"}, 0, 1, 1, 1, nil},
			4:  {ParameterKey{4, "Mono"}, 0, 2, 0, 1, []string{"OFF", "ON", "to Left"}},
			5:  {ParameterKey{5, "Width"}, -100, 100, 100, 1, nil}, // Considering values are in percentage (-1.00 to 1.00)
			6:  {ParameterKey{6, "M/S-Proc"}, 0, 1, 0, 1, nil},
			7:  {ParameterKey{7, "Polarity"}, 0, 3, 0, 1, []string{"OFF", "Both", "Left", "Right"}},
			8:  {ParameterKey{8, "Crossfeed"}, 0, 5, 0, 1, nil},
			9:  {ParameterKey{9, "DA Filter"}, 0, 6, 2, 1, []string{"SD Sharp", "SD Slow", "Sharp", "Slow", "NOS", "SD LD"}},
			10: {ParameterKey{10, "De-Emphasis"}, 0, 2, 0, 1, []string{"Auto", "ON", "OFF"}},
			11: {ParameterKey{11, "Dual EQ"}, 0, 1, 0, 1, nil},
			12: {ParameterKey{12, "Volume"}, -1145, 60, -100, 5, nil}, // Considering values are in decibels (-114.5 to 6.0)
			13: {ParameterKey{13, "Lock Volume"}, 0, 1, 0, 1, nil},
			14: {ParameterKey{14, "Balance"}, -100, 100, 0, 1, nil}, // Considering values are in percentage (-1.00 to 1.00)
			15: {ParameterKey{15, "Mute"}, 0, 1, 0, 1, nil},
			16: {ParameterKey{16, "Dim"}, 0, 1, 0, 1, nil},
			17: {ParameterKey{17, "Loopback to USB"}, 0, 9, 0, 1, []string{"OFF", "pre FX to 1/2", "post FX to 1/2", "post 1/2 -6dB", "pre FX to 3/4", "post FX to 3/4", "post 3/4 -6dB", "pre FX to 5/6", "post FX to 5/6", "post 5/6 -6dB"}},
			18: {ParameterKey{18, "Dig. DC Protection"}, 0, 2, 2, 1, []string{"OFF", "ON", "Filter"}},
			19: {ParameterKey{19, "Rear TRS Source"}, 0, 1, 0, 1, []string{"Line 1/2", "Ph. 3/4"}},
			20: {ParameterKey{20, "Loudness Enable"}, 0, 1, 0, 1, nil},
			21: {ParameterKey{21, "Bass Gain"}, 10, 100, 70, 5, nil},           // Considering values are in decibels (1.0 to 10.0)
			22: {ParameterKey{22, "Treble Gain"}, 10, 100, 70, 5, nil},         // Considering values are in decibels (1.0 to 10.0)
			23: {ParameterKey{23, "Low Vol Ref"}, -9000, -2000, -3000, 5, nil}, // Considering values are in decibels (-90.0 to -20.0)
		},
	}

	EQParameters := map[int]map[int]Parameter{
		1: {
			1: {ParameterKey{1, "Source"}, 0, 6, 0, 1, []string{"Auto", "AES", "SPDIF", "Analog", "USB 1/2", "USB 3/4"}},

			2: {ParameterKey{1, "EQ Enable"}, 0, 6, 0, 1, nil},
			// 2:  {ParameterKey{2, "Ref Level"}, 0, 3, 2, 1, []string{"+4 dBu", "+13 dBu", "+19 dBu", "+24 dBu"}},
			3:  {ParameterKey{3, "Band 1 Type"}, 0, 3, 1, 1, []string{"Peak", "Shelf", "Hi Pass", "Hi Cut"}},
			4:  {ParameterKey{4, "Band 1 Gain"}, -12, 12, 0, 5, nil},
			5:  {ParameterKey{5, "Band 1 Freq"}, 20, 20000, 100, 1, nil},
			6:  {ParameterKey{6, "Band 1 Q"}, 5, 99, 10, 1, nil},
			7:  {ParameterKey{4, "Band 2 Gain"}, -12, 12, 0, 5, nil},
			8:  {ParameterKey{5, "Band 2 Freq"}, 20, 20000, 100, 1, nil},
			9:  {ParameterKey{6, "Band 2 Q"}, 5, 99, 10, 1, nil},
			10: {ParameterKey{4, "Band 3 Gain"}, -12, 12, 0, 5, nil},
			11: {ParameterKey{5, "Band 3 Freq"}, 20, 20000, 100, 1, nil},
			12: {ParameterKey{6, "Band 3 Q"}, 5, 99, 10, 1, nil},
			13: {ParameterKey{4, "Band 4 Gain"}, -12, 12, 0, 5, nil},
			14: {ParameterKey{5, "Band 4 Freq"}, 20, 20000, 100, 1, nil},
			15: {ParameterKey{6, "Band 4 Q"}, 5, 99, 10, 1, nil},
			16: {ParameterKey{3, "Band 5 Type"}, 0, 3, 1, 1, []string{"Peak", "Shelf", "Hi Pass", "Hi Cut"}},
			17: {ParameterKey{4, "Band 5 Gain"}, -12, 12, 0, 5, nil},
			18: {ParameterKey{5, "Band 5 Freq"}, 20, 20000, 100, 1, nil},
			19: {ParameterKey{6, "Band 5 Q"}, 5, 99, 10, 1, nil},
		},
	}
	EQParameters[4] = maps.Clone(EQParameters[1])
	EQParameters[7] = maps.Clone(EQParameters[1])
	EQParameters[10] = maps.Clone(EQParameters[1])
	EQParameters[2] = maps.Clone(EQParameters[1])
	EQParameters[5] = maps.Clone(EQParameters[1])
	EQParameters[8] = maps.Clone(EQParameters[1])

	EQParameters[11] = maps.Clone(EQParameters[1])

	ChannelParameters[6] = maps.Clone(ChannelParameters[3])
	ChannelParameters[9] = maps.Clone(ChannelParameters[3])
	Parameters = ChannelParameters
	maps.Copy(Parameters, DeviceParameters)
	maps.Copy(Parameters, EQParameters)

}
