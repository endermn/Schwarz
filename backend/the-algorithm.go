package main

import (
	"maps"
	"math"
	"slices"
)

func extractPoints(grid [][]square, products set[int]) ([]point, []point, point) {
	items := []point{}
	checkouts := []point{}
	var end point
	for y, row := range grid {
		for x, cell := range row {
			if cell.Kind == productSquare && products.contains(cell.ProductID) { // Required products to be visited
				items = append(items, point{x, y})
			} else if cell.Kind.isCheckout() {
				checkouts = append(checkouts, point{x, y})
			} else if cell.Kind == endSquare {
				end = point{x, y}
			}
		}
	}
	return items, checkouts, end
}

func bfs(grid [][]square, start point, pointsWithItems set[point]) ([][]float64, [][]point) {
	directions := []point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	width := getWidth(grid)
	height := len(grid)
	dist := makeGrid[float64](width, height)
	for y, r := range dist {
		for x := range r {
			dist[y][x] = math.Inf(1)
		}
	}
	prev := makeGrid[point](width, height)
	dist[start.Y][start.X] = 0
	queue := []point{start}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, d := range directions {
			n := point{p.X + d.X, p.Y + d.Y}
			if 0 <= n.X && n.X < width && 0 <= n.Y && n.Y < height && !(grid[n.Y][n.X].Kind == wallSquare || grid[n.Y][n.X].Kind == productSquare && !pointsWithItems.contains(n)) && dist[n.Y][n.X] == math.Inf(1) {
				if pointsWithItems.contains(p) && pointsWithItems.contains(n) {
					continue // skip direct paths between items
				}
				dist[n.Y][n.X] = dist[p.Y][p.X] + 1
				prev[n.Y][n.X] = p
				queue = append(queue, n)
			}
		}
	}

	return dist, prev
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
	width := getWidth(grid)
	height := len(grid)

	distMatrix = makeGrid[float64](len(points), len(points))
	pathMatrix = makeGrid[[]point](width, height)

	for i, p := range points {
		pset := set[point]{}
		for _, q := range points {
			pset.insert(q)
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

func tsp_with_one_checkout_greedy(dist_matrix [][]float64, path_matrix [][][]point, items []point, checkouts []point, entrance point, exit point) []point {
	n_items := len(items)
	entrance_idx := 0
	exit_idx := len(items) + len(checkouts) + 1

	all_points := slices.Concat([]point{entrance}, items, checkouts, []point{exit})
	n := len(all_points)

	visited := make([]bool, n)
	visited[entrance_idx] = true
	visited_items := 0
	visited_checkout := false

	current_pos := entrance_idx
	path_order := []int{current_pos}

	for visited_items < n_items || !visited_checkout {
		next_pos := -1
		min_dist := math.Inf(1)

		for i := range n {
			if i == 0 {
				continue
			}
			if !visited[i] {
				if i <= n_items { // item
					if dist_matrix[current_pos][i] < min_dist {
						min_dist = dist_matrix[current_pos][i]
						next_pos = i
					}
				} else if i > n_items && visited_items == n_items { // checkout
					if dist_matrix[current_pos][i] < min_dist {
						min_dist = dist_matrix[current_pos][i]
						next_pos = i
					}
				}
			}
		}

		if next_pos == -1 {
			break
		}

		visited[next_pos] = true
		path_order = append(path_order, next_pos)
		current_pos = next_pos

		if next_pos <= n_items {
			visited_items += 1
		} else if next_pos > n_items {
			visited_checkout = true
		}
	}

	path_order = append(path_order, exit_idx)

	full_path := []point{}
	for i := range len(path_order) - 1 {
		x := path_matrix[path_order[i]][path_order[i+1]]
		full_path = append(full_path, x[:len(x)-1]...)
	}
	full_path = append(full_path, path_matrix[path_order[len(path_order)-1]][exit_idx]...)

	return full_path
}

func solve(grid [][]square, start point, products set[int]) []point {
	items, checkouts, end := extractPoints(grid, products)
	points := slices.Concat([]point{start}, items, checkouts, []point{end})
	dist_matrix, path_matrix := createDistanceAndPathMatrix(grid, points)
	path := tsp_with_one_checkout_greedy(dist_matrix, path_matrix, items, checkouts, start, end)
	return path
}

func theAlgorithm(grid [][]square, start point, products set[int]) []point {
	bestPath := []point{}
	// best_egg := -1
	for _, egg := range eggs {
		new_required := maps.Clone(products)
		new_required.insert(egg)
		// grid = generate_grid(data_dict, new_required)
		path := solve(grid, start, new_required)
		if len(path) < len(bestPath) || len(bestPath) == 0 {
			bestPath = path
			// best_egg = egg
		}
	}

	return bestPath
}
