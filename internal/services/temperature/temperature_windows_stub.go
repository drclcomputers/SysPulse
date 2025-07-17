//go:build !windows
// +build !windows

package temperature

func GetWindowsTemperatureInfo() (*TemperatureData, error) {
	return &TemperatureData{
		Sensors: make([]TemperatureSensor, 0),
	}, nil
}
