package easylpac

import (
	"errors"
	"fmt"
	"image/color"
	"strings"

	"EasyLPAC/internal/easylpac/activation"
	appconfig "EasyLPAC/internal/easylpac/config"
	appi18n "EasyLPAC/internal/easylpac/i18n"
	"EasyLPAC/internal/easylpac/model"
	"EasyLPAC/internal/easylpac/resources"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	nativeDialog "github.com/sqweek/dialog"
	"golang.design/x/clipboard"
)

var WMain fyne.Window
var spacer *canvas.Rectangle

func buildLanguageSelect() *widget.Select {
	options, labelToPreference := appi18n.LanguagePreferenceOptions()
	selectedPreference := appi18n.CurrentLanguagePreference(App.Preferences())

	selectWidget := widget.NewSelect(options, nil)
	selectWidget.SetSelected(appi18n.LanguagePreferenceLabel(selectedPreference))
	selectWidget.OnChanged = func(label string) {
		preference, ok := labelToPreference[label]
		if !ok || preference == appi18n.CurrentLanguagePreference(App.Preferences()) {
			return
		}
		appi18n.SetLanguagePreference(App.Preferences(), preference)
		App.Settings().SetTheme(resources.NewTheme(appi18n.CurrentLanguageTag()))
		ReloadMainWindowContent()
	}
	return selectWidget
}

func localizeActivationError(err error) error {
	if errors.Is(err, activation.ErrInvalidCodeFormat) {
		return errors.New(appi18n.TR.Trans("message.qr_code_format_error"))
	}
	return err
}

func buildMainContent() fyne.CanvasObject {
	statusBar := container.NewGridWrap(fyne.Size{
		Width:  100,
		Height: DownloadButton.MinSize().Height,
	}, StatusLabel, StatusProcessBar)

	spacer = canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(1, 1))

	topToolBar := container.NewBorder(
		layout.NewSpacer(),
		nil,
		container.New(layout.NewHBoxLayout(), OpenLogButton, spacer, RefreshButton, spacer),
		FreeSpaceLabel,
		container.NewBorder(
			nil,
			nil,
			widget.NewLabel(appi18n.TR.Trans("label.card_reader")),
			nil,
			container.NewHBox(container.NewGridWrap(fyne.Size{
				Width:  280,
				Height: ApduDriverSelect.MinSize().Height,
			}, ApduDriverSelect), ApduDriverRefreshButton)),
	)

	profileTabContent := container.NewBorder(
		topToolBar,
		container.NewBorder(
			nil,
			nil,
			nil,
			container.NewHBox(ProfileMaskCheck, DownloadButton,
				// spacer, DiscoveryButton,
				spacer, SetNicknameButton,
				spacer, SwitchStateButton,
				spacer, DeleteProfileButton),
			statusBar),
		nil,
		nil,
		ProfileList)
	ProfileTab = container.NewTabItem(appi18n.TR.Trans("tab_bar.profile"), profileTabContent)

	notificationTabContent := container.NewBorder(
		topToolBar,
		container.NewBorder(
			nil,
			nil,
			nil,
			container.NewHBox(NotificationMaskCheck,
				spacer, ProcessNotificationButton,
				spacer, ProcessAllNotificationButton,
				spacer, BatchRemoveNotificationButton,
				spacer, RemoveNotificationButton),
			statusBar),
		nil,
		nil,
		NotificationList)
	NotificationTab = container.NewTabItem(appi18n.TR.Trans("tab_bar.notification"), notificationTabContent)

	chipInfoTabContent := container.NewBorder(
		topToolBar,
		container.NewBorder(
			nil,
			nil,
			nil,
			nil,
			statusBar),
		nil,
		nil,
		container.NewBorder(
			container.NewVBox(
				container.NewHBox(
					EidLabel, CopyEidButton, layout.NewSpacer(), EUICCManufacturerLabel),
				container.NewHBox(
					DefaultDpAddressLabel, SetDefaultSmdpButton, layout.NewSpacer(), ViewCertInfoButton),
				container.NewHBox(
					RootDsAddressLabel, layout.NewSpacer(), CopyEuiccInfo2Button)),
			nil,
			nil,
			nil,
			container.NewScroll(EuiccInfo2Entry),
		))
	ChipInfoTab = container.NewTabItem(appi18n.TR.Trans("tab_bar.chip_info"), chipInfoTabContent)

	aidEntryHint := &widget.Label{Text: appi18n.TR.Trans("label.aid_valid")}
	aidEntry := &widget.Entry{
		Text: appconfig.Instance.LpacAID,
		Validator: validation.NewAllStrings(
			validation.NewRegexp(`^.{32}$`, appi18n.TR.Trans("message.aid_length_illegal")),
			validation.NewRegexp(`[[:xdigit:]]{32}`, appi18n.TR.Trans("message.aid_not_hex")),
		),
	}
	aidEntry.OnChanged = func(s string) {
		val := aidEntry.Validate()
		if val != nil {
			aidEntryHint.SetText(val.Error())
		} else {
			// Use last known good value only
			appconfig.Instance.LpacAID = s
			aidEntryHint.SetText(appi18n.TR.Trans("label.aid_valid"))
		}
	}
	setToDefaultAidButton := widget.NewButton(
		appi18n.TR.Trans("label.aid_default_button"),
		func() {
			aidEntry.SetText(appconfig.AIDDefault)
		})
	setTo5berAidButton := widget.NewButton(
		appi18n.TR.Trans("label.aid_5ber_button"),
		func() {
			aidEntry.SetText(appconfig.AID5BER)
		})
	setToEsimmeAidButton := widget.NewButton(
		appi18n.TR.Trans("label.aid_esimme_button"),
		func() {
			aidEntry.SetText(appconfig.AIDESIMME)
		})
	setToXesimAidButton := widget.NewButton(
		appi18n.TR.Trans("label.aid_xesim_button"),
		func() {
			aidEntry.SetText(appconfig.AIDXESIM)
		})
	LanguageSelect = buildLanguageSelect()

	settingsTabContent := container.NewVBox(
		&widget.Label{Text: appi18n.TR.Trans("label.interface_language"), TextStyle: fyne.TextStyle{Bold: true}},
		container.NewGridWrap(
			fyne.Size{
				Width:  220,
				Height: LanguageSelect.MinSize().Height,
			},
			LanguageSelect),

		&widget.Label{Text: appi18n.TR.Trans("label.lpac_isdr_aid"), TextStyle: fyne.TextStyle{Bold: true}},
		container.NewHBox(container.NewGridWrap(
			fyne.Size{
				Width:  320,
				Height: aidEntry.MinSize().Height,
			}, aidEntry),
			setToDefaultAidButton,
			setTo5berAidButton,
			setToEsimmeAidButton,
			setToXesimAidButton),
		aidEntryHint,

		&widget.Label{Text: appi18n.TR.Trans("label.lpac_debug_output"), TextStyle: fyne.TextStyle{Bold: true}},
		&widget.Check{
			Text:    appi18n.TR.Trans("label.enable_env_LIBEUICC_DEBUG_HTTP_check"),
			Checked: appconfig.Instance.DebugHTTP,
			OnChanged: func(b bool) {
				appconfig.Instance.DebugHTTP = b
			},
		},
		&widget.Check{
			Text:    appi18n.TR.Trans("label.enable_env_LIBEUICC_DEBUG_APDU_check"),
			Checked: appconfig.Instance.DebugAPDU,
			OnChanged: func(b bool) {
				appconfig.Instance.DebugAPDU = b
			},
		},

		&widget.Label{Text: appi18n.TR.Trans("label.easylpac_settings"), TextStyle: fyne.TextStyle{Bold: true}},
		&widget.Check{
			Text:    appi18n.TR.Trans("label.auto_process_notification_check"),
			Checked: appconfig.Instance.AutoMode,
			OnChanged: func(b bool) {
				appconfig.Instance.AutoMode = b
			},
		})
	SettingsTab = container.NewTabItem(appi18n.TR.Trans("tab_bar.settings"), settingsTabContent)

	thankstoText := widget.NewRichTextFromMarkdown(appi18n.TR.Trans("thanks_to"))

	aboutText := widget.NewRichTextFromMarkdown(appi18n.TR.Trans("about"))

	aboutTabContent := container.NewBorder(
		nil,
		container.NewBorder(nil, nil,
			container.NewHBox(
				widget.NewLabel(fmt.Sprintf(appi18n.TR.Trans("label.version")+" %s", Version)),
				LpacVersionLabel),
			widget.NewLabel(fmt.Sprintf(appi18n.TR.Trans("label.euicc_data")+" %s", EUICCDataVersion))),
		nil,
		nil,
		container.NewCenter(container.NewVBox(thankstoText, aboutText)))
	AboutTab = container.NewTabItem(appi18n.TR.Trans("tab_bar.about"), aboutTabContent)

	Tabs = container.NewAppTabs(ProfileTab, NotificationTab, ChipInfoTab, SettingsTab, AboutTab)

	return Tabs
}

func ReloadMainWindowContent() {
	if WMain == nil {
		return
	}

	driverOptions := append([]string(nil), ApduDriverSelect.Options...)
	selectedDriver := ApduDriverSelect.Selected
	driverIFID := appconfig.Instance.DriverIFID
	refreshNeeded := RefreshNeeded

	SelectedProfile = Unselected
	SelectedNotification = Unselected
	ProfileStateAllowDisable = false

	InitWidgets()
	WMain.SetContent(buildMainContent())

	ApduDriverSelect.SetOptions(driverOptions)
	if selectedDriver != "" {
		ApduDriverSelect.SetSelected(selectedDriver)
	} else {
		ApduDriverSelect.ClearSelected()
	}
	appconfig.Instance.DriverIFID = driverIFID
	RefreshNeeded = refreshNeeded

	RefreshCurrentUIState()
	if Tabs != nil && len(Tabs.Items) > 3 {
		Tabs.SelectIndex(3)
	}
}

func InitMainWindow() fyne.Window {
	w := App.NewWindow("EasyLPAC")
	w.Resize(fyne.Size{
		Width:  850,
		Height: 545,
	})
	w.SetMaster()
	w.SetContent(buildMainContent())
	RefreshCurrentUIState()
	return w
}

func InitDownloadDialog() dialog.Dialog {
	smdpEntry := &widget.Entry{PlaceHolder: appi18n.TR.Trans("label.smdp_entry_placeholder")}
	matchIDEntry := &widget.Entry{PlaceHolder: appi18n.TR.Trans("label.match_id_entry_placeholder")}
	confirmCodeEntry := &widget.Entry{PlaceHolder: appi18n.TR.Trans("label.confirm_code_entry_placeholder")}
	imeiEntry := &widget.Entry{PlaceHolder: appi18n.TR.Trans("label.imei_entry_placeholder")}

	formItems := []*widget.FormItem{
		{Text: appi18n.TR.Trans("label.smdp"), Widget: smdpEntry},
		{Text: appi18n.TR.Trans("label.match_id"), Widget: matchIDEntry},
		{Text: appi18n.TR.Trans("label.confirm_code"), Widget: confirmCodeEntry},
		{Text: appi18n.TR.Trans("label.imei"), Widget: imeiEntry},
	}

	form := widget.NewForm(formItems...)
	var d dialog.Dialog
	showConfirmCodeNeededDialog := func() {
		dialog.ShowInformation(appi18n.TR.Trans("dialog.confirm_code_required"),
			appi18n.TR.Trans("message.confirm_code_required"), WMain)
	}
	cancelButton := &widget.Button{
		Text: appi18n.TR.Trans("dialog.cancel"),
		Icon: theme.CancelIcon(),
		OnTapped: func() {
			d.Hide()
		},
	}
	downloadButton := &widget.Button{
		Text:       appi18n.TR.Trans("label.download_profile_button"),
		Icon:       theme.ConfirmIcon(),
		Importance: widget.HighImportance,
		OnTapped: func() {
			d.Hide()
			pullConfig := model.PullInfo{
				SMDP:        strings.TrimSpace(smdpEntry.Text),
				MatchID:     strings.TrimSpace(matchIDEntry.Text),
				ConfirmCode: strings.TrimSpace(confirmCodeEntry.Text),
				IMEI:        strings.TrimSpace(imeiEntry.Text),
			}
			go func() {
				err := RefreshNotification()
				if err != nil {
					ShowLpacErrDialog(err)
					return
				}
				LpacProfileDownload(pullConfig)
			}()
		},
	}
	// 回调函数需要操作这两个 Button，预先声明
	var selectQRCodeButton *widget.Button
	var pasteFromClipboardButton *widget.Button
	disableButtons := func() {
		cancelButton.Disable()
		downloadButton.Disable()
		selectQRCodeButton.Disable()
		pasteFromClipboardButton.Disable()
	}
	enableButtons := func() {
		cancelButton.Enable()
		downloadButton.Enable()
		selectQRCodeButton.Enable()
		pasteFromClipboardButton.Enable()
	}

	selectQRCodeButton = &widget.Button{
		Text: appi18n.TR.Trans("label.select_qrcode_button"),
		Icon: theme.FileImageIcon(),
		OnTapped: func() {
			go func() {
				disableButtons()
				defer enableButtons()
				fileBuilder := nativeDialog.File().Title(appi18n.TR.Trans("dialog.select_qrcode"))
				fileBuilder.Filters = []nativeDialog.FileFilter{
					{
						Desc:       appi18n.TR.Trans("dialog.image_desc") + " (*.PNG, *.png, *.JPG, *.jpg, *.JPEG, *.jpeg)",
						Extensions: []string{"PNG", "png", "JPG", "jpg", "JPEG", "jpeg"},
					},
					{
						Desc:       appi18n.TR.Trans("dialog.all_files_desc") + " (*.*)",
						Extensions: []string{"*"},
					},
				}

				filename, err := fileBuilder.Load()
				if err != nil {
					if err.Error() != "Cancelled" {
						panic(err)
					}
				} else {
					code, err := activation.DecodeQRCodeFile(filename)
					if err != nil {
						dialog.ShowError(err, WMain)
					} else {
						pullInfo, confirmCodeNeeded, err2 := activation.DecodeCode(code)
						if err2 != nil {
							dialog.ShowError(localizeActivationError(err2), WMain)
						} else {
							smdpEntry.SetText(pullInfo.SMDP)
							matchIDEntry.SetText(pullInfo.MatchID)
							if confirmCodeNeeded {
								go showConfirmCodeNeededDialog()
							}
						}
					}
				}
			}()
		},
	}
	pasteFromClipboardButton = &widget.Button{
		Text: appi18n.TR.Trans("label.paste_from_clipboard_button"),
		Icon: theme.ContentPasteIcon(),
		OnTapped: func() {
			go func() {
				disableButtons()
				defer enableButtons()
				var err error
				var pullInfo model.PullInfo
				var confirmCodeNeeded bool

				format, result, err := activation.ReadClipboard()
				if err != nil {
					dialog.ShowError(err, WMain)
					return
				}
				switch format {
				case clipboard.FmtImage:
					code, decodeErr := activation.DecodeQRCodeBytes(result)
					if decodeErr != nil {
						dialog.ShowError(decodeErr, WMain)
						return
					}
					pullInfo, confirmCodeNeeded, err = activation.DecodeCode(code)
				case clipboard.FmtText:
					pullInfo, confirmCodeNeeded, err = activation.DecodeCode(activation.CompleteCode(string(result)))
				default:
					panic("unexpected clipboard format")
				}
				if err != nil {
					dialog.ShowError(localizeActivationError(err), WMain)
					return
				}
				smdpEntry.SetText(pullInfo.SMDP)
				matchIDEntry.SetText(pullInfo.MatchID)
				if confirmCodeNeeded {
					go showConfirmCodeNeededDialog()
				}
			}()
		},
	}
	d = dialog.NewCustomWithoutButtons(appi18n.TR.Trans("label.download_profile_button"), container.NewBorder(
		nil,
		container.NewVBox(spacer, container.NewCenter(selectQRCodeButton), spacer,
			container.NewCenter(pasteFromClipboardButton), spacer,
			container.NewCenter(container.NewHBox(cancelButton, spacer, downloadButton))),
		nil,
		nil,
		form), WMain)
	d.Resize(fyne.Size{
		Width:  520,
		Height: 380,
	})
	return d
}

func InitSetNicknameDialog() dialog.Dialog {
	entry := &widget.Entry{PlaceHolder: appi18n.TR.Trans("label.set_nickname_entry_placeholder")}
	form := []*widget.FormItem{
		{Text: appi18n.TR.Trans("label.set_nickname_button"), Widget: entry},
	}
	d := dialog.NewForm(appi18n.TR.Trans("label.set_nickname_form"), appi18n.TR.Trans("dialog.submit"), appi18n.TR.Trans("dialog.cancel"), form, func(b bool) {
		if b {
			if err := LpacProfileNickname(model.Profiles[SelectedProfile].Iccid, entry.Text); err != nil {
				ShowLpacErrDialog(err)
			}
			err := RefreshProfile()
			if err != nil {
				ShowLpacErrDialog(err)
			}
		}
	}, WMain)
	d.Resize(fyne.Size{
		Width:  400,
		Height: 180,
	})
	return d
}

func InitSetDefaultSmdpDialog() dialog.Dialog {
	entry := &widget.Entry{PlaceHolder: appi18n.TR.Trans("label.set_default_smdp_entry_placeholder")}
	form := []*widget.FormItem{
		{Text: appi18n.TR.Trans("label.default_smdp"), Widget: entry},
	}
	d := dialog.NewForm(appi18n.TR.Trans("label.set_default_smdp_form"), appi18n.TR.Trans("dialog.submit"), appi18n.TR.Trans("dialog.cancel"), form, func(b bool) {
		if b {
			if err := LpacChipDefaultSmdp(entry.Text); err != nil {
				ShowLpacErrDialog(err)
			}
			err := RefreshChipInfo()
			if err != nil {
				ShowLpacErrDialog(err)
			}
		}
	}, WMain)
	d.Resize(fyne.Size{
		Width:  510,
		Height: 200,
	})
	return d
}

func ShowLpacErrDialog(err error) {
	fyne.Do(func() {
		l := &widget.Label{Text: fmt.Sprintf("%v", err)}
		content := container.NewVBox(
			container.NewCenter(container.NewHBox(
				widget.NewIcon(theme.ErrorIcon()),
				widget.NewLabel(appi18n.TR.Trans("dialog.lpac_error")))),
			container.NewCenter(l),
			container.NewCenter(widget.NewLabel(appi18n.TR.Trans("message.lpac_error"))))
		dialog.ShowCustom(appi18n.TR.Trans("dialog.error"), appi18n.TR.Trans("dialog.ok"), content, WMain)
	})
}

func ShowSelectItemDialog() {
	fyne.Do(func() {
		d := dialog.NewInformation(appi18n.TR.Trans("dialog.info"), appi18n.TR.Trans("message.select_item"), WMain)
		d.Show()
	})
}

func ShowSelectCardReaderDialog() {
	fyne.Do(func() {
		dialog.ShowInformation(appi18n.TR.Trans("dialog.info"), appi18n.TR.Trans("message.select_card_reader"), WMain)
	})
}

func ShowRefreshNeededDialog() {
	fyne.Do(func() {
		dialog.ShowInformation(appi18n.TR.Trans("dialog.info"), appi18n.TR.Trans("message.refresh_required")+"\n", WMain)
	})
}
