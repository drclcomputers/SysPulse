package ui

import (
	"fmt"
	"os"

	"syspulse/internal/export"
	loggerv2 "syspulse/internal/logger/v2"
	"syspulse/internal/utils"
)

var (
	log        *loggerv2.Logger
	dataPoints []export.DataPoint
)

func init() {
	var err error
	log, err = loggerv2.New("logs", loggerv2.INFO)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	dataPoints = make([]export.DataPoint, 0)
}

func NewDashboard() *utils.Dashboard {
	return newDashboard()
}

func Run(d *utils.Dashboard) error {
	defer log.Close()

	quit := make(chan struct{})

	startWorkers(d, quit)

	startExportWorker(d, quit)

	d.App.EnableMouse(true)
	d.App.EnablePaste(true)

	log.Info("Starting SysPulse application")
	err := d.App.Run()
	close(quit)

	performFinalExport(d)

	log.Info("Shutting down SysPulse application")
	return err
}

func StartUI() {
	if err := Run(NewDashboard()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
