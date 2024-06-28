package main

import (
	"maps"
	"math"
	"slices"

	. "github.com/stoyan-kukev/team-project/backend/util"
)

func extractPoints(grid [][]square, products Set[int]) (items []point, checkouts []point, end point) {
	for y, row := range grid {
		for x, cell := range row {
			if cell.Kind == productSquare && products.Contains(cell.ProductID) { // Required products to be visited
				items = append(items, point{x, y})
			} else if cell.Kind.isCheckout() {
				checkouts = append(checkouts, point{x, y})
			} else if cell.Kind == endSquare {
				end = point{x, y}
			}
		}
	}
	return
}

func bfs(grid [][]square, start point, targetPoints Set[point]) ([][]float64, [][]point) {
	startingFromCheckout := grid[start.Y][start.X].Kind.isCheckout()
	directions := []point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	width := getWidth(grid)
	height := len(grid)

	distFromStart := makeGrid[float64](width, height)
	for y := range height {
		for x := range width {
			distFromStart[y][x] = math.Inf(1)
		}
	}

	prev := makeGrid[point](width, height)
	distFromStart[start.Y][start.X] = 0
	queue := []point{start}

	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		for _, dir := range directions {
			nextPos := point{pos.X + dir.X, pos.Y + dir.Y}
			if !(0 <= nextPos.X && nextPos.X < width && 0 <= nextPos.Y && nextPos.Y < height) {
				continue
			}
			if !(grid[nextPos.Y][nextPos.X].Kind == wallSquare || grid[nextPos.Y][nextPos.X].Kind == productSquare && !targetPoints.Contains(nextPos) || grid[nextPos.Y][nextPos.X].Kind.isCheckout() && startingFromCheckout) && distFromStart[nextPos.Y][nextPos.X] == math.Inf(1) {
				if targetPoints.Contains(pos) && targetPoints.Contains(nextPos) {
					continue // skip direct paths between items
				}
				distFromStart[nextPos.Y][nextPos.X] = distFromStart[pos.Y][pos.X] + 1
				prev[nextPos.Y][nextPos.X] = pos
				queue = append(queue, nextPos)
			}
		}
	}

	return distFromStart, prev
}

func reconstructPath(prev [][]point, start point, end point) []point {
	path := []point{}
	at := end
	for at != start {
		path = append(path, at)
		if prev[at.Y][at.X] == (point{0, 0}) {
			return nil
		}
		at = prev[at.Y][at.X]
	}
	path = append(path, start)
	slices.Reverse(path)
	return path
}

func createDistanceAndPathMatrix(grid [][]square, points []point) (distMatrix [][]float64, pathMatrix [][][]point) {
	distMatrix = makeGrid[float64](len(points), len(points))
	pathMatrix = makeGrid[[]point](len(points), len(points))

	for i, p := range points {
		pset := Set[point]{}
		for _, q := range points {
			pset.Insert(q)
		}
		dist, prev := bfs(grid, p, pset)
		for j, q := range points {
			if i != j {
				distMatrix[i][j] = dist[q.Y][q.X]
				pathMatrix[i][j] = reconstructPath(prev, p, q)
			}
		}
	}

	return
}

func tspWithOneCheckoutGreedy(distMatrix [][]float64, pathMatrix [][][]point, items []point, checkouts []point, entrance point, exit point) []point {
	allPoints := slices.Concat([]point{entrance}, items, checkouts, []point{exit})
	iBegin := 0
	iEnd := len(items) + len(checkouts) + 1

	visited := make([]bool, len(allPoints))
	visited[iBegin] = true
	visitedItems := 0
	visitedCheckout := false

	currentPos := iBegin
	pathOrder := []int{currentPos}

	for visitedItems < len(items) || !visitedCheckout {
		nextPos := -1
		minDist := math.Inf(1)

		for iPoint := range len(allPoints) {
			if iPoint == 0 {
				continue
			}
			if !visited[iPoint] {
				if iPoint <= len(items) || visitedItems == len(items) { // item or checkout
					if distMatrix[currentPos][iPoint] < minDist {
						minDist = distMatrix[currentPos][iPoint]
						nextPos = iPoint
					}
				}
			}
		}

		if nextPos == -1 {
			break
		}

		visited[nextPos] = true
		pathOrder = append(pathOrder, nextPos)
		currentPos = nextPos

		if nextPos <= len(items) {
			visitedItems += 1
		} else {
			visitedCheckout = true
		}
	}

	pathOrder = append(pathOrder, iEnd)

	fullPath := []point{}
	for i := range len(pathOrder) - 1 {
		pathToNextPoint := pathMatrix[pathOrder[i]][pathOrder[i+1]]
		fullPath = append(fullPath, pathToNextPoint[:len(pathToNextPoint)-1]...)
	}
	fullPath = append(fullPath, pathMatrix[pathOrder[len(pathOrder)-1]][iEnd]...)
	fullPath = append(fullPath, exit)

	return fullPath
}

func solve(grid [][]square, start point, products Set[int]) []point {
	items, checkouts, end := extractPoints(grid, products)
	points := slices.Concat([]point{start}, items, checkouts, []point{end})
	distMatrix, pathMatrix := createDistanceAndPathMatrix(grid, points)
	path := tspWithOneCheckoutGreedy(distMatrix, pathMatrix, items, checkouts, start, end)
	return path
}

func theAlgorithm(grid [][]square, start point, products Set[int]) []point {
	bestPath := []point{}
	// best_egg := -1
	for _, egg := range eggs {
		newRequired := maps.Clone(products)
		newRequired.Insert(egg)
		// grid = generate_grid(data_dict, new_required)
		path := solve(grid, start, newRequired)
		if len(path) < len(bestPath) || len(bestPath) == 0 {
			bestPath = path
			// best_egg = egg
		}
	}

	return bestPath
}
