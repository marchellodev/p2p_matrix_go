package p2p_matrix

import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Network struct {
	//NodeInstance NodeInstanceInterface
	NodeFactory func(NodeInstance) NodeInstanceInterface
}

func deserializeScriptFile(scriptPath string) ScriptModel {
	data, _ := ioutil.ReadFile(scriptPath)

	var script ScriptModel
	_ = json.Unmarshal(data, &script)

	return script
}

// here we create new node based on the id from script
func newNode(n Network, script ScriptModel, id int, adapter NetworkAdapter) NodeInstanceInterface {
	return n.NodeFactory(NodeInstance{
		Id:   id,
		Node: script.Nodes[id],
		// todo NetworkAdapter
		networkAdapter: adapter,
	})
}

func runStory(script ScriptModel, network Network) ResultData {

	resultData := ResultData{}

	allNodes := make(map[int]NodeInstanceInterface)
	var activeNodes []int

	// profiling
	var currentOperationConnections []int
	var currentOperation StoryElementOperation
	var currentOperationTime float64
	readRight := 0
	readWrong := 0

	// here we need to know what is going on in the network: reading, writing and detailed info
	var networkAdapter NetworkAdapter
	networkAdapter.sendString = func(from NodeInstance, to int, data string) bool {
		active := false

		for _, node := range activeNodes {
			if node == to {
				active = true
			}
		}
		if !active {
			return false
		}

		allNodes[to].Listen(from.Id, data)
		return true
	}
	networkAdapter.sendRead = func(from NodeInstance, to int, file int) (bool, float64) {
		active := false

		for _, node := range activeNodes {
			if node == to {
				active = true
			}
		}
		if !active {
			return false, 0
		}

		if currentOperation.Type == "read" {
			// means that we are communicating with other nodes when reading
			currentOperationConnections = append(currentOperationConnections, to)
		}

		// here we need to calculate how much time the request will take

		fileData := allNodes[to].Read(from.Id, file)

		networkSpeed := (script.Nodes[from.Id].Speed + script.Nodes[to].Speed) / 2

		location1 := script.Nodes[from.Id].Location
		location2 := script.Nodes[to].Location
		var ping float64

		for _, pair := range script.Pings.Pings {
			if (pair.Location1 == location1 && pair.Location2 == location2) || (pair.Location1 == location2 && pair.Location2 == location1) {
				ping = pair.Ping
				break
			}
		}

		time := fileData/networkSpeed + (ping / 1000)
		currentOperationTime += time

		return true, fileData
	}
	networkAdapter.sendWrite = func(from NodeInstance, to int, file int, size float64) bool {
		active := false

		for _, node := range activeNodes {
			if node == to {
				active = true
			}
		}
		if !active {
			return false
		}

		allNodes[to].Write(from.Id, file, size)
		return true
	}

	//networkAdapter.sendString = func(a int, b string) {}
	// todo

	for id, storyElement := range script.Story {

		fmt.Println(id)

		for _, action := range storyElement.NodeActions {

			if action.Action == "on" {

				// if node does not yet exists, create it
				exists := false
				for key := range allNodes {
					if key == action.Node {
						exists = true
						break
					}
				}

				activeNodes = append(activeNodes, action.Node)

				if !exists {

					//for i := range activeNodes {
					//	j := rand.Intn(i + 1)
					//	activeNodes[i], activeNodes[j] = activeNodes[j], activeNodes[i]
					//}

					//rand.Shuffle(len(activeNodes), func(i, j int) {
					//	activeNodes[i], activeNodes[j] = activeNodes[j], activeNodes[i]
					//})

					// create the node
					allNodes[action.Node] = newNode(network, script, action.Node, networkAdapter)

					if len(activeNodes) == 1 {

						allNodes[action.Node].Activate(-1)
					} else {
						// get random node to bootstrap from
						//minus := 1
						//bNode := action.Node
						//for bNode == action.Node {
						//	bNode = activeNodes[len(activeNodes)-minus]
						//	minus++
						//}

						allNodes[action.Node].Activate(script.Nodes[action.Node].Bootstrap)
					}

				} else {
					allNodes[action.Node].Activate(-2)
				}

				//if id == 0 {
				//	fmt.Println("on 0")
				//}
				//network.NodeInstance.Activate(id)
			} else {
				//if id == 0 {
				//	fmt.Println("off 0")
				//}
				activeNodes = remove(activeNodes, action.Node)
			}
		}

		for _, operation := range storyElement.Operations {
			currentOperationConnections = []int{}
			currentOperationTime = 0

			currentOperation = operation

			if operation.Type == "write" {
				allNodes[operation.Node].Write(-1, operation.File, script.Files[operation.File].Size)
				// todo check if the node is active ?
			} else {
				read := allNodes[operation.Node].Read(-1, operation.File)
				if read == script.Files[operation.File].Size {
					readRight++
				} else {
					readWrong++
				}

				resultData.WriteOperation(currentOperation, currentOperationConnections)
				resultData.WriteOperationTime(currentOperation, currentOperationTime)

			}

			currentOperation = StoryElementOperation{}
		}

		resultData.WriteStorage(id, allNodes)
	}

	resultData.FileNotFound = float64(readWrong) / (float64(readWrong) + float64(readRight))

	return resultData
}

func remove(s []int, obj int) []int {
	var i int
	for n, el := range s {
		if el == obj {
			i = n
		}
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (n Network) RunScript(scriptPath string, resultPath string) {

	script := deserializeScriptFile(scriptPath)

	result := runStory(script, n)

	result.SaveToJson(resultPath)

	//var adapter = NetworkAdapter{
	//	sendString: func(i int, s string) {
	//
	//	},
	//	sendRead: func(i int, i2 int) float64 {
	//		return 0
	//	},
	//	sendWrite: func(i int, i2 int, f float64) {
	//
	//	},
	//}
	//
	//fmt.Println(adapter)

	//n.Activate(0, adapter)

	//data, err := ioutil.ReadFile(script)

	//if err != nil {
	//	panic(err)
	//}

	//fmt.Println(data)

	//d1 := []byte("hello\ngo\n")
	//err := ioutil.WriteFile(path, d1, 0644)
	//if err != nil {
	//	panic(err)
	//}
}

// new stuff

type NodeInstanceInterface interface {

	// NodeInstance - the bootstrap node
	// NodeInstance contains only Id and Node
	Activate(int)

	// int - sender; string - message
	Listen(int, string)

	// Here we (int = 0) or other node asks to read file with id int
	Read(int, int) float64

	// Here we (int = 0) or other node asks to read file with id int
	Write(int, int, float64)

	SysGetPeers() []int
	SysGetStorage() map[int]float64
}

// add list of peers & other stuff
type NodeInstance struct {
	Node
	Id             int
	networkAdapter NetworkAdapter
}

func (n NodeInstance) NetworkSendString(to int, message string) bool {
	return n.networkAdapter.SendString(n, to, message)
}

func (n NodeInstance) NetworkSendWrite(to int, file int, size float64) bool {
	return n.networkAdapter.sendWrite(n, to, file, size)
}

func (n NodeInstance) NetworkSendRead(to int, file int) (bool, float64) {
	return n.networkAdapter.sendRead(n, to, file)
}
