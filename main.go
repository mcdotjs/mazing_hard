package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	lastDrawTime time.Time
	count        int
	stopCount    bool
	countdva     int
	refreshRate  string
}

var newpath [][]int

func (g *Game) Update() error {
	if g.count == len(maze.finalPath) {
		g.stopCount = true
		newpath = maze.finalPath
	} else {
		newpath = maze.finalPath[0 : g.count%len(maze.finalPath)-1]
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	drawCells(screen, Cells)

	for i := 0; i < len(newpath)-1; i++ {
		f := newpath[i]
		s := newpath[i+1]
		Cells[f[0]][f[1]].drawMove(screen, Cells[s[0]][s[1]])
	}

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
	numberOfCols:  100,
	numberOfRows:  70,
	padding:       30,
	drawingFertig: false,
}

func main() {
	ebiten.SetWindowSize(1200, 800)
	ebiten.SetWindowTitle("Hello, World!")
	game := &Game{
		refreshRate: "33ms",
		stopCount:   false,
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
	maze.finalPath = append(maze.finalPath, []int{maze.numberOfCols - 1, maze.numberOfRows - 1})
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
