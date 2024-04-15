package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

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
	axisType bool // x:true , y:false
}

type Game struct {
	width, height int
	scale         float64
	points        []Point
}

func NewGame(width, height int, scl scale, p []Point) *Game {
	var NewGameScale float64
	//0.8 coeficient is added, so larger point won't be located on the frame of the screen
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
	vector.StrokeLine(screen, 0, float32(offset.y), float32(g.width), float32(offset.y), 3, color.White, false)
	vector.StrokeLine(screen, float32(offset.x), 0, float32(offset.x), float32(g.height), 3, color.White, false)
	for _, v := range g.points {
		vector.DrawFilledCircle(screen, float32((v.x*g.scale + offset.x)), float32((v.y*(-1)*g.scale + offset.y)), 5, color.RGBA{255, 0, 0, 255}, false)
	}
	a, b, c := approximationQuad(g.points)
	fmt.Println(a, b, c)
	for x := -200.0 / g.scale; x < 200/g.scale; x += 0.01 / g.scale {
		screen.Set(int(x*g.scale)+int(offset.x), int((a*x*x+b*x+c)*g.scale*(-1))+int(offset.x), color.White)
	}

}

func approximationQuad(p []Point) (a, b, c float64) {
	var sumX, sumY, sumXY, sumX2, sumX2Y float64
	for _, v := range p {
		sumX += v.x
		sumY += v.y
		sumXY += v.x * v.y
		sumX2 += v.x * v.x
		sumX2Y += v.x * v.x * v.y
	}
	n := float64(len(p))
	a = (n*sumXY - sumX*sumY) / (n*sumX2 - math.Pow(sumX, 2))
	b = sumY - a*(sumX)
	c = (sumY / n) - (b * sumX / n) - (a * sumX2 / n)
	return
}

func main() {
	var n int
	var x, y float64
	fmt.Println("Enter number of points:")
	fmt.Scan(&n)
	p := make([]Point, n)
	maxPoint := scale{-1, false}
	fmt.Println("Enter coordinates of points:")
	// p = []Point{{10000, 12350}, {-10000, -13210}, {100, -13210}, {-3200, -13210}, {-12300, -13210}, {-1020, -1323}}
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
