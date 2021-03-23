package main

// todo write all errors to the log file
import (
	"C"
	"encoding/json"
	"fmt"
	"p2p_matrix_go/p2p_matrix"
)

// todo think about potential for abuse (you can't build olympiads with that lol)
// todo are you sure boy lol ???
var storage map[int]map[int]float64
var peers map[int][]int

var cNetwork = p2p_matrix.Network{Storage: &storage, Peers: &peers, Activate: activate, Listen: listen, Read: nil, Write: nil}

func activate(bootstrap int, node int, adapter p2p_matrix.NetworkAdapter) {
	if bootstrap == 0 {
		return
	}

	peers[node] = append(peers[node], bootstrap)
	adapter.SendString(bootstrap, "give_me_your_peers")
}

func listen(from int, to int, data string, adapter p2p_matrix.NetworkAdapter) {
	if data == "give_me_your_peers" {
		jsonString, err := json.Marshal(peers)
		if err != nil {
			panic(err)
		}
		adapter.SendString(from, string(jsonString))
	} else {
		var results []int
		err := json.Unmarshal([]byte(data), &results)
		if err != nil {
			panic(err)
		}

		peers = append(peers, results...)
	}

}

//func read(from int, file int, adapter p2p_matrix.NetworkAdapter) float64 {
//	if from == 0 {
//		for _, peer := range peers {
//			if data := adapter.SendRead(peer, file); data != 0 {
//				return data
//			}
//		}
//		return 0
//	} else {
//		if val, ok := storage[file]; ok {
//			return val
//		} else {
//			return 0
//		}
//	}
//}
//
//func write(from int, file int, data float64, adapter p2p_matrix.NetworkAdapter) {
//	storage[file] = data
//
//	if from == 0 {
//		for _, peer := range peers {
//			adapter.SendWrite(peer, file, data)
//		}
//	}
//}

func main() {
	fmt.Println("running test")

	cNetwork.RunScript("D:\\ctemp\\IdeaProjects\\p2p_model\\storage\\scripts\\test2.json", "")
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
