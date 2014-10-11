package main

import (
	"os"
	"fmt"
	"strconv"
	"errors"
	"math"
	"math/rand"
)

func main() {
    if len(os.Args) != 4 {
		fmt.Println("Error: worng number of arguments")
		return
	}

	r, err := strconv.ParseFloat(os.Args[1], 64)
	if nil != err || r < 0 {
		fmt.Println("Error: Invalid population rate " + os.Args[1])
	}

	stepSize, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil || stepSize <= 0 {
		fmt.Println("Error: Stepsize is not a valid number " + os.Args[2])
		return
	}

	rule, err := parseRule(os.Args[3]) // rule for Cellular automata, a slice of length 8
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	drawPopSize(r)
	drawRandomWalk(stepSize)
	drawCelluarAutomata(rule)
	drawCoolPicture()
}


/*===========================================================================
 *  Draw PopSize picture
 *=========================================================================*/
 func  drawPopSize(r float64) {
 	pic := CreateNewCanvas(500,100)
    pic.SetStrokeColor(MakeColor(0,0,255))
    pic.SetLineWidth(1)
 	x_0 := 0.1
 	max_t := 100
 	popsize := PopSize(r, x_0, max_t)
 	pic.MoveTo(5.0*0, 100 - 100*popsize[0])
 	for i := 1; i < len(popsize); i++ {
 		pic.LineTo(5*float64(i), 100 - 100*popsize[i])
 	}
 	pic.Stroke()
 	pic.SaveToPNG("PopSize.png")
 }

/**
 * Growth of a Population
 *
 * The size at time t of a population with a birth rate r can be modeled as:
 *
 *      x_t = r*x_{t-1}(1 - x_{t-1})
 * return slice of population number
 */
 func PopSize(r, x_0 float64, max_t int) []float64 {
 	var popsize []float64
 	popsize = make([]float64, 0)
 	for max_t != 0 {
 		popsize = append(popsize, x_0)
	 	x_1 := r * x_0 * (1 - x_0);
	 	if x_1 < 0 {
	 		x_1 = 0.0
	 	}
	 	if x_1 > 1 {
	 		x_1 = 1.0
	 	}
	 	x_0 = x_1
	 	//popsize = append(popsize, x_0)
	 	max_t--
 	}
 	//fmt.Println(popsize)
 	return popsize
 }




/*=========================================================================
 *  Draw RandomWalk picture
 *========================================================================*/
func drawRandomWalk(stepSize float64) {
	var width float64 = 500
	var height float64 = 500
	pic := CreateNewCanvas(500, 500)
	pic.SetStrokeColor(MakeColor(0,0,0))
	pic.SetLineWidth(1)
	
	var steps int = 1000
	var seed int64 = 12345
	rand.Seed(seed) //initialize ranodm
	var x, y = width/2, height/2
	pic.MoveTo(x ,y)
 	for i := 0; i < steps; i++ {
 		x,y = randStep(x, y, width, height, stepSize)
 		pic.LineTo(x, y)
 	}
 	pic.Stroke()
 	pic.SaveToPNG("RandomWalk.png")
}


// generate a new step
func randStep(x, y, width, height, stepSize float64) (nx, ny float64) {
	var deltaX, deltaY float64// store the newly generated moving distance in x and y direction
	nx = x
	ny = y
	for (math.Abs(nx - x) <= math.SmallestNonzeroFloat64 && math.Abs(ny - y) <= math.SmallestNonzeroFloat64) || !inField(nx, width) || !inField(ny, height) {
	// if new destination is same as the original one or out of boards
	// generate a new destination repeatly
		deltaX, deltaY = randDelta(stepSize)
		nx, ny = x + deltaX, y + deltaY
	}
	return
}

/*generate relative moving distance*/
func randDelta(stepSize float64) (deltaX, deltaY float64) {
	theta := rand.Float64() * 2 * math.Pi // range [0, 2*Pi)
	deltaX = stepSize * math.Cos(theta)
	deltaY = stepSize * math.Sin(theta)
	return deltaX, deltaY
}

/* check whether the new location has go over board*/
func inField(coord, board float64) bool {
	return coord >= 0 && coord < board
}



/*=========================================================================
 *  Draw Cellular Automata picture
 *========================================================================*/
 func drawCelluarAutomata(rule []int) {
 	var width int = 100
 	pic := CreateNewCanvas(width*5, 255)
 	pic.SetLineWidth(1)
	var prevLv, currLv, temp []int
	prevLv = make([]int, width) // previous level, states of cells at time t-1
	prevLv[width/2] = 1
	currLv = make([]int, width) // current level, states of cells at time t
	drawCells(pic, prevLv, 0)
	var steps = 50
	for i := 0; i < steps; i++ {
		currLv = calcCellState(rule, prevLv, currLv) // update states of cells at time t
		drawCells(pic, currLv, i + 1)
		temp = prevLv // update previous level
		prevLv = currLv
		currLv = temp
	}
	pic.SaveToPNG("CA.png")
}

/*
* draw all cell states at time t
* @param: cells []int
*           array that represents current states of the cells
*           "#" stands for filled
*           " " (white space) represents empty
*/
func drawCells(pic Canvas, cells []int, row int) {
	for i, cell := range cells {
		if cell == 1 {
			drawSquare(pic, row, i, false)  // if the cell is filled, draw black sqaure
		} else {
			drawSquare(pic, row, i, true)   // if the cell if empty, draw yellow sqaure
		}
	}
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
func drawSquare(pic Canvas, r, c int, defaultColor bool) {
	y1, x1 := float64(r*5), float64(c*5)
	y2, x2 := y1 + 5, x1 + 5
	if defaultColor {
		pic.SetFillColor(MakeColor(255,255,0))
		pic.SetStrokeColor(MakeColor(225,255,0))
	} else {
		pic.SetFillColor(MakeColor(0,0,0))
		pic.SetStrokeColor(MakeColor(0,0,0))
	}
	pic.MoveTo(x1, y1)
	pic.LineTo(x1, y2)
	pic.LineTo(x2, y2)
	pic.LineTo(x2, y1)
	pic.LineTo(x1, y1)
	pic.Fill()
}

 /* parse the rule provide by string
*  @param: str string
*           input string in binary/dec format
*           Default format : 0/1 binary string
*  return: slice to store the rule, 1--full 0--empty  
*/
func parseRule(str string) ([]int, error) {
	var rule = make([]int, 0)
	var err error = nil
	num, _err := strconv.ParseInt(str, 2, 32)  // try to parse the string as binary string
	if ( _err != nil) {
		num, _err = strconv.ParseInt(str, 10, 32) // if error occurs, re-prase it as decimal string
		if _err != nil {
			return rule, _err
		}
	} else { // it's a binary string, check its length
		if(len([]rune(str)) != 8) {
			err = errors.New("RULE is not valid")
		}
	}

	if num < 0  || num > 255 {
		err = errors.New("RULE is not valid")
	}
	// now index is a decimal number, convert it into an array which simulates a binary number
	var index int  = 0
	for index < 8 {
		rule = append(rule, int(num%2))
		num = num / 2
		index++
	}

	// reverse []rule
	var reversed_rule = make([]int, 8)
	for index = 0; index < 8; index++ {
		reversed_rule[index] = rule[7 - index]
	}
	return reversed_rule, err
}

/*
* calculate next state of cells from the state of previous state
* @param: rule
*           rule that produce next state
*         prevLv
*           cell states of time t-1
*         currLv
*           cell states of time t
* return: current state of all cells
*/
func calcCellState(rule, prevLv, currLv []int) []int {
	var num int = 0  // convert states of 3 cells into a single decimal number
	for i := 0; i < len(prevLv); i++ {
		if i == 0 && i != len(prevLv) - 1 { // edge case : left boundary
			num = prevLv[i]*10 + prevLv[i + 1]  
		} else if  i == len(prevLv) - 1 && i != 0 { // edge case : right boundary
			num = prevLv[i - 1]*100 + prevLv[i]*10
		} else { 	// normal case
			num = prevLv[i - 1]*100 + prevLv[i]*10 + prevLv[i + 1]
		}
		switch num {
			case 111 : currLv[i] = rule[0]
			case 110 : currLv[i] = rule[1]
			case 101 : currLv[i] = rule[2]
			case 100 : currLv[i] = rule[3]
			case 11 : currLv[i] = rule[4]
			case 10 : currLv[i] = rule[5]
			case 1 : currLv[i] = rule[6]
			case 0 : currLv[i] = rule[7]
		}
	}
	//printCells(currLv)
	return currLv
}


/*=========================================================================
 *  Draw interesting picture
 *========================================================================*/
func drawCoolPicture() {
	width, height := 500, 500
	pic := CreateNewCanvas(width, height)
	pic.SetStrokeColor(MakeColor(0,0,0))
	pic.SetLineWidth(1)

	pic.MoveTo(calcPostion(0, width, height))
	var steps int = 100
	for i := 0; i < steps; i++ {
		var theta float64 = 2 * math.Pi / float64(steps) * float64(i)
		pic.LineTo(calcPostion(theta, width, height))
	}
	pic.LineTo(calcPostion(0, width, height))
	pic.SetFillColor(MakeColor(255,0,0))
	pic.Fill()
	pic.SaveToPNG("MyCoolPicture.png")
}

func calcPostion(theta float64, width, height int) (x, y float64) {
	x = 8*(16*math.Pow(math.Sin(theta), 3)) + float64(width/2)
	y = 8*(13*math.Cos(theta) - 5*math.Cos(2*theta)-2*math.Cos(3*theta)-math.Cos(4*theta)) + float64(height/2)
	return float64(width) - x, float64(height) - y
}