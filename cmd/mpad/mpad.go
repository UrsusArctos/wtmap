package main

import (
	"UrsusArctos/wtmap/internal/pkg/wtmapfonts"
	"UrsusArctos/wtmap/internal/pkg/wtmapmain"
	"UrsusArctos/wtmap/internal/pkg/wtmapobj"
	"UrsusArctos/wtmap/internal/pkg/wtmapstbar"
	"UrsusArctos/wtmap/pkg/wtapi"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ProjectName     = "WT Map"
	defaultHostname = "localhost"
)

var (
	WTMap wtmapmain.TWTMapMain
)

func Update() error {
	return nil
}

func Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	var mapDimension float64 = 2048
	if WTMap.MapImage != nil {
		mapDimension = float64(WTMap.MapImage.Bounds().Dx())
	}
	// calculate downscale factor
	WTMap.MapDownscaleFactor = math.Min(float64(outsideWidth), float64(outsideHeight)) / mapDimension
	// calculate render offsets
	WTMap.MapRenderOffset = image.Point{X: 0, Y: 0}
	if outsideHeight > outsideWidth {
		// Portrait orientation
		WTMap.MapRenderOffset.Y = (outsideHeight - int(mapDimension*WTMap.MapDownscaleFactor)) / 2
	} else {
		// Landscape orientation
		WTMap.MapRenderOffset.X = (outsideWidth - int(mapDimension*WTMap.MapDownscaleFactor)) / 2
	}
	// prepare draw image options
	WTMap.ImageOptions = &ebiten.DrawImageOptions{}
	WTMap.ImageOptions.GeoM.Scale(WTMap.MapDownscaleFactor, WTMap.MapDownscaleFactor)
	WTMap.ImageOptions.GeoM.Translate(float64(WTMap.MapRenderOffset.X), float64(WTMap.MapRenderOffset.Y))
	// return
	return outsideWidth, outsideHeight
}

func Draw(screen *ebiten.Image) {
	WTMap.UIMutex.Lock()
	defer WTMap.UIMutex.Unlock()
	// Update map image
	if WTMap.MapImage != nil {
		screen.DrawImage(WTMap.MapImage, WTMap.ImageOptions)
		wtmapobj.DrawMapObjects(screen, &WTMap)
	} else {
		screen.Clear()
	}
	// Update satus bar
	if (WTMap.ColorSB != nil) && (len(WTMap.TextSB) > 0) {
		wtmapstbar.DrawStatusBar(screen, WTMap.FontSB, WTMap.ColorSB, WTMap.TextSB)
	}
}

func main() {
	// 1. Initialize Ebitengine
	ebiten.SetFullscreen(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle(ProjectName)
	// 2. Initialize WTAPI client
	WTMap.WTAPIC = wtapi.NewClient(defaultHostname)
	defer WTMap.WTAPIC.Close()
	WTMap.MapGen = 0
	// 3. Initialize resources
	WTMap.FontSB = wtmapfonts.GetFontFace(wtmapfonts.TTFont[wtmapfonts.FontSBTEXT], wtmapstbar.DefaultTextSize)
	WTMap.FontVG = wtmapfonts.GetFontFace(wtmapfonts.TTFont[wtmapfonts.FontICONS], wtmapobj.DefaultGlyphSize)
	WTMap.ImageDirected = ebiten.NewImage(wtmapobj.DefaultDirectedImageSize, wtmapobj.DefaultDirectedImageSize)
	// 4. Set handlers
	WTMap.HandlerUpdate = Update
	WTMap.HandlerLayout = Layout
	WTMap.HandlerDraw = Draw
	// 5. Run
	go WTMap.WTMapMainWorker()
	ebiten.RunGame(&WTMap)
}
