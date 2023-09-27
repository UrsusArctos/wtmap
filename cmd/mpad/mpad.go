package main

import (
	"UrsusArctos/wtmap/internal/pkg/wtmapfonts"
	"UrsusArctos/wtmap/internal/pkg/wtmapstbar"
	"UrsusArctos/wtmap/pkg/wtapi"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	ProjectName     = "WT Map"
	defaultHostname = "localhost"
)

type (
	TWTMapMain struct {
		// UI Sync
		uiMutex sync.Mutex
		// UI Elements
		sbFont  font.Face
		sbText  string
		sbColor color.Color
		// API Client
		wtapic wtapi.TWTAPIClient
	}
)

var (
	WTMap TWTMapMain
)

func (wtmap *TWTMapMain) Draw(screen *ebiten.Image) {
	wtmap.uiMutex.Lock()
	defer wtmap.uiMutex.Unlock()
	// screen.Clear()
	// Update satus bar
	if (wtmap.sbColor != nil) && (len(wtmap.sbText) > 0) {
		wtmapstbar.DrawStatusBar(screen, wtmap.sbFont, wtmap.sbColor, wtmap.sbText)
	}
}

func (wtmap *TWTMapMain) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (wtmap *TWTMapMain) Update() error {
	return nil
}

func main() {
	// 1. Initialize Ebitengine
	ebiten.SetFullscreen(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle(ProjectName)
	// 2. Initialize WTAPI client
	WTMap.wtapic = wtapi.NewClient(defaultHostname)
	defer WTMap.wtapic.Close()
	// 3. Initialize resources
	WTMap.sbFont = wtmapfonts.GetFontFace(wtmapfonts.TTFont[wtmapfonts.FontSBTEXT], wtmapstbar.DefaultTextSize)
	// 4. Run
	go WTMap.WTMapMainWorker()
	ebiten.RunGame(&WTMap)
}
