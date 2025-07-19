package ui

import (
	"fmt"
	"os"
	"sort"
	"syspulse/internal/plugins"
	"syspulse/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

type widgetPosition struct {
	widget tview.Primitive
	row    int
	column int
}

type TerminalSize struct {
	Width  int
	Height int
}

func GetTerminalSize() (int, int, error) {
	if width, height, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
		return width, height, nil
	}

	if width, height, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		return width, height, nil
	}

	if width, height, err := term.GetSize(int(os.Stderr.Fd())); err == nil {
		return width, height, nil
	}

	return 80, 24, nil
}

func (d *Dashboard) CheckTerminalSize() {
	width, height, err := GetTerminalSize()
	if err != nil {
		return
	}

	availableHeight := height - 3
	minWidth := 80
	minHeight := 24

	if width < minWidth || availableHeight < minHeight {
		d.ShowTerminalSizeWarning(width, height)
	}
}

func (d *Dashboard) ShowTerminalSizeWarning(width, height int) {
	warning := fmt.Sprintf("Terminal Size: %dx%d\n\nTerminal is too small for optimal display.\n\nFor best experience:\n• Increase terminal size to at least 80x24\n• Use fullscreen mode\n• Zoom out if needed\n\nContinue anyway or quit?", width, height)

	modal := tview.NewModal().
		SetText(warning).
		AddButtons([]string{"Continue", "Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				d.App.Stop()
			} else {
				d.App.SetRoot(d.MainWidget, true)
			}
		})

	d.App.SetRoot(modal, false).SetFocus(modal)
}

func (d *Dashboard) initMainLayout() {
	d.rebuildLayoutGrid()
	d.CheckTerminalSize()
}

func (d *Dashboard) rebuildLayoutGrid() {
	grid := tview.NewGrid().
		SetRows(make([]int, d.Theme.Layout.Rows)...).
		SetColumns(make([]int, d.Theme.Layout.Columns)...).
		SetMinSize(1, 1).
		SetGap(d.Theme.Layout.Spacing, d.Theme.Layout.Spacing)

	focusableWidgets := d.buildFocusableWidgets()
	d.addWidgetsToGrid(grid)
	d.setupInputHandlers(focusableWidgets)
	d.createMainWidget(grid)
}

func (d *Dashboard) buildFocusableWidgets() []tview.Primitive {
	widgetPositions := make([]widgetPosition, 0)

	widgetPositions = append(widgetPositions, widgetPosition{
		widget: d.HeaderWidget,
		row:    -1,
		column: 0,
	})

	if d.CpuWidget != nil && d.Theme.Layout.CPU.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.CpuWidget,
			row:    d.Theme.Layout.CPU.Row,
			column: d.Theme.Layout.CPU.Column,
		})
	}
	if d.MemWidget != nil && d.Theme.Layout.Memory.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.MemWidget,
			row:    d.Theme.Layout.Memory.Row,
			column: d.Theme.Layout.Memory.Column,
		})
	}
	if d.DiskWidget != nil && d.Theme.Layout.Disk.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.DiskWidget,
			row:    d.Theme.Layout.Disk.Row,
			column: d.Theme.Layout.Disk.Column,
		})
	}
	if d.NetWidget != nil && d.Theme.Layout.Network.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.NetWidget,
			row:    d.Theme.Layout.Network.Row,
			column: d.Theme.Layout.Network.Column,
		})
	}
	if d.ProcessWidget != nil && d.Theme.Layout.Process.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.ProcessWidget,
			row:    d.Theme.Layout.Process.Row,
			column: d.Theme.Layout.Process.Column,
		})
	}
	if d.GPUWidget != nil && d.Theme.Layout.GPU.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.GPUWidget,
			row:    d.Theme.Layout.GPU.Row,
			column: d.Theme.Layout.GPU.Column,
		})
	}
	if d.LoadWidget != nil && d.Theme.Layout.Load.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.LoadWidget,
			row:    d.Theme.Layout.Load.Row,
			column: d.Theme.Layout.Load.Column,
		})
	}
	if d.TemperatureWidget != nil && d.Theme.Layout.Temperature.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.TemperatureWidget,
			row:    d.Theme.Layout.Temperature.Row,
			column: d.Theme.Layout.Temperature.Column,
		})
	}
	if d.NetworkConnsWidget != nil && d.Theme.Layout.NetworkConns.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.NetworkConnsWidget,
			row:    d.Theme.Layout.NetworkConns.Row,
			column: d.Theme.Layout.NetworkConns.Column,
		})
	}
	if d.DiskIOWidget != nil && d.Theme.Layout.DiskIO.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.DiskIOWidget,
			row:    d.Theme.Layout.DiskIO.Row,
			column: d.Theme.Layout.DiskIO.Column,
		})
	}
	if d.ProcessTreeWidget != nil && d.Theme.Layout.ProcessTree.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.ProcessTreeWidget,
			row:    d.Theme.Layout.ProcessTree.Row,
			column: d.Theme.Layout.ProcessTree.Column,
		})
	}
	if d.BatteryWidget != nil && d.Theme.Layout.Battery.Enabled {
		widgetPositions = append(widgetPositions, widgetPosition{
			widget: d.BatteryWidget,
			row:    d.Theme.Layout.Battery.Row,
			column: d.Theme.Layout.Battery.Column,
		})
	}

	if d.PluginManager != nil {
		if pluginManager, ok := d.PluginManager.(*plugins.PluginManager); ok {
			pluginInfo := pluginManager.GetPluginInfo()

			for _, info := range pluginInfo {
				if info.Config.Enabled && info.Widget != nil {
					widgetPositions = append(widgetPositions, widgetPosition{
						widget: info.Widget,
						row:    info.Config.Layout.Row,
						column: info.Config.Layout.Column,
					})
				}
			}
		}
	}

	sort.Slice(widgetPositions, func(i, j int) bool {
		if widgetPositions[i].row == widgetPositions[j].row {
			return widgetPositions[i].column < widgetPositions[j].column
		}
		return widgetPositions[i].row < widgetPositions[j].row
	})

	focusableWidgets := make([]tview.Primitive, 0, len(widgetPositions))
	for _, wp := range widgetPositions {
		focusableWidgets = append(focusableWidgets, wp.widget)
	}

	return focusableWidgets
}

func (d *Dashboard) addWidgetsToGrid(grid *tview.Grid) {
	if d.Theme.Layout.CPU.Enabled && d.CpuWidget != nil {
		grid.AddItem(d.CpuWidget,
			d.Theme.Layout.CPU.Row, d.Theme.Layout.CPU.Column,
			d.Theme.Layout.CPU.RowSpan, d.Theme.Layout.CPU.ColSpan,
			d.Theme.Layout.CPU.MinWidth, 0, false)
	}

	if d.Theme.Layout.Memory.Enabled && d.MemWidget != nil {
		grid.AddItem(d.MemWidget,
			d.Theme.Layout.Memory.Row, d.Theme.Layout.Memory.Column,
			d.Theme.Layout.Memory.RowSpan, d.Theme.Layout.Memory.ColSpan,
			d.Theme.Layout.Memory.MinWidth, 0, false)
	}

	if d.Theme.Layout.Disk.Enabled && d.DiskWidget != nil {
		grid.AddItem(d.DiskWidget,
			d.Theme.Layout.Disk.Row, d.Theme.Layout.Disk.Column,
			d.Theme.Layout.Disk.RowSpan, d.Theme.Layout.Disk.ColSpan,
			d.Theme.Layout.Disk.MinWidth, 0, false)
	}

	if d.Theme.Layout.Network.Enabled && d.NetWidget != nil {
		grid.AddItem(d.NetWidget,
			d.Theme.Layout.Network.Row, d.Theme.Layout.Network.Column,
			d.Theme.Layout.Network.RowSpan, d.Theme.Layout.Network.ColSpan,
			d.Theme.Layout.Network.MinWidth, 0, false)
	}

	if d.Theme.Layout.Process.Enabled && d.ProcessWidget != nil {
		grid.AddItem(d.ProcessWidget,
			d.Theme.Layout.Process.Row, d.Theme.Layout.Process.Column,
			d.Theme.Layout.Process.RowSpan, d.Theme.Layout.Process.ColSpan,
			d.Theme.Layout.Process.MinWidth, 0, false)
	}

	if d.Theme.Layout.GPU.Enabled && d.GPUWidget != nil {
		grid.AddItem(d.GPUWidget,
			d.Theme.Layout.GPU.Row, d.Theme.Layout.GPU.Column,
			d.Theme.Layout.GPU.RowSpan, d.Theme.Layout.GPU.ColSpan,
			d.Theme.Layout.GPU.MinWidth, 0, false)
	}

	if d.Theme.Layout.Load.Enabled && d.LoadWidget != nil {
		grid.AddItem(d.LoadWidget,
			d.Theme.Layout.Load.Row, d.Theme.Layout.Load.Column,
			d.Theme.Layout.Load.RowSpan, d.Theme.Layout.Load.ColSpan,
			d.Theme.Layout.Load.MinWidth, 0, false)
	}

	if d.Theme.Layout.Temperature.Enabled && d.TemperatureWidget != nil {
		grid.AddItem(d.TemperatureWidget,
			d.Theme.Layout.Temperature.Row, d.Theme.Layout.Temperature.Column,
			d.Theme.Layout.Temperature.RowSpan, d.Theme.Layout.Temperature.ColSpan,
			d.Theme.Layout.Temperature.MinWidth, 0, false)
	}

	if d.Theme.Layout.NetworkConns.Enabled && d.NetworkConnsWidget != nil {
		grid.AddItem(d.NetworkConnsWidget,
			d.Theme.Layout.NetworkConns.Row, d.Theme.Layout.NetworkConns.Column,
			d.Theme.Layout.NetworkConns.RowSpan, d.Theme.Layout.NetworkConns.ColSpan,
			d.Theme.Layout.NetworkConns.MinWidth, 0, false)
	}

	if d.Theme.Layout.DiskIO.Enabled && d.DiskIOWidget != nil {
		grid.AddItem(d.DiskIOWidget,
			d.Theme.Layout.DiskIO.Row, d.Theme.Layout.DiskIO.Column,
			d.Theme.Layout.DiskIO.RowSpan, d.Theme.Layout.DiskIO.ColSpan,
			d.Theme.Layout.DiskIO.MinWidth, 0, false)
	}

	if d.Theme.Layout.ProcessTree.Enabled && d.ProcessTreeWidget != nil {
		grid.AddItem(d.ProcessTreeWidget,
			d.Theme.Layout.ProcessTree.Row, d.Theme.Layout.ProcessTree.Column,
			d.Theme.Layout.ProcessTree.RowSpan, d.Theme.Layout.ProcessTree.ColSpan,
			d.Theme.Layout.ProcessTree.MinWidth, 0, false)
	}

	if d.Theme.Layout.Battery.Enabled && d.BatteryWidget != nil {
		grid.AddItem(d.BatteryWidget,
			d.Theme.Layout.Battery.Row, d.Theme.Layout.Battery.Column,
			d.Theme.Layout.Battery.RowSpan, d.Theme.Layout.Battery.ColSpan,
			d.Theme.Layout.Battery.MinWidth, 0, false)
	}

	if d.PluginManager != nil {
		plugins.AddPluginWidgetsToGrid((*utils.Dashboard)(d), grid)
	}
}

func (d *Dashboard) setupInputHandlers(focusableWidgets []tview.Primitive) {
	currentFocus := -1

	d.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !d.InModalState {
			isMainWidgetActive := !d.InModalState

			currentFocused := d.App.GetFocus()
			isInputField := false
			if _, ok := currentFocused.(*tview.InputField); ok {
				isInputField = true
			}

			shouldProcessGlobalKeys := isMainWidgetActive && !isInputField

			switch event.Key() {
			case tcell.KeyTab:
				if isMainWidgetActive && len(focusableWidgets) > 0 {
					currentFocus = (currentFocus + 1) % len(focusableWidgets)
					d.App.SetFocus(focusableWidgets[currentFocus])
				}
				return nil
			case tcell.KeyBacktab:
				if isMainWidgetActive && len(focusableWidgets) > 0 {
					currentFocus--
					if currentFocus < 0 {
						currentFocus = len(focusableWidgets) - 1
					}
					d.App.SetFocus(focusableWidgets[currentFocus])
				}
				return nil
			}

			switch event.Rune() {
			case 'q', 'Q':
				d.quitModal()
				return nil
			case 'c', 'C':
				if shouldProcessGlobalKeys && d.CpuWidget != nil {
					d.App.SetFocus(d.CpuWidget)
				}
			case 'm', 'M':
				if shouldProcessGlobalKeys && d.MemWidget != nil {
					d.App.SetFocus(d.MemWidget)
				}
			case 'd', 'D':
				if shouldProcessGlobalKeys && d.DiskWidget != nil {
					d.App.SetFocus(d.DiskWidget)
				}
			case 'n', 'N':
				if shouldProcessGlobalKeys && d.NetWidget != nil {
					d.App.SetFocus(d.NetWidget)
				}
			case 'p', 'P':
				if shouldProcessGlobalKeys && d.ProcessWidget != nil {
					d.App.SetFocus(d.ProcessWidget)
				}
			case 'h', 'H':
				if shouldProcessGlobalKeys {
					d.showHelpModal()
				}
				return nil
			}
		}

		return event
	})
}

func (d *Dashboard) createMainWidget(grid *tview.Grid) {
	d.FooterWidget = tview.NewTextView().
		SetText("Press 'h' for help | TAB to cycle widgets | 'q' to quit").
		SetTextColor(tview.Styles.PrimaryTextColor)
	d.FooterWidget.SetBackgroundColor(utils.GetColorFromName(d.Theme.Background))

	d.MainWidget = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.HeaderWidget, 2, 1, false).
		AddItem(grid, 0, 8, true).
		AddItem(d.FooterWidget, 1, 1, false)

	d.MainWidget.SetBackgroundColor(utils.GetColorFromName(d.Theme.Background))

	grid.SetBackgroundColor(utils.GetColorFromName(d.Theme.Background))

	d.App.SetRoot(d.MainWidget, true)
}
