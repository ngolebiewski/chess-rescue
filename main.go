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
    screenWidth     = 256
    screenHeight    = 256
    scale           = 3   
    frameWidth      = 16
    frameHeight     = 16
)

var keyImage *ebiten.Image
var count int = 0
var direction bool
var playerY float64 = 50.0

type Game struct{
	keys []ebiten.Key
}

func (g *Game) Update() error {
	if count % screenWidth * scale == 0 {
		direction = !direction
	}
	if direction{
		count++}
	if !direction{
		count--}
	// if count > screenWidth * scale {
	// }
	
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		playerY--
		log.Println("Arrow Up pressed - do something")
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		playerY++
		log.Println("Arrow Down pressed - do something else")
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Chess Rescue!")
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	// op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(float64(count), playerY)
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
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

