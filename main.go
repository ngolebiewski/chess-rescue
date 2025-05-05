package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var _ embed.FS // keeps embed imported even if not directly used
//go:embed assets/images/key.png
var keyPNG []byte

const (
	screenWidth  = 256
	screenHeight = 256
	scale        = 3
	frameWidth   = 16
	frameHeight  = 16
)

var (
	keyImage  *ebiten.Image
	// count     int = 0
	// direction bool
	// playerY   float64 = 50.0
)

type Game struct {
	keys []ebiten.Key
	count     int
	direction bool
	playerY   float64
}

func (g *Game) Update() error {
	if g.count%screenWidth*scale == 0 {
		g.direction = !g.direction
	}
	if g.direction {
		g.count++
	}
	if !g.direction {
		g.count--
	}
	// if count > screenWidth * scale {
	// }

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.playerY--
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.playerY++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Chess Rescue!")

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	// op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(float64(g.count), g.playerY)
	screen.DrawImage(keyImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Decode an image from the image file's byte slice.

	img, _, err := image.Decode(bytes.NewReader(keyPNG))
	if err != nil {
		log.Fatal(err)
	}
	keyImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth*scale, screenHeight*scale)
	ebiten.SetWindowTitle("Chess Rescue")

	game := &Game{
		count:     0,
		direction: false,
		playerY:   50.0,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
