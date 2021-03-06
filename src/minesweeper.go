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

// board struct holds all data members necessary to describe the board and
// the actual state of the game.
// eBoard - encoded board, actual solution
// dBoard - decoded board, represents current game board with
type board struct {
	side            int
	numMines        int
	remainingFields int
	eBoard          [10][10]int
	dBoard          [10][10]int
	gConst          gameConstants
}

// game constant variables
type gameConstants struct {
	mine         int
	undiscovered int
	flag         int
}

func main() {
	playGame()
}

// fillDecodedBoard initializes decoded board to all fields containing
// undiscovered value.
func (b *board) fillDecodedBoard() {
	for i := 0; i < b.side; i++ {
		for j := 0; j < b.side; j++ {
			(*b).dBoard[i][j] = b.gConst.undiscovered
		}
	}
}

// shuffleMines place randomly mines on the board.
func (b *board) shuffleMines() {
	// create truly random number generator using current time as seed
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// generate random numbers
	a := make([]int, 100)
	for i := range a {
		a[i] = i
	}
	r.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

	// assign mines to the encoded board
	mines := a[:b.numMines]
	for _, val := range mines {
		(*b).eBoard[val/b.side][val%b.side] = b.gConst.mine
	}
}

// fillEncodedBoard calculates values of each field of the board based on the
// position of the mines.
func (b *board) fillEncodedBoard() {
	for i := 0; i < b.side; i++ {
		for j := 0; j < b.side; j++ {
			// skip fields were are mines
			if (*b).eBoard[i][j] == b.gConst.mine {
				continue
			}
			total := 0
			// check all 8 neighbouring fields for existence of mines
			if i-1 >= 0 && j-1 >= 0 && (*b).eBoard[i-1][j-1] == b.gConst.mine {
				total++
			}
			if i-1 >= 0 && (*b).eBoard[i-1][j] == b.gConst.mine {
				total++
			}
			if i-1 >= 0 && j+1 < (*b).side && (*b).eBoard[i-1][j+1] == b.gConst.mine {
				total++
			}
			if j-1 >= 0 && (*b).eBoard[i][j-1] == b.gConst.mine {
				total++
			}
			if j+1 < (*b).side && (*b).eBoard[i][j+1] == b.gConst.mine {
				total++
			}
			if i+1 < (*b).side && j-1 >= 0 && (*b).eBoard[i+1][j-1] == b.gConst.mine {
				total++
			}
			if i+1 < (*b).side && (*b).eBoard[i+1][j] == b.gConst.mine {
				total++
			}
			if i+1 < (*b).side && j+1 < (*b).side && (*b).eBoard[i+1][j+1] == b.gConst.mine {
				total++
			}
			(*b).eBoard[i][j] = total
		}
	}
}

// setUp fills in both decoded and encoded boards with initial values.
func (b *board) setUp() {
	(*b).fillDecodedBoard()
	(*b).shuffleMines()
	(*b).fillEncodedBoard()
}

// printBoard prints the board.
func printBoard(board [10][10]int, constants gameConstants) {
	println("   A B C D E F G H I J")
	for i := 0; i < len(board); i++ {
		str := fmt.Sprint(i) + " "
		for j := 0; j < len(board[0]); j++ {
			switch board[i][j] {
			case constants.mine:
				str += " M"
			case constants.undiscovered:
				str += " -"
			case constants.flag:
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

// userInput asks for the user input and validates it.
// Returns valide user input.
func (b *board) userInput() []string {
	// create an instance of buffered I/O
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("-----------------------------------------")
		fmt.Println("To quit the game press Q and press Enter.")
		fmt.Println("Please enter a field to undiscover (eg. A0) or set/unset a flag (eg. A0 F): ")
		fmt.Print("-> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)
		// converting to upper letters simplifies necessary number of comparisons
		input = strings.ToUpper(input)
		words := strings.Split(input, " ")
		// validation of the input
		if (len(words) == 0 || len(words) > 2 || len(words[0]) > 2) || (len(words) == 2 && len(words[1]) != 1) {
			fmt.Println("Invalid input!")
			continue
		}
		first := words[0]
		if len(first) == 1 && strings.Compare(first, "Q") == 0 {
			b.gameExit(false)
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

// gameExit exits the game loop.
// If result is true, congratulions are displayed on the screen.
func (b *board) gameExit(result bool) {
	if result {
		fmt.Println("CONGRATULATION YOU WON!!!")
	}
	fmt.Println("---------------------")
	fmt.Println("Thank you very much for the game!")
	fmt.Println("---------------------")
	fmt.Println("Solution:")
	printBoard(b.eBoard, b.gConst)
	fmt.Println("---------------------")
	for {
		fmt.Println("Please press any key to leave the game.")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if len(input) != 0 {
			os.Exit(0)
		}
	}
}

// inputConverter converts users input into position in the matrix.
// Returns (x, y) where x is the number of row and y is the number of column.
func inputConverter(field string) (int, int) {
	startValue := int('A')
	f := int(field[0])
	// strconv.Atoi can be used when string is a number representation
	n, _ := strconv.Atoi(field[1:])
	return f - startValue, n
}

// checkValue verifies if user guess is correct or not.
// If user guess is incorrect and falls into mine the gameExit is called.
func (b *board) checkValue(row int, col int, checkFlag bool, print bool) {
	if checkFlag {
		if (*b).dBoard[row][col] == b.gConst.flag {
			(*b).dBoard[row][col] = b.gConst.undiscovered
			fmt.Println("Flag was added!")
		} else {
			(*b).dBoard[row][col] = b.gConst.flag
			fmt.Println("Flag was removed!")
		}
	} else {
		(*b).dBoard[row][col] = (*b).eBoard[row][col]
		if (*b).dBoard[row][col] != b.gConst.mine {
			(*b).remainingFields -= 1
			if print {
				fmt.Println("Good pick!")
			}
			// update neighbouring fields if the guess field value is zero
			if (*b).dBoard[row][col] == 0 {
				b.updateNeihbours(row, col)
			}
			// check winning conditions
			if (*b).remainingFields == 0 {
				b.gameExit(true)
			}
		} else {
			fmt.Println("MINE!!! The game lost")
			fmt.Println("Solution:")
			printBoard(b.eBoard, b.gConst)
			b.gameExit(false)
		}
	}
}

// updateNeighbours displays values of neihbouring fields if the guessed field
// value equals zero.
func (b *board) updateNeihbours(row int, col int) {
	toCheck := []int{(row * b.side) + col}
	for len(toCheck) != 0 {
		i, j := toCheck[0]/b.side, toCheck[0]%b.side
		n, m := i-1, j-1
		if n >= 0 && m >= 0 && (*b).dBoard[n][m] == b.gConst.undiscovered {
			b.checkValue(n, m, false, false)
			if (*b).eBoard[n][m] == 0 {
				toCheck = append(toCheck, (n*b.side)+m)
			}
		}
		n, m = i-1, j
		if n >= 0 && (*b).dBoard[n][m] == b.gConst.undiscovered {
			b.checkValue(n, m, false, false)
			if (*b).eBoard[n][m] == 0 {
				toCheck = append(toCheck, (n*b.side)+m)
			}
		}
		n, m = i-1, j+1
		if n >= 0 && m < b.side && (*b).dBoard[n][m] == b.gConst.undiscovered {
			b.checkValue(n, m, false, false)
			if (*b).eBoard[n][m] == 0 {
				toCheck = append(toCheck, (n*b.side)+m)
			}
		}
		n, m = i, j-1
		if m >= 0 && (*b).dBoard[n][m] == b.gConst.undiscovered {
			b.checkValue(n, m, false, false)
			if (*b).eBoard[n][m] == 0 {
				toCheck = append(toCheck, (n*b.side)+m)
			}
		}
		n, m = i, j+1
		if m < b.side && (*b).dBoard[n][m] == b.gConst.undiscovered {
			b.checkValue(n, m, false, false)
			if (*b).eBoard[n][m] == 0 {
				toCheck = append(toCheck, (n*b.side)+m)
			}
		}
		n, m = i+1, j-1
		if n < b.side && m >= 0 && (*b).dBoard[n][m] == b.gConst.undiscovered {
			b.checkValue(n, m, false, false)
			if (*b).eBoard[n][m] == 0 {
				toCheck = append(toCheck, (n*b.side)+m)
			}
		}
		n, m = i+1, j
		if n < b.side && (*b).dBoard[n][m] == b.gConst.undiscovered {
			b.checkValue(n, m, false, false)
			if (*b).eBoard[n][m] == 0 {
				toCheck = append(toCheck, (n*b.side)+m)
			}
		}
		n, m = i+1, j+1
		if n < b.side && m < b.side && (*b).dBoard[n][m] == b.gConst.undiscovered {
			b.checkValue(n, m, false, false)
			if (*b).eBoard[n][m] == 0 {
				toCheck = append(toCheck, (n*b.side)+m)
			}
		}
		toCheck = append(toCheck[1:])
	}
}

// playGame starts new minesweeper game.
func playGame() {
	side, numMines := 10, 10
	// create instance of minesweeper board and initialize it
	b := board{
		side:            side,
		numMines:        numMines,
		remainingFields: side*side - numMines,
		eBoard:          [10][10]int{},
		dBoard:          [10][10]int{},
		gConst: gameConstants{
			mine:         -1,
			undiscovered: -2,
			flag:         -3,
		},
	}
	b.setUp()

	fmt.Println("-----------------------------------------")
	fmt.Println("WELCOME IN MINESWEEPER")
	fmt.Println("Prepare for a lot of fun! :)\n")
	// Infinite game loop
	for {
		fmt.Println("Current board:")
		printBoard(b.dBoard, b.gConst)
		field := b.userInput()
		col, row := inputConverter(field[0])
		var checkFlag bool = false
		if b.dBoard[row][col] != b.gConst.undiscovered && b.dBoard[row][col] != b.gConst.flag {
			fmt.Println("The field", field[0], "was already revealed")
			continue
		}
		if len(field) == 2 {
			checkFlag = true
		}
		b.checkValue(row, col, checkFlag, true)
	}
}
