package main

// todo write all errors to the log file
// todo round all floats
import (
	"C"
	"encoding/json"
	"fmt"
	"p2p_matrix_go/p2p_matrix"
)

var cNetwork = p2p_matrix.Network{NodeFactory: func(n p2p_matrix.NodeInstance) p2p_matrix.NodeInstanceInterface {
	return &MyNode{NodeInstance: n}
}}

func main() {
	fmt.Println("running test")

	cNetwork.RunScript("/home/mark/IdeaProjects/p2p_matrix/storage/scripts/f2.json", "result.json")

}

//export GetModelName
func GetModelName() *C.char {
	return C.CString("test model")
}

//export Run
// script - path to .json script
// path - path to the result file
func Run(script *C.char, path *C.char) {
	// DO NOT EDIT

	go func() {
		cNetwork.RunScript(C.GoString(script), C.GoString(path))
	}()
}

type MyNode struct {
	p2p_matrix.NodeInstance

	Peers        []int
	Storage      map[int]float64
	Bootstrapper int
}

// safe method of adding peers
func (n *MyNode) addPeer(peers ...int) bool {

	addedNew := false

	for _, peer := range peers {

		if n.NodeInstance.Id == peer || n.NodeInstance.Id == 0 {
			continue
		}

		exists := false

		for _, el := range n.Peers {
			if el == peer {
				exists = true
				break
			}
		}

		if !exists {
			n.Peers = append(n.Peers, peer)

			addedNew = true
		}

	}

	return addedNew
}

func (n *MyNode) Activate(bootstrap int) {
	if n.Storage == nil {
		n.Storage = make(map[int]float64)
	}

	if bootstrap == -1 {
		return
	}

	if bootstrap != -2 {
		n.Bootstrapper = bootstrap

		n.addPeer(bootstrap)

		n.NetworkSendString(bootstrap, "give_me_your_peers")

	}

	if bootstrap == -2 {
		success := false

		id := 0
		for success == false {
			success = n.NetworkSendString(n.Peers[id], "give_me_your_peers")
			id++
		}
	}

}

// todo make sure there are no two nodes with the same id
func (n *MyNode) Listen(from int, message string) {

	if message == "give_me_your_peers" {

		// if we received it - we bootstrapped someone
		// we give that node our peers and share this node with our bootstrapper

		n.addPeer(from)
		//n.Peers = append(n.Peers, from)

		jsonString, err := json.Marshal(n.Peers)
		if err != nil {
			panic(err)
		}

		n.NetworkSendString(from, string(jsonString))

		success := false
		id := 0
		for success == false {
			success = n.NetworkSendString(id, string(jsonString))
			id++
		}

		for el, data := range n.Storage {
			n.NetworkSendWrite(from, el, data)
		}

	} else {

		var results []int
		err := json.Unmarshal([]byte(message), &results)
		if err != nil {
			panic(err)
		}

		added := n.addPeer(results...)

		if !added {
			return
		}

		jsonString, err := json.Marshal(n.Peers)
		if err != nil {
			panic(err)
		}

		for _, node := range n.Peers {
			n.NetworkSendString(node, string(jsonString))
		}

	}
}

// todo wtf
func (n MyNode) Read(from int, file int) float64 {

	if from != -1 {
		return 0
	}

	for _, node := range n.Peers {
		n.NetworkSendRead(node, file)
	}

	for el, val := range n.Storage {
		if el == file {
			return val
		}
	}

	return 0
}

/*
Sig:
213005.21219999684
9702

Sig:
213005.21219999678
9702

213005.2121999969
9702

*/
// todo make sure we use pointer when we want to change the object
// todo wtf
func (n *MyNode) Write(from int, file int, data float64) {

	alreadyHave := false

	for el, _ := range n.Storage {
		if el == file {
			alreadyHave = true
			// todo use brake for performance reasons
			break
		}
	}

	if !alreadyHave {
		n.Storage[file] = data
	}

	if from == -1 {
		for _, node := range n.Peers {
			n.NetworkSendWrite(node, file, data)
		}
	}

	//n.addPeer(from)
	//
	//if n.Storage == nil {
	//	n.Storage = make(map[int]float64)
	//}
	//
	//for key, _ := range n.Storage {
	//	if key == file {
	//		return
	//	}
	//}
	//
	//n.Storage[file] = data
	//
	//// todo check if already exists
	//
	//for _, peer := range n.Peers {
	//	n.NetworkSendWrite(peer, file, data)
	//}

}

func (n MyNode) SysGetPeers() []int {
	return n.Peers
}

func (n MyNode) SysGetStorage() map[int]float64 {
	//var result float64
	//
	//for _, size := range n.Storage {
	//	result += size
	//}

	return n.Storage
}
