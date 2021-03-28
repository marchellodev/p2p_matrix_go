package p2p_matrix

import (
	"encoding/json"
	"github.com/montanaflynn/stats"
	"io/ioutil"
)

type ResultData struct {
	StorageHistory      []ResultStorageStep `json:"storageHistory"`
	StorageHistoryStats StatsData           `json:"storageHistoryStats"`

	UsedNodes      []UsedNodes `json:"usedNodes"`
	UsedNodesStats StatsData   `json:"usedNodesStats"`

	OperationTime      []OperationTime `json:"operationTime"`
	OperationTimeStats StatsData       `json:"operationTimeStats"`

	FileNotFound float64 `json:"fileNotFound"`
}

type StatsData struct {
	Average           float64 `json:"average"`
	Median            float64 `json:"median"`
	Range             float64 `json:"range"`
	StandardDeviation float64 `json:"standardDeviation"`
}

type ResultStorageStep struct {
	Step  int                           `json:"step"`
	State map[int]ResultStorageStepNode `json:"state"`
}

type ResultStorageStepNode struct {
	Storage map[int]float64 `json:"storage"`
	Sum     float64         `json:"sum"`
}

type UsedNodes struct {
	Operation   StoryElementOperation `json:"operation"`
	Connections []int                 `json:"connections"`
}

type OperationTime struct {
	Operation StoryElementOperation `json:"operation"`
	Time      float64               `json:"time"`
}

func (result *ResultData) WriteStorage(step int, allNodes map[int]NodeInstanceInterface) {
	storageStep := ResultStorageStep{}

	storageStep.Step = step
	storageStep.State = make(map[int]ResultStorageStepNode)

	for id, node := range allNodes {
		data := node.SysGetStorage()
		var sum float64

		for _, data := range data {
			sum += data
		}

		storageStep.State[id] = ResultStorageStepNode{Storage: data, Sum: sum}
	}

	result.StorageHistory = append(result.StorageHistory, storageStep)
}

func (result *ResultData) WriteOperation(operation StoryElementOperation, connections []int) {
	result.UsedNodes = append(result.UsedNodes, UsedNodes{Operation: operation, Connections: connections})
}

func (result *ResultData) WriteOperationTime(operation StoryElementOperation, time float64) {
	result.OperationTime = append(result.OperationTime, OperationTime{Operation: operation, Time: time})
}

func (result *ResultData) computeStats() {

	func() {
		var elements []float64
		for _, el := range result.StorageHistory {
			for _, data := range el.State {
				elements = append(elements, data.Sum)
			}
		}
		mean, _ := stats.Mean(elements)
		median, _ := stats.Median(elements)
		max, _ := stats.Max(elements)
		min, _ := stats.Min(elements)
		standardDeviation, _ := stats.StandardDeviation(elements)
		result.StorageHistoryStats = StatsData{
			Average:           mean,
			Median:            median,
			Range:             max - min,
			StandardDeviation: standardDeviation,
		}
	}()

	func() {

		var elements []int
		for _, el := range result.UsedNodes {
			elements = append(elements, len(el.Connections))

		}

		elementsData := stats.LoadRawData(elements)
		mean, _ := stats.Mean(elementsData)
		median, _ := stats.Median(elementsData)
		max, _ := stats.Max(elementsData)
		min, _ := stats.Min(elementsData)
		standardDeviation, _ := stats.StandardDeviation(elementsData)
		result.UsedNodesStats = StatsData{
			Average:           mean,
			Median:            median,
			Range:             max - min,
			StandardDeviation: standardDeviation,
		}
	}()

	func() {

		var elements []float64
		for _, el := range result.OperationTime {
			elements = append(elements, el.Time)
		}

		mean, _ := stats.Mean(elements)
		median, _ := stats.Median(elements)
		max, _ := stats.Max(elements)
		min, _ := stats.Min(elements)
		standardDeviation, _ := stats.StandardDeviation(elements)
		result.OperationTimeStats = StatsData{
			Average:           mean,
			Median:            median,
			Range:             max - min,
			StandardDeviation: standardDeviation,
		}
	}()
}

func (result ResultData) SaveToJson(path string) {

	result.computeStats()

	jsonString, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	_ = ioutil.WriteFile(path, jsonString, 0644)
}
