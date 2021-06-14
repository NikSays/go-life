package main

import (
	"image"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/drawille"
)

const SIZE = 212
const CELLCHAR = '#'
const FPS = 3000
const FILENAME = "field.in"

var temp image.Rectangle
var CLEARBUF = ui.NewBuffer(temp)

type FieldT [SIZE][SIZE]bool

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var field FieldT
	field.read(FILENAME)

	window := ui.NewCanvas()
	window.SetRect(0, 0, SIZE, SIZE)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second / FPS).C
	i := 0
	for {
		i++
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			window.Canvas.CellMap = make(map[image.Point]drawille.Cell)
			field.createMap(window)
			field.conwayStep()
			ui.Render(window)
		}
	}
}

func (field *FieldT) displacedDiagonal(displacement int) {
	updateCell := func(field *FieldT, x int, y int) {
		if x+displacement == y {
			field[y][x] = true
		} else {
			field[y][x] = false
		}
	}
	field.iterateField(updateCell)
}
func (field *FieldT) conwayStep() {
	var newField FieldT
	//Create torus coordinates
	//Example: for size:=20, coordinates are from 0 to 19
	//         warpCoordinates turns -1 to 19, 20 to 0
	wrapCoordinates := func(x, y int) (int, int) {
		warpx := (x + SIZE) % SIZE
		warpy := (y + SIZE) % SIZE
		return warpx, warpy
	}
	//Count all live cells around given x,y; excluding x,y cell itself
	//                                       including those around screen corners
	countNeighbors := func(field *FieldT, x, y int) (sum int) {
		for j := y - 1; j <= y+1; j++ {
			for i := x - 1; i <= x+1; i++ {
				if (i == x) && (j == y) {
					continue
				}
				warpi, warpj := wrapCoordinates(i, j)
				if field[warpj][warpi] {
					sum++
				}
			}
		}
		return
	}
	//
	newField = *field
	updateCell := func(field *FieldT, x, y int) {
		switch count := countNeighbors(field, x, y); {
		case count < 2:
			newField[y][x] = false
		case count == 3:
			newField[y][x] = true
		case count > 3:
			newField[y][x] = false
		}
	}
	field.iterateField(updateCell)
	*field = newField
}
