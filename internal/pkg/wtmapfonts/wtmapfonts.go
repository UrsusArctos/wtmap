package wtmapfonts

import (
	"UrsusArctos/wtmap/fonts"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	defaultDPI = 72

	FontICONS  = 0
	FontSBTEXT = 1
)

var (
	TTFontName = [...]string{"icons.ttf", "nk57.ttf"}
	TTFont     []*opentype.Font
)

func loadFont(fontName string) (otf *opentype.Font) {
	fontTTF, err := fonts.EmbedFonts.ReadFile(fontName)
	if (len(fontTTF) > 0) && (err == nil) {
		otf, err = opentype.Parse(fontTTF)
		if err == nil {
			return otf
		}
	}
	return nil
}

func GetFontFace(fontTT *opentype.Font, fontSize float64) (fontFace font.Face) {
	if fontTT != nil {
		var err error
		fontFace, err = opentype.NewFace(fontTT, &opentype.FaceOptions{Size: fontSize, DPI: defaultDPI, Hinting: font.HintingNone})
		if err == nil {
			return fontFace
		}
	}
	return nil
}

func init() {
	TTFont = make([]*opentype.Font, len(TTFontName))
	for i := range TTFontName {
		TTFont[i] = loadFont(TTFontName[i])
	}
}
