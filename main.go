package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var _ embed.FS // keeps embed imported to fool the Go Linter
// this is where the image gets loaded in, the embed is *supposed* to be commented. It's special and not really a comment.

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
	keyImage *ebiten.Image
)

type Game struct {
	keys 			[]ebiten.Key
	count     int
	direction bool
	playerY   float64
}

func edgeCheck(x, y int, g *Game) bool{
	if x == 0 || x == screenWidth*scale {
		fmt.Println("On x edge")
		return true
	}
	if y == 1 || y == screenHeight - 1 {
		fmt.Println("On y edge")
		if y < frameHeight {g.playerY += frameHeight}
		if y > screenHeight-frameHeight{g.playerY -= frameHeight}
		return true
	}
	return false
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

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && !edgeCheck(g.count, int(g.playerY), g){
		g.playerY--
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && !edgeCheck(g.count, int(g.playerY), g){
		g.playerY++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	msg := fmt.Sprintf(`Chess Rescue!
	player y: %v
	count/x: %v`, g.playerY, g.count)
	ebitenutil.DebugPrint(screen, msg)

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
		playerY:   screenHeight/2,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
