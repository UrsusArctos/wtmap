package wtmapstbar

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

const (
	textOffset      = 8
	DefaultTextSize = 18
)

func DrawStatusBar(canvas *ebiten.Image, fontFace font.Face, barColor color.Color, barText string) {
	b, _ := font.BoundString(fontFace, barText)
	barHeight := b.Max.Y.Ceil() - b.Min.Y.Ceil() + textOffset
	// draw background
	vector.DrawFilledRect(canvas,
		float32(canvas.Bounds().Min.X),
		float32(canvas.Bounds().Max.Y-barHeight),
		float32(canvas.Bounds().Dx()),
		float32(barHeight),
		barColor, false)
	// draw text
	text.Draw(canvas, barText, fontFace, canvas.Bounds().Min.X+textOffset/2, canvas.Bounds().Max.Y-textOffset, color.RGBA{0xff, 0xff, 0xff, 0xff})
}
