package observable

import (
	"fmt"
)

type MyObserver struct{}

func (o *MyObserver) Update(message string) {
	fmt.Printf("Received message: %s\n", message)
}