package utils

import (
	"container/heap"
	"fmt"
	"log"
	"slices"
	"strings"
)

type Neighbor[T comparable] struct {
	Value T
	Cost  int
}
type PositionGraph = Graph[Position]

type Graph[T comparable] map[T][]Neighbor[T]

func (graph *Graph[T]) String() string {
	var sb strings.Builder
	for from, tos := range *graph {
		sb.WriteString(fmt.Sprintf("%v -> %v\n", from, tos))
	}
	return sb.String()
}

func (graph *Graph[T]) AddEdge(from T, to T, Cost int) {
	//fmt.Println("Adding edge from", from, "to", to, "with Cost", Cost)
	(*graph)[from] = append((*graph)[from], Neighbor[T]{Value: to, Cost: Cost})
}

func makePathFromPrevMap[T comparable](prev map[T]T, froms []T, to T) [][]T {
	ret := make([][]T, 0)
	for _, from := range froms {
		path := make([]T, 0)
		path = append(path, to)
		curr := to
		for curr != from {
			prevNode, ok := prev[curr]
			if !ok {
				break
			}
			curr = prevNode
			path = append(path, curr)
		}
		slices.Reverse(path)
		ret = append(ret, path)
	}
	return ret
}

func (graph *Graph[T]) FindDistanceAndPath(froms []T, ends []T) (int, [][]T) {
	//fmt.Println("Finding path from", froms, "to", ends)
	toVisit := make(PriorityQueue[T], len(froms))
	distances := make(map[T]int)
	for idx, from := range froms {
		toVisit[idx] = &Item[T]{value: from, priority: 0}
		distances[from] = 0
	}
	heap.Init(&toVisit)
	prev := make(map[T]T)

	for toVisit.Len() > 0 {
		currItem := heap.Pop(&toVisit).(*Item[T])
		currPos := currItem.value
		currDistance := currItem.priority

		//log.Printf("Visiting %v with cost %d ", currPos, currDistance)

		if slices.Contains(ends, currPos) {
			pathsToCurr := makePathFromPrevMap(prev, froms, currPos)
			return currDistance, pathsToCurr
		}

		for _, neighbor := range (*graph)[currPos] {
			prevCost, hasPrevCost := distances[neighbor.Value]
			itemInToVisit := toVisit.Get(neighbor.Value)
			newCostToNeighbor := currDistance + neighbor.Cost
			if hasPrevCost && prevCost < newCostToNeighbor {
				continue
			}
			if itemInToVisit != nil && itemInToVisit.priority < newCostToNeighbor {
				continue
			}

			prev[neighbor.Value] = currPos
			distances[neighbor.Value] = newCostToNeighbor
			//log.Printf("Adding %v to toVisit with cost %d", neighbor.Value, currDistance+neighbor.Cost)
			if itemInToVisit != nil {
				toVisit.update(itemInToVisit, neighbor.Value, newCostToNeighbor)
			} else {
				heap.Push(&toVisit, &Item[T]{value: neighbor.Value, priority: currDistance + neighbor.Cost})
			}
		}
	}
	log.Fatal("No path found", graph)
	return -1, nil
}

func (graph *Graph[T]) FindCycleBFS(from T) []T {
	visited := make(map[T]bool)
	pathsToNode := make(map[T][]T)
	toVisit := make([]T, 0)
	toVisit = append(toVisit, from)
	for len(toVisit) > 0 {
		curr := toVisit[len(toVisit)-1]

		pathToCurr := pathsToNode[curr]

		toVisit = toVisit[:len(toVisit)-1]
		if curr == from && len(pathToCurr) > 0 {
			return pathToCurr
		}

		for _, neighbor := range (*graph)[curr] {
			if visited[neighbor.Value] && neighbor.Value != from {
				// Don't go backwards along the path
				continue
			}

			toVisit = append(toVisit, neighbor.Value)
			pathsToNode[neighbor.Value] = make([]T, len(pathsToNode[curr]))
			copy(pathsToNode[neighbor.Value], pathsToNode[curr])
			pathsToNode[neighbor.Value] = append(pathsToNode[neighbor.Value], curr)
		}

		visited[curr] = true
	}
	log.Fatal("No path found", graph)
	return nil
}
