package main

import (
	"encoding/csv"
	"io"
	"strconv"
)

type squareKind uint8

const (
	emptySquare         squareKind = 0
	endSquare           squareKind = 1
	wallSquare          squareKind = 2
	productSquare       squareKind = 3
	humanCheckoutSquare squareKind = 4
	selfCheckoutSquare  squareKind = 5
)

func (k squareKind) isCheckout() bool {
	return k == humanCheckoutSquare || k == selfCheckoutSquare
}

// blob encoding:
// 0000000000000000 empty
// 0001000000000000 end
// 0010000000000000 wall
// 0011pppppppppppp product p=productID
// 0100nnnnnnnnnnnn human checkout n=number
// 0101nnnnnnnnnnnn self checkout n=number

type square struct {
	Kind      squareKind `json:"kind"`
	ProductID int        `json:"productId"`
}

func makeGrid[x any](width int, height int) [][]x {
	grid := make([][]x, height)
	for y := range grid {
		grid[y] = make([]x, width)
	}
	return grid
}

func getWidth[x any](grid [][]x) int {
	if len(grid) == 0 {
		return 0
	}
	return len(grid[0])
}

func encodeGrid(grid [][]square) []byte {
	height := len(grid)
	width := getWidth(grid)
	bytes := make([]byte, height*width*2)
	for y, row := range grid {
		for x, sq := range row {
			word := uint16(sq.Kind) << 12
			if sq.Kind == productSquare {
				word |= uint16(sq.ProductID)
			}
			bytes[(y*width+x)*2] = byte(word)
			bytes[(y*width+x)*2+1] = byte(word >> 8)
		}
	}
	return bytes
}

func decodeGrid(bytes []byte, width int) [][]square {
	height := len(bytes) / 2 / width
	grid := make([][]square, height)
	for y := range height {
		grid[y] = make([]square, width)
		for x := range width {
			word := uint16(bytes[(y*width+x)*2]) | uint16(bytes[(y*width+x)*2+1])<<8
			number := int(word & 0xfff)
			grid[y][x] = square{Kind: squareKind(word >> 12), ProductID: number}
		}
	}
	return grid
}

func parseStoreCSV(input io.Reader) (grid [][]square, start point, err error) {
	reader := csv.NewReader(input)
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	minX := int(^uint(0) >> 1)
	maxX := 0
	minY := int(^uint(0) >> 1)
	maxY := 0
	for _, record := range records {
		var x, y int
		x, err = strconv.Atoi(record[1])
		if err != nil {
			return
		}
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}

		y, err = strconv.Atoi(record[2])
		if err != nil {
			return
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	width := maxX - minX + 1
	height := maxY - minY + 1
	grid = make([][]square, height)
	for y := range height {
		grid[y] = make([]square, width)
	}

	for _, record := range records {
		name := record[0]

		x, _ := strconv.Atoi(record[1])
		x -= minX
		y, _ := strconv.Atoi(record[2])
		y -= minY

		if name == "BL" {
			grid[y][x].Kind = wallSquare
		} else if name[0] == 'P' {
			var num int
			num, err = strconv.Atoi(record[0][1:])
			if err != nil {
				return
			}
			grid[y][x] = square{Kind: productSquare, ProductID: num}
		} else if name[0] == 'C' && name[1] == 'A' {
			grid[y][x] = square{Kind: selfCheckoutSquare}
		} else if name[0] == 'S' {
			grid[y][x] = square{Kind: humanCheckoutSquare}
		} else if name == "EX" {
			grid[y][x].Kind = endSquare
		} else if name == "EN" {
			start = point{x, y}
		}
	}

	return
}

var eggs = []int{170, 130, 240, 119, 239} // TODO: add others

// func (s square) isEgg() bool {
// 	return s.Kind == productSquare && slices.Contains(eggs, s.ProductID)
// }

// func findRouteToOne(start point, isAcceptable func(square) bool) ([]point, square) {
// 	var explored [gridHeight][gridWidth]bool
// 	explored[start.Y][start.X] = true

// 	var prevSquare [gridHeight][gridWidth]point // only for visited non-start

// 	pendingSquares := []point{start}
// 	add := func(oldPos point, pos point) {
// 		if pos.X >= 0 && pos.X < gridWidth && pos.Y >= 0 && pos.Y < gridHeight && !explored[pos.Y][pos.X] {
// 			pendingSquares = append(pendingSquares, pos)
// 			explored[pos.Y][pos.X] = true
// 			prevSquare[pos.Y][pos.X] = oldPos
// 		}
// 	}

// 	for len(pendingSquares) > 0 {
// 		pos := pendingSquares[0]
// 		pendingSquares = pendingSquares[1:]

// 		square := grid[pos.Y][pos.X]
// 		if isAcceptable(square) {
// 			path := []point{}
// 			for pos != start {
// 				path = append(path, pos)
// 				pos = prevSquare[pos.Y][pos.X]
// 			}
// 			slices.Reverse(path)
// 			return path, square
// 		} else if pos == start || square.Kind == emptySquare || square.isEgg() {
// 			adjacentOffsets := []point{
// 				{-1, -1}, {0, -1}, {1, -1},
// 				{-1, 0}, {1, 0},
// 				{-1, 1}, {0, 1}, {1, 1},
// 			}
// 			for _, offset := range adjacentOffsets {
// 				add(pos, point{pos.X + offset.X, pos.Y + offset.Y})
// 			}
// 		}
// 	}

// 	panic("NO PATH?!")
// }

// func findRoute(products set[int]) []point {
// 	pos := start
// 	path := []point{}
// 	for len(products) > 0 {
// 		newPathSegment, productSquare := findRouteToOne(pos, func(s square) bool {
// 			return s.Kind == productSquare && products.contains(s.ProductID)
// 		})
// 		pos = newPathSegment[len(newPathSegment)-1]
// 		delete(products, productSquare.ProductID)
// 		path = append(path, newPathSegment...)
// 	}

// 	newPathSegment, _ := findRouteToOne(pos, func(s square) bool {
// 		return s.Kind == checkoutSquare
// 	})
// 	path = append(path, newPathSegment...)
// 	pos = newPathSegment[len(newPathSegment)-1]

// 	newPathSegment, _ = findRouteToOne(pos, func(s square) bool {
// 		return s.Kind == endSquare
// 	})
// 	path = append(path, newPathSegment...)
// 	return path
// }
