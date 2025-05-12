package lechain

/*
messages between lechain and other entities
*/

// return a handle to listen to the chain block changes
// (signals when the head of the chain changes)
func (lc *LeChain) RegisterListener() chan bool {
	signal := make(chan bool, 1)
	lc.handles = append(lc.handles, signal)
	return signal
}

func (lc *LeChain) Broadcast() {
	for _, c := range lc.handles {
		c <- true
	}
}

// TODO: (extension) method to remove handles
