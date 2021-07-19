package main

import (
	"fmt"
	"github.com/oddg/hungarian-algorithm"
)

func main() {

	//Eredeti maátrix és a megoldás \-el jelölve
	//parcelID: 1  2  3  4
	//				  ------------------
	//droneID: 1  |	\4  2  5  7
	//droneID: 2  |	 8 \3  10 8
	//droneID: 3  |	 12 5  4 \5
	//droneID: 4  |	 6  3 \7  14
	b := [][]int{{4, 8, 12, 6}, {2, 3, 5, 4}, {5, 10, 4, 7}, {7, 8, 5, 14}}
	fmt.Println(hungarianAlgorithm.Solve(b))
	//A megoldása az algoritmus szerint: [0 1 3 2], ami az indexeket jelöli.
	//A megoldás tábla:
	//parcelID: 1  2  3  4
	//				  ------------------
	//droneID: 1  |	\0  1  0  1
	//droneID: 2  |	 2 \0  3  0
	//droneID: 3  |	 9  5  0 \0
	//droneID: 4  |	 6  3 \7  14

}
