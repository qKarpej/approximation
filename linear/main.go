package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 800
)

type Point struct {
	x, y float64
}
type scale struct {
	maxVal   float64
	axisType bool
}

type Game struct {
	width, height int
	scale         float64
	points        []Point
}

func NewGame(width, height int, scl scale, p []Point) *Game {
	var NewGameScale float64
	if scl.axisType {
		NewGameScale = (float64(width/2) / scl.maxVal) * 0.8
	} else {
		NewGameScale = (float64(height/2) / scl.maxVal) * 0.8
	}
	return &Game{
		width:  width,
		height: height,
		scale:  NewGameScale,
		points: p,
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA64{255, 255, 255, 30})
	offset := Point{
		x: float64(g.width / 2),
		y: float64(g.height / 2),
	}
	for _, v := range g.points {
		vector.DrawFilledCircle(screen, float32((v.x*g.scale + offset.x)), float32((v.y*(-1)*g.scale + offset.y)), 5, color.RGBA{255, 0, 0, 255}, false)
	}
	k, b := approximation(g.points)
	vector.StrokeLine(screen, 0, float32(((k*(-offset.x/g.scale)+b)*(-1))*g.scale+offset.y), float32(g.width), float32(((k*(offset.x/g.scale)+b)*(-1))*g.scale+offset.y), 2, color.RGBA{255, 165, 0, 150}, false)
	vector.StrokeLine(screen, 0, float32(offset.y), float32(g.width), float32(offset.y), 3, color.White, false)
	vector.StrokeLine(screen, float32(offset.x), 0, float32(offset.x), float32(g.height), 3, color.White, false)

}

func approximation(p []Point) (k, b float64) {
	var sumX, sumY, sumXY, sumX2 float64
	for _, v := range p {
		sumX += v.x
		sumY += v.y
		sumXY += v.x * v.y
		sumX2 += v.x * v.x
	}
	return (float64(len(p))*sumXY - sumX*sumY) / (float64(len(p))*sumX2 - sumX*sumX), (sumY - k*sumX) / float64(len(p))
}

func main() {
	var n int
	var x, y float64
	fmt.Println("Enter number of points:")
	fmt.Scan(&n)
	p := make([]Point, n)
	maxPoint := scale{-1, false}
	fmt.Println("Enter coordinates of points:")
	for i := 0; i < n; i++ {
		fmt.Scan(&x, &y)
		p[i] = Point{x, y}
		if abs(x) > maxPoint.maxVal {
			maxPoint.maxVal = abs(x)
			maxPoint.axisType = true
		}
		if abs(y) > maxPoint.maxVal {
			maxPoint.maxVal = abs(y)
			maxPoint.axisType = false
		}
	}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g := NewGame(screenWidth, screenHeight, maxPoint, p)
	fmt.Println(g.scale)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func abs(a float64) float64 {
	if a > 0 {
		return a
	} else {
		return a * (-1)
	}
}
