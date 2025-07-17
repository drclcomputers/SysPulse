package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"syspulse/internal/errors"
	"syspulse/internal/utils"

	"github.com/rivo/tview"
)

func newDashboard() *utils.Dashboard {
	d := &utils.Dashboard{
		App: tview.NewApplication(),
	}
	if err := (*Dashboard)(d).loadTheme(); err != nil {
		log.Fatal(fmt.Sprintf("Failed to load theme: %v", err))
	}
	(*Dashboard)(d).applyThemeColors()
	(*Dashboard)(d).initWidgets()
	return d
}

func (d *Dashboard) loadTheme() error {
	var themeData utils.Theme

	dataFile, err := os.ReadFile("config.json")
	if err != nil {
		dataFileBackup, err := os.ReadFile("internal/services/UI/default.json")
		if err != nil {
			return errors.NewAppError(errors.ConfigError,
				"Cannot load neither 'config.json' nor 'default.json'", err)
		}
		if err = json.Unmarshal(dataFileBackup, &themeData); err != nil {
			return errors.NewAppError(errors.ConfigError,
				"Failed to parse default.json", err)
		}
	} else {
		if err = json.Unmarshal(dataFile, &themeData); err != nil {
			return errors.NewAppError(errors.ConfigError,
				"Failed to parse config.json", err)
		}
	}

	if err = utils.Validate(themeData); err != nil {
		return errors.Wrap(err, "Theme validation failed")
	}

	d.Theme = themeData

	return nil
}

func (d *Dashboard) applyThemeColors() {
	backgroundColor := utils.GetColorFromName(d.Theme.Background)
	foregroundColor := utils.GetColorFromName(d.Theme.Foreground)

	tview.Styles.PrimitiveBackgroundColor = backgroundColor
	tview.Styles.ContrastBackgroundColor = backgroundColor
	tview.Styles.PrimaryTextColor = foregroundColor
	tview.Styles.SecondaryTextColor = utils.GetColorFromName(d.Theme.Altforeground)
	tview.Styles.TertiaryTextColor = foregroundColor
	tview.Styles.InverseTextColor = backgroundColor
	tview.Styles.ContrastSecondaryTextColor = foregroundColor

	tview.Styles.TitleColor = foregroundColor
	tview.Styles.BorderColor = foregroundColor
	tview.Styles.GraphicsColor = foregroundColor
}
