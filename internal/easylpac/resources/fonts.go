package resources

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed fonts/DroidSansFallback.ttf
var droidSansFallback []byte

//go:embed fonts/DroidSansMono.ttf
var droidSansMono []byte

//go:embed fonts/DroidSansBold.ttf
var droidSansBold []byte

//go:embed fonts/NotoSansJP.ttf
var notoSansJP []byte

//go:embed fonts/NotoSansJP-Bold.ttf
var notoSansJPBold []byte

//go:embed fonts/NotoSansSC-Regular.otf
var notoSansSC []byte

//go:embed fonts/NotoSansSC-Bold.otf
var notoSansSCBold []byte

//go:embed fonts/NotoSansTC.ttf
var notoSansTC []byte

//go:embed fonts/NotoSansTC-Bold.ttf
var notoSansTCBold []byte

var DroidSansFallback = &fyne.StaticResource{
	StaticName:    "DroidSansFallback.ttf",
	StaticContent: droidSansFallback,
}

var DroidSansMono = &fyne.StaticResource{
	StaticName:    "DroidSansMono.ttf",
	StaticContent: droidSansMono,
}

var DroidSansBold = &fyne.StaticResource{
	StaticName:    "DroidSansBold.ttf",
	StaticContent: droidSansBold,
}

var NotoSansJP = &fyne.StaticResource{
	StaticName:    "NotoSansJP.ttf",
	StaticContent: notoSansJP,
}

var NotoSansJPBold = &fyne.StaticResource{
	StaticName:    "NotoSansJP-Bold.ttf",
	StaticContent: notoSansJPBold,
}

var NotoSansSC = &fyne.StaticResource{
	StaticName:    "NotoSansSC-Regular.otf",
	StaticContent: notoSansSC,
}

var NotoSansSCBold = &fyne.StaticResource{
	StaticName:    "NotoSansSC-Bold.otf",
	StaticContent: notoSansSCBold,
}

var NotoSansTC = &fyne.StaticResource{
	StaticName:    "NotoSansTC.ttf",
	StaticContent: notoSansTC,
}

var NotoSansTCBold = &fyne.StaticResource{
	StaticName:    "NotoSansTC-Bold.ttf",
	StaticContent: notoSansTCBold,
}
