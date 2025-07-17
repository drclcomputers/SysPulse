//go:build !linux
// +build !linux

package temperature

func GetLinuxTemperatureInfo() (*TemperatureData, error) {
	return &TemperatureData{
		Sensors: make([]TemperatureSensor, 0),
	}, nil
}
