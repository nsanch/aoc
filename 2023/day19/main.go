package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Bounds struct {
	minX, maxX int
	minM, maxM int
	minA, maxA int
	minS, maxS int
}

func NewBounds() Bounds {
	return Bounds{
		minX: 1,
		maxX: 4000,
		minM: 1,
		maxM: 4000,
		minA: 1,
		maxA: 4000,
		minS: 1,
		maxS: 4000,
	}
}

func (b Bounds) NumPossibleValues() int {
	return (b.maxX - b.minX + 1) *
		(b.maxM - b.minM + 1) *
		(b.maxA - b.minA + 1) *
		(b.maxS - b.minS + 1)
}

func (b *Bounds) Clone() Bounds {
	return Bounds{
		minX: b.minX,
		maxX: b.maxX,
		minM: b.minM,
		maxM: b.maxM,
		minA: b.minA,
		maxA: b.maxA,
		minS: b.minS,
		maxS: b.maxS,
	}
}

type Rule struct {
	autoAccept         bool
	conditionReads     string
	conditionIsGT      bool
	conditionThreshold int
	destination        string
}

func (r Rule) String() string {
	if r.autoAccept {
		return fmt.Sprintf("Rule: autoAccept %s", r.destination)
	} else if r.conditionIsGT {
		return fmt.Sprintf("Rule: %s > %d -> %s", r.conditionReads, r.conditionThreshold, r.destination)
	} else {
		return fmt.Sprintf("Rule: %s < %d -> %s", r.conditionReads, r.conditionThreshold, r.destination)
	}
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

func (r Rule) ConstrainAcceptableRange(b *Bounds) bool {
	if r.autoAccept {
		return true
	}
	var minV *int
	var maxV *int
	switch r.conditionReads {
	case "x":
		minV = &b.minX
		maxV = &b.maxX
	case "m":
		minV = &b.minM
		maxV = &b.maxM
	case "a":
		minV = &b.minA
		maxV = &b.maxA
	case "s":
		minV = &b.minS
		maxV = &b.maxS
	default:
		log.Panicf("Invalid condition reads %s", r.conditionReads)
	}

	if r.conditionIsGT {
		if *maxV != 0 && *maxV < r.conditionThreshold {
			// we know that V must be less than the threshold. So this rule can never pass.
			return false
		}
		*minV = max(*minV, r.conditionThreshold+1)
		return true
	} else {
		if *minV != 0 && *minV >= r.conditionThreshold {
			// we know that V must be >= than the threshold. So this rule can never pass.
			return false
		}
		*maxV = min(*maxV, r.conditionThreshold-1)
		return true
	}
}

func (r Rule) ConstrainToFailureRange(b *Bounds) bool {
	if r.autoAccept {
		return false
	}

	var minV *int
	var maxV *int
	switch r.conditionReads {
	case "x":
		minV = &b.minX
		maxV = &b.maxX
	case "m":
		minV = &b.minM
		maxV = &b.maxM
	case "a":
		minV = &b.minA
		maxV = &b.maxA
	case "s":
		minV = &b.minS
		maxV = &b.maxS
	default:
		log.Panicf("Invalid condition reads %s", r.conditionReads)
	}

	if r.conditionIsGT {
		// we know we fail this rule, so V must be <= threshold.
		if *minV != 0 && *minV > r.conditionThreshold {
			// we know that V must be greater than the threshold. So this rule can never pass.
			return false
		}
		*maxV = min(*maxV, r.conditionThreshold)
		return true
	} else {
		// similar to above. we can only pass if we're >= threshold. so if max is < than threshold,
		// this can never pass.
		if *maxV != 0 && *maxV < r.conditionThreshold {
			// we know that V must be less than the threshold. So this rule can never pass.
			return false
		}
		*minV = max(*minV, r.conditionThreshold)
		return true
	}
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

func WalkPaths(workflowMap map[string]Workflow, currWorkflow string, bounds Bounds) []Bounds {
	if currWorkflow == "A" {
		// we've reached the end-state, return the bounds.
		return []Bounds{bounds}
	}

	ret := make([]Bounds, 0)
	workflow := workflowMap[currWorkflow]
	for _, rule := range workflow.rules {
		// ignore the terminal state reject rules.
		if rule.GetDestination() != "R" {
			fmt.Printf("evaluating rule for success %s\n", rule.String())
			b := bounds.Clone()
			if !rule.ConstrainAcceptableRange(&b) {
				fmt.Printf("rule %s cannot pass, bounds: %v\n", rule.String(), bounds)
			} else {
				fmt.Printf("Recursing! bounds after rule %s: %v\n", rule.String(), b)
				ret = slices.Concat(ret, WalkPaths(workflowMap, rule.GetDestination(), b))
				// cannot break here because there could be multiple ways in this ruleset to
				// get to the same workflow, but we must've failed this rule to keep going, so continue
				// with the logic below.
			}
		}

		fmt.Printf("evaluating rule for failure %s\n", rule.String())
		// we must fail this rule to get to the next one.
		if !rule.ConstrainToFailureRange(&bounds) {
			fmt.Printf("rule %s cannot fail, bounds: %v\n", rule.String(), bounds)
			break
		}
		fmt.Printf("bounds after rule %s: %v\n", rule.String(), bounds)
	}
	return ret
}

func part2(fname string) int {
	_, workflows := ParseFile(fname)
	workflowMap := make(map[string]Workflow)
	for _, workflow := range workflows {
		workflowMap[workflow.name] = workflow
	}
	allPossibleBounds := WalkPaths(workflowMap, "in", NewBounds())
	ret := 0
	for _, bounds := range allPossibleBounds {
		fmt.Println(bounds)
		ret += bounds.NumPossibleValues()
	}
	return ret
}

func main() {
	fmt.Println(part1("day19-input-easy.txt"))
	fmt.Println(part1("day19-input.txt"))

	fmt.Println(part2("day19-input-easy.txt"))
	fmt.Println(part2("day19-input.txt"))
}
