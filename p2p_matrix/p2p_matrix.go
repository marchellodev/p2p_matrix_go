package p2p_matrix

import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
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

func runStory(script ScriptModel, network Network) {

	allNodes := make(map[int]NodeInstanceInterface)
	var activeNodes []int

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

		return true, allNodes[to].Read(from.Id, file)
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
					}
				}

				activeNodes = append(activeNodes, action.Node)

				if !exists {

					//for i := range activeNodes {
					//	j := rand.Intn(i + 1)
					//	activeNodes[i], activeNodes[j] = activeNodes[j], activeNodes[i]
					//}

					//rand.Shuffle(len(activeNodes), func(i, j int){
					//	activeNodes[i], activeNodes[j] = activeNodes[j], activeNodes[i]
					//})

					// create the node
					allNodes[action.Node] = newNode(network, script, action.Node, networkAdapter)

					if len(activeNodes) == 1 {

						allNodes[action.Node].Activate(-1)
					} else {
						// get random node to bootstrap from

						allNodes[action.Node].Activate(activeNodes[len(activeNodes)-2])
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
			if operation.Type == "write" {
				allNodes[operation.Node].Write(-1, operation.File, script.Files[operation.File].Size)
				// todo check if the node is active ?
			} else {
				read := allNodes[operation.Node].Read(0, operation.File)
				if read != script.Files[operation.File].Size {
					fmt.Println("FUCK YOU")
					fmt.Println(operation.File)
					fmt.Println(read)
					fmt.Println(script.Files[operation.File].Size)
					fmt.Println(operation.Node)

					for node, _ := range allNodes {
						fmt.Print(strconv.Itoa(node) + " ")
						fmt.Print(allNodes[node].SysGetStorage())
						fmt.Print(" ")
						fmt.Println(allNodes[node].SysGetPeers())
						//
					}
					return
				}

			}
		}

		// 23 has 82
		// 82 has
		if id == -2 {
			for _, node := range activeNodes {
				fmt.Print(strconv.Itoa(node) + " ")
				fmt.Print(allNodes[node].SysGetPeers())
				fmt.Print(" ")
				fmt.Println(allNodes[node].SysGetStorage())
				//
			}
			return
		}

		// todo: fix an example DHT algorithm & finish the runner + adapter limitations + statistics

	}

	/*

	 */
	fmt.Println("Running Done")
	//fmt.Println(activeNodes)
	//
	var storage float64
	var peers int
	//
	for _, node := range allNodes {
		st := node.SysGetStorage()
		for _, el := range st {
			storage += el
		}

		peers += len(node.SysGetPeers())
		//sum += node.SysGetStorage()

		//fmt.Print(strconv.Itoa(id) + " ")
		//fmt.Println(node.SysGetStorage())
	}
	// 2.4172943228197068e+07
	// 2.4172943228196982e+07
	fmt.Println("Sig:")
	fmt.Println(storage)
	fmt.Println(peers)

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

func (n Network) RunScript(scriptPath string, ResultPath string) {

	script := deserializeScriptFile(scriptPath)

	runStory(script, n)

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
