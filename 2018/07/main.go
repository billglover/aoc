package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type node struct {
	id         string
	children   []*node
	parents    []*node
	startTime  int
	inProgress bool
	visited    bool
	weight     int
	lock       sync.RWMutex
}

type graph map[string]*node

func (g graph) NextTask() *node {
	for _, n := range g {
		if n.IsAvailable() {
			return n
		}
	}
	return nil
}

type worker int

func (n *node) Start(t int) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.startTime = t
	n.inProgress = true
}

func (n *node) Complete() {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.visited = true
}

func (n *node) Completed() bool {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.visited
}

func (n *node) Started() bool {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.inProgress
}

func (n *node) IsAvailable() bool {
	n.lock.RLock()
	defer n.lock.RUnlock()

	if n.inProgress || n.visited {
		return false
	}

	if n.inProgress == false && n.parents == nil {
		return true
	}

	for p := range n.parents {
		if n.parents[p].Completed() == false {
			return n.inProgress
		}
	}
	return true
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("unable to open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("unable to read file:", err)
		os.Exit(1)
	}

	depends := map[string][]string{}
	nodes := map[string]bool{}

	re := regexp.MustCompile("^Step ([A-Z]) must be finished before step ([A-Z]) can begin.$")
	lines := strings.Split(string(data), "\n")
	for l := range lines {
		matches := re.FindStringSubmatch(lines[l])
		if matches == nil {
			fmt.Println("unable to parse line:", lines[l])
			os.Exit(1)
		}
		parent, child := matches[1], matches[2]
		nodes[parent], nodes[child] = true, true

		if _, ok := depends[parent]; ok == false {
			depends[parent] = []string{child}
			continue
		}
		depends[parent] = append(depends[parent], child)
	}

	tree, answerOne := partOne(nodes, depends)
	fmt.Println("Part One:", answerOne)

	for n := range tree {
		tree[n].visited = false
	}

	answerTwo := partTwo(5, 0, tree, answerOne)
	fmt.Println("Part Two:", answerTwo)
}

// PartOne solves the first part of the puzzle.
func partOne(vs map[string]bool, ds map[string][]string) (map[string]*node, string) {
	var answer string

	nodes := map[string]*node{}
	for v := range vs {
		n := node{id: v}
		nodes[v] = &n
	}

	for d := range ds {
		parentNode := nodes[d]
		for _, c := range ds[d] {
			childNode := nodes[c]
			childNode.parents = append(childNode.parents, parentNode)
			parentNode.children = append(parentNode.children, childNode)
		}
	}

	// look for a root nodes
	candidates := []string{}
	for n := range nodes {
		if nodes[n].parents == nil {
			candidates = append(candidates, n)
		}
	}

	// walk the tree identifying sequence of nodes
	for {
		if len(candidates) == 0 {
			break
		}

		var nextNode string
		candidates, nextNode = visitNextNode(candidates, nodes)
		answer += nextNode
	}

	return nodes, answer
}

func visitNextNode(candidates []string, nodes map[string]*node) ([]string, string) {
	sort.Strings(candidates)
	nextNode := candidates[0]
	nodes[nextNode].visited = true

	for n := range nodes {
		if nodes[n].visited == true {
			continue
		}

		unblocked := true
		for p := range nodes[n].parents {
			if nodes[n].parents[p].visited == false {
				unblocked = false
			}
		}

		if unblocked && inSlice(candidates, n) == false {
			candidates = append(candidates, n)
		}
	}

	newCandidates := make([]string, len(candidates)-1)
	copy(newCandidates, candidates[1:])
	return newCandidates, nextNode
}

// PartTwo solves the second part of the puzzle.
func partTwo(workerCount, baseDuration int, tree map[string]*node, sequence string) int {

	// add job weights
	for id := range tree {
		tree[id].weight = 60 + int(id[0]) - 65
	}

	// declare our graph so we can use methods on it
	g := graph(tree)

	// declare our workers
	workers := make([]*node, workerCount)

	t := 0
	for {
		//fmt.Println("\ntime:", t)

		// if there are no more jobs to complete, end
		complete := true
		for _, n := range g {
			if n.Completed() == false {
				complete = false
			}
		}
		if complete == true {
			break
		}

		// at each time, t, update our worker status
		for w := range workers {

			// if our worker has finished the task, mark it as complete
			if workers[w] != nil {
				if t > workers[w].startTime+workers[w].weight {
					//fmt.Println(w, "complete", workers[w].id)
					workers[w].Complete()
					workers[w] = nil
				}
			}

			// if our worker has unfinished work then just continue
			if workers[w] != nil && workers[w].Completed() == false {
				//fmt.Println(w, "working on:", workers[w].id)
				continue
			}

			// if our worker has no task, get the next one from the graph
			if workers[w] == nil {
				workers[w] = g.NextTask()
				if workers[w] == nil {
					//fmt.Println(w, "no work available")
					continue
				}
				workers[w].Start(t)
				//fmt.Println(w, "starting work on:", workers[w].id)
			}
		}
		if t > 20000 {
			break
		}
		w1, w2 := ".", "."
		if workers[0] != nil {
			w1 = workers[0].id
		}
		if workers[1] != nil {
			w2 = workers[1].id
		}

		fmt.Printf("%2d:\t%s\t%s\n", t, w1, w2)

		t++
	}

	fmt.Println("ended at t:", t-1)
	return t - 1
}

func startClock(master chan int, chans []chan int) {
	for t := range master {
		for c := range chans {
			chans[c] <- t
		}
	}
	for c := range chans {
		close(chans[c])
	}
}

func inSlice(slice []string, s string) bool {
	for i := range slice {
		if slice[i] == s {
			return true
		}
	}
	return false
}

func nextTask(tree map[string]*node, nextTask string) string {
	if tree[nextTask].parents == nil {
		return nextTask
	}

	for p := range tree[nextTask].parents {
		if tree[nextTask].parents[p].visited == false {
			return ""
		}
	}
	return nextTask

}
