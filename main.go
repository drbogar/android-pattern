/*
The goal of this program is to find all possible patterns
that can be implemented on Android's well-known nine-point screen lock.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

/*
Contains is a function that determines whether the given array contains
the given integer.
*/
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type DotType int

const (
	CENTER DotType = iota
	CORNER
	SIDE
)

/*
0 1 2
3 4 5
6 7 8
*/

/*
getDotType tells you whether the given point with index is a corner point,
a side point or a center point.
*/
func getDotType(index int) DotType {
	x := int(index / 3)
	y := int(index % 3)
	/*
		The dexter diagonal would be the one from upper left to lower right, and
		the sinister diagonal the other one.
	*/
	onDexterDiagonal := x == y
	onSinisterDiagonal := (x + y) == 2
	if onDexterDiagonal && onSinisterDiagonal {
		return CENTER
	} else if onDexterDiagonal || onSinisterDiagonal {
		return CORNER
	} else {
		return SIDE
	}
}

/*
getMiddle tells you the index of the point between the point with index x and
the point with index y.
*/
func getMiddle(x int, y int) int {
	if x > y {
		return ((x - y) / 2) + y
	} else {
		return ((y - x) / 2) + x
	}
}

// The walker function walks the possible routes in a concurrent manner.
func walker(state []int, results chan<- []int, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(state) > 3 {
		results <- state
	}

	/*
		lastDot is the index of the last dot in the current walker function
		state array.
	*/
	lastDot := state[len(state)-1]

	switch getDotType(lastDot) {
	case CORNER:
		for i := 0; i < 9; i++ {
			if !contains(state, i) {
				ns := append(append(make([]int, 0), state...), i)
				if getDotType(i) != CORNER {
					wg.Add(1)
					go walker(ns, results, wg)
				} else {
					middleDot := getMiddle(lastDot, i)
					if contains(state, middleDot) {
						wg.Add(1)
						go walker(ns, results, wg)
					}
				}
			}
		}
	case SIDE:
		for i := 0; i < 9; i++ {
			if !contains(state, i) {
				ns := append(append(make([]int, 0), state...), i)
				if getDotType(i) != SIDE {
					wg.Add(1)
					go walker(ns, results, wg)
				} else if lastDot+i == 8 {
					if contains(state, 4) {
						wg.Add(1)
						go walker(ns, results, wg)
					}
				} else {
					wg.Add(1)
					go walker(ns, results, wg)
				}
			}
		}
	case CENTER:
		for i := 0; i < 9; i++ {
			if !contains(state, i) {
				ns := append(append(make([]int, 0), state...), i)
				wg.Add(1)
				go walker(ns, results, wg)
			}
		}
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	defer timeTrack(time.Now(), "main")
	// The wg is a WaitGroup, which indicates when the routes are finished.
	wg := &sync.WaitGroup{}
	nineFact := 9 * 8 * 7 * 6 * 5 * 4 * 3 * 2 * 1
	resultChanBuffSize :=
		nineFact +
			nineFact/1 +
			nineFact/2/1 +
			nineFact/3/2/1 +
			nineFact/4/3/2/1 +
			nineFact/5/4/3/2/1
	// The results channel collects possible unlocking patterns.
	results := make(chan []int, resultChanBuffSize)
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go walker([]int{i}, results, wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	i := 0
	f, _ := os.Create("dat1.csv")
	defer f.Close()
	for v := range results {
		i++
		f.WriteString(fmt.Sprint(i, ";", v, "\n"))
	}
	fmt.Println("resultCount", i)
}
