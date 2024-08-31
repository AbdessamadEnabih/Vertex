package state

import (
	"fmt"
	"sync"
)

type State struct {
	data map[string]string
	mu   sync.Mutex
}

func Insert() {
	fmt.Print("state")
}
