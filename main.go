package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var inputLimit = 9
var startTime = time.Now()
var lastEnemySpawn = time.Now()

const (
	screenWidth  = 640
	screenHeight = 480
	centerX      = screenWidth / 2
	centerY      = screenHeight / 2
)

type Game struct {
	count       int
	userInput   UserInput
	mouseX      int
	mouseY      int
	projectiles []Projectile
	player      Player
	enemies     []Enemy
}

type Player struct {
	xPosition   float64
	yPosition   float64
	pSquare     *ebiten.Image
	pSquareHalo *ebiten.Image
}

type Projectile struct {
	xPosition float64
	yPosition float64
	xVelocity float64
	yVelocity float64
	pSquare   *ebiten.Image
}

func (p *Projectile) isOffScreen() bool {
	return p.xPosition < 0 || p.xPosition > screenWidth || p.yPosition < 0 || p.yPosition > screenHeight
}

func getVelocity(mouseX int, mouseY int) (xVelocity float64, yVelocity float64) {
	denominator := math.Abs(float64(mouseX-centerX)) + math.Abs(float64(mouseY-centerY))/2
	return (float64(mouseX-centerX) / denominator), (float64(mouseY-centerY) / denominator)
}

func (g *Game) spawnProjectile() {
	xVelocity, yVelocity := getVelocity(g.mouseX, g.mouseY)
	i := ebiten.NewImage(2, 2)
	i.Fill(color.White)
	p := Projectile{centerX, centerY, xVelocity, yVelocity, i}
	g.projectiles = append(g.projectiles, p)
}

type UserInput struct {
	currentInput []rune
}

func (u *UserInput) append() {
	if len(u.currentInput) < inputLimit {
		u.currentInput = ebiten.AppendInputChars(u.currentInput)
	}
}

func (u *UserInput) delete() {
	if len(u.currentInput) > 0 {
		u.currentInput = u.currentInput[:len(u.currentInput)-1]
	}
}

func (g *Game) Update() error {
	g.count++
	g.userInput.append()
	if time.Since(lastEnemySpawn) > 5*time.Second {
		lastEnemySpawn = time.Now()
		//random border location
		xLocation := 0
		yLocation := 0
		isOnYAxis := rand.Intn(2) == 1
		if isOnYAxis {
			isXMin := rand.Intn(2) == 1
			if isXMin {
				xLocation = 0
			} else {
				xLocation = screenWidth
			}
			yLocation = rand.Intn(screenHeight + 1)
		} else {
			isYMin := rand.Intn(2) == 1
			if isYMin {
				yLocation = 0
			} else {
				yLocation = screenHeight
			}
			xLocation = rand.Intn(screenWidth + 1)
		}

		e := Enemy{
			xPosition: float64(xLocation),
			yPosition: float64(yLocation),
			eSquare:   ebiten.NewImage(10, 10),
		}
		e.eSquare.Fill(color.White)
		g.enemies = append(g.enemies, e)
	}
	for i := range g.enemies {
		g.enemies[i].Update(g.player.xPosition, g.player.yPosition)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.userInput.delete()
	}
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		//if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.spawnProjectile()
	}
	for i, p := range g.projectiles {
		g.projectiles[i].xPosition += p.xVelocity
		g.projectiles[i].yPosition += p.yVelocity
	}
	i := 0
	for _, p := range g.projectiles {
		if !p.isOffScreen() {
			g.projectiles[i] = p
			i++
		}
	}
	g.projectiles = g.projectiles[:i]
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	geo := ebiten.GeoM{}
	geo.Translate(g.player.xPosition-5, g.player.yPosition-5)
	hgeo := ebiten.GeoM{}
	hgeo.Translate(g.player.xPosition-10, g.player.yPosition-10)

	chalo := ebiten.ColorScale{}
	chalo.SetA(.3)
	chalo.SetR(.471)
	chalo.SetG(.29)
	chalo.SetB(.557)

	c := ebiten.ColorScale{}
	c.SetA(1)
	c.SetR(.353)
	c.SetG(.161)
	c.SetB(.443)
	//
	hop := &ebiten.DrawImageOptions{
		GeoM:       hgeo,
		ColorScale: chalo,
	}
	op := &ebiten.DrawImageOptions{
		GeoM:       geo,
		ColorScale: c,
	}
	for _, e := range g.enemies {
		geo := ebiten.GeoM{}
		geo.Translate(e.xPosition, e.yPosition)
		eop := &ebiten.DrawImageOptions{
			GeoM: geo,
		}
		screen.DrawImage(e.eSquare, eop)
	}

	screen.DrawImage(g.player.pSquareHalo, hop)
	screen.DrawImage(g.player.pSquare, op)

	for _, p := range g.projectiles {
		pgeo := ebiten.GeoM{}
		pgeo.Translate(p.xPosition, p.yPosition)
		piop := &ebiten.DrawImageOptions{
			GeoM: pgeo,
		}
		screen.DrawImage(p.pSquare, piop)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Mx: %d My: %d Frame Count: %d TPS: %0.2f %s", g.mouseX, g.mouseY, g.count, ebiten.ActualTPS(), string(g.userInput.currentInput)))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	p := Player{
		xPosition:   centerX,
		yPosition:   centerY,
		pSquare:     ebiten.NewImage(10, 10),
		pSquareHalo: ebiten.NewImage(20, 20),
	}
	p.pSquare.Fill(color.White)
	p.pSquareHalo.Fill(color.White)
	if err := ebiten.RunGame(&Game{
		player: p,
	}); err != nil {
		log.Fatal(err)
	}
}
