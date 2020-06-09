package main

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width          = 1000
	height         = 800
	numCellX       = 50
	numCellY       = 50
	dimensionCellX = float64(width / numCellX)
	dimensionCellY = float64(height / numCellY)
	gameTitle      = "Game of life"
)

var (
	backgroundColor = color.RGBA{R: 25, G: 25, B: 25}
)

func main() {
	pixelgl.Run(run)
}

func run() {
	window := createWindow()
	imd := imdraw.New(nil)
	gameState := [numCellX][numCellY]int{}
	paused := false

	//Initial game state, remove if you don't want it
	gameState[10][1] = 1
	gameState[10][2] = 1
	gameState[10][3] = 1

	gameState[21][21] = 1
	gameState[22][22] = 1
	gameState[22][23] = 1
	gameState[21][23] = 1
	gameState[20][23] = 1

	for !window.Closed() {

		//Pause
		if window.JustPressed(pixelgl.KeySpace) && !paused {
			paused = true
		} else if window.JustPressed(pixelgl.KeySpace) && paused {
			paused = false
		}

		if !paused {

			paused = false
			//Restart Game Logic
			window.Clear(backgroundColor)
			imd.Clear()
			time.Sleep(time.Millisecond * 100)
			newGameState := gameState

			cellX, cellY, state := getStateIfClicked(window)
			if state != -1 {
				newGameState[cellX][cellY] = state
			}

			for x := 0; x < numCellX; x++ {
				for y := 0; y < numCellY; y++ {

					//Control logic
					controlledX := x
					controlledY := y
					if x-1 == -1 {
						controlledX = numCellX - 1
					} else {
						controlledX = (x - 1) % numCellX
					}
					if y-1 == -1 {
						controlledY = numCellY - 1
					} else {
						controlledY = (y - 1) % numCellY
					}
					//Game logic
					numberOfNeighs := gameState[controlledX][controlledY] +
						gameState[x][controlledY] +
						gameState[(x+1)%numCellX][controlledY] +
						gameState[controlledX][y%numCellY] +
						gameState[(x+1)%numCellX][y%numCellY] +
						gameState[controlledX][(y+1)%numCellY] +
						gameState[x%numCellX][(y+1)%numCellY] +
						gameState[(x+1)%numCellX][(y+1)%numCellY]

					if gameState[x][y] == 0 && numberOfNeighs == 3 {
						newGameState[x][y] = 1
					} else if gameState[x][y] == 1 && (numberOfNeighs < 2 || numberOfNeighs > 3) {
						newGameState[x][y] = 0
					}

					//Drawing logic
					poly := getPolygonVect(x, y)
					if newGameState[x][y] == 0 {
						//uncomment if you want to see a grid
						//drawBlankCell(imd, poly)
					} else {
						drawFilledCell(imd, poly)
					}
				}
				imd.Draw(window)
			}
			//Updating logic
			gameState = newGameState
		} else {
			//TODO: Extract to method
			cellX, cellY, state := getStateIfClicked(window)
			if state != -1 {
				poly := getPolygonVect(cellX, cellY)
				if state == 0 {
					drawBlankCell(imd, poly)
				} else {
					drawFilledCell(imd, poly)
				}
				imd.Draw(window)
				gameState[cellX][cellY] = state
			}
		}

		window.Update()
	}
}

func createWindow() *pixelgl.Window {
	config := pixelgl.WindowConfig{
		Title:  gameTitle,
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}
	window.Clear(backgroundColor)
	return window
}

func getPolygonVect(x, y int) *[]pixel.Vec {
	return &[]pixel.Vec{
		{X: (float64(x) * dimensionCellX), Y: (float64(y) * dimensionCellY)},
		{X: ((float64(x) + 1) * dimensionCellX), Y: (float64(y) * dimensionCellY)},
		{X: ((float64(x) + 1) * dimensionCellX), Y: (float64(y+1) * dimensionCellY)},
		{X: (float64(x) * dimensionCellX), Y: (float64(y+1) * dimensionCellY)},
	}
}

func drawFilledCell(imd *imdraw.IMDraw, poly *[]pixel.Vec) {
	imd.Color = color.White
	imd.Push(*poly...)
	imd.Polygon(0)
}

func drawBlankCell(imd *imdraw.IMDraw, poly *[]pixel.Vec) {
	imd.Color = color.Black
	imd.Push(*poly...)
	imd.Polygon(0)
}

func getStateIfClicked(window *pixelgl.Window) (int, int, int) {
	if window.Pressed(pixelgl.MouseButtonLeft) {
		celX, celY := getCellXYByMousePosition(window)
		return celX, celY, 1
	} else if window.Pressed(pixelgl.MouseButtonRight) {
		celX, celY := getCellXYByMousePosition(window)
		return celX, celY, 0
	}
	return -1, -1, -1
}

func getCellXYByMousePosition(window *pixelgl.Window) (int, int) {
	mouseX, mouseY := window.MousePosition().XY()
	celX := int(mouseX / dimensionCellX)
	celY := int(mouseY / dimensionCellY)
	return celX, celY
}
