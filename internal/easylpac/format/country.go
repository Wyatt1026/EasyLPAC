package format

import "strings"

func CountryCodeToEmoji(countryCode string) string {
	if len(countryCode) != 2 {
		return "🌎"
	}
	countryCode = strings.ToUpper(countryCode)
	rune1 := rune(countryCode[0]-'A') + 0x1F1E6
	rune2 := rune(countryCode[1]-'A') + 0x1F1E6
	return string([]rune{rune1, rune2})
}
