package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/*===============================================================
 * Functions to manipulate a "field" of cells --- the main data
 * that must be managed by this program.
 *==============================================================*/

// The data stored in a single cell of a field
type Cell struct {
	kind  string
	score float64
}

// createField should create a new field of the ysize rows and xsize columns,
// so that field[r][c] gives the Cell at position (r,c).
func createField(rsize, csize int) [][]Cell {
	f := make([][]Cell, rsize)
	for i := range f {
		f[i] = make([]Cell, csize)
	}
	return f
}

// inField returns true iff (row,col) is a valid cell in the field
func inField(field [][]Cell, row, col int) bool {
	return row >= 0 && row < len(field) && col >= 0 && col < len(field[0])
}

// readFieldFromFile should open the given file and read the initial
// values for the field. The first line of the file will contain
// two space-separated integers saying how many rows and columns
// the field should have:
//    10 15
// each subsequent line will consist of a string of Cs and Ds, which
// are the initial strategies for the cells:
//    CCCCCCDDDCCCCCC
//
// If there is ever an error reading, this function should cause the
// program to quit immediately.
func readFieldFromFile(filename string) [][]Cell {
    // WRITE YOUR CODE HERE
    file, err := os.Open(filename)
    if err != nil {
    	fmt.Println("Error: something went wrong opening the file.")
    }
    scanner := bufio.NewScanner(file)
    // read the first line, get row and col information
    var rows, cols int
    if scanner.Scan() {
    	fmt.Sscanf(scanner.Text(), "%v %v", &rows, &cols)
    }
    //fmt.Println(rows)
    //fmt.Println(cols)
    var text string
    var cells [][]Cell = make([][]Cell, rows)
    for i := 0; i < rows; i++ {
    	cells[i] = make([]Cell, cols)
    	scanner.Scan()
    	text = scanner.Text()
    	//fmt.Println(text)
    	for  j := 0; j < cols; j++ {
    		//fmt.Print(string(text[j]))
    		cells[i][j] = Cell{string(text[j]), 0}
    	}
    	//fmt.Print("\n")
    }

    if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(3)
	}

	return cells // This is included only so this template will compile
}

// drawField should draw a representation of the field on a canvas and save the
// canvas to a PNG file with a name given by the parameter filename.  Each cell
// in the field should be a 5-by-5 square, and cells of the "D" kind should be
// drawn red and cells of the "C" kind should be drawn blue.
func drawField(field [][]Cell, filename string) {
    // WRITE YOUR CODE HERE
    rows := len(field)
    cols := len(field[0])
    //fmt.Println("rows: " + string(rows))
    //fmt.Println("cols: " + string(cols))
    pic := CreateNewCanvas(rows*5, cols*5)
 	pic.SetLineWidth(1)
    for i := 0; i < rows; i++ {
    	for j := 0; j < cols; j++ {
    		if field[i][j].kind == "C" {
			drawSquare(pic, i, j, "blue")  // if the cell is "C", draw blue sqaure
		} else {
			drawSquare(pic, i, j, "red")   // if the cell if "D", draw red sqaure
		}

    	}
	}
	pic.SaveToPNG(filename)
}

/* drawSquare 
 * @param: pic Canvas
 *  @param: r   int
 *			row index
 *  @param  c   int
 *			coloum index
 *  @param  default boolean
 *			true -- yellow
 *			false -- black
 */
func drawSquare(pic Canvas, r, c int, color string) {
	y1, x1 := float64(r*5), float64(c*5)
	y2, x2 := y1 + 5, x1 + 5
	if color == "red" {
		pic.SetFillColor(MakeColor(255,0,0))
		pic.SetStrokeColor(MakeColor(255,0,0))
	} else if color == "blue" {
		pic.SetFillColor(MakeColor(0,0,255))
		pic.SetStrokeColor(MakeColor(0,0,255))
	} else if color == "yellow" {
		pic.SetFillColor(MakeColor(255,255,0))
		pic.SetStrokeColor(MakeColor(255,255,0))
	} else if color == "green" {
		pic.SetFillColor(MakeColor(0,255,0))
		pic.SetStrokeColor(MakeColor(0,255,0))
	}
	pic.MoveTo(x1, y1)
	pic.LineTo(x1, y2)
	pic.LineTo(x2, y2)
	pic.LineTo(x2, y1)
	pic.LineTo(x1, y1)
	pic.Fill()
}

/*===============================================================
 * Functions to simulate the spatial games
 *==============================================================*/

// play a game between a cell of type "me" and a cell of type "them" (both me
// and them should be either "C" or "D"). This returns the reward that "me"
// gets when playing against them.
func gameBetween(me, them string, b float64) float64 {
	if me == "C" && them == "C" {
		return 1
	} else if me == "C" && them == "D" {
		return 0
	} else if me == "D" && them == "C" {
		return b
	} else if me == "D" && them == "D" {
		return 0
	} else {
		fmt.Println("type ==", me, them)
		panic("This shouldn't happen")
	}
}

// updateScores goes through every cell, and plays the Prisoner's dilema game
// with each of it's in-field nieghbors (including itself). It updates the
// score of each cell to be the sum of that cell's winnings from the game.
func updateScores(field [][]Cell, b float64) {
    // WRITE YOUR CODE HERE
    rows := len(field)
    cols := len(field[0])
    var sum float64
    var me string
    for i := 0; i < rows; i++ {
    	for j := 0; j < cols; j++ {
    		me = field[i][j].kind
    		sum = 0.0
    		for m := -1; m < 2; m++ {
    			for n := -1; n < 2; n++ {
    				if inField(field, i + m, j + n) {
    					sum += gameBetween(me, field[i + m][j + n].kind, b)
    				}
    			}
    		}
    		field[i][j].score = sum
    	}
    }
}

// updateStrategies create a new field by going through every cell (r,c), and
// looking at each of the cells in its neighborhood (including itself) and the
// setting the kind of cell (r,c) in the new field to be the kind of the
// neighbor with the largest score
func updateStrategies(field [][]Cell) [][]Cell {
    // WRITE YOUR CODE HERE
    rows := len(field)
    cols := len(field[0])
    var cells [][]Cell = make([][]Cell, rows)
    var maxScoreKind string
    for i := 0; i < rows; i++ {
    	cells[i] = make([]Cell, cols)
    	for  j := 0; j < cols; j++ {
 			maxScoreKind = getMaxScoreKind(field, i , j)
    		cells[i][j] = Cell{maxScoreKind, 0}
    	}
    	//fmt.Print("\n")
    }

	return cells // This is included only so this template will compile
}

func getMaxScoreKind(field [][]Cell, i, j int) string {
	var cellWithMaxScore Cell = Cell{"C", -1}
   	for m := -1; m < 2; m++ {
    	for n := -1; n < 2; n++ {
    		if inField(field, i + m, j + n) {
    			if field[i+m][j+n].score > cellWithMaxScore.score {
    				cellWithMaxScore = field[i+m][j+n]
    			}
    		}
    	}
    }
    return cellWithMaxScore.kind
}

// evolve takes an intial field and evolves it for nsteps according to the game
// rule. At each step, it should call "updateScores()" and the updateStrategies
func evolve(field [][]Cell, nsteps int, b float64) [][]Cell {
	for i := 0; i < nsteps; i++ {
		updateScores(field, b)
		field = updateStrategies(field)
	}
	return field
}

// evolve takes an intial field and evolves it for nsteps according to the game
// rule. At each step, it should call "updateScores()" and the updateStrategies
func evolveExtra(field [][]Cell, nsteps int, b float64) (prevField, currField [][]Cell) {
	rows := len(field)
	cols := len(field[0])
	prevField = make([][]Cell, rows)
	for i := 0; i < rows; i++ {
		prevField[i] = make([]Cell, cols)
	}

	for i := 0; i < nsteps; i++ {
		if i == nsteps - 1 {
			/*for m := 0; m < rows; m++ {
    			for  n := 0; n < cols; n++ {
    				prevField[m][n] = Cell{field[m][n].kind, 0}
    			}
   			}*/
   			copy(prevField, field)
		}
		updateScores(field, b)
		field = updateStrategies(field)
	}
	return prevField, field
}

// drawField should draw a representation of the field on a canvas and save the
// canvas to a PNG file with a name given by the parameter filename.  Each cell
// in the field should be a 5-by-5 square, and cells of the "D" kind should be
// drawn red and cells of the "C" kind should be drawn blue.
func drawFieldExtra(prevField, field [][]Cell, filename string) {
    // WRITE YOUR CODE HERE
    rows := len(field)
    cols := len(field[0])
    //fmt.Println("prev rows: " + string(len(prevField)))
    //fmt.Println("prev cols: " + string(len(prevField[0])))
    pic := CreateNewCanvas(rows*5, cols*5)
 	pic.SetLineWidth(1)
    for i := 0; i < rows; i++ {
    	for j := 0; j < cols; j++ {
    		if prevField[i][j].kind == "C" && field[i][j].kind == "C"  {
				drawSquare(pic, i, j, "blue")  // if the cell is "C", draw blue sqaure
			} else if prevField[i][j].kind == "D" && field[i][j].kind == "D"  {
				drawSquare(pic, i, j, "red")  // if the cell is "C", draw blue sqaure
			} else if prevField[i][j].kind == "C" && field[i][j].kind == "D"  {
				drawSquare(pic, i, j, "yellow")  // if the cell is "C", draw blue sqaure
			} else if prevField[i][j].kind == "D" && field[i][j].kind == "C"  {
				drawSquare(pic, i, j, "green")  // if the cell is "C", draw blue sqaure
			}
		}
	}
	pic.SaveToPNG(filename)
}


// Implements a Spatial Games version of prisoner's dilemma. The command-line
// usage is:
//     ./spatial field_file b nsteps
// where 'field_file' is the file continaing the initial arrangment of cells, b
// is the reward for defecting against a cooperator, and nsteps is the number
// of rounds to update stategies.
//
func main() {
	// parse the command line
	if len(os.Args) != 4 {
		fmt.Println("Error: should spatial field_file b nsteps")
		return
	}

	fieldFile := os.Args[1]

	b, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil || b <= 0 {
		fmt.Println("Error: bad b parameter.")
		return
	}

	nsteps, err := strconv.Atoi(os.Args[3])
	if err != nil || nsteps < 0 {
		fmt.Println("Error: bad number of steps.")
		return
	}

    // read the field
	field := readFieldFromFile(fieldFile)
    fmt.Println("Field dimensions are:", len(field), "by", len(field[0]))

    // evolve the field for nsteps and write it as a PNG
	newfield := evolve(field, nsteps, b)
	drawField(newfield, "Prisoners.png")

	// evolve the field for nsteps and write it as a PNG
	prevField, currField := evolveExtra(field, nsteps, b)
	drawFieldExtra(prevField, currField, "PrisonersExtra.png")
}