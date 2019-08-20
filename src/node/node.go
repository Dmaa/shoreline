package node

import (
	"io/ioutil"
	"os"
	"strconv"
)

const NUM_NODES uint64 = 1024

type Node struct {
	lastID                uint64
	nodeId                int
	myStorage             string
	checkpointCounter     int
	checkpointGranularity int
}

// Helper function to write a node ID checkpoint to persistent storage.
func (nde *Node) write_node_id() {
	f, err := os.Create(nde.myStorage)
	check(err)

	_, err = f.WriteString(strconv.FormatUint(nde.lastID, 10))
	check(err)
	f.Close()
}

// Return a unique uint64 ID.
func (nde *Node) Get_id() uint64 {
	nde.lastID += NUM_NODES
	if nde.checkpointCounter == nde.checkpointGranularity {
		nde.write_node_id()
		nde.checkpointCounter = 0
	} else {
		nde.checkpointCounter++
	}
	return nde.lastID
}

// Assume already implemented helper function
// Returns node_id
func (nde *Node) node_id() int {
	return 1
}

// Assume already implemented helper function
// Returns num_millis since epoch
func (nde *Node) timestamp() {

}

// Function copied from gobyexample.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Function to create a node (equivalent of constructor) Must call this to start a node.
func Make(nodeId int, checkpointGranularity int, storageDir string) Node {
	nde := &Node{}
	nde.nodeId = nodeId
	nde.checkpointGranularity = checkpointGranularity

	nde.myStorage = storageDir + strconv.Itoa(nde.nodeId)

	if _, err := os.Stat(nde.myStorage); os.IsNotExist(err) {
		// We are starting up for the first time. Persistent storage starts with a 0.
		f, err := os.Create(nde.myStorage)
		check(err)

		_, err = f.WriteString(strconv.Itoa(nde.nodeId))
		check(err)
		f.Close()
		nde.lastID = uint64(nodeId)
		// Value is already default-initialized to 0, but make it more clear.
		nde.checkpointCounter = 0
	} else if err == nil {
		// File already exists, meaning we are resuming from crash.
		dat, err := ioutil.ReadFile(nde.myStorage)
		check(err)
		nde.lastID, err = strconv.ParseUint(string(dat), 10, 64)
		nde.lastID += 1024 * uint64(nde.checkpointGranularity)
		nde.checkpointCounter = 0
	}
	return *nde
}
