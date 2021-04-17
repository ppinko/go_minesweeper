package main

import (
	"fmt"
	"math/rand"
	"time"
)

type board struct {
	eBoard [10][10]int // encoded board
	dBoard [10][10]int // decoded board
}

func main() {
	var b board
	b.setUp()
	printBoard(b.eBoard)
}

const mine int = -1
const undiscovered int = -2
const flag int = -3

func (b *board) setUp() {
	b.fillDecodedBoard()
	b.shuffleMines()
	b.fillEncodedBoard()
}

func (b *board) shuffleMines() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	a := make([]int, 100)
	for i := range a {
		a[i] = i
	}
	r.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	mines := a[:21]
	for _, val := range mines {
		(*b).eBoard[val/10][val%10] = mine
	}
}

func (b *board) fillDecodedBoard() {
	var i int = 0
	for i < 10 {
		var j int = 0
		for j < 10 {
			(*b).dBoard[i][j] = undiscovered
			j++
		}
		i++
	}
}

func (b *board) fillEncodedBoard() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if (*b).eBoard[i][j] == mine {
				continue
			}
			total := 0
			if i-1 >= 0 && j-1 >= 0 && (*b).eBoard[i-1][j-1] == mine {
				total++
			}
			if i-1 >= 0 && (*b).eBoard[i-1][j] == mine {
				total++
			}
			if i-1 >= 0 && j+1 < 10 && (*b).eBoard[i-1][j+1] == mine {
				total++
			}
			if j-1 >= 0 && (*b).eBoard[i][j-1] == mine {
				total++
			}
			if j+1 < 10 && (*b).eBoard[i][j+1] == mine {
				total++
			}
			if i+1 < 10 && j-1 >= 0 && (*b).eBoard[i+1][j-1] == mine {
				total++
			}
			if i+1 < 10 && (*b).eBoard[i+1][j] == mine {
				total++
			}
			if i+1 < 10 && j+1 < 10 && (*b).eBoard[i+1][j+1] == mine {
				total++
			}
			(*b).eBoard[i][j] = total
		}
	}
}

func printBoard(board [10][10]int) {
	println("   A B C D E F G H I J")
	for i := 0; i < 10; i++ {
		var str string
		str += fmt.Sprint(i) + " "
		for j := 0; j < 10; j++ {
			switch board[i][j] {
			case mine:
				str += " M"
			case undiscovered:
				str += " -"
			case flag:
				str += " F"
			case 0:
				str += "  "
			case 1, 2, 3, 4, 5, 6, 7, 8:
				str += " " + fmt.Sprint(board[i][j])
			}
		}
		fmt.Println(str)
	}
}
