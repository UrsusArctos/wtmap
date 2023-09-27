package main

import "image/color"

const (
	trWaitingForWT      = "Waiting for WT client to start..."
	trWaitingForSession = "Waiting for game session to begin..."
	trInSession         = "In game session!"
)

var (
	colorWaitingForWT      = color.RGBA{0x88, 0x00, 0x00, 0xff}
	colorWaitingForSession = color.RGBA{0xb0, 0x50, 0x00, 0xff}
	colorInSession         = color.RGBA{0x20, 0x70, 0x00, 0xff}
)
