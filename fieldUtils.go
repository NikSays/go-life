package main

import (
	"bufio"
	"image"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
)

func (field *FieldT) iterateField(action func(field *FieldT, x, y int)) {
	for y := 0; y < SIZE; y++ {
		for x := 0; x < SIZE; x++ {
			action(field, x, y)
		}
	}
}
func (field *FieldT) read(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Couldn't open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > SIZE {
			log.Print("Too wide, turncating")
			time.Sleep(time.Second)
		}
		for x := 0; (x < SIZE) && (x < len(line)); x++ {
			if line[x] == CELLCHAR {
				field[y][x] = true
			}
		}
		y++
		if y > SIZE {
			log.Print("Too high, turncating")
			time.Sleep(time.Second)
			break
		}
	}
}
func (field *FieldT) createMap(c *ui.Canvas) {
	printOne := func(field *FieldT, x, y int) {
		if field[y][x] {
			c.SetPoint(image.Pt(x, y), ui.ColorWhite)
		}
	}
	field.iterateField(printOne)
}
