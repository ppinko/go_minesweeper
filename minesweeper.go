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
	b.shuffle()
	fmt.Println(b.dBoard)
	fmt.Println(b.eBoard)
}

func (b *board) setUp() {
	// initial all values to unknown (-2)
	var i int = 0
	for i < 10 {
		var j int = 0
		for j < 10 {
			(*b).dBoard[i][j] = -2 // -2 equals to '-' (not discovered yet)
			j++
		}
		i++
	}
}

func (b *board) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	a := make([]int, 100)
	for i := range a {
		a[i] = i
	}
	r.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	mines := a[:21]
	for _, val := range mines {
		(*b).eBoard[val/10][val%10] = -1
	}
}
