package main

import (
	"maps"
	"math"
	"slices"
)

func extractPoints(products set[int]) ([]point, []point) {
	items := []point{}
	checkouts := []point{}
	for y, row := range grid {
		for x, cell := range row {
			if cell.Kind == productSquare && products.contains(cell.ProductID) { // Required products to be visited
				items = append(items, point{x, y})
			} else if cell.Kind == checkoutSquare {
				checkouts = append(checkouts, point{x, y})
			}
		}
	}
	return items, checkouts
}

func bfs(start point, pointsWithItems set[point]) ([gridHeight][gridWidth]float64, [gridHeight][gridWidth]point) {
	directions := []point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	var dist [gridHeight][gridWidth]float64
	for y, r := range dist {
		for x, _ := range r {
			dist[y][x] = math.Inf(1)
		}
	}
	var prev [gridHeight][gridWidth]point
	dist[start.Y][start.X] = 0
	queue := []point{start}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, d := range directions {
			n := point{p.X + d.X, p.Y + d.Y}
			if 0 <= n.X && n.X < gridWidth && 0 <= n.Y && n.Y < gridHeight && !(grid[n.Y][n.X].Kind == wallSquare || grid[n.Y][n.X].Kind == productSquare && !pointsWithItems.contains(n)) && dist[n.Y][n.X] == math.Inf(1) {
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

func reconstructPath(prev [gridHeight][gridWidth]point, start point, end point) []point {
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

func create_distance_and_path_matrix(points []point) ([][]float64, [gridHeight][gridWidth][]point) {
	distMatrix := make([][]float64, len(points))
	for i := range distMatrix {
		distMatrix[i] = make([]float64, len(points))
	}
	var pathMatrix [gridHeight][gridWidth][]point

	for i, p := range points {
		pset := set[point]{}
		for _, q := range points {
			pset.insert(q)
		}
		dist, prev := bfs(p, pset)
		for j, q := range points {
			if i != j {
				distMatrix[i][j] = dist[q.Y][q.X]
				pathMatrix[i][j] = reconstructPath(prev, p, q)
			}
		}
	}

	return distMatrix, pathMatrix
}

func tsp_with_one_checkout_greedy(dist_matrix [][]float64, path_matrix [gridHeight][gridWidth][]point, items []point, checkouts []point, entrance point, exit point) []point {
	n_items := len(items)
	// n_checkouts := len(checkouts)
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

func solve(products set[int]) []point {
	items, checkouts := extractPoints(products)
	points := slices.Concat([]point{start}, items, checkouts, []point{end})
	dist_matrix, path_matrix := create_distance_and_path_matrix(points)
	path := tsp_with_one_checkout_greedy(dist_matrix, path_matrix, items, checkouts, start, end)
	return path
}

// def open_csv(filename):
// 	data_dict = {}
// 	with open(filename, 'r') as file:
// 		csv_reader = csv.reader(file)

// 		next(csv_reader)

// 		# Iterate over each row in the CSV file
// 		for row in csv_reader:
// 			if "B" in row[0]:
// 				# Blokadite nqmat unikalno ime - BL,2,3 BL,3,4, za tova te sa s pulno
// 				# ime za da e unikalen klucha
// 				data_dict[row[0] + row[1] + row[2]] = (int(row[1]), int(row[2]))
// 			else:
// 				data_dict[row[0]] = (int(row[1]), int(row[2]))

// 	return data_dict

// def generate_grid(data_dict, required_products):
// 	grid = [[0 for _ in range(41)] for _ in range(41)]

// 	for entry in data_dict:
// 		x, y = data_dict[entry][0], data_dict[entry][1]

// 		if entry in required_products:
// 			grid[x][y] = 6  // Required products
// 		elif "P" in entry:
// 			grid[x][y] = 2  // Non-required products
// 		elif "CA" in entry or "S" in entry:
// 			grid[x][y] = 5
// 		elif "B" in entry:  // Blockade
// 			grid[x][y] = 1
// 		elif "EN" in entry: // Entrance
// 			grid[x][y] = 3
// 		elif "EX" in entry: // Exit
// 			grid[x][y] = 4
// 		else:
// 			grid[x][y] = 0  // Empty

// 	return grid

func theAlgorithm(products set[int]) []point {
	// data_dict = open_csv("placement.csv")

	bestPath := []point{}
	// best_egg := -1
	for _, egg := range eggs {
		new_required := maps.Clone(products)
		new_required.insert(egg)
		// grid = generate_grid(data_dict, new_required)
		path := solve(new_required)
		if len(path) < len(bestPath) || len(bestPath) == 0 {
			bestPath = path
			// best_egg = egg
		}
	}

	return bestPath

	// products.add(best_egg)

	// final_grid = generate_grid(data_dict, products)
	// print(final_grid, best_path)
}
