package main

import (
	"fmt"
	"node"
	"time"
)

// Throughput test - can I handle 100,000 requests/sec?
func checkThroughput(myNode *node.Node) {
	start := time.Now()
	for i := 0; i < 100000; i++ {

		myNode.Get_id()
	}
	elapsed := time.Since(start)
	fmt.Printf("Took %s\n", elapsed)
	if elapsed > time.Second {
		panic("Took too long!")
	}
}

// Check to make sure all nums for a particular node ID are unique.
func checkCorrectness(myNode *node.Node, unique_nums *map[uint64]int) {
	for i := 0; i < 100000; i++ {
		n := myNode.Get_id()
		if _, ok := (*unique_nums)[n]; ok {
			panic("Non unique uid")
		} else {
			(*unique_nums)[n] = 1
		}
	}
}

// Check if a single node can recover from a crash.
func simulateCrash(myNode *node.Node, nodeId int, unique_nums *map[uint64]int) {
	for i := 0; i < 50000; i++ {
		n := myNode.Get_id()
		if _, ok := (*unique_nums)[n]; ok {
			panic("Non unique uid")
		} else {
			(*unique_nums)[n] = 1
		}
	}
	secondNode := node.Make(0, 25, "/Users/dharma/Documents/shoreline/persistent_storage/")
	for i := 0; i < 50000; i++ {
		n := secondNode.Get_id()
		if _, ok := (*unique_nums)[n]; ok {
			panic("Non unique uid")
		} else {
			(*unique_nums)[n] = 1
		}
	}
}

// Check if 100,000 requests to nodes with different nodeIDs result in different get_IDs.
func multipleNodes() {
	unique_nums := make(map[uint64]int)
	n := node.Make(2, 25, "/Users/dharma/Documents/shoreline/persistent_storage/")
	n1 := node.Make(3, 25, "/Users/dharma/Documents/shoreline/persistent_storage/")
	n2 := node.Make(4, 25, "/Users/dharma/Documents/shoreline/persistent_storage/")

	for i := 0; i < 100000; i++ {
		nI := n.Get_id()
		n1I := n1.Get_id()
		n2I := n2.Get_id()

		if _, ok := unique_nums[nI]; ok {
			panic("Non unique uid")
		} else {
			unique_nums[nI] = 1
		}

		if _, ok := unique_nums[n1I]; ok {
			panic("Non unique uid")
		} else {
			unique_nums[n1I] = 1
		}

		if _, ok := unique_nums[n2I]; ok {
			panic("Non unique uid")
		} else {
			unique_nums[n2I] = 1
		}
	}
}

func main() {
	// TODO(user): Replace this with your own storage dir.
	n := node.Make(0, 25, "/Users/dharma/Documents/shoreline/persistent_storage/")
	checkThroughput(&n)
	unique_nums := make(map[uint64]int)
	checkCorrectness(&n, &unique_nums)
	simulateCrash(&n, 0, &unique_nums)
	multipleNodes()
}
