package main

import (
	"k8s.io/klog/v2"
)

func parseStatus(parameterBytes []byte) map[int]map[int]int {
	parsedValues := make(map[int]map[int]int)

	// Parse each set of 3 bytes
	for i := 0; i < len(parameterBytes)-2; i += 3 {
		address, index, value := parseParameterBytes(parameterBytes[i : i+3])
		if index == 0 {
			klog.Info("0", address, index, value)

			continue
		}
		_, exists := parsedValues[address]
		if !exists {
			parsedValues[address] = make(map[int]int)
		}
		if address == 3 || address == 9 || address == 6 {
			klog.V(3).Info(address, index, value)
			address1, index1, value1 := parseChannelBytes(parameterBytes[i : i+3])
			klog.V(3).Info(address, index, value)
			parsedValues[address1][index1] = value1
		} else if address == 12 {
			parsedValues[address][index] = value
		} else if address == 1 || address == 4 || address == 7 || address == 10 {
			klog.V(3).Info(address, index, value)
			address1, index1, value1 := parseEQBytes(parameterBytes[i : i+3])
			klog.V(3).Info(address, index, value)
			parsedValues[address1][index1] = value1
		}

	}
	klog.V(3).Info(parsedValues)
	return parsedValues
}
func parseParameterBytes(parameterBytes []byte) (address, index, value int) {
	if len(parameterBytes) != 3 {
		// Ensure the input has exactly three bytes
		return -1, -1, -1
	}

	// Extracting values from the bytes
	address = int((parameterBytes[0] >> 3) & 0x0F)
	index = int(((parameterBytes[0] & 0x07) << 3) | ((parameterBytes[1] >> 4) & 0x07))

	// Combine upper bits and lower bits for parameter value
	upperBits := int(parameterBytes[1] & 0x0F)
	lowerBits := int(parameterBytes[2])
	value = (upperBits << 4) | lowerBits

	return address, index, value
}

func generateChannelCommand(channelAddress, parameterIndex, parameterValue int) []byte {
	if channelAddress < 0 || channelAddress > 15 || parameterIndex < 0 || parameterIndex > 31 || parameterValue < -2048 || parameterValue > 2047 {
		return nil
	}
	byte1 := byte((channelAddress << 3) | (parameterIndex >> 2))
	klog.V(3).Infof("%#v", byte1)

	byte2 := byte(((parameterIndex & 0x03) << 5) | ((parameterValue >> 7) & 0x1F))
	klog.V(3).Infof("%#v", byte2)

	byte3 := byte(parameterValue & 0x7F)
	klog.V(3).Infof("%#v", byte3)

	return []byte{byte1, byte2, byte3}
}

func generateDeviceCommand(parameterIndex, parameterValue int) []byte {
	address := 12 << 3
	klog.V(3).Infof("%#v", address)

	upperBits := ((parameterIndex >> 3) & 0x7)
	klog.V(3).Infof("%#v", upperBits)

	// Combining bits to form Byte 1
	byte1 := byte(address | upperBits)
	klog.V(3).Infof("%#v", byte1)

	// Byte 2: Lower 3 bits of Parameter Index and Upper 4 bits of Parameter Value
	byte2 := byte(((parameterIndex & 0x7) << 4) | ((parameterValue >> 3) & 0xF))
	klog.V(3).Infof("%#v", byte2)

	// Byte 3: Lower 7 bits of Parameter Value
	byte3 := byte(parameterValue & 0x7F)
	klog.V(3).Infof("%#v", byte3)

	return []byte{byte1, byte2, byte3}
}
func parseChannelBytes(data []byte) (channelAddress, parameterIndex, parameterValue int) {
	if len(data) != 3 {
		return 0, 0, 0
	}
	// Parsing Byte 1
	channelAddress = int((data[0] >> 3) & 0x0F)
	parameterIndex = int(((data[0] & 0x07) << 2) | ((data[1] >> 5) & 0x03))

	parameterValue = int(((int(data[1]) & 0x1F) << 7) | int(data[2]))

	if parameterValue > 2047 {
		parameterValue -= 4096
	}

	return channelAddress, parameterIndex, parameterValue
}

func parseEQBytes(data []byte) (channelAddress, parameterIndex, parameterValue int) {
	if len(data) != 3 {
		return 0, 0, 0
	}

	channelAddress = int((data[0] >> 3) & 0x0F)
	parameterIndex = int(((data[0] & 0x07) << 2) | ((data[1] >> 5) & 0x07))

	parameterValue = int(((int(data[1]) & 0x0F) << 4) | int(data[2]&0x0F))

	if (data[1] & 0x08) != 0 {
		parameterValue *= 10
	}

	if parameterValue > 2047 {
		parameterValue -= 4096
	}

	return channelAddress, parameterIndex, parameterValue
}
