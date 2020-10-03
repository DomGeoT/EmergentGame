package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth, screenHeight = 640, 480
)

type entity struct {
	image         *ebiten.Image
	xPos, yPos    float64
	width, height float64
	speed         float64
}

func (e *entity) moveEntity(xDir float64, yDir float64) {
	e.xPos += xDir * e.speed
	e.yPos += yDir * e.speed
}

func (e entity) drawEntity(screen *ebiten.Image) {
	entOp := &ebiten.DrawImageOptions{}
	entOp.GeoM.Translate(e.xPos, e.yPos)
	screen.DrawImage(e.image, entOp)
}

type momentumEntity struct {
	baseEntity entity
	xDir, yDir float64
}

type rect struct {
	x, y          float64
	width, height float64
}

func (e momentumEntity) getNextRect() rect {
	return rect{
		e.baseEntity.xPos + e.xDir,
		e.baseEntity.yPos + e.yDir,
		e.baseEntity.width,
		e.baseEntity.height}
}

func (e entity) getRect() rect {
	return rect{
		e.xPos,
		e.yPos,
		e.width,
		e.height}
}

func (r1 rect) intersect(r2 rect) bool {
	if r1.x >= r2.x+r2.width || r2.x >= r1.x+r1.width {
		return false
	}

	if r1.y+r1.height <= r2.y || r2.y+r2.height <= r1.y {
		return false
	}

	return true
}

func (e *momentumEntity) moveEntity(xDir float64, yDir float64) {
	e.xDir += xDir
	e.yDir += yDir
	e.baseEntity.xPos += e.xDir * e.baseEntity.speed
	e.baseEntity.yPos += e.yDir * e.baseEntity.speed
}

func (e *momentumEntity) bounceYDir() {
	e.yDir *= -1
}

func (e *momentumEntity) bounceXDir() {
	e.xDir *= -1
}

func willHit(e entity, m momentumEntity) bool {
	var eHitbox, mNextHitbox rect
	eHitbox = e.getRect()
	mNextHitbox = m.getNextRect()

	return eHitbox.intersect(mNextHitbox)
}

var (
	err        error
	background *ebiten.Image
	ball       *ebiten.Image
	paddle     *ebiten.Image
	playerOne  entity
	playerTwo  entity
	ballEnt    momentumEntity
)

func init() {
	background, _, err = ebitenutil.NewImageFromFile("../assets/space.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	ball, _, err = ebitenutil.NewImageFromFile("../assets/ball.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	paddle, _, err = ebitenutil.NewImageFromFile("../assets/paddle.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	playerOne = entity{paddle, screenWidth / 2, 16, 64, 16, 4}
	playerTwo = entity{paddle, screenWidth / 2, screenHeight - 32, 64, 16, 4}

	ballEnt = momentumEntity{entity{ball, screenWidth / 2, screenHeight / 2, 32, 32, 4}, 0, 0}
	ballEnt.xDir = 1
	ballEnt.yDir = 1
}

func update(screen *ebiten.Image) error {
	var p1XDir float64
	p1XDir = 0

	var p2XDir float64
	p2XDir = 0

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p1XDir--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p1XDir++
	}

	if ballEnt.baseEntity.xPos > playerTwo.xPos {
		p2XDir = 1
	} else {
		p2XDir = -1
	}

	if ballEnt.baseEntity.xPos <= 0 || ballEnt.baseEntity.xPos+ballEnt.baseEntity.width >= screenWidth {
		ballEnt.bounceXDir()
	}

	if willHit(playerOne, ballEnt) {
		ballEnt.bounceYDir()
		ballEnt.baseEntity.speed *= 1.01
	}
	if willHit(playerTwo, ballEnt) {
		ballEnt.bounceYDir()
		ballEnt.baseEntity.speed *= 1.01
	}

	playerOne.moveEntity(p1XDir, 0)
	playerTwo.moveEntity(p2XDir, 0)
	ballEnt.moveEntity(0, 0)

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if ballEnt.baseEntity.yPos > screenHeight || ballEnt.baseEntity.yPos < -ballEnt.baseEntity.width {
		ballEnt.baseEntity.xPos = screenWidth / 2
		ballEnt.baseEntity.yPos = screenHeight / 2
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)

	playerOne.drawEntity(screen)
	playerTwo.drawEntity(screen)
	ballEnt.baseEntity.drawEntity(screen)

	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Hello, World!"); err != nil {
		log.Fatal(err)
	}
}
