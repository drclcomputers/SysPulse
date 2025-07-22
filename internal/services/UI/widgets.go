package ui

import (
	"fmt"
	"syspulse/internal/plugins"
	"syspulse/internal/services/battery"
	"syspulse/internal/services/disk"
	"syspulse/internal/services/gpu"
	"syspulse/internal/services/load"
	"syspulse/internal/services/memory"
	"syspulse/internal/services/network"
	"syspulse/internal/services/processes"
	"syspulse/internal/services/sysinfo"
	"syspulse/internal/services/temperature"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (d *Dashboard) initWidgets() {
	d.initHeaderWidget()
	d.initCPUWidget()
	d.initMemoryWidget()
	d.initDiskWidget()
	d.initNetworkWidget()
	d.initProcessWidget()
	d.initGPUWidget()
	d.initLoadWidget()
	d.initTemperatureWidget()
	d.initNetworkConnsWidget()
	d.initDiskIOWidget()
	d.initProcessTreeWidget()
	d.initBatteryWidget()
	d.initPluginSystem()
	d.initMainLayout()
}

func (d *Dashboard) initHeaderWidget() {
	d.HeaderWidget = tview.NewBox()
	d.HeaderWidget.SetTitle(createHeaderTitle()).
		SetTitleColor(tview.Styles.PrimaryTextColor)
	d.HeaderWidget.SetTitleAlign(tview.AlignCenter)
	utils.SetBorderStyle(d.HeaderWidget)
	d.HeaderWidget.SetBackgroundColor(utils.GetColorFromName(d.Theme.Background))
	d.HeaderWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Rune()
		switch key {
		case 'q', 'Q':
			d.quitModal()
			return nil
		case 'i', 'I', rune(tcell.KeyEnter):
			modal := tview.NewModal().
				SetText(sysinfo.GetSystemInfo()).
				AddButtons([]string{"Close"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					d.App.SetRoot(d.MainWidget, true).SetFocus(d.HeaderWidget)
				})
			modal.SetTitle("System Overview")
			d.App.SetRoot(modal, true).SetFocus(modal)
		}
		return nil
	})
}

func (d *Dashboard) initCPUWidget() {
	if d.Theme.Layout.CPU.Enabled {
		d.CpuWidget = tview.NewBox()
		utils.SetBorderStyle(d.CpuWidget)
		d.CpuWidget.SetTitle(fmt.Sprint("CPU | ", sysinfo.GetCpuName())).
			SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				key := event.Rune()
				switch key {
				case 'q', 'Q':
					d.quitModal()
					return nil
				case 'i', 'I', rune(tcell.KeyEnter):
					textView := tview.NewTextView().
						SetDynamicColors(true).
						SetRegions(true).
						SetWordWrap(true).
						SetScrollable(true).
						SetText(sysinfo.GetCpuFormattedInfo())

					utils.SetBorderStyle(textView.Box)
					textView.SetTitle("CPU Information (Arrow keys to scroll, ESC to close)").
						SetTitleAlign(tview.AlignCenter)

					textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						switch event.Key() {
						case tcell.KeyEscape:
							d.App.SetRoot(d.MainWidget, true).SetFocus(d.CpuWidget)
							return nil
						}

						switch event.Rune() {
						case 'q', 'Q':
							d.App.SetRoot(d.MainWidget, true).SetFocus(d.CpuWidget)
							return nil
						}

						return event
					})

					flex := tview.NewFlex().
						AddItem(nil, 0, 1, false).
						AddItem(tview.NewFlex().
							SetDirection(tview.FlexRow).
							AddItem(nil, 0, 1, false).
							AddItem(textView, 0, 10, true).
							AddItem(nil, 0, 1, false), 0, 10, true).
						AddItem(nil, 0, 1, false)

					d.App.SetRoot(flex, true).SetFocus(textView)
				}
				return nil
			})

		if d.Theme.Layout.CPU.BorderColor != "" {
			d.CpuWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.CPU.BorderColor))
		}
		if d.Theme.Layout.CPU.ForegroundColor != "" {
			d.CpuWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.CPU.ForegroundColor))
		}
	}
}

func (d *Dashboard) initMemoryWidget() {
	if d.Theme.Layout.Memory.Enabled {
		d.MemWidget = tview.NewBox()
		utils.SetBorderStyle(d.MemWidget)
		d.MemWidget.SetTitle(fmt.Sprint("Memory Usage | ", memory.GetRAM())).
			SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				key := event.Rune()
				switch key {
				case 'q', 'Q':
					d.quitModal()
					return nil
				case 'i', 'I', rune(tcell.KeyEnter):
					modal := tview.NewModal().
						SetText(memory.GetMemoryFormattedInfo()).
						AddButtons([]string{"Close"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							d.App.SetRoot(d.MainWidget, true).SetFocus(d.MemWidget)
						})
					modal.SetTitle("Memory Usage Information")
					d.App.SetRoot(modal, true).SetFocus(modal)
				}
				return nil
			})

		if d.Theme.Layout.Memory.BorderColor != "" {
			d.MemWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.Memory.BorderColor))
		}
		if d.Theme.Layout.Memory.ForegroundColor != "" {
			d.MemWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.Memory.ForegroundColor))
		}
	}
}

func (d *Dashboard) initDiskWidget() {
	if d.Theme.Layout.Disk.Enabled {
		d.DiskWidget = tview.NewBox()
		utils.SetBorderStyle(d.DiskWidget)
		d.DiskWidget.SetTitle(fmt.Sprint("Disk Usage | ", disk.GetNumberofPartitions())).
			SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				key := event.Rune()
				switch key {
				case 'q', 'Q':
					d.quitModal()
					return nil
				case 'i', 'I', rune(tcell.KeyEnter):
					textView := tview.NewTextView().
						SetDynamicColors(true).
						SetRegions(true).
						SetWordWrap(true).
						SetScrollable(true).
						SetText(disk.GetDiskFormattedInfo())

					utils.SetBorderStyle(textView.Box)
					textView.SetTitle("Disk Information (Arrow keys to scroll, ESC to close)").
						SetTitleAlign(tview.AlignCenter)

					textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						switch event.Key() {
						case tcell.KeyEscape:
							d.App.SetRoot(d.MainWidget, true).SetFocus(d.DiskWidget)
							return nil
						}

						switch event.Rune() {
						case 'q', 'Q':
							d.App.SetRoot(d.MainWidget, true).SetFocus(d.DiskWidget)
							return nil
						}

						return event
					})

					flex := tview.NewFlex().
						AddItem(nil, 0, 1, false).
						AddItem(tview.NewFlex().
							SetDirection(tview.FlexRow).
							AddItem(nil, 0, 1, false).
							AddItem(textView, 0, 5, true).
							AddItem(nil, 0, 1, false), 0, 5, true).
						AddItem(nil, 0, 1, false)

					d.App.SetRoot(flex, true).SetFocus(textView)
				}
				return nil
			})

		if d.Theme.Layout.Disk.BorderColor != "" {
			d.DiskWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.Disk.BorderColor))
		}
		if d.Theme.Layout.Disk.ForegroundColor != "" {
			d.DiskWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.Disk.ForegroundColor))
		}
	}
}

func (d *Dashboard) initNetworkWidget() {
	if d.Theme.Layout.Network.Enabled {
		d.NetWidget = tview.NewBox()
		utils.SetBorderStyle(d.NetWidget)
		d.NetWidget.SetTitle("Network Activity / Interfaces").
			SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				key := event.Rune()
				switch key {
				case 'q', 'Q':
					d.quitModal()
					return nil
				case 'i', 'I', rune(tcell.KeyEnter):
					textView := tview.NewTextView().
						SetDynamicColors(true).
						SetRegions(true).
						SetWordWrap(true).
						SetScrollable(true).
						SetText(network.GetNetworkFormattedInfo())

					utils.SetBorderStyle(textView.Box)
					textView.SetTitle("Network Activity & Interfaces (Arrow keys to scroll, ESC to close)").
						SetTitleAlign(tview.AlignCenter)

					textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						switch event.Key() {
						case tcell.KeyEscape:
							d.App.SetRoot(d.MainWidget, true).SetFocus(d.NetWidget)
							return nil
						}

						switch event.Rune() {
						case 'q', 'Q':
							d.App.SetRoot(d.MainWidget, true).SetFocus(d.NetWidget)
							return nil
						}

						return event
					})

					flex := tview.NewFlex().
						AddItem(nil, 0, 1, false).
						AddItem(tview.NewFlex().
							SetDirection(tview.FlexRow).
							AddItem(nil, 0, 1, false).
							AddItem(textView, 0, 6, true).
							AddItem(nil, 0, 1, false), 0, 6, true).
						AddItem(nil, 0, 1, false)

					d.App.SetRoot(flex, true).SetFocus(textView)
				}
				return nil
			})

		if d.Theme.Layout.Network.BorderColor != "" {
			d.NetWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.Network.BorderColor))
		}
		if d.Theme.Layout.Network.ForegroundColor != "" {
			d.NetWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.Network.ForegroundColor))
		}
	}
}

func (d *Dashboard) initGPUWidget() {
	if d.Theme.Layout.GPU.Enabled {
		d.GPUWidget = tview.NewBox()
		utils.SetBorderStyle(d.GPUWidget)
		d.GPUWidget.SetTitle(fmt.Sprint("GPU | ", gpu.GetGPUTitle())).
			SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				key := event.Rune()
				switch key {
				case 'q', 'Q':
					d.quitModal()
					return nil
				case 'i', 'I', rune(tcell.KeyEnter):
					gpuInfoView := tview.NewTextView().
						SetText(gpu.GetGPUFormattedInfo()).
						SetScrollable(true).
						SetWrap(true)
					utils.SetBorderStyle(gpuInfoView.Box)
					gpuInfoView.SetTitle("GPU Information (Arrow keys to scroll, ESC to close)").
						SetTitleAlign(tview.AlignCenter)
					gpuInfoView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						if event.Key() == tcell.KeyEscape {
							d.App.SetRoot(d.MainWidget, true).SetFocus(d.GPUWidget)
							return nil
						}
						return event
					})

					flex := tview.NewFlex().
						AddItem(nil, 0, 1, false).
						AddItem(tview.NewFlex().
							SetDirection(tview.FlexRow).
							AddItem(nil, 0, 1, false).
							AddItem(gpuInfoView, 0, 3, true).
							AddItem(nil, 0, 1, false), 0, 3, true).
						AddItem(nil, 0, 1, false)

					d.App.SetRoot(flex, true).SetFocus(gpuInfoView)
				}
				return nil
			})

		if d.Theme.Layout.GPU.BorderColor != "" {
			d.GPUWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.GPU.BorderColor))
		}
		if d.Theme.Layout.GPU.ForegroundColor != "" {
			d.GPUWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.GPU.ForegroundColor))
		}
	}
}

func (d *Dashboard) initProcessWidget() {
	if d.Theme.Layout.Process.Enabled {
		d.ProcessWidget = tview.NewList()
		utils.SetBorderStyle(d.ProcessWidget.Box)
		d.ProcessWidget.SetTitle(fmt.Sprint("Processes - ", processes.GetNrProcesses()))
		d.ProcessWidget.ShowSecondaryText(false)
		d.ProcessWidget.SetSelectedBackgroundColor(tcell.ColorDarkBlue)
		d.ProcessWidget.SetSelectedTextColor(utils.GetColorFromName(d.Theme.Altforeground))
		d.ProcessWidget.SetWrapAround(true)
		d.ProcessWidget.SetMainTextColor(utils.GetColorFromName(d.Theme.Layout.Process.ForegroundColor))
		d.ProcessWidget.SetInputCapture(d.getProcessInputHandler())

		if d.Theme.Layout.Process.BorderColor != "" {
			d.ProcessWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.Process.BorderColor))
		}
		if d.Theme.Layout.Process.ForegroundColor != "" {
			d.ProcessWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Foreground))
		}
	}
}

func (d *Dashboard) getProcessInputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			current := d.ProcessWidget.GetCurrentItem()
			if current > 0 {
				d.ProcessWidget.SetCurrentItem(current - 1)
			}
			return nil
		case tcell.KeyDown:
			current := d.ProcessWidget.GetCurrentItem()
			if current < d.ProcessWidget.GetItemCount()-1 {
				d.ProcessWidget.SetCurrentItem(current + 1)
			}
			return nil
		case tcell.KeyEnter:
			processes.ShowProcessDetails((*utils.Dashboard)(d))
			return nil
		}

		switch event.Rune() {
		case 'y', 'Y':
			if d.Theme.Sorting == "cpu" {
				d.Theme.Sorting = "mem"
			} else {
				d.Theme.Sorting = "cpu"
			}
		case 'i', 'I':
			processes.ShowProcessDetails((*utils.Dashboard)(d))
			return nil
		case 'f', 'F':
			d.showProcessSearch()
			return nil
		case 'q', 'Q':
			d.quitModal()
			return nil
		case 'w', 'W', 'o', 'O':
			current := d.ProcessWidget.GetCurrentItem()
			if current > 0 {
				d.ProcessWidget.SetCurrentItem(current - 1)
			}
			return nil
		case 's', 'S', 'l', 'L':
			current := d.ProcessWidget.GetCurrentItem()
			if current < d.ProcessWidget.GetItemCount()-1 {
				d.ProcessWidget.SetCurrentItem(current + 1)
			}
			return nil
		case 'k', 'K':
			var selectedPID int32
			currentItem := d.ProcessWidget.GetCurrentItem()
			if currentItem >= 0 && currentItem < d.ProcessWidget.GetItemCount() {
				text, _ := d.ProcessWidget.GetItemText(currentItem)
				fmt.Sscanf(text, "%s (PID: %d)", new(string), &selectedPID)
			}
			d.showProcessKillModal(selectedPID)
		}
		return event
	}
}

func (d *Dashboard) initLoadWidget() {
	d.LoadWidget = tview.NewBox()
	utils.SetBorderStyle(d.LoadWidget)
	d.LoadWidget.SetTitle("Load Average").
		SetTitleAlign(tview.AlignCenter)
	d.LoadWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Rune()
		switch key {
		case 'i', 'I', rune(tcell.KeyEnter):
			modal := tview.NewModal().
				SetText(load.GetLoadFormattedInfo()).
				AddButtons([]string{"Close"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					d.App.SetRoot(d.MainWidget, true).SetFocus(d.LoadWidget)
				})
			modal.SetTitle("System Load Average")
			d.App.SetRoot(modal, true).SetFocus(modal)
		}
		return nil
	})

	if d.Theme.Layout.Load.BorderColor != "" {
		d.LoadWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.Load.BorderColor))
	}
	if d.Theme.Layout.Load.ForegroundColor != "" {
		d.LoadWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.Load.ForegroundColor))
	}
}

func (d *Dashboard) initTemperatureWidget() {
	d.TemperatureWidget = tview.NewBox()
	utils.SetBorderStyle(d.TemperatureWidget)
	d.TemperatureWidget.SetTitle("Temperature").
		SetTitleAlign(tview.AlignCenter)
	d.TemperatureWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Rune()
		switch key {
		case 'i', 'I', rune(tcell.KeyEnter):
			tempInfoView := tview.NewTextView().
				SetText(temperature.GetTemperatureFormattedInfo()).
				SetScrollable(true).
				SetWrap(true)
			utils.SetBorderStyle(tempInfoView.Box)
			tempInfoView.SetTitle("System Temperature Monitoring (Arrow keys to scroll, ESC to close)").
				SetTitleAlign(tview.AlignCenter)
			tempInfoView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEscape {
					d.App.SetRoot(d.MainWidget, true).SetFocus(d.TemperatureWidget)
					return nil
				}
				return event
			})

			flex := tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(nil, 0, 1, false).
					AddItem(tempInfoView, 0, 3, true).
					AddItem(nil, 0, 1, false), 0, 3, true).
				AddItem(nil, 0, 1, false)

			d.App.SetRoot(flex, true).SetFocus(tempInfoView)
		}
		return nil
	})

	if d.Theme.Layout.Temperature.BorderColor != "" {
		d.TemperatureWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.Temperature.BorderColor))
	}
	if d.Theme.Layout.Temperature.ForegroundColor != "" {
		d.TemperatureWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.Temperature.ForegroundColor))
	}
}

func (d *Dashboard) initNetworkConnsWidget() {
	d.NetworkConnsWidget = tview.NewBox()
	utils.SetBorderStyle(d.NetworkConnsWidget)
	d.NetworkConnsWidget.SetTitle("Network Connections").
		SetTitleAlign(tview.AlignCenter)
	d.NetworkConnsWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Rune()
		switch key {
		case 'i', 'I', rune(tcell.KeyEnter):
			modal := network.CreateNetworkConnectionsModal(d.App, d.MainWidget)
			d.App.SetRoot(modal, true).SetFocus(modal)
		}
		return nil
	})

	if d.Theme.Layout.NetworkConns.BorderColor != "" {
		d.NetworkConnsWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.NetworkConns.BorderColor))
	}
	if d.Theme.Layout.NetworkConns.ForegroundColor != "" {
		d.NetworkConnsWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.NetworkConns.ForegroundColor))
	}
}

func (d *Dashboard) initDiskIOWidget() {
	d.DiskIOWidget = tview.NewBox()
	utils.SetBorderStyle(d.DiskIOWidget)
	d.DiskIOWidget.SetTitle("Disk I/O").
		SetTitleAlign(tview.AlignCenter)
	d.DiskIOWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Rune()
		switch key {
		case 'i', 'I', rune(tcell.KeyEnter):
			textView := tview.NewTextView().
				SetDynamicColors(true).
				SetRegions(true).
				SetWordWrap(true).
				SetScrollable(true).
				SetText(disk.GetDiskIOFormattedInfo())

			utils.SetBorderStyle(textView.Box)
			textView.SetTitle("Disk I/O Information (Arrow keys to scroll, ESC to close)").
				SetTitleAlign(tview.AlignCenter)

			textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Key() {
				case tcell.KeyEscape:
					d.App.SetRoot(d.MainWidget, true).SetFocus(d.DiskIOWidget)
					return nil
				}

				switch event.Rune() {
				case 'q', 'Q':
					d.App.SetRoot(d.MainWidget, true).SetFocus(d.DiskIOWidget)
					return nil
				}

				return event
			})

			flex := tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(nil, 0, 1, false).
					AddItem(textView, 0, 7, true).
					AddItem(nil, 0, 1, false), 0, 7, true).
				AddItem(nil, 0, 1, false)

			d.App.SetRoot(flex, true).SetFocus(textView)
		}
		return nil
	})

	if d.Theme.Layout.DiskIO.BorderColor != "" {
		d.DiskIOWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.DiskIO.BorderColor))
	}
	if d.Theme.Layout.DiskIO.ForegroundColor != "" {
		d.DiskIOWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.DiskIO.ForegroundColor))
	}
}

func (d *Dashboard) initProcessTreeWidget() {
	d.ProcessTreeWidget = tview.NewBox()
	utils.SetBorderStyle(d.ProcessTreeWidget)
	d.ProcessTreeWidget.SetTitle("Process Tree").
		SetTitleAlign(tview.AlignCenter)
	d.ProcessTreeWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Rune()
		switch key {
		case 'i', 'I', rune(tcell.KeyEnter):
			d.showProcessTreeModal()
		}
		return nil
	})

	if d.Theme.Layout.ProcessTree.BorderColor != "" {
		d.ProcessTreeWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.ProcessTree.BorderColor))
	}
	if d.Theme.Layout.ProcessTree.ForegroundColor != "" {
		d.ProcessTreeWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.ProcessTree.ForegroundColor))
	}
}

func (d *Dashboard) initBatteryWidget() {
	d.BatteryWidget = tview.NewBox()
	utils.SetBorderStyle(d.BatteryWidget)
	d.BatteryWidget.SetTitle("Battery").
		SetTitleAlign(tview.AlignCenter)
	d.BatteryWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Rune()
		switch key {
		case 'i', 'I', rune(tcell.KeyEnter):
			modal := tview.NewModal().
				SetText(battery.GetBatteryFormattedInfo()).
				AddButtons([]string{"Close"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					d.App.SetRoot(d.MainWidget, true).SetFocus(d.BatteryWidget)
				})
			modal.SetTitle("Battery Information")
			d.App.SetRoot(modal, true).SetFocus(modal)
		}
		return nil
	})

	if d.Theme.Layout.Battery.BorderColor != "" {
		d.BatteryWidget.SetBorderColor(utils.GetColorFromName(d.Theme.Layout.Battery.BorderColor))
	}
	if d.Theme.Layout.Battery.ForegroundColor != "" {
		d.BatteryWidget.SetTitleColor(utils.GetColorFromName(d.Theme.Layout.Battery.ForegroundColor))
	}
}

func (d *Dashboard) initPluginSystem() {
	if err := plugins.InitializePluginSystem((*utils.Dashboard)(d)); err != nil {
		fmt.Printf("Failed to initialize plugin system: %v\n", err)
	}
}
