package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	lastDrawTime time.Time
	count        int
	countdva     int
	refreshRate  string
}

var Newpath [][]int

func (g *Game) Update() error {
	makeCellsVisible(Cells, g, &maze)

	if Cells[maze.numberOfCols-1][maze.numberOfRows-1].visible == true {

		num := maze.numberOfCols * maze.numberOfRows
		if g.count < num {
			Newpath = Path[0:g.countdva]
		} else {
			Newpath = Path
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	drawCells(screen, Cells)
	//BUG: rf
	//drawMoveBase(screen, Cells, Newpath, g.countdva, &maze)

	if Cells[maze.numberOfCols-1][maze.numberOfRows-1].visible == true {
		g.updatingCounterDva()
	}
	num := maze.numberOfCols * maze.numberOfRows
	fmt.Println(num, maze.finalPath)
	temp := Cells[0][0]
	for i := range len(maze.finalPath) - 1 {
		f := maze.finalPath[i]
		s := maze.finalPath[i+1]
		Cells[f[0]][f[1]].drawMove(screen, Cells[s[0]][s[1]])
		temp = Cells[s[0]][s[1]]
	}

	temp.drawMove(screen, Cells[maze.end.col][maze.end.row])
	// beg := GridItem{col: 0, row: 0}
	// seen := [][]int{}
	// path := [][]int{}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1200, 800
}

var Cells [][]*Cell
var Path [][]int

type Maze struct {
	numberOfCols     int
	numberOfRows     int
	padding          int
	start            GridItem
	end              GridItem
	drawingRowNumber int
	drawingFertig    bool
	finalPath        [][]int
}

var maze = Maze{
	numberOfCols:  13,
	numberOfRows:  8,
	padding:       80,
	drawingFertig: false,
}

func main() {
	ebiten.SetWindowSize(1200, 800)
	ebiten.SetWindowTitle("Hello, World!")
	game := &Game{
		refreshRate: "33ms",
	}
	game.updatingStuff()
	Cells = createCells(maze.numberOfCols, maze.numberOfRows, float32(maze.padding))

	maze.start = GridItem{col: 0, row: 0}
	maze.end = GridItem{col: maze.numberOfCols - 1, row: maze.numberOfRows - 1}
	Cells[maze.start.col][maze.start.row].topBorder = false
	Cells[maze.end.col][maze.end.row].bottomBorder = false
	beg := GridItem{col: 0, row: 0}

	seen := [][]int{}
	path := [][]int{}
	_, Path = removeWalls(Cells, beg, &seen, &path, &maze)

	start := GridItem{col: 0, row: 0}
	end := GridItem{col: maze.numberOfCols - 1, row: maze.numberOfRows - 1}
	maze.finalPath = solve(&maze, Cells, start, end)
	fmt.Println("main", maze.finalPath)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
