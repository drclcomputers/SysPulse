//go:build !darwin
// +build !darwin

package temperature

func GetDarwinTemperatureInfo() (*TemperatureData, error) {
	return &TemperatureData{
		Sensors: make([]TemperatureSensor, 0),
	}, nil
}
