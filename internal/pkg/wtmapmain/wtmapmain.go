package wtmapmain

import (
	"UrsusArctos/wtmap/pkg/wtapi"
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
)

const (
	defaultSleep               = 200
	sleepMultiplierWTWait      = 10
	sleepMultiplierSessionWait = 5

	trWaitingForWT      = "Waiting for WT client to start..."
	trWaitingForSession = "Waiting for game session to begin..."
	trInSession         = "In game session!"
)

type (
	TGenericHandler func() error
	TLayoutHandler  func(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	TDrawHandler    func(screen *ebiten.Image)

	TWTMapMain struct {
		// Handler forwarders
		HandlerUpdate TGenericHandler
		HandlerLayout TLayoutHandler
		HandlerDraw   TDrawHandler
		// UI Sync
		UIMutex sync.Mutex
		// UI Elements
		FontSB             font.Face   // font face for statusbar
		FontVG             font.Face   // font face for vehicle glyphs
		TextSB             string      // text for statusbar
		ColorSB            color.Color // background color for statusbar
		MapDownscaleFactor float64
		MapRenderOffset    image.Point
		ImageOptions       *ebiten.DrawImageOptions
		// Preloaded image for directional markers
		ImageDirected *ebiten.Image
		// API Client
		WTAPIC wtapi.TWTAPIClient
		// In-game data
		MapGen   int64
		MapFile  []byte
		MapImage *ebiten.Image
		MapObj   wtapi.TMapObjects
	}
)

var (
	colorWaitingForWT      = color.RGBA{0x88, 0x00, 0x00, 0xff}
	colorWaitingForSession = color.RGBA{0xb0, 0x50, 0x00, 0xff}
	colorInSession         = color.RGBA{0x20, 0x70, 0x00, 0xff}
)

func (wtmap *TWTMapMain) SetStatusBarPolitely(sbT string, sbC color.Color) {
	wtmap.UIMutex.Lock()
	defer wtmap.UIMutex.Unlock()
	wtmap.TextSB = sbT
	wtmap.ColorSB = sbC
}

func (wtmap *TWTMapMain) ClearMapImagePolitely() {
	wtmap.UIMutex.Lock()
	defer wtmap.UIMutex.Unlock()
	if wtmap.MapImage != nil {
		wtmap.MapImage.Clear()
	}
	wtmap.MapImage = nil
}

func (wtmap *TWTMapMain) WTMapMainWorker() {
	for {
		actualSleep := defaultSleep * time.Millisecond
		switch wtmap.WTAPIC.GetWTClientState() {
		case wtapi.WTClientStateNotStarted:
			{
				wtmap.ClearMapImagePolitely()
				wtmap.SetStatusBarPolitely(trWaitingForWT, colorWaitingForWT)
				wtmap.MapObj = nil
				actualSleep = defaultSleep * sleepMultiplierWTWait * time.Millisecond
			}
		case wtapi.WTClientStateIdle:
			{
				wtmap.ClearMapImagePolitely()
				wtmap.SetStatusBarPolitely(trWaitingForSession, colorWaitingForSession)
				wtmap.MapObj = nil
				actualSleep = defaultSleep * sleepMultiplierSessionWait * time.Millisecond
			}
		case wtapi.WTClientStateInSession:
			{
				// =================
				// 1. Map Generation
				if wtmap.WTAPIC.GetMapGeneration() != wtmap.MapGen {
					// 1.1. Map generation have changed, reload map
					tmpMapFile, err := wtmap.WTAPIC.GetMapFile()
					if err == nil {
						wtmap.MapFile = make([]byte, len(tmpMapFile))
						if copy(wtmap.MapFile, tmpMapFile) == len(tmpMapFile) {
							// 1.2. Try to create drawable image from it
							wtmap.ClearMapImagePolitely()
							// Create new map image
							var err2 error
							wtmap.UIMutex.Lock()
							wtmap.MapImage, _, err2 = ebitenutil.NewImageFromReader(bytes.NewReader(wtmap.MapFile))
							wtmap.UIMutex.Unlock()
							if err2 != nil {
								wtmap.ClearMapImagePolitely()
							}
						} // if copy is ok
					} // if getmap file is ok
					// update mapgen
					wtmap.MapGen = wtmap.WTAPIC.GetMapGeneration()
				} // if mapgen changed
				// =================
				// 2. Update in-game objects
				wtmap.UIMutex.Lock()
				wtmap.MapObj = wtmap.WTAPIC.GetMapObjects()
				wtmap.UIMutex.Unlock()
				// 3. Update status bar for in-game session message
				wtmap.SetStatusBarPolitely(trInSession, colorInSession)
			}
		}
		// Yield
		time.Sleep(actualSleep)
	}
}

func (wtmap *TWTMapMain) Draw(screen *ebiten.Image) {
	if wtmap.HandlerDraw != nil {
		wtmap.HandlerDraw(screen)
	}
}

func (wtmap *TWTMapMain) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if wtmap.HandlerLayout != nil {
		return wtmap.HandlerLayout(outsideWidth, outsideHeight)
	} else {
		return outsideWidth, outsideHeight
	}
}

func (wtmap *TWTMapMain) Update() error {
	if wtmap.HandlerUpdate != nil {
		return wtmap.HandlerUpdate()
	} else {
		return nil
	}
}
