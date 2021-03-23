package p2p_matrix

type NetworkAdapter struct {
	sendString func(int, int, string)
	sendRead   func(int, int) float64
	sendWrite  func(int, int, float64)
}

func (n NetworkAdapter) SendString(from int, to int, data string) {
	n.sendString(from, to, data)
}

func (n NetworkAdapter) SendRead(user int, file int) float64 {
	return n.sendRead(user, file)
}

func (n NetworkAdapter) SendWrite(user int, file int, data float64) {
	n.sendWrite(user, file, data)
}
