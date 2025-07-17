package sysinfo

import (
	"fmt"

	"github.com/shirou/gopsutil/host"
)

func Sensors() {
	temps, _ := host.SensorsTemperatures()
	for _, t := range temps {
		fmt.Printf("Sensor: %s, Temp: %.1f°C\n", t.SensorKey, t.Temperature)
	}
}
