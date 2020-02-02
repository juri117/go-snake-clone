package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth      = 320 * 2
	screenHeight     = 320 * 1.5
	tileSize         = 32
	fontSize         = 32
	smallFontSize    = fontSize / 2
	pipeWidth        = tileSize * 2
	pipeStartOffsetX = 8
	pipeIntervalX    = 8
	pipeGapY         = 5
	frq              = 15
)

type BodyPart struct {
	posX int
	posY int
}

type Game struct {
	direction     int
	prevDirection int
	counter       int
	headI         int
	tail          []BodyPart
	food          BodyPart
	gameOver      int
	bodyImg       *ebiten.Image
	foodImg       *ebiten.Image
}

func NewGame() *Game {
	g := &Game{}
	g.bodyImg = LoadImg("assets/ball.png")
	g.foodImg = LoadImg("assets/food.png")
	g.init()
	return g
}

func LoadImg(fileName string) *ebiten.Image {
	reader, _ := ebitenutil.OpenFile(fileName)
	img, _, _ := image.Decode(reader)
	gameImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	return gameImg
}

func (g *Game) init() {
	var new_tail []BodyPart
	g.tail = new_tail
	g.tail = append(g.tail, BodyPart{posX: 5, posY: 5})
	g.tail = append(g.tail, BodyPart{posX: 4, posY: 5})
	g.tail = append(g.tail, BodyPart{posX: 3, posY: 5})
	g.headI = 0
	g.direction = 1
	g.prevDirection = g.direction
	g.gameOver = 0
	g.counter = 0
	g.place_food()
}

func (g *Game) place_food() {
	tryX := rand.Intn(int(screenWidth/tileSize) - 1)
	tryY := rand.Intn(int(screenHeight/tileSize) - 1)
	for g.get_crash_index(tryX, tryY) >= 0 {
		tryX = rand.Intn(int(screenWidth/tileSize) - 1)
		tryY = rand.Intn(int(screenHeight/tileSize) - 1)
	}
	g.food = BodyPart{posX: tryX, posY: tryY}
}

func (g *Game) get_crash_index(posX int, posY int) int {
	for i, part := range g.tail {
		if part.posX == posX && part.posY == posY {
			return i
		}
	}
	return -1
}

func (g *Game) append_body_part(new_part BodyPart) {
	var new_tail []BodyPart
	for i, part := range g.tail {
		if i == g.headI {
			new_tail = append(new_tail, new_part)
		}
		new_tail = append(new_tail, part)
	}
	g.tail = new_tail
	g.headI++
}

func (g *Game) Update(screen *ebiten.Image) error {
	// check for keyboad input
	var pressed []ebiten.Key
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			pressed = append(pressed, k)
		}
	}

	if g.gameOver > 0 {
		screen.Fill(color.NRGBA{0xFF, 0x00, 0x00, 0xff})
		g.counter = -99
	}

	// update direction
	if len(pressed) > 0 {
		switch pressed[0] {
		case ebiten.KeyUp:
			if g.prevDirection != 2 {
				g.direction = 0
			}
		case ebiten.KeyRight:
			if g.prevDirection != 3 {
				g.direction = 1
			}
		case ebiten.KeyDown:
			if g.prevDirection != 0 {
				g.direction = 2
			}
		case ebiten.KeyLeft:
			if g.prevDirection != 1 {
				g.direction = 3
			}
		case ebiten.KeyR:
			g.init()
		}
	}

	// update snake
	if g.counter >= frq {
		g.counter = 0

		lastI := g.headI - 1
		if lastI < 0 {
			lastI = len(g.tail) - 1
		}

		lastPart := g.tail[lastI]
		switch g.direction {
		case 0:
			g.tail[lastI].posX = g.tail[g.headI].posX
			g.tail[lastI].posY = g.tail[g.headI].posY - 1
			if g.tail[lastI].posY < 0 {
				g.tail[lastI].posY = int(screenHeight/tileSize) - 1
			}
		case 1:
			g.tail[lastI].posX = g.tail[g.headI].posX + 1
			g.tail[lastI].posY = g.tail[g.headI].posY
			if g.tail[lastI].posX > (screenWidth/tileSize)-1 {
				g.tail[lastI].posX = 0
			}
		case 2:
			g.tail[lastI].posX = g.tail[g.headI].posX
			g.tail[lastI].posY = g.tail[g.headI].posY + 1
			if g.tail[lastI].posY > (screenHeight/tileSize)-1 {
				g.tail[lastI].posY = 0
			}
		case 3:
			g.tail[lastI].posX = g.tail[g.headI].posX - 1
			g.tail[lastI].posY = g.tail[g.headI].posY
			if g.tail[lastI].posX < 0 {
				g.tail[lastI].posX = int(screenWidth/tileSize) - 1
			}
		}
		g.headI = lastI
		g.prevDirection = g.direction

		// ceck if food is hit
		if g.tail[g.headI].posX == g.food.posX && g.tail[g.headI].posY == g.food.posY {
			g.append_body_part(lastPart)
			g.place_food()
		}

		// check für canibalism
		for i, part := range g.tail {
			if i != g.headI {
				if g.tail[g.headI].posX == part.posX && g.tail[g.headI].posY == part.posY {
					g.gameOver = 1
				}
			}
		}
	}
	g.counter++

	for _, part := range g.tail {
		g.draw_body_part(screen, part.posX, part.posY)
	}
	g.draw_food(screen, g.food.posX, g.food.posY)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f, SCORE: %d", ebiten.CurrentTPS(), len(g.tail)))

	g.counter++
	return nil
}

func (g *Game) draw_body_part(screen *ebiten.Image, posX int, posY int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(float64(posX*tileSize), float64(posY*tileSize))
	screen.DrawImage(g.bodyImg, op)
}

func (g *Game) draw_food(screen *ebiten.Image, posX int, posY int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(float64(posX*tileSize), float64(posY*tileSize))
	screen.DrawImage(g.foodImg, op)
}

func main() {
	g := NewGame()
	if err := ebiten.Run(g.Update, screenWidth, screenHeight, 1, "Sanke"); err != nil {
		panic(err)
	}
}
