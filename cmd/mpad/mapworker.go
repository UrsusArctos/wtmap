package main

import (
	"bytes"
	"image/color"
	"time"

	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	defaultSleep               = 500
	sleepMultiplierWTWait      = 4
	sleepMultiplierSessionWait = 2
)

func (wtmap *TWTMapMain) SetStatusBarPolitely(sbT string, sbC color.Color) {
	wtmap.uiMutex.Lock()
	defer WTMap.uiMutex.Unlock()
	wtmap.sbText = sbT
	wtmap.sbColor = sbC
}

func (wtmap *TWTMapMain) ClearMapImagePolitely() {
	wtmap.uiMutex.Lock()
	defer WTMap.uiMutex.Unlock()
	if wtmap.mapImage != nil {
		wtmap.mapImage.Clear()
	}
	wtmap.mapImage = nil
}

func (wtmap *TWTMapMain) WTMapMainWorker() {
	for {
		actualSleep := defaultSleep * time.Millisecond
		if wtmap.wtapic.IsWTRunning() {
			if wtmap.wtapic.IsInSession() {
				// 1. Watch map generation as map_info is now updated
				if wtmap.wtapic.GetMapGeneration() != wtmap.mapGen {
					// 1.1. Map generation have changed, reload map
					tmpMapFile, err := wtmap.wtapic.GetMapFile()
					if err == nil {
						wtmap.mapFile = make([]byte, len(tmpMapFile))
						if copy(wtmap.mapFile, tmpMapFile) == len(tmpMapFile) {
							// 1.2. Try to create drawable image from it
							wtmap.ClearMapImagePolitely()
							// Create new map image
							var err2 error
							wtmap.uiMutex.Lock()
							wtmap.mapImage, _, err2 = ebitenutil.NewImageFromReader(bytes.NewReader(wtmap.mapFile))
							WTMap.uiMutex.Unlock()
							if err2 != nil {
								wtmap.ClearMapImagePolitely()
							}
						} // if copy is ok
					} // if getmap file is ok
					// update mapgen
					wtmap.mapGen = wtmap.wtapic.GetMapGeneration()
				} // if mapgen changed
				// 2. Update status bar for in-game session message
				wtmap.SetStatusBarPolitely(trInSession, colorInSession)
			} else {
				wtmap.ClearMapImagePolitely()
				wtmap.SetStatusBarPolitely(trWaitingForSession, colorWaitingForSession)
				actualSleep = defaultSleep * sleepMultiplierSessionWait * time.Millisecond
			}
		} else {
			wtmap.ClearMapImagePolitely()
			wtmap.SetStatusBarPolitely(trWaitingForWT, colorWaitingForWT)
			actualSleep = defaultSleep * sleepMultiplierWTWait * time.Millisecond
		}
		// Yield
		time.Sleep(actualSleep)
	}
}
