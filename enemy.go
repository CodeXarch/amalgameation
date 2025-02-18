package main

import "github.com/hajimehoshi/ebiten/v2"

type Enemy struct {
	xPosition float64
	yPosition float64
	eSquare   *ebiten.Image
}

func (e *Enemy) Update(playerX float64, playerY float64) {
	// Move the enemy towards the player
	if e.xPosition < playerX {
		e.xPosition += .1
	} else {
		e.xPosition -= .1
	}
	if e.yPosition < playerY {
		e.yPosition += .1
	} else {
		e.yPosition -= .1
	}
}
