package main

/*
The program as it is currently written does not support values of n greater than 120-ish
*/

import (
	"fmt"
	"math/rand"
	"sync"
	// "time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

/*
There are two places where we had to primarily use the two constants mentioned
below: firstly to determine the size of the window, which demands that it be int,
and secondly to describe the position of the gridlines, which demand the values to
be float64. This is why the following two values are kept as untyped constants.
*/
const WIN_WIDTH = 768
const WIN_HEIGHT = 768

type State int

/*
DEAD is given 0 and ALIVE 1, instead of the other way round, because in the `process`
function we have to find the number of live neighbours and this assignment makes it
cleaner to write
*/
const (
	DEAD = iota
	ALIVE
)

type Cell struct {
	lower_left pixel.Vec
	length     float64
	rect       *imdraw.IMDraw
	state      State
}

func (c *Cell) germinate() {
	c.state = ALIVE
	c.rect.Clear()
	c.rect.Color = pixel.RGB(0, 0, 0)
	c.rect.Push(pixel.V(c.lower_left.X, c.lower_left.Y))
	c.rect.Push(pixel.V(c.lower_left.X+c.length, c.lower_left.Y+c.length))
	c.rect.Rectangle(0)
}

func (c *Cell) depart_this_life() {
	c.state = DEAD
	c.rect.Clear()
	c.rect.Color = pixel.RGB(1, 1, 1)
	c.rect.Push(pixel.V(c.lower_left.X, c.lower_left.Y))
	c.rect.Push(pixel.V(c.lower_left.X+c.length, c.lower_left.Y+c.length))
	c.rect.Rectangle(0.5)
}

func process(cur_board [][]Cell, row, col int, next_board [][]Cell, wg *sync.WaitGroup) {
	defer wg.Done()

	dy := [8]int{-1, 0, 1, 1, 1, 0, -1, -1}
	dx := [8]int{-1, -1, -1, 0, 1, 1, 1, 0}

	count := 0
	for i := range dx {
		var nrow int = row + dx[i]
		var ncol int = col + dy[i]
		if is_valid(cur_board, nrow, ncol) {
			count += int(cur_board[nrow][ncol].state)
		}
	}

	if (count > 3 || count < 2) && cur_board[row][col].state == ALIVE {
		next_board[row][col].depart_this_life()
	} else if count == 3 && cur_board[row][col].state == DEAD {
		next_board[row][col].germinate()
	} else {
		if cur_board[row][col].state == ALIVE {
			next_board[row][col].germinate()
		} else {
			next_board[row][col].depart_this_life()
		}
	}
}

func run() {
	var n float64
	fmt.Print("Enter the number of cells ")
	fmt.Scan(&n)
	var INIT_PROPORTION_FILLED float64
	if n < 30{
		INIT_PROPORTION_FILLED = 0.20
	} else{
		INIT_PROPORTION_FILLED = 0.10
	}

	cfg := pixelgl.WindowConfig{
		Title:     "Game of Life",
		Bounds:    pixel.R(0, 0, WIN_WIDTH, WIN_HEIGHT),
		Icon:      []pixel.Picture{pixel.PictureDataFromImage(getIcon())},
		VSync:     true,
		Resizable: false,
		Maximized: false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	im_draw := imdraw.New(nil)

	// Let us draw the gridlines!
	im_draw.Color = pixel.RGB(1, 1, 1)
	for i := 1; i < int(n); i++ {
		im_draw.Push(pixel.V(float64(i)*(WIN_WIDTH/n), 0),
			pixel.V(float64(i)*(WIN_WIDTH/n), WIN_HEIGHT))
		im_draw.Line(1)
	}
	for j := 1; j < int(n); j++ {
		im_draw.Push(pixel.V(0, float64(j)*(WIN_HEIGHT/n)),
			pixel.V(WIN_WIDTH, float64(j)*(WIN_HEIGHT/n)))
		im_draw.Line(1)
	}

	cur_board := make([][]Cell, int(n))
	for row := range cur_board {
		cur_board[row] = make([]Cell, int(n))
	}
	for row := range cur_board {
		for col := range cur_board[row] {
			cur_board[row][col].lower_left = pixel.V(float64(row)*(WIN_WIDTH/n),
				float64(col)*(WIN_HEIGHT/n))
			cur_board[row][col].length = min(WIN_HEIGHT, WIN_WIDTH) / n
			cur_board[row][col].rect = imdraw.New(nil)
			cur_board[row][col].state = DEAD // initialization
			rect := cur_board[row][col].rect
			rect.Color = pixel.RGB(1, 1, 1)

			/*
				Interesting thing is that you do not need to specify all the 4 points of
				the rectangle: a fact that I overlooked in the first attempt and only
				discovered it serendipitously. Read the documentation, folks.
			*/
			rect.Push(pixel.V(cur_board[row][col].lower_left.X, cur_board[row][col].lower_left.Y))
			rect.Push(pixel.V(cur_board[row][col].lower_left.X+cur_board[row][col].length, cur_board[row][col].lower_left.Y+cur_board[row][col].length))
			rect.Rectangle(0.5)
		}
	}

	// Deep copy of cur_board slice 
	next_board := make([][]Cell, int(n))
	for row := range next_board {
		next_board[row] = make([]Cell, int(n))
		for col := range next_board[row] {
			next_board[row][col].lower_left = cur_board[row][col].lower_left
			next_board[row][col].length = cur_board[row][col].length
			next_board[row][col].rect = imdraw.New(nil)
			next_board[row][col].state = DEAD
			rect := next_board[row][col].rect
			rect.Color = pixel.RGB(1, 1, 1)
			rect.Push(pixel.V(next_board[row][col].lower_left.X, next_board[row][col].lower_left.Y))
			rect.Push(pixel.V(next_board[row][col].lower_left.X+next_board[row][col].length, 
				next_board[row][col].lower_left.Y+next_board[row][col].length))
			rect.Rectangle(0.5)
		}
	}

	for row := range cur_board {
		for col := range cur_board[row] {
			r := rand.Float64()
			if 0.00 <= r && r <= INIT_PROPORTION_FILLED {
				next_board[row][col].germinate()
			}
		}
	}

	// The animation loop in which the frame updates take place
	for !win.Closed() {
		win.Clear(colornames.Skyblue)

		cur_board, next_board = next_board, cur_board

		/*
			Because we are inevitably going to have a large number of cells, I thought
			why not explore this exciting feature in Go I just learnt about? So here goes
			concurrency:
		*/
		var wg sync.WaitGroup
		for row := range cur_board {
			for col := range cur_board[row] {
				wg.Add(1)
				go process(cur_board, row, col, next_board, &wg)
			}
		}
		wg.Wait()

		var draw_wg sync.WaitGroup
		im_draw.Draw(win)
		for row := range cur_board {
			for col := range cur_board[row] {
				draw_wg.Add(1)
				go func(){
					defer draw_wg.Done()
					cell := cur_board[row][col]
					cell.rect.Draw(win)
				}()
			}
		}
		draw_wg.Wait()

		win.Update()
		// time.Sleep(50 * time.Millisecond) 
		//Uncomment the above line for better visualization for small values of n
	}
}

func main() {
	pixelgl.Run(run)
}
