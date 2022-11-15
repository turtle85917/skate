package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenSizeX = tileSize * (width + 2)
	screenSizeY = tileSize * (height + 2)
	width       = 15
	height      = 12
	tileSize    = 50
)

var (
	background = color.RGBA{120, 120, 120, 255}
	blue       = color.RGBA{0, 112, 255, 255}
	green      = color.RGBA{0, 255, 46, 255}
	brown      = color.RGBA{95, 85, 60, 255}
	mint       = color.RGBA{0, 255, 166, 255}

	player = Player{x: 3, y: 6}

	board = [height][width]int{}
)

type Game struct{}

type Player struct {
	x int
	y int
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(background)

	for y := -1; y < height+1; y++ {
		for x := -1; x < width+1; x++ {
			fx := float64(tileSize * (x + 1))
			fy := float64(tileSize * (y + 1))

			if x == -1 || y == -1 || x == width || y == height {
				ebitenutil.DrawRect(screen, fx, fy, tileSize, tileSize, brown)
			} else if x == player.x && y == player.y {
				ebitenutil.DrawRect(screen, fx, fy, tileSize, tileSize, mint)
			} else {
				ebitenutil.DrawRect(screen, fx, fy, tileSize, tileSize, getBlock(board[y][x]))
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}

func getBlock(tile int) color.RGBA {
	switch tile {
	case 1:
		return blue
	case 2:
		return green
	default:
		return blue
	}
}

func main() {
	ebiten.SetWindowTitle("Skate")
	ebiten.SetWindowSize(screenSizeX, screenSizeY)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
