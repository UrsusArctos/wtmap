package main

import (
	"UrsusArctos/wtmap/internal/pkg/wtmapfonts"
	"UrsusArctos/wtmap/internal/pkg/wtmapstbar"
	"UrsusArctos/wtmap/pkg/wtapi"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	ProjectName     = "WT Map"
	defaultHostname = "localhost"
	mapDimension    = 2048
)

type (
	TWTMapMain struct {
		// UI Sync
		uiMutex sync.Mutex
		// UI Elements
		sbFont             font.Face
		sbText             string
		sbColor            color.Color
		mapDownscaleFactor float64
		mapRenderOffset    image.Point
		imageOptions       *ebiten.DrawImageOptions
		// API Client
		wtapic wtapi.TWTAPIClient
		// In-game data
		mapGen   int64
		mapFile  []byte
		mapImage *ebiten.Image
	}
)

var (
	WTMap TWTMapMain
)

func (wtmap *TWTMapMain) Draw(screen *ebiten.Image) {
	wtmap.uiMutex.Lock()
	defer wtmap.uiMutex.Unlock()
	// Update map image
	if wtmap.mapImage != nil {
		screen.DrawImage(wtmap.mapImage, wtmap.imageOptions)
	} else {
		screen.Clear()
	}
	// Update satus bar
	if (wtmap.sbColor != nil) && (len(wtmap.sbText) > 0) {
		wtmapstbar.DrawStatusBar(screen, wtmap.sbFont, wtmap.sbColor, wtmap.sbText)
	}
}

func (wtmap *TWTMapMain) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// calculate downscale factor
	wtmap.mapDownscaleFactor = math.Min(float64(outsideWidth), float64(outsideHeight)) / float64(mapDimension)
	// calculate render offsets
	wtmap.mapRenderOffset = image.Point{X: 0, Y: 0}
	if outsideHeight > outsideWidth {
		// Portrait orientation
		wtmap.mapRenderOffset.Y = (outsideHeight - int(mapDimension*wtmap.mapDownscaleFactor)) / 2
	} else {
		// Landscape orientation
		wtmap.mapRenderOffset.X = (outsideWidth - int(mapDimension*wtmap.mapDownscaleFactor)) / 2
	}
	// prepare draw image options
	wtmap.imageOptions = &ebiten.DrawImageOptions{}
	wtmap.imageOptions.GeoM.Scale(wtmap.mapDownscaleFactor, wtmap.mapDownscaleFactor)
	wtmap.imageOptions.GeoM.Translate(float64(wtmap.mapRenderOffset.X), float64(wtmap.mapRenderOffset.Y))
	// return
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
	WTMap.mapGen = 0
	// 3. Initialize resources
	WTMap.sbFont = wtmapfonts.GetFontFace(wtmapfonts.TTFont[wtmapfonts.FontSBTEXT], wtmapstbar.DefaultTextSize)
	// 4. Run
	go WTMap.WTMapMainWorker()
	ebiten.RunGame(&WTMap)
}
