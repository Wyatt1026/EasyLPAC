package easylpac

import (
	"fmt"
	"os"
	"path/filepath"

	appconfig "EasyLPAC/internal/easylpac/config"
	appi18n "EasyLPAC/internal/easylpac/i18n"
	"EasyLPAC/internal/easylpac/model"
	"EasyLPAC/internal/easylpac/registry"
	"EasyLPAC/internal/easylpac/resources"
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

const Version = "development"
const EUICCDataVersion = "20260227"

var App fyne.App

func init() {
	registry.InitCiRegistry()
	registry.InitEumRegistry()
	App = fyneapp.New()
	appi18n.Init(App.Preferences())
	App.Settings().SetTheme(resources.NewTheme(appi18n.CurrentLanguageTag()))

	if err := appconfig.Load(); err != nil {
		panic(err)
	}
	if _, err := os.Stat(appconfig.Instance.LogDir); os.IsNotExist(err) {
		err := os.Mkdir(appconfig.Instance.LogDir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func Run() {
	var err error
	appconfig.Instance.LogFile, err = os.Create(filepath.Join(appconfig.Instance.LogDir, appconfig.Instance.LogFilename))
	if err != nil {
		panic(err)
	}
	defer appconfig.Instance.LogFile.Close()

	InitWidgets()
	go UpdateStatusBarListener()
	go LockButtonListener()

	WMain = InitMainWindow()

	_, err = os.Stat(filepath.Join(appconfig.Instance.LpacDir, appconfig.Instance.EXEName))
	if err != nil {
		d := dialog.NewError(fmt.Errorf(" %s", appi18n.TR.Trans("message.lpac_not_found")), WMain)
		d.SetOnClosed(func() {
			os.Exit(127)
		})
		d.Show()
	} else {
		if version, err2 := LpacVersion(); err2 != nil {
			DetectedLpacVersion = ""
			updateLpacVersionLabel()
		} else {
			DetectedLpacVersion = version
			updateLpacVersionLabel()
		}
		RefreshApduDriver()
		if model.ApduDrivers != nil {
			ApduDriverSelect.SetSelectedIndex(0)
		}
	}

	WMain.Show()
	App.Run()
}
