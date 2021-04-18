package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type board struct {
	side            int
	numMines        int
	remainingFields int
	eBoard          [10][10]int // encoded board
	dBoard          [10][10]int // decoded board

}

func main() {
	playGame()
}

const mine int = -1
const undiscovered int = -2
const flag int = -3

func (b *board) setUp() {
	(*b).fillDecodedBoard()
	(*b).shuffleMines()
	(*b).fillEncodedBoard()
}

func (b *board) shuffleMines() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	a := make([]int, 100)
	for i := range a {
		a[i] = i
	}
	r.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	mines := a[:b.numMines+1]
	for _, val := range mines {
		(*b).eBoard[val/b.side][val%b.side] = mine
	}
}

func (b *board) fillDecodedBoard() {
	var i int = 0
	for i < b.side {
		var j int = 0
		for j < b.side {
			(*b).dBoard[i][j] = undiscovered
			j++
		}
		i++
	}
}

func (b *board) fillEncodedBoard() {
	for i := 0; i < b.side; i++ {
		for j := 0; j < b.side; j++ {
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

func userInput() []string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("-----------------------------------------")
		fmt.Println("To quit the game press Q and press Enter.")
		fmt.Println("Please enter a field to undiscover (eg. A0) or set/unset a flag (eg. A0 F): ")
		fmt.Print("-> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)
		input = strings.ToUpper(input)
		words := strings.Split(input, " ")
		if (len(words) == 0 || len(words) > 2 || len(words[0]) > 2) || (len(words) == 2 && len(words[1]) != 1) {
			fmt.Println("Invalid input!")
			continue
		}
		first := words[0]
		if len(first) == 1 && strings.Compare(first, "Q") == 0 {
			gameExit(false)
		}
		if len(first) == 2 && first[:1] >= "A" && first[:1] <= "J" && first[1:2] >= "0" && first[1:2] <= "9" {
			if len(words) == 2 {
				if words[1] != "F" {
					fmt.Println("Invalid input!")
					continue
				}
			}
			return words
		}
	}
}

// how to set default values
func gameExit(result bool) {
	if result {
		fmt.Println("Congratulation you won!!!")
	}
	fmt.Println("---------------------")
	fmt.Println("Thank you very much for the game!")
	for {
		fmt.Println("Please press any key to leave the game.")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if len(input) != 0 {
			os.Exit(0)
		}
	}
}

func playGame() {
	side, numMines := 10, 20
	// create instance of minesweeper board and initialize it
	b := board{
		side:            side,
		numMines:        numMines,
		remainingFields: side*side - numMines,
		eBoard:          [10][10]int{},
		dBoard:          [10][10]int{},
	}
	b.setUp()

	fmt.Println("-----------------------------------------")
	fmt.Println("WELCOME IN MINESWEEPER")
	fmt.Println("Prepare for a lot of fun! :)\n")
	for {
		fmt.Println("Current board:")
		printBoard(b.dBoard)
		field := userInput()
		col, row := inputConverter(field[0])
		var checkFlag bool = false
		if b.dBoard[row][col] != undiscovered && b.dBoard[row][col] != flag {
			fmt.Println("The field", field[0], "was already revealed")
			continue
		}
		if len(field) == 2 {
			checkFlag = true
		}
		b.checkValue(row, col, checkFlag)
	}
}

func inputConverter(field string) (int, int) {
	startValue := int('A')
	f := int(field[0])
	n, _ := strconv.Atoi(field[1:])
	return f - startValue, n
}

func (b *board) checkValue(row int, col int, checkFlag bool) {
	if checkFlag {
		if (*b).dBoard[row][col] == flag {
			(*b).dBoard[row][col] = undiscovered
			fmt.Println("Flag was added!")
		} else {
			(*b).dBoard[row][col] = flag
			fmt.Println("Flag was removed!")
		}
	} else {
		(*b).dBoard[row][col] = (*b).eBoard[row][col]
		if (*b).dBoard[row][col] != mine {
			(*b).remainingFields -= 1
			fmt.Println("Good pick!")
			if (*b).remainingFields == 0 {
				gameExit(true)
			}
		} else {
			fmt.Println("MINE!!! The game lost")
			fmt.Println("Solution:")
			printBoard(b.eBoard)
			gameExit(false)
		}
	}
}
