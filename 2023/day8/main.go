package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Node struct {
	name  string
	left  string
	right string
}

func (node Node) String() string {
	return fmt.Sprintf("Node{name: %s, left: %s, right: %s}", node.name, node.left, node.right)
}

/*
	func (node Node) getLeftNode(nodeMap map[string]Node) Node {
		return nodeMap[node.left]
	}

	func (node Node) getRightNode(nodeMap map[string]Node) Node {
		return nodeMap[node.right]
	}
*/
func parseFile(fname string) (string, map[string]Node) {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	nodeMap := make(map[string]Node)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	instructions := scanner.Text()

	// consume empty line
	scanner.Scan()
	scanner.Text()

	for scanner.Scan() {
		line := scanner.Text()
		var node Node
		nodeRE := regexp.MustCompile(`(\w+) += +\((\w+), +(\w+)\)`)
		matches := nodeRE.FindStringSubmatch(line)
		node.name = matches[1]
		node.left = matches[2]
		node.right = matches[3]
		//log.Print(node)
		nodeMap[node.name] = node
	}
	return instructions, nodeMap
}

func navigateMap(nodeMap map[string]Node, instructions string) int {
	currentNode := nodeMap["AAA"]
	numSteps := 0
	for {
		for _, c := range instructions {
			if currentNode.name == "ZZZ" {
				return numSteps
			}
			if c == 'L' {
				currentNode = nodeMap[currentNode.left]
			} else if c == 'R' {
				currentNode = nodeMap[currentNode.right]
			} else {
				log.Fatal("Invalid instruction", c)
			}
			numSteps++
		}
	}
}

func part1(fname string) int {
	instructions, nodeMap := parseFile(fname)
	return navigateMap(nodeMap, instructions)
}

func main() {
	fmt.Println(part1("day8-input-easy.txt"))
	fmt.Println(part1("day8-input-easy2.txt"))
	fmt.Println(part1("day8-input.txt"))
}
