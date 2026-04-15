package resources

import (
	"image/color"

	appi18n "EasyLPAC/internal/easylpac/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type AppTheme struct {
	languageTag string
}

func NewTheme(languageTag string) fyne.Theme {
	return AppTheme{languageTag: languageTag}
}

func (t AppTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0xe6, G: 0x77, B: 0x2e, A: 0xff}
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 0xe6, G: 0x77, B: 0x2e, A: 0xff}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xf5, G: 0x65, B: 0x08, A: 0x2a}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0xf5, G: 0x65, B: 0x08, A: 0x2a}
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}

func (t AppTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Italic || s.Symbol {
		return theme.DefaultTheme().Font(s)
	}
	if s.Monospace {
		return DroidSansMono
	}
	if s.Bold {
		switch t.languageTag {
		case appi18n.LanguageJapanese:
			return NotoSansJPBold
		case appi18n.LanguageSimplifiedZH:
			return NotoSansSCBold
		case appi18n.LanguageTraditionalZH:
			return NotoSansTCBold
		default:
			return DroidSansBold
		}
	}
	switch t.languageTag {
	case appi18n.LanguageJapanese:
		return NotoSansJP
	case appi18n.LanguageSimplifiedZH:
		return NotoSansSC
	case appi18n.LanguageTraditionalZH:
		return NotoSansTC
	default:
		return DroidSansFallback
	}
}

func (t AppTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t AppTheme) Size(s fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(s)
}
