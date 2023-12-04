// Basic example of a REST server with several routes, using only the standard library.
package main

import (
	"fmt"

	"github.com/k-zehnder/gophersignal/internals/taskstore"
)

func main() {
	fmt.Println("is working:", taskstore.IsWorking())
}
