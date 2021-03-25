package p2p_matrix

type NetworkAdapter struct {
	sendString func(NodeInstance, int, string) bool
	sendRead   func(NodeInstance, int, int) (bool, float64)
	sendWrite  func(NodeInstance, int, int, float64) bool
}

func (n NetworkAdapter) SendString(node NodeInstance, to int, data string) bool {

	return n.sendString(node, to, data)
}

func (n NetworkAdapter) SendRead(node NodeInstance, user int, file int) (bool, float64) {
	return n.sendRead(node, user, file)
}

func (n NetworkAdapter) SendWrite(node NodeInstance, user int, file int, data float64) bool {
	return n.sendWrite(node, user, file, data)
}
