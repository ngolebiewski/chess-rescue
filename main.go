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

var _ embed.FS // keeps embed imported to fool the Go Linter in VS Code
// this is where the image gets loaded in, the embed is *supposed* to be commented. It's special and not really a comment.

//go:embed assets/images/key.png
var keyPNG []byte

const (
	screenWidth  = 480
	screenHeight = 320
	scale        = 2
	frameWidth   = 16
	frameHeight  = 16
	spacer       = 10
)

var (
	keyImage *ebiten.Image
)

type position struct {
	pX float64
	pY float64
}

type Game struct {
	keys         []ebiten.Key
	count        int
	direction    bool
	playerY      float64
	playerX      float64
	playerTrails bool
	positions    []position
}

func (g *Game) AddPosition(x, y float64) {
	g.positions = append(g.positions, position{pX: x, pY: y})

	if len(g.positions) > spacer*5 {
		g.positions = g.positions[len(g.positions)-spacer*5:]
	}

}

// Essentially, and crudely, checks if the player is hitting the top or bottom of the screen, and then bounces them back a bit. Would be better to reverse their direction.
func edgeCheck(x, y int, g *Game) bool {
	if x == 0 || x == screenWidth*scale {
		fmt.Println("On x edge")
		return true
	}
	if y == 1 || y == screenHeight-1 {
		fmt.Println("On y edge")
		if y < frameHeight {
			g.playerY += frameHeight
		}
		if y > screenHeight-frameHeight {
			g.playerY -= frameHeight
		}
		return true
	}
	return false
}

func (g *Game) Update() error {
	g.count++
	if g.count%screenWidth*scale == 0 {
		g.direction = !g.direction
	}
	if g.direction {
		g.playerX--
	}
	if !g.direction {
		g.playerX++
	}

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && !edgeCheck(g.count, int(g.playerY), g) {
		g.playerY--
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && !edgeCheck(g.count, int(g.playerY), g) {
		g.playerY++
	}

	g.AddPosition(g.playerX, g.playerY)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	msg := fmt.Sprintf(`Chess Rescue!
	player y: %v
	player x: %v
	count: %v
	TPS: %v`,
		g.playerY, g.playerX, g.count, ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)

	//Draw Trails
	if g.playerTrails && len(g.positions) >= spacer*4 {
		alpha := .5
		for i := range 4 {
			pos := g.positions[len(g.positions)-((i+1)*spacer)]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(pos.pX, pos.pY)
			op.ColorScale.ScaleAlpha(float32(alpha))
			screen.DrawImage(keyImage, op)
			alpha -= .1
		}
	}

		//Draw Main Image
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(g.playerX), g.playerY)
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
		count:        0,
		direction:    false,
		playerY:      screenHeight / 2,
		playerTrails: true,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
