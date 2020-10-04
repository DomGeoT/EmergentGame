package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth, screenHeight = 640, 480
)

const (
	tileWidth, tileHeight = 64, 64
)

const (
	worldWidth, worldHeight = 20, 20
)

type Game struct{}

var (
	grassTileImg, rockTileImg *ebiten.Image
	game                      Game
	world                     tilemap
)

type tile struct {
	img  *ebiten.Image
	x, y float64
}

func (t tile) drawTile(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(t.x, t.y)
	screen.DrawImage(t.img, op)
}

type tilemap struct {
	tiles [][]tile
}

func init() {
	var err error
	grassTileImg, _, err = ebitenutil.NewImageFromFile("../assets/grassTile.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	rockTileImg, _, err = ebitenutil.NewImageFromFile("../assets/rockTile.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	world = tilemap{generateTileMap()}
}

func (g *Game) Update(screen *ebiten.Image) error {

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var i, j int

	for i = 0; i < worldWidth; i++ {
		for j = 0; j < worldHeight; j++ {
			world.tiles[i][j].drawTile(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func generateTileMap() [][]tile {
	var tileMap [][]tile
	var i, j int

	tileMap = make([][]tile, worldHeight)
	for i := range tileMap {
		tileMap[i] = make([]tile, worldWidth)
	}

	for i = 0; i < worldWidth; i++ {
		for j = 0; j < worldHeight; j++ {
			if i == 0 || j == 0 || (i-1) == worldWidth || (j-1) == worldHeight {
				tileMap[i][j] = tile{rockTileImg, (float64)(i * tileWidth), (float64)(j * tileHeight)}
			} else {
				tileMap[i][j] = tile{grassTileImg, (float64)(i * tileWidth), (float64)(j * tileHeight)}
			}
		}
	}

	return tileMap
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Your game's title")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
