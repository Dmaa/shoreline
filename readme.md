To run, modify line 90 in main.go to be your own chosen storage directory and type "go run main.go" in the shoreline directory.

Code Overview: 
Node: node.go
Tests and runner: main.go

Solution Overview:

Per the spec, I chose to create a solution that relies completely on persistent storage and without coordination. Other solutions
I chose included Paxos/Primary Backup, but those distributed system paradigms are slow and require many RPCs. Assumption: I assume every Node has a single client. My solution could be modified to serve multiple clients with one node fairly quickly by locking within critical sections (functions that give out Node_IDs), but I chose not to implement that due to the above assumption.

1. Solution to Get_id():
I chose to initialize each node with a uid of its own ID, and increment the ID by 1024 every time a request to get_id came in. This solution is globally unique because there are max 1024 nodes, so every node will have the section of the id space of (id + (1024 * n)), where id is the node's starting ID and n is any positive integer (and is = to the number of requests the node received). I did not test or create a solution that would handle ID overflow, as the amount of time before overflow (assuming 100,000 requests per second per node) = 2^64/(1024 * 100000 * 60 * 60 * 24 * 365) = 5712 years, which is longer than our system should last.

No coordination is required, as every node controls its own section of the ID space, and crashes are handled using persitent storage checkpoints (written to disk). I could have incremented the ID by less than 1024 (by n where n is the number of nodes, as that is known at startup) to save on UIDs, but I decided to use 1024 for a simpler solution (also because as proved above, we have no shortage of IDs).

2. Performance:
My node is able to achieve 1000 qps by checkpointing only every 25 IDs. (Note you can change this granularity yourself).

In main.go, I wrote a series of tests to verify performance. To ensure my node can handle 100,000 qps, I wrote a test measuring the amount of time it took for a node to handle 1000 requests (checkThroughput). I realize this disregards potential problems such as disk write speeds and RPC latencies/throttling being different for different machines, but as that was not specified in the spec I just ensured the code worked on my machine. (On my macbook pro, 100,000 requests took about .8-.9 seconds. If it's taking longer on your machine, increase the checkpoint parameter.)

3. Example Cases and Testing

Some cases I tested:

- Verifying correctness with multiple machines. Logically, since every node controls its own section of the ID space, multiple nodes should not result in duplicate IDs. I verified this in multipleNodes() in main.go.


- Verifying correctness with a crash and restart. Logically, since on startup every node reads its persistent storage, a node should take into consideration which IDs it has previously given out before serving new Get_id() requests. (My API requires that nodes are brought up through Make() on startup.) I tested this with the function SimulateCrash().

- Verifying correctness that a single node has unique IDs. Logically, this should be true because nodes increment an internal counter before giving out IDs. I verified this in checkCorrectness(). 