package main 

import (
	"fmt"
	"os"
	"strconv"
)



//Board
type Board struct {
	size int  // number of rows and cols in the board
	cell [][]int // 
	st stack // 
}

//
type Cell struct {
	r int
	c int
}


type stack []Cell

func (s stack) Empty() bool { return len(s) == 0 }
func (s stack) Peek() Cell   { return s[len(s)-1] }
func (s *stack) Push(i Cell)  { (*s) = append((*s), i) }
func (s *stack) Pop() Cell {
  d := (*s)[len(*s)-1]
  (*s) = (*s)[:len(*s)-1]
  return d
}


//Create a new Board, initilize it with efualt configuration
func createBoard(size, numOfSandpiles int) *Board {
	cell := make([][]int, size, size);
	for i := 0; i < size; i++ {
		cell[i] = make([]int, size)
		for j := 0; j < size; j++ {
			cell[i][j] = 0
		}
	}
	cell[size/2][size/2] = numOfSandpiles
	st := make(stack, 0)
	if numOfSandpiles >= 4 {
		st.Push(Cell{size/2, size/2})
	}
	return &Board{size, cell, st}
}

// returns true if ( r, c) is within the field.
// otherwise return false
func (b *Board) Contains(r, c int) bool {
	if r >= 0 && c >= 0 && r < b.size && c < b.size { return true }
	return false
}

// sets the value of cell ( r, c)
func (b *Board) Set(r, c, value int) {
	if b.Contains(r, c) {
		b.cell[r][c] = value
	}
}

// returns the value of the cell ( r, c).
func (b *Board) Cell(r, c int) int {
	if b.Contains(r, c) {
		return b.cell[r][c]
	}
	return -1;
}

func (b *Board) isConverged() bool {
	if b.st.Empty() {
		return true
	}
	return false
}

func (b *Board) NumRows() int {
	return b.size
}

func (b *Board) NumCols() int {
	return b.size
}

func (b *Board) Topple(r, c int) {
	value := b.Cell(r, c)
	b.Set(r, c, value - 4)
	if value - 4 >= 4 {
		b.st.Push(Cell{r,c})
	}
	b.UpdateCell(r-1, c)
	b.UpdateCell(r+1, c)
	b.UpdateCell(r, c-1)
	b.UpdateCell(r, c+1)
}

func (b *Board) UpdateCell(r, c int) {
	if b.Contains(r, c) {
		b.Set(r, c, b.Cell(r, c) + 1)
		if b.Cell(r, c) >= 4 {
			b.Topple(r, c)
		}
	}
}

func ComputeSteadyState(b *Board) {
	for !b.st.Empty() {
		//fmt.Println(len(b.st))
		cell := b.st.Pop()
		b.Topple(cell.r, cell.c)
	}
}

func DrawBoard(b *Board) {
	pic := CreateNewCanvas(b.size, b.size)
	pic.SetLineWidth(1)
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.cell[i][j] == 0 {
				// if the cell is 0, draw black sqaure
				drawSquare(pic, i, j, "balck") 
			} else if b.cell[i][j] == 1 {
				// if the cell if 1, draw dark gray sqaure
				drawSquare(pic, i, j, "darkGray") 
			} else if b.cell[i][j] == 2 {
				// if the cell if 1, draw dark gray sqaure
				drawSquare(pic, i, j, "lightGray") 
			} else if b.cell[i][j] == 3 {
				// if the cell if 1, draw dark gray sqaure
				drawSquare(pic, i, j, "white") 
			}
		}
	}
	pic.SaveToPNG("sandpile.png")
}

/*===============================================================
 * Functions to draw square
 *==============================================================
 *  @param: pic Canvas
 *  @param: r   int
 *			row index
 *  @param  c   int
 *			coloum index
 *  @param  color string
 *			support red, blue, yellow, green
 */
func drawSquare(pic Canvas, r, c int, color string) {
	y1, x1 := float64(r), float64(c)
	y2, x2 := y1+1, x1+1
	if color == "balck" {
		pic.SetFillColor(MakeColor(0, 0, 0))
		pic.SetStrokeColor(MakeColor(0, 0, 0))
	} else if color == "darkGray" {
		pic.SetFillColor(MakeColor(85, 85, 85))
		pic.SetStrokeColor(MakeColor(85, 85, 85))
	} else if color == "lightGray" {
		pic.SetFillColor(MakeColor(170, 170, 170))
		pic.SetStrokeColor(MakeColor(170, 170, 170))
	} else if color == "white" {
		pic.SetFillColor(MakeColor(255, 255, 255))
		pic.SetStrokeColor(MakeColor(255, 255, 255))
	}
	pic.MoveTo(x1, y1)
	pic.LineTo(x1, y2)
	pic.LineTo(x2, y2)
	pic.LineTo(x2, y1)
	pic.LineTo(x1, y1)
	pic.Fill()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error: command should be: sandpile SIZE PILE")
		return
	}
	size, err := strconv.Atoi(os.Args[1]) // get board size
	if err != nil {
		fmt.Println("Error: Board size should be an integer")
		return
	}

	numOfSandpiles, err := strconv.Atoi(os.Args[2]) // get board size
	if err != nil {
		fmt.Println("Error: Number of sandpiles should be an integer")
		return
	}

	b := createBoard(size, numOfSandpiles)
	ComputeSteadyState(b)
	DrawBoard(b)
}