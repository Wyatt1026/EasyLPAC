package i18n

import (
	"embed"
	"strings"

	"github.com/Xuanwo/go-locale"
	"github.com/fullpipe/icu-mf/mf"
	"golang.org/x/text/language"
)

//go:embed *.yaml
var localeFiles embed.FS

type PreferenceStore interface {
	String(key string) string
	SetString(key, value string)
}

var TR mf.Translator
var LanguageTag string

const (
	LanguagePreferenceKey = "language"
	LanguageSystemDefault = "system"
	LanguageEnglish       = "en"
	LanguageSimplifiedZH  = "zh-CN"
	LanguageTraditionalZH = "zh-TW"
	LanguageJapanese      = "ja-JP"
)

func normalizeLanguageTag(tag string) string {
	tag = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(tag)), "_", "-")
	switch {
	case strings.HasPrefix(tag, "ja"):
		return LanguageJapanese
	case strings.HasPrefix(tag, "zh-hant"),
		strings.Contains(tag, "-tw"),
		strings.Contains(tag, "-hk"),
		strings.Contains(tag, "-mo"):
		return LanguageTraditionalZH
	case strings.HasPrefix(tag, "zh"):
		return LanguageSimplifiedZH
	default:
		return LanguageEnglish
	}
}

func normalizeLanguagePreference(preference string) string {
	preference = strings.TrimSpace(strings.ToLower(preference))
	switch {
	case preference == "":
		return LanguageSystemDefault
	case preference == LanguageSystemDefault:
		return LanguageSystemDefault
	case strings.HasPrefix(preference, "en"):
		return LanguageEnglish
	case strings.HasPrefix(preference, "zh-cn"),
		strings.HasPrefix(preference, "zh-hans"),
		preference == "zh":
		return LanguageSimplifiedZH
	case strings.HasPrefix(preference, "zh-tw"),
		strings.HasPrefix(preference, "zh-hant"),
		strings.HasPrefix(preference, "zh-hk"),
		strings.HasPrefix(preference, "zh-mo"):
		return LanguageTraditionalZH
	case strings.HasPrefix(preference, "ja"):
		return LanguageJapanese
	default:
		return LanguageSystemDefault
	}
}

func detectSystemLanguage() string {
	tag, err := locale.Detect()
	if err != nil {
		return LanguageEnglish
	}
	return normalizeLanguageTag(tag.String())
}

func CurrentLanguagePreference(store PreferenceStore) string {
	if store == nil {
		return LanguageSystemDefault
	}
	return normalizeLanguagePreference(store.String(LanguagePreferenceKey))
}

func resolveLanguageTag(preference string) string {
	preference = normalizeLanguagePreference(preference)
	if preference == LanguageSystemDefault {
		return detectSystemLanguage()
	}
	return normalizeLanguageTag(preference)
}

func SetLanguagePreference(store PreferenceStore, preference string) {
	preference = normalizeLanguagePreference(preference)
	if store != nil {
		store.SetString(LanguagePreferenceKey, preference)
	}
	Init(store)
}

func LanguagePreferenceLabel(preference string) string {
	switch normalizeLanguagePreference(preference) {
	case LanguageEnglish:
		return TR.Trans("label.language_english")
	case LanguageSimplifiedZH:
		return TR.Trans("label.language_simplified_chinese")
	case LanguageTraditionalZH:
		return TR.Trans("label.language_traditional_chinese")
	case LanguageJapanese:
		return TR.Trans("label.language_japanese")
	default:
		return TR.Trans("label.language_system")
	}
}

func LanguagePreferenceOptions() ([]string, map[string]string) {
	preferences := []string{
		LanguageSystemDefault,
		LanguageEnglish,
		LanguageSimplifiedZH,
		LanguageTraditionalZH,
		LanguageJapanese,
	}
	options := make([]string, 0, len(preferences))
	labelToPreference := make(map[string]string, len(preferences))
	for _, preference := range preferences {
		label := LanguagePreferenceLabel(preference)
		options = append(options, label)
		labelToPreference[label] = preference
	}
	return options, labelToPreference
}

func CurrentLanguageTag() string {
	return LanguageTag
}

func Init(store PreferenceStore) {
	bundle, err := mf.NewBundle(
		mf.WithDefaultLangFallback(language.English),
		mf.WithYamlProvider(localeFiles),
	)
	if err != nil {
		panic(err)
	}
	LanguageTag = resolveLanguageTag(CurrentLanguagePreference(store))
	TR = bundle.Translator(LanguageTag)
}
