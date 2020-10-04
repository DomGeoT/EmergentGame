package camera

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/math/f64"
)

// Camera for viewing a world
type Camera struct {
	Viewport   f64.Vec2
	Position   f64.Vec2
	ZoomFactor int
	Rotation   int
}

// String generates a string representation of a camera
func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.Viewport[0] * 0.5,
		c.Viewport[1] * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

// Render draws provided world image to the screen according to the camera
func (c *Camera) Render(world, screen *ebiten.Image) error {
	return screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

// ScreenToWorld converts screen coords to world position
func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	}
	// When scaling it can happend that matrix is not invertable
	return math.NaN(), math.NaN()
}

// Reset moves
func (c *Camera) Reset() {
	c.Position[0] = 0
	c.Position[1] = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}

// Game contains
type Game struct {
	layers [][]int
	world  *ebiten.Image
	camera Camera
}
