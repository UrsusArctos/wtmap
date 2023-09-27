package main

import (
	"image/color"
	"time"
)

const (
	defaultSleep               = 500
	sleepMultiplierWTWait      = 4
	sleepMultiplierSessionWait = 2
)

func (wtmap *TWTMapMain) SetStatusBarPolite(sbT string, sbC color.Color) {
	wtmap.uiMutex.Lock()
	defer WTMap.uiMutex.Unlock()
	wtmap.sbText = sbT
	wtmap.sbColor = sbC
}

func (wtmap *TWTMapMain) WTMapMainWorker() {
	for {
		actualSleep := defaultSleep * time.Millisecond
		if wtmap.wtapic.IsWTRunning() {
			if wtmap.wtapic.IsInSession() {
				// do some actual api calls maybe?
				//
				wtmap.SetStatusBarPolite(trInSession, colorInSession)
			} else {
				wtmap.SetStatusBarPolite(trWaitingForSession, colorWaitingForSession)
				actualSleep = defaultSleep * sleepMultiplierSessionWait * time.Millisecond
			}
		} else {
			wtmap.SetStatusBarPolite(trWaitingForWT, colorWaitingForWT)
			actualSleep = defaultSleep * sleepMultiplierWTWait * time.Millisecond
		}
		// Yield
		time.Sleep(actualSleep)
	}
}
