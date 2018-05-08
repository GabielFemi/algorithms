package graph

import "math"

func floydWarShall(g weightedGraph, init func(weightedGraph) ([][][]int, []interface{}), handler func(*[][][]int, int, int, int),
	rebuild func([]interface{}, [][]int)) {
	array, vertices := init(g)

	for k := range array[:len(array)-1] {
		for i := range array[k] {
			for j := range array[k][i] {
				handler(&array, k, i, j)
			}
		}
	}
	rebuild(vertices, array[len(array)-1])
}

func distFloydInit(g weightedGraph) ([][][]int, []interface{}) {
	vertices := g.AllVertices()
	array := make([][][]int, len(vertices)+1, len(vertices)+1)

	for k := range array {
		array[k] = make([][]int, len(vertices), len(vertices))
		for i := range array[k] {
			array[k][i] = make([]int, len(vertices), len(vertices))
			if k == 0 {
				for j := range array[k][i] {
					currentEdge := edge{vertices[i], vertices[j]}
					if i == j {
						array[k][i][j] = 0
					} else if !g.CheckEdge(currentEdge) {
						array[k][i][j] = math.MaxInt32
					} else {
						array[k][i][j] = g.Weight(currentEdge)
					}
				}
			}
		}
	}

	return array, vertices
}

func distFloydHandler(array *[][][]int, k, i, j int) {
	(*array)[k+1][i][j] = (*array)[k][i][j]
	if (*array)[k][i][k]+(*array)[k][k][j] < (*array)[k+1][i][j] {
		(*array)[k+1][i][j] = (*array)[k][i][k] + (*array)[k][k][j]
	}
}

func distFloydWarShall(g weightedGraph) weightedGraph {
	newG := createGraphByType(g).(weightedGraph)
	rebuild := func(vertices []interface{}, array [][]int) {
		for i := range vertices {
			for j := range vertices {
				if array[i][j] < math.MaxInt32 {
					newG.AddEdgeWithWeight(edge{vertices[i], vertices[j]}, array[i][j])
				}
			}
		}
	}

	floydWarShall(g, distFloydInit, distFloydHandler, rebuild)
	return newG
}

func pathFloydWarShall(g weightedGraph) map[interface{}]weightedGraph {
	var pathArray [][]int
	var pathForest map[interface{}]weightedGraph
	init := func(g weightedGraph) ([][][]int, []interface{}) {
		distArray, vertices := distFloydInit(g)
		pathArray = make([][]int, len(vertices), len(vertices))
		for i := range pathArray {
			pathArray[i] = make([]int, len(vertices), len(vertices))
			for j := range pathArray[i] {
				currentEdge := edge{vertices[i], vertices[j]}
				if i == j || !g.CheckEdge(currentEdge) {
					pathArray[i][j] = math.MaxInt32
				} else {
					pathArray[i][j] = i
				}
			}
		}
		return distArray, vertices
	}

	handler := func(array *[][][]int, k, i, j int) {
		if (*array)[k][i][j] > (*array)[k][k][j]+(*array)[k][i][k] {
			pathArray[i][j] = pathArray[k][j]
		}
		distFloydHandler(array, k, i, j)
	}

	rebuild := func(vertices []interface{}, array [][]int) {
		pathForest = make(map[interface{}]weightedGraph)
		for i := range vertices {
			pathForest[vertices[i]] = createGraphByType(g).(weightedGraph)
			for j := range vertices {
				if pathArray[i][j] < math.MaxInt32 {
					currentEdge := edge{vertices[pathArray[i][j]], vertices[j]}
					pathForest[vertices[i]].AddEdgeWithWeight(currentEdge, g.Weight(currentEdge))
				}
			}
		}
	}

	floydWarShall(g, init, handler, rebuild)
	return pathForest
}
