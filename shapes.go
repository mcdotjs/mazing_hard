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

func solve(screen *ebiten.Image, maze *Maze, cells [][]*Cell, current GridItem, path *[][]int, seen *[][]int, row int) (bool, [][]int) {
	fmt.Println("seenn", seen, path, current)
	if current == maze.end {
		return true, *path
	}
	if current.col > maze.numberOfCols-1 || current.row > maze.numberOfRows-1 || current.col < 0 || current.row < 0 {
		fmt.Println("mimo", current)
		return true, *path
	}

	for _, s := range *seen {
		if s[0] == current.col && s[1] == current.row {
			return false, *path
		}
	}

	*seen = append(*seen, []int{current.col, current.row})
	*path = append(*path, []int{current.col, current.row})
	//_____________________________________________________________________
	// NOTE: pre vsetky styry susedov musim zavolat
	// toto setovanie rowu mi ide zaradom a to neni dobre

	fmt.Println(row)
	directions := [][]int{
		{-1, 0}, // Up
		{1, 0},  // Down
		{0, -1}, // Left
		{0, 1},  // Right
	}
	for i := range directions {
		next := GridItem{
			col: (current.col + directions[i][0]) % maze.numberOfCols,
			row: (current.row + directions[i][1]) % maze.numberOfRows,
		}
		if next.col < 0 || next.row < 0 {
			continue
		}

		currentCell := cells[current.col][current.row]
		nextCell := cells[next.col][next.row]
		fmt.Println("curr", currentCell)
		fmt.Println("next", nextCell)
		if current.col < next.col && (!nextCell.leftBorder && !currentCell.rightBorder) { // Next is above
			currentCell.drawMove(screen, nextCell)

		} else if current.row < next.row && (!nextCell.topBorder && !currentCell.bottomBorder) { // Next is above
			currentCell.drawMove(screen, nextCell)

		}

		if current.row > next.row && !nextCell.bottomBorder && !currentCell.topBorder { // Next is above
			currentCell.drawMove(screen, nextCell)
		}
		if current.col > next.col && (!nextCell.rightBorder && !currentCell.leftBorder) { // Next is above
			currentCell.drawMove(screen, nextCell)
		}
		// }
		// if !nextCell.topBorder && !currentCell.bottomBorder { // Next is to the right
		// 	currentCell.drawMove(screen, nextCell)
		// }
		//
		// else if !nextCell.topBorder && !currentCell.bottomBorder { // Next is to the right
		// 	currentCell.drawMove(screen, nextCell)
		// }
		//}
		//__________________________________________________________________________
		if g, p := solve(screen, maze, cells, next, seen, path, row); g {
			return true, p
		}

	}

	//BUG: here !!! musis mu dat na vyber so susedov ... nie setnut dalsieho
	// next := GridItem{
	// 	col: (current.col) % maze.numberOfCols,
	// 	row: row,
	// }
	// //NOTE: tu urobit possible neighnbors?
	// a foreach???

	// var chosenDirection []int
	// directions := map[string][]int{
	// 	"up":    {-1, 0}, // Up
	// 	"down":  {1, 0},  // Down
	// 	"left":  {0, -1}, // Left
	// 	"right": {0, 1},  // Right
	// }
	//
	// currentCell := cells[current.col][current.row]
	// if !currentCell.rightBorder {
	// 	chosenDirection = directions["right"]
	// } else if !currentCell.leftBorder {
	// 	chosenDirection = directions["left"]
	// } else if !currentCell.bottomBorder {
	// 	chosenDirection = directions["down"]
	// } else if !currentCell.topBorder {
	// 	chosenDirection = directions["up"]
	// }
	// //for _, dir := range directions {
	// c := current.col + chosenDirection[0]
	// r := current.row + chosenDirection[1]
	// fmt.Println(c, r, chosenDirection)
	// if c < 1 || r < 1 {
	// 	return false, *path
	// }
	// next := GridItem{
	// 	col: c,
	// 	row: r,
	// }
	// if next.col < 0 || next.col > maze.numberOfCols-1 || next.row < 0 || next.row > maze.numberOfRows-1 {
	// 	// if HERE try another direction
	// 	return false, *path
	// }
	// for _, s := range *seen {
	// 	if s[0] == next.col && s[1] == next.row {
	// 		return false, *path
	// 		// nasla sa zhoda chod von
	// 	}
	// }
	//
	// Remove walls between current and next
	// if next.col < 0 || next.row < 0 {
	// 	return false, *path
	// }

	// Determine which walls to remove
	*path = (*path)[:len(*path)-1]
	return false, *path
}

func solvei(screen *ebiten.Image, maze *Maze, cells [][]*Cell, current GridItem, path *[][]int, seen *[][]int, row int) (bool, [][]int) {
	if current == maze.end {
		return true, *path
	}

	if current.col > maze.numberOfCols-1 || current.row > maze.numberOfRows-1 || current.col < 0 || current.row < 0 {
		fmt.Println("mimo", current)
		return false, *path
	}
	//fmt.Println("nextCellPosition", seen)

	for _, s := range *seen {
		if s[0] == current.col && s[1] == current.row {
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
	// directionsMap := map[string][]int{
	// 	"top":    {-1, 0}, // Up
	// 	"bottom": {1, 0},  // Down
	// 	"left":   {0, -1}, // Left
	// 	"right":  {0, 1},  // Right
	// }

	*path = append(*path, []int{current.col, current.row})

	currentCell := cells[current.col][current.row]
	var currentDirections []string
	if currentCell.topBorder == false {
		currentDirections = append(currentDirections, "top")
	}

	if currentCell.bottomBorder == false {
		currentDirections = append(currentDirections, "bottom")
	}

	if currentCell.leftBorder == false {
		currentDirections = append(currentDirections, "left")
	}
	if currentCell.rightBorder == false {
		currentDirections = append(currentDirections, "right")
	}
	fmt.Println(len(currentDirections))
	//for _, val := range directionsMap {
	for i := range directions {
		//for i := 4; i < 4; i++ {
		// next := GridItem{
		// 	col: current.col + val[0],
		// 	row: current.row + val[1],
		// }
		fmt.Println(directions[i], currentDirections[i])
		next := GridItem{
			col: (current.col + directions[i][0]) % maze.numberOfCols,
			row: (current.row + directions[i][1]) % maze.numberOfRows,
		}
		if next.col < 0 || next.row < 0 {
			continue
		}

		//

		// Remove walls between current and next
		nextCell := cells[next.col][next.row]
		// if currentDirections[i] == "top" && nextCell.bottomBorder == false {
		// 	currentCell.drawMove(screen, nextCell)
		// }
		//
		// if currentDirections[i] == "bottom" && nextCell.topBorder == false {
		// 	currentCell.drawMove(screen, nextCell)
		// }
		//
		// if currentDirections[i] == "left" && nextCell.rightBorder == false {
		// 	currentCell.drawMove(screen, nextCell)
		// }

		if !nextCell.topBorder {
			currentCell.drawMove(screen, nextCell)
		}

		if !nextCell.bottomBorder {
			currentCell.drawMove(screen, nextCell)
		}
		// // if currentCell.rightBorder && nextCell.leftBorder {
		// 	fmt.Println("LLLLLL")
		//
		// 	return false, *path
		// }
		//
		// if currentCell.leftBorder && nextCell.rightBorder {
		// 	fmt.Println(",,,,,LLLLLL")
		// 	return false, *path
		// }
		fmt.Println("curr", currentCell)
		fmt.Println("next", nextCell)
		// //
		if g, p := solvei(screen, maze, cells, next, seen, path, row); g {
			return true, p
		}
	}
	return false, *path
}
func removeWalls(cells [][]*Cell, currentItem GridItem, seen *[][]int, path *[][]int, maze *Maze) (bool, [][]int) {
	if currentItem.col > maze.numberOfCols-1 || currentItem.row > maze.numberOfRows-1 || currentItem.col < 0 || currentItem.row < 0 {
		fmt.Println("mimo", currentItem)
		return false, *path
	}
	//fmt.Println("nextCellPosition", seen)

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
				//game.countdva++
				// Do something on each tick
				//fmt.Println("going", game.count)
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
				//fmt.Println("going", game.count)
			}
		}
	}()
}
func RandBool() bool {
	return rand.Intn(2) == 1
}
