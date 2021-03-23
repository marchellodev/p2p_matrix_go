package p2p_matrix

import "C"
import (
	"encoding/json"
	"io/ioutil"
)

type Network struct {
	// here we store files and their sizes
	Storage *map[int]map[int]float64

	// an array of all nodes we know about
	Peers *map[int][]int

	// when node joins the network
	// bootstrap_id, node_id, adapter
	Activate func(int, int, NetworkAdapter)

	// this is called when someone sends something to this node via NetworkAdapter
	// from, to, data, adapter
	Listen func(int, int, string, NetworkAdapter)

	// here we (int = 0) or other node asks to read file with id int
	// the function returns the size of the file
	Read func(int, int, NetworkAdapter) float64

	// here we (int = 0) or other node asks to write a file with an id int
	// the function returns the size of the file
	Write func(int, int, float64, NetworkAdapter)
}

func deserializeScriptFile(scriptPath string) ScriptModel {
	data, _ := ioutil.ReadFile(scriptPath)

	var script ScriptModel
	_ = json.Unmarshal(data, &script)

	return script
}

func runStory(script ScriptModel, network Network) {
	var activeNodes []int
	var networkAdapter NetworkAdapter

	// 	adapter.SendString(bootstrap, "give_me_your_peers")

	networkAdapter.sendString = func(from int, user int, data string) {
		// todo maybe should be a pointer ?
		network.Listen(user, data, networkAdapter)
	}

	//networkAdapter.sendString = func(a int, b string) {}
	// todo

	for _, storyElement := range script.Story {

		for id, action := range storyElement.NodeActions {
			if action == "on" {
				network.Activate(id, networkAdapter)
				activeNodes = append(activeNodes, id)
			} else {
				activeNodes = remove(activeNodes, id)
			}
		}

	}
}

func remove(s []int, i int) []int {
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
