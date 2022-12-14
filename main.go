package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	skyblue    = color.RGBA{27, 229, 235, 255}
	green      = color.RGBA{0, 255, 46, 255}
	brown      = color.RGBA{95, 85, 60, 255}
	mint       = color.RGBA{0, 255, 166, 255}

	player    = Player{Point{x: 3, y: 6}}
	direction = Point{}
	goal      = Point{x: 10, y: 11}
	block     = []Block{}

	board = [height][width]int{}
)

type Game struct {
	clear bool
}

type Point struct {
	x int
	y int
}

type Player struct {
	Point
}

type Block struct {
	Point
}

func (g *Game) Update() error {
	direction = Point{}
	if g.clear {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		direction.x = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		direction.x = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		direction.y = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		direction.y = 1
	}

	skip := len(blockFilter(block, func(b Block) bool {
		return b.x == player.x+direction.x && b.y == player.y
	})) != 0

	for idx := 0; idx < width; idx++ {
		if direction.x != 0 && !skip {
			tx := direction.x * (idx + 1)
			bf := blockFilter(block, func(b Block) bool {
				return (b.x-direction.x == player.x+tx && b.y == player.y)
			})

			if len(bf) > 0 || player.x+tx < 1 || player.x+tx > width-2 {
				player.x += direction.x * (idx + 1)
				break
			}
		}
	}

	skip = len(blockFilter(block, func(b Block) bool {
		return b.y == player.y+direction.y && b.x == player.x
	})) != 0

	for idx := 0; idx < height; idx++ {
		if direction.y != 0 && !skip {
			ty := direction.y * (idx + 1)
			bf := blockFilter(block, func(b Block) bool {
				return (b.y-direction.y == player.y+ty && b.x == player.x)
			})

			if len(bf) > 0 || player.y+ty < 1 || player.y+ty > height-2 {
				player.y += direction.y * (idx + 1)
				break
			}
		}
	}

	if player.x == goal.x && player.y == goal.y {
		g.clear = true
	}

	if player.x < 0 {
		player.x = 0
	}
	if player.x > width-1 {
		player.x = width - 1
	}
	if player.y < 0 {
		player.y = 0
	}
	if player.y > height-1 {
		player.y = height - 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	board = [height][width]int{}
	screen.Fill(background)

	for _, blo := range block {
		board[blo.y][blo.x] = 1
	}

	board[goal.y][goal.x] = 3

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

	image, _, err := ebitenutil.NewImageFromFile("./assets/game-clear.png")
	if err != nil {
		log.Print(err)
	} else if g.clear {
		sizeX := float64(image.Bounds().Size().X / 2)
		sizeY := float64(image.Bounds().Size().Y / 2)

		screenCenterX := float64(screenSizeX / 2)
		screenCenterY := float64(screenSizeY / 2)

		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenCenterX-sizeX, screenCenterY-sizeY)

		screen.DrawImage(image, &op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}

func getBlock(tile int) color.RGBA {
	switch tile {
	case 1:
		return brown
	case 2:
		return green
	case 3:
		return skyblue
	default:
		return blue
	}
}

func blockFilter(v []Block, f func(b Block) bool) []Block {
	var nb []Block
	for _, bv := range v {
		if f(bv) {
			nb = append(nb, bv)
		}
	}

	return nb
}

func init() {
	block = append(block,
		Block{Point{x: 9, y: 11}},
		Block{Point{x: 0, y: 1}})
}

func main() {
	ebiten.SetWindowTitle("Skate")
	ebiten.SetWindowSize(screenSizeX, screenSizeY)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
