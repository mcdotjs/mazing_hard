package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	x float32
	y float32
}

type Line struct {
	start Point
	end   Point
}

func drawLine(screen *ebiten.Image, a, b Point, color color.RGBA) {
	vector.StrokeLine(screen, a.x, a.y, b.x, b.y, 1, color, true)
}

type Cell struct {
	topLeft      Point
	bottomRight  Point
	middlePoint  Point
	leftBorder   bool
	rightBorder  bool
	topBorder    bool
	bottomBorder bool
	visited      bool
	visible      bool
}

func (cell *Cell) drawCell(screen *ebiten.Image, borderColor color.RGBA) {
	if cell.topBorder {
		drawLine(screen, cell.topLeft, Point{x: cell.bottomRight.x, y: cell.topLeft.y}, borderColor)
	}

	if cell.rightBorder {
		drawLine(screen, Point{x: cell.bottomRight.x, y: cell.topLeft.y}, cell.bottomRight, borderColor)
	}
	if cell.bottomBorder {
		drawLine(screen, cell.bottomRight, Point{x: cell.topLeft.x, y: cell.bottomRight.y}, borderColor)
	}

	if cell.leftBorder {
		drawLine(screen, Point{x: cell.topLeft.x, y: cell.bottomRight.y}, cell.topLeft, borderColor)
	}
}

type Grid struct {
	cols  int
	rows  int
	cells [][]Cell
}

func (c *Cell) drawMove(screen *ebiten.Image, another *Cell) {
	drawLine(screen, Point{x: c.middlePoint.x, y: c.middlePoint.y}, Point{x: another.middlePoint.x, y: another.middlePoint.y}, color.RGBA{255, 0, 0, 255})
}

func createCells(cols int, rows int, padding float32) [][]*Cell {
	cells := make([][]*Cell, 0)
	for c := range cols {

		fmt.Println("in")
		colsSlice := make([]*Cell, 0)
		for r := range rows {
			aa := float32(c) * cellXSize
			bb := float32(r) * cellYSize
			topLeft := Point{x: padding + aa, y: padding + bb}
			bottomRight := Point{x: padding + aa + cellXSize, y: padding + bb + cellYSize}
			newCell := Cell{}
			newCell.topLeft = topLeft
			newCell.bottomRight = bottomRight
			newCell.topBorder = true
			newCell.rightBorder = true
			newCell.bottomBorder = true
			newCell.leftBorder = true
			newCell.visited = false
			newCell.middlePoint = Point{x: (topLeft.x + bottomRight.x) / 2, y: (topLeft.y + bottomRight.y) / 2}
			newCell.visible = false
			colsSlice = append(colsSlice, &newCell)
		}
		cells = append(cells, colsSlice)
	}
	return cells
}

type GridItem struct {
	col int
	row int
}

func drawMoveBase(screen *ebiten.Image, cells [][]*Cell, path [][]int, count int, maze *Maze) {
	//BUG: rf
	// nechapem
	num := maze.numberOfCols * maze.numberOfRows
	fmt.Println(Newpath)
	if len(Newpath) == 0 {
		return
	}
	var newpath [][]int
	if count < num {
		newpath = Newpath[0:count]
	} else {
		newpath = Newpath
	}
	for i := range newpath {
		if i < num-1 && path[i+1] != nil {
			f := path[i]
			s := path[i+1]
			cells[f[0]][f[1]].drawMove(screen, cells[s[0]][s[1]])
		}
	}
}

func removeWalls(cells [][]*Cell, currentItem GridItem, seen *[][]int, path *[][]int, maze *Maze) (bool, [][]int) {
	if currentItem.col > maze.numberOfCols-1 || currentItem.row > maze.numberOfRows-1 || currentItem.col < 0 || currentItem.row < 0 {
		fmt.Println("mimo", currentItem)
		return false, *path
	}
	fmt.Println("nextCellPosition", seen)

	for _, s := range *seen {
		if s[0] == currentItem.col && s[1] == currentItem.row {
			return false, *path
		}
	}
	//INFO: bud pointer alebo deep copy
	// newPath := make([][]int, len(path))
	// copy(newPath, path) // Create a deep copy to avoid modifying the original
	// newPath = append(newPath, []int{current.col, current.row})

	directions := [][]int{
		{-1, 0}, // Up
		{1, 0},  // Down
		{0, -1}, // Left
		{0, 1},  // Right
	}

	// Shuffle or randomize directions if desired
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	*path = append(*path, []int{currentItem.col, currentItem.row})
	*seen = append(*seen, []int{currentItem.col, currentItem.row})

	for _, direction := range directions {
		next := GridItem{
			col: currentItem.col + direction[0],
			row: currentItem.row + direction[1],
		}

		if next.col < 0 || next.col > maze.numberOfCols-1 || next.row < 0 || next.row > maze.numberOfRows-1 {
			// if HERE try another direction
			continue
		}
		alreadyVisited := false
		for _, s := range *seen {
			if s[0] == next.col && s[1] == next.row {
				alreadyVisited = true
				break
				// nasla sa zhoda chod von
			}
		}
		if alreadyVisited {
			// je visited teda chod na dalsi direction
			continue
		}

		// Remove walls between current and next
		currentCell := cells[currentItem.col][currentItem.row]
		nextCell := cells[next.col][next.row]

		// Determine which walls to remove
		if next.row < currentItem.row { // Next is above
			currentCell.topBorder = false
			nextCell.bottomBorder = false
		} else if next.row > currentItem.row { // Next is below
			currentCell.bottomBorder = false
			nextCell.topBorder = false
		} else if next.col < currentItem.col { // Next is to the left
			currentCell.leftBorder = false
			nextCell.rightBorder = false
		} else if next.col > currentItem.col { // Next is to the right
			currentCell.rightBorder = false
			nextCell.leftBorder = false
		}

		if g, p := removeWalls(cells, next, seen, path, maze); g {
			return true, p
		}
	}
	return false, *path
}

func drawCells(screen *ebiten.Image, cells [][]*Cell) {
	for _, r := range cells {
		for _, cell := range r {
			if cell.visible {
				cell.drawCell(screen, color.RGBA{255, 255, 255, 255})
			}
		}
	}

}

func makeCellsVisible(cells [][]*Cell, game *Game, maze *Maze) {
	currentCell := cells[game.count%maze.numberOfCols][(game.count/maze.numberOfCols)%maze.numberOfRows]
	currentCell.visible = true
}

func (game *Game) updatingCounterDva() {
	p, e := time.ParseDuration(game.refreshRate)
	if e != nil {
		fmt.Println(e)
	}
	tic := time.NewTicker(p)
	go func() {
		for {
			select {
			case <-tic.C:
				game.countdva++
				// Do something on each tick
				fmt.Println("going", game.count)
			}
		}
	}()
}
func (game *Game) updatingStuff() {
	p, e := time.ParseDuration(game.refreshRate)
	if e != nil {
		fmt.Println(e)
	}
	tic := time.NewTicker(p)
	go func() {
		for {
			select {
			case <-tic.C:
				game.count++
				// Do something on each tick
				fmt.Println("going", game.count)
			}
		}
	}()
}
func RandBool() bool {
	return rand.Intn(2) == 1
}
