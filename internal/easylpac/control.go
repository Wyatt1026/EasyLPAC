package easylpac

import (
	appconfig "EasyLPAC/internal/easylpac/config"
	appformat "EasyLPAC/internal/easylpac/format"
	appi18n "EasyLPAC/internal/easylpac/i18n"
	"EasyLPAC/internal/easylpac/model"
	"EasyLPAC/internal/easylpac/registry"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const StatusProcess = 1
const StatusReady = 0
const Unselected = -1

var SelectedProfile = Unselected
var SelectedNotification = Unselected

var RefreshNeeded = true
var ProfileMaskNeeded bool
var NotificationMaskNeeded bool
var ProfileStateAllowDisable bool
var CurrentStatus = StatusReady
var ControlsLocked bool
var DetectedLpacVersion string

var StatusChan = make(chan int)
var LockButtonChan = make(chan bool)

func applyStatusState() {
	if StatusLabel == nil || StatusProcessBar == nil {
		return
	}
	switch CurrentStatus {
	case StatusProcess:
		StatusLabel.SetText(appi18n.TR.Trans("label.status_processing"))
		StatusProcessBar.Start()
		StatusProcessBar.Show()
	default:
		StatusLabel.SetText(appi18n.TR.Trans("label.status_ready"))
		StatusProcessBar.Stop()
		StatusProcessBar.Hide()
	}
}

func updateLpacVersionLabel() {
	if LpacVersionLabel == nil {
		return
	}
	if DetectedLpacVersion == "" {
		LpacVersionLabel.SetText(appi18n.TR.Trans("label.lpac_version_unknown"))
		return
	}
	LpacVersionLabel.SetText(appi18n.TR.Trans("label.lpac_version") + " " + DetectedLpacVersion)
}

func updateChipInfoView() {
	if model.ChipInfo == nil {
		return
	}

	convertToString := func(value interface{}) string {
		if value == nil {
			return appi18n.TR.Trans("label.not_set")
		}
		if str, ok := value.(string); ok {
			return str
		}
		return appi18n.TR.Trans("label.not_set")
	}

	EidLabel.SetText(fmt.Sprintf(appi18n.TR.Trans("label.info_eid")+" %s", model.ChipInfo.EidValue))
	DefaultDpAddressLabel.SetText(fmt.Sprintf(appi18n.TR.Trans("label.default_smdp_address")+"  %s", convertToString(model.ChipInfo.EuiccConfiguredAddresses.DefaultDpAddress)))
	RootDsAddressLabel.SetText(fmt.Sprintf(appi18n.TR.Trans("label.root_smds_address")+"  %s", convertToString(model.ChipInfo.EuiccConfiguredAddresses.RootDsAddress)))
	if eum := registry.GetEUM(model.ChipInfo.EidValue); eum != nil {
		manufacturer := fmt.Sprint(eum.Manufacturer, " ", appformat.CountryCodeToEmoji(eum.Country))
		EUICCManufacturerLabel.SetText(appi18n.TR.Trans("label.manufacturer") + " " + manufacturer)
	} else {
		EUICCManufacturerLabel.SetText(appi18n.TR.Trans("label.manufacturer_unknown"))
	}
	bytes, err := json.MarshalIndent(model.ChipInfo.EUICCInfo2, "", "  ")
	if err != nil {
		ShowLpacErrDialog(fmt.Errorf(appi18n.TR.Trans("message.failed_to_decode_euiccinfo2")+"\n%s", err))
	}
	EuiccInfo2Entry.SetText(string(bytes))
	freeSpace := float64(model.ChipInfo.EUICCInfo2.ExtCardResource.FreeNonVolatileMemory) / 1024
	FreeSpaceLabel.SetText(fmt.Sprintf(appi18n.TR.Trans("label.free_space")+" %.2f KiB", math.Round(freeSpace*100)/100))

	CopyEidButton.Show()
	SetDefaultSmdpButton.Show()
	EuiccInfo2Entry.Show()
	ViewCertInfoButton.Show()
	EUICCManufacturerLabel.Show()
	CopyEuiccInfo2Button.Show()
}

func currentButtons() []*widget.Button {
	return []*widget.Button{
		RefreshButton, DownloadButton, SetNicknameButton, SwitchStateButton, DeleteProfileButton,
		ProcessNotificationButton, ProcessAllNotificationButton, RemoveNotificationButton, BatchRemoveNotificationButton,
		SetDefaultSmdpButton, ApduDriverRefreshButton,
	}
}

func currentChecks() []*widget.Check {
	return []*widget.Check{
		ProfileMaskCheck, NotificationMaskCheck,
	}
}

func currentSelects() []*widget.Select {
	return []*widget.Select{
		ApduDriverSelect, LanguageSelect,
	}
}

func applyLockState() {
	for _, button := range currentButtons() {
		if button == nil {
			continue
		}
		if ControlsLocked {
			button.Disable()
		} else {
			button.Enable()
		}
	}
	for _, check := range currentChecks() {
		if check == nil {
			continue
		}
		if ControlsLocked {
			check.Disable()
		} else {
			check.Enable()
		}
	}
	for _, selectWidget := range currentSelects() {
		if selectWidget == nil {
			continue
		}
		if ControlsLocked {
			selectWidget.Disable()
		} else {
			selectWidget.Enable()
		}
	}
}

func RefreshCurrentUIState() {
	applyStatusState()
	updateLpacVersionLabel()
	if ProfileList != nil {
		ProfileList.Refresh()
	}
	if NotificationList != nil {
		NotificationList.Refresh()
	}
	updateChipInfoView()
	applyLockState()
}

func RefreshProfile() error {
	var err error
	model.Profiles, err = LpacProfileList()
	if err != nil {
		return err
	}
	// 刷新 List
	fyne.Do(func() {
		ProfileList.Refresh()
		ProfileList.UnselectAll()
		SwitchStateButton.SetText(appi18n.TR.Trans("label.switch_state_button_enable"))
		SwitchStateButton.SetIcon(theme.ConfirmIcon())
	})
	return nil
}

func RefreshNotification() error {
	var err error
	model.Notifications, err = LpacNotificationList()
	if err != nil {
		return err
	}
	sort.Slice(model.Notifications, func(i, j int) bool {
		return model.Notifications[i].SeqNumber < model.Notifications[j].SeqNumber
	})
	// 刷新 List
	fyne.Do(func() {
		NotificationList.Refresh()
		NotificationList.UnselectAll()
	})
	return nil
}

func RefreshChipInfo() error {
	var err error
	model.ChipInfo, err = LpacChipInfo()
	if err != nil {
		return err
	}
	if model.ChipInfo == nil {
		return nil
	}

	fyne.Do(func() {
		updateChipInfoView()
	})
	return nil
}

func RefreshApduDriver() {
	var err error
	model.ApduDrivers, err = LpacDriverApduList()
	if err != nil {
		ShowLpacErrDialog(err)
	}
	var options []string
	for _, d := range model.ApduDrivers {
		// exclude YubiKey and CanoKey
		if strings.Contains(d.Name, "canokeys.org") || strings.Contains(d.Name, "YubiKey") {
			continue
		}
		// Workaround: lpac shows an empty driver when no card reader inserted under macOS
		if d.Name == "" {
			continue
		}
		options = append(options, d.Name)
	}
	ApduDriverSelect.SetOptions(options)
	ApduDriverSelect.ClearSelected()
	appconfig.Instance.DriverIFID = ""
	ApduDriverSelect.Refresh()
}

func OpenLog() {
	if err := OpenProgram(appconfig.Instance.LogDir); err != nil {
		d := dialog.NewError(err, WMain)
		d.Show()
	}
}

func OpenProgram(name string) error {
	var launcher string
	switch runtime.GOOS {
	case "windows":
		launcher = "explorer"
	case "darwin":
		launcher = "open"
	case "linux":
		launcher = "xdg-open"
	}
	if launcher == "" {
		return fmt.Errorf("unsupported platform, failed to open")
	}
	return exec.Command(launcher, name).Start()
}

func Refresh() {
	if appconfig.Instance.DriverIFID == "" {
		ShowSelectCardReaderDialog()
		return
	}
	err := RefreshProfile()
	if err != nil {
		ShowLpacErrDialog(err)
		return
	}
	err = RefreshNotification()
	if err != nil {
		ShowLpacErrDialog(err)
		return
	}
	err = RefreshChipInfo()
	if err != nil {
		ShowLpacErrDialog(err)
		return
	}
	RefreshNeeded = false
}

func UpdateStatusBarListener() {
	for {
		status := <-StatusChan
		CurrentStatus = status
		fyne.Do(func() {
			applyStatusState()
		})
	}
}

func LockButtonListener() {
	for {
		lock := <-LockButtonChan
		ControlsLocked = lock
		fyne.Do(func() {
			applyLockState()
		})
	}
}

func SetDriverIFID(name string) {
	for _, d := range model.ApduDrivers {
		if name == d.Name {
			// 未选择过读卡器
			if appconfig.Instance.DriverIFID == "" {
				appconfig.Instance.DriverIFID = d.Env
			} else {
				// 选择过读卡器，要求刷新
				appconfig.Instance.DriverIFID = d.Env
				RefreshNeeded = true
			}
		}
	}
}
