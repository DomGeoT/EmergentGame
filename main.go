package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth, screeHeight = 640, 480
)

type player struct {
	image      *ebiten.Image
	xPos, yPos float64
	speed      float64
}

var (
	err        error
	background *ebiten.Image
	spaceship  *ebiten.Image
	playerOne  player
)

func init() {
	background, _, err = ebitenutil.NewImageFromFile("assets/space.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	spaceship, _, err = ebitenutil.NewImageFromFile("assets/spaceship.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	playerOne = player{spaceship, screenWidth / 2, screeHeight / 2, 4}
}

func movePlayer(xDir float64, yDir float64) {
	playerOne.xPos += xDir * playerOne.speed
	playerOne.yPos += yDir * playerOne.speed
}

func update(screen *ebiten.Image) error {
	var xDir, yDir float64
	xDir, yDir = 0, 0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		yDir--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		yDir++
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		xDir--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		xDir++
	}

	movePlayer(xDir, yDir)

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)

	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(playerOne.xPos, playerOne.yPos)
	screen.DrawImage(playerOne.image, playerOp)

	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screeHeight, 1, "Hello, World!"); err != nil {
		log.Fatal(err)
	}
}
