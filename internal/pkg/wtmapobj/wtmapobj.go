package wtmapobj

import (
	"UrsusArctos/wtmap/internal/pkg/wtmapmain"
	"UrsusArctos/wtmap/pkg/wtapi"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

const (
	typeAirfield      = "airfield"
	typeAirRespawn    = "respawn_base_fighter"
	typeGroundRespawn = "respawn_base_tank"
	typeCaptureZone   = "capture_zone"
	typeGroundModel   = "ground_model"
	typeAircraft      = "aircraft"

	// ground vehicles
	iconPlayer        = "Player"
	iconGroundRespawn = "respawn_base_tank"
	iconCaptureZone   = "capture_zone"
	iconStructure     = "Structure"
	iconTankDestroyer = "TankDestroyer"
	iconSPAA          = "SPAA"
	iconAirdefence    = "Airdefence"
	iconWheeled       = "Wheeled"
	iconLightTank     = "LightTank"
	iconMediumTank    = "MediumTank"
	iconHeavyTank     = "HeavyTank"
	iconTorpedoBoat   = "TorpedoBoat"
	iconShip          = "Ship"
	// aircraft
	iconAirRespawn = "respawn_base_fighter"
	iconFighter    = "Fighter"
	iconBomber     = "Bomber"

	playerCircleRadius       = 6
	headingMarkerLength      = 16
	defaultStrokeWidth       = 2
	airfieldStrokeWidth      = 5
	DefaultDirectedImageSize = 32

	// vehicle glyphs
	DefaultGlyphSize = 32
	vgTankDestroyer  = "b"
	vgSPAA           = "4"
	vgWheeled        = "5"
	vgLightTank      = "t"
	vgMediumTank     = "g"
	vgHeavyTank      = "0"
	vgTorpedoBoat    = "s"
	vgShip           = "w"
	vgAirRespawn     = "6"
	vgCaptureZone    = "7"
	vgStructure      = "l"
	vgFighter        = "."
	vgBomber         = ":"
)

func calcHeadingMarker(obj wtapi.TMapObject) ( /*deltaX, deltaY*/ float64, float64) {
	unitLen := math.Sqrt((*obj.DX)*(*obj.DX) + (*obj.DY)*(*obj.DY))
	return float64(headingMarkerLength) * (*obj.DX) / unitLen, float64(headingMarkerLength) * (*obj.DY) / unitLen
}

func remapCoordinates(x float64, y float64, wtmain *wtmapmain.TWTMapMain) ( /*canvasX, canvasY*/ float64, float64) {
	mapDimensionReal := float64(wtmain.MapImage.Bounds().Dx()) * wtmain.MapDownscaleFactor
	return float64(wtmain.MapRenderOffset.X) + mapDimensionReal*x, float64(wtmain.MapRenderOffset.Y) + mapDimensionReal*y
}

func areCoordinatesValid(x, y float64) bool {
	return ((x > 0) && (x < 1) && (y > 0) && (y < 1))
}

func drawLinearAirfield(canvas *ebiten.Image, wtmain *wtmapmain.TWTMapMain, obj wtapi.TMapObject) {
	sx, sy := remapCoordinates(*obj.SX, *obj.SY, wtmain)
	ex, ey := remapCoordinates(*obj.EX, *obj.EY, wtmain)
	if areCoordinatesValid(*obj.SX, *obj.SY) && areCoordinatesValid(*obj.EX, *obj.EY) {
		vector.StrokeLine(canvas, float32(sx), float32(sy), float32(ex), float32(ey), airfieldStrokeWidth, obj.AdaptColor(), true)
	}
}

func drawVehicleGlyph(canvas *ebiten.Image, wtmain *wtmapmain.TWTMapMain, obj wtapi.TMapObject, glyphChar string, ofX int, ofY int, encircled bool) {
	b, _ := font.BoundString(wtmain.FontVG, glyphChar)
	glyphHeight, glyphWidth := b.Max.Y.Ceil()-b.Min.Y.Ceil(), b.Max.X.Ceil()-b.Min.X.Ceil()
	if areCoordinatesValid(*obj.X, *obj.Y) {
		cx, cy := remapCoordinates(*obj.X, *obj.Y, wtmain)
		// This circle is used to adjust new glyph placement, see offsets in drawVehicleGlyph()
		if encircled {
			vector.StrokeCircle(canvas, float32(cx), float32(cy), 16, defaultStrokeWidth, obj.AdaptColor(), true)
		}
		// Draw actual glyph
		text.Draw(canvas, glyphChar, wtmain.FontVG, int(cx)-glyphWidth+ofX, int(cy)+glyphHeight+ofY, obj.AdaptColor())
	}
}

func drawVehicleGlyphRotated(canvas *ebiten.Image, wtmain *wtmapmain.TWTMapMain, obj wtapi.TMapObject, glyphChar string, ofX int, ofY int, encircled bool) {
	if areCoordinatesValid(*obj.X, *obj.Y) {
		b, _ := font.BoundString(wtmain.FontVG, glyphChar)
		glyphHeight, glyphWidth := b.Max.Y.Ceil()-b.Min.Y.Ceil(), b.Max.X.Ceil()-b.Min.X.Ceil()
		wtmain.ImageDirected.Clear()
		if encircled {
			vector.StrokeCircle(wtmain.ImageDirected, DefaultDirectedImageSize/2, DefaultDirectedImageSize/2, DefaultDirectedImageSize/2, defaultStrokeWidth, obj.AdaptColor(), true)
		}
		text.Draw(wtmain.ImageDirected, glyphChar, wtmain.FontVG, (DefaultDirectedImageSize/2)-glyphWidth+9, (DefaultDirectedImageSize/2)+glyphHeight-2, obj.AdaptColor())
		io := &ebiten.DrawImageOptions{}
		cx, cy := remapCoordinates(*obj.X, *obj.Y, wtmain)
		io.GeoM.Reset()
		io.GeoM.Translate(-(DefaultDirectedImageSize / 2), -(DefaultDirectedImageSize / 2))
		// PI/2 is added because screen rotational axis coincide with in-game Y axis
		io.GeoM.Rotate((math.Pi / 2) + math.Atan2(*obj.DY, *obj.DX))
		io.GeoM.Translate(cx-float64(glyphWidth)+float64(ofX), cy+float64(glyphHeight)+float64(ofY))
		canvas.DrawImage(wtmain.ImageDirected, io)
	}
}

func drawAirRespawnDirected(canvas *ebiten.Image, wtmain *wtmapmain.TWTMapMain, obj wtapi.TMapObject, encircled bool) {
	if areCoordinatesValid(*obj.X, *obj.Y) {
		drawVehicleGlyph(canvas, wtmain, obj, vgAirRespawn, 5, -5, encircled)
		cx, cy := remapCoordinates(*obj.X, *obj.Y, wtmain)
		dx, dy := calcHeadingMarker(obj) // remember, these are not coordinates
		vector.StrokeLine(canvas, float32(cx), float32(cy), float32(cx+dx), float32(cy+dy), defaultStrokeWidth, obj.AdaptColor(), true)
	}
}

func DrawMapObjects(canvas *ebiten.Image, wtmain *wtmapmain.TWTMapMain) {
	if len(wtmain.MapObj) > 0 {
		for i := range wtmain.MapObj {
			// 1. Draw player regardless of vehicle type
			if wtmain.MapObj[i].Icon == iconPlayer {
				// Player: circle and heading marker
				if areCoordinatesValid(*wtmain.MapObj[i].X, *wtmain.MapObj[i].Y) {
					cx, cy := remapCoordinates(*wtmain.MapObj[i].X, *wtmain.MapObj[i].Y, wtmain)
					dx, dy := calcHeadingMarker(wtmain.MapObj[i]) // remember, these are not coordinates
					vector.StrokeCircle(canvas, float32(cx), float32(cy), playerCircleRadius, defaultStrokeWidth, wtmain.MapObj[i].AdaptColor(), true)
					vector.StrokeLine(canvas, float32(cx), float32(cy), float32(cx+dx), float32(cy+dy), defaultStrokeWidth, wtmain.MapObj[i].AdaptColor(), true)
				}
			}
			// 2. Draw other vehicles
			switch wtmain.MapObj[i].Type {
			case typeAirfield:
				{
					drawLinearAirfield(canvas, wtmain, wtmain.MapObj[i])
				}
			case typeAirRespawn:
				{
					drawAirRespawnDirected(canvas, wtmain, wtmain.MapObj[i], false)
				}
			case typeGroundRespawn:
				{
					// There is no point in drawing this, as they are too many
				}
			case typeCaptureZone:
				{
					drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgCaptureZone, 8, -8, false)
				}
			case typeAircraft:
				switch wtmain.MapObj[i].Icon {
				case iconFighter:
					{
						drawVehicleGlyphRotated(canvas, wtmain, wtmain.MapObj[i], vgFighter, 9, -1, false)
					}
				case iconBomber:
					{
						drawVehicleGlyphRotated(canvas, wtmain, wtmain.MapObj[i], vgBomber, 9, -1, false)
					}
				case iconPlayer:
					{
					}
				default:
					{
						// fmt.Printf("Unknown aircraft icon '%s'\n", wtmain.MapObj[i].Icon)
					}
				}
			case typeGroundModel:
				{ // Ground vehicles
					switch wtmain.MapObj[i].Icon {
					case iconTankDestroyer:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgTankDestroyer, 2, 0, false)
						}
					case iconSPAA, iconAirdefence:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgSPAA, 3, -1, false)
						}
					case iconWheeled:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgWheeled, -2, 4, false)
						}
					case iconLightTank:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgLightTank, 4, 8, false)
						}
					case iconMediumTank:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgMediumTank, 2, 0, false)
						}
					case iconHeavyTank:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgHeavyTank, 9, 2, false)
						}
					case iconTorpedoBoat:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgTorpedoBoat, 3, 0, false)
						}
					case iconShip:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgShip, 1, 2, false)
						}
					case iconStructure:
						{
							drawVehicleGlyph(canvas, wtmain, wtmain.MapObj[i], vgStructure, 0, 0, false)
						}
					case iconPlayer:
						{
						}
					default:
						{
							// fmt.Printf("Unknown ground model icon '%s'\n", wtmain.MapObj[i].Icon)
						}
					}
				}
			default:
				{
					// fmt.Printf("Unknown object type '%s'\n", wtmain.MapObj[i].Type)
				}
			}
		}
	}
}
