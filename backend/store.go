package main

import (
	"slices"
)

type squareKind uint8

const (
	emptySquare squareKind = iota
	endSquare
	wallSquare
	productSquare
	checkoutSquare
)

type square struct {
	Kind         squareKind `json:"kind"`
	ProductID    int        `json:"productId"`
	CheckoutName string     `json:"checkoutName"`
}

var eggs = []int{204} // TODO: add others

func (s square) isEgg() bool {
	return s.Kind == productSquare && slices.Contains(eggs, s.ProductID)
}

func findRouteToOne(start point, isAcceptable func(square) bool) ([]point, square) {
	var explored [gridHeight][gridWidth]bool
	explored[start.Y][start.X] = true

	var prevSquare [gridHeight][gridWidth]point // only for visited non-start

	pendingSquares := []point{start}
	add := func(oldPos point, pos point) {
		if pos.X >= 0 && pos.X < gridWidth && pos.Y >= 0 && pos.Y < gridHeight && !explored[pos.Y][pos.X] {
			pendingSquares = append(pendingSquares, pos)
			explored[pos.Y][pos.X] = true
			prevSquare[pos.Y][pos.X] = oldPos
		}
	}

	for len(pendingSquares) > 0 {
		pos := pendingSquares[0]
		pendingSquares = pendingSquares[1:]

		square := grid[pos.Y][pos.X]
		if isAcceptable(square) {
			path := []point{}
			for pos != start {
				path = append(path, pos)
				pos = prevSquare[pos.Y][pos.X]
			}
			slices.Reverse(path)
			return path, square
		} else if pos == start || square.Kind == emptySquare || square.isEgg() {
			adjacentOffsets := []point{
				{-1, -1}, {0, -1}, {1, -1},
				{-1, 0}, {1, 0},
				{-1, 1}, {0, 1}, {1, 1},
			}
			for _, offset := range adjacentOffsets {
				add(pos, point{pos.X + offset.X, pos.Y + offset.Y})
			}
		}
	}

	panic("NO PATH?!")
}

func findRoute(products set[int]) []point {
	pos := start
	path := []point{}
	for len(products) > 0 {
		newPathSegment, productSquare := findRouteToOne(pos, func(s square) bool {
			return s.Kind == productSquare && products.contains(s.ProductID)
		})
		pos = newPathSegment[len(newPathSegment)-1]
		delete(products, productSquare.ProductID)
		path = append(path, newPathSegment...)
	}

	newPathSegment, _ := findRouteToOne(pos, func(s square) bool {
		return s.Kind == checkoutSquare
	})
	path = append(path, newPathSegment...)
	pos = newPathSegment[len(newPathSegment)-1]

	newPathSegment, _ = findRouteToOne(pos, func(s square) bool {
		return s.Kind == endSquare
	})
	path = append(path, newPathSegment...)
	return path
}
