package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	autoAccept         bool
	conditionReads     string
	conditionIsGT      bool
	conditionThreshold int
	destination        string
}

func NewRule(conditionStr string) Rule {
	if !strings.Contains(conditionStr, ":") {
		return Rule{autoAccept: true, destination: conditionStr}
	}
	re := regexp.MustCompile(`(x|m|a|s)(>|<)(\d+):(\w+)`)
	matches := re.FindStringSubmatch(conditionStr)
	if len(matches) != 5 {
		log.Panicf("Invalid condition string %s", conditionStr)
	}
	threshold, err := strconv.Atoi(matches[3])
	if err != nil {
		log.Panicf("Invalid threshold value %s", matches[3])
	}

	return Rule{
		conditionReads:     matches[1],
		conditionIsGT:      matches[2] == ">",
		conditionThreshold: threshold,
		destination:        matches[4],
	}
}

func (r Rule) ShouldApply(part Part) bool {
	if r.autoAccept {
		return true
	}
	value := 0
	switch r.conditionReads {
	case "x":
		value = part.x
	case "m":
		value = part.m
	case "a":
		value = part.a
	case "s":
		value = part.s
	default:
		log.Panicf("Invalid condition reads %s", r.conditionReads)
	}
	if r.conditionIsGT {
		return value > r.conditionThreshold
	}
	return value < r.conditionThreshold
}

func (r Rule) GetDestination() string {
	return r.destination
}

type Workflow struct {
	name  string
	rules []Rule
}

func NewWorkflow(workflowStr string) Workflow {
	re := regexp.MustCompile(`^([a-z0-9]+){([^}]+)}$`)
	matches := re.FindStringSubmatch(workflowStr)
	if len(matches) != 3 {
		log.Panicf("Invalid workflow string %s", workflowStr)
	}
	name := matches[1]
	rulesStr := matches[2]
	rulesStrings := strings.Split(rulesStr, ",")
	rules := make([]Rule, len(rulesStrings))
	for i, ruleStr := range rulesStrings {
		rules[i] = NewRule(ruleStr)
	}
	return Workflow{
		name:  name,
		rules: rules,
	}
}

type Part struct {
	x, m, a, s int
}

func NewPart(line string) Part {
	re := regexp.MustCompile(`{x=(\d+),m=(\d+),a=(\d+),s=(\d+)}`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 5 {
		log.Panicf("Invalid part string %s", line)
	}
	x, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Panicf("Invalid x value %s", matches[1])
	}
	m, err := strconv.Atoi(matches[2])
	if err != nil {
		log.Panicf("Invalid m value %s", matches[2])
	}
	a, err := strconv.Atoi(matches[3])
	if err != nil {
		log.Panicf("Invalid a value %s", matches[3])
	}
	s, err := strconv.Atoi(matches[4])
	if err != nil {
		log.Panicf("Invalid s value %s", matches[4])
	}
	return Part{
		x: x,
		m: m,
		a: a,
		s: s,
	}
}

func (p Part) PartScore() int {
	return p.x + p.m + p.a + p.s
}

func ParseFile(fname string) ([]Part, []Workflow) {
	file, err := os.Open(fname)
	if err != nil {
		log.Panicf("Failed to open file %s: %v", fname, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	workflows := make([]Workflow, 0)
	for scanner.Scan() {
		t := scanner.Text()
		t = strings.TrimSpace(t)
		if t == "" {
			break
		}
		w := NewWorkflow(t)
		//fmt.Println(w)
		workflows = append(workflows, w)
	}

	parts := make([]Part, 0)
	for scanner.Scan() {
		t := scanner.Text()
		t = strings.TrimSpace(t)
		parts = append(parts, NewPart(t))
	}

	return parts, workflows
}

func applyWorkflows(part Part, workflowMap map[string]Workflow, currWorkflow string) string {
	workflow := workflowMap[currWorkflow]
	for _, rule := range workflow.rules {
		if rule.ShouldApply(part) {
			dest := rule.GetDestination()
			//fmt.Printf("Rule %v applied to part %v, going to %s\n", rule, part, dest)
			if dest == "A" || dest == "R" {
				return dest
			}
			return applyWorkflows(part, workflowMap, dest)
		}
	}
	log.Panic("Should not reach here", part)
	return ""
}

func part1(fname string) int {
	parts, workflows := ParseFile(fname)
	workflowMap := make(map[string]Workflow)
	for _, workflow := range workflows {
		workflowMap[workflow.name] = workflow
	}
	out := 0
	for _, part := range parts {
		result := applyWorkflows(part, workflowMap, "in")
		if result == "A" {
			out += part.PartScore()
		}
	}
	return out
}

func main() {
	fmt.Println(part1("day19-input-easy.txt"))
	fmt.Println(part1("day19-input.txt"))
}
