package main

import (
	"flag"
	"fmt"

	"github.com/kzahedi/goc3d"
)

func main() {
	inputPtr := flag.String("c3d", "", "Input filename.")

	flag.Parse()
	fmt.Println(*inputPtr)

	goc3d.ReadC3D(*inputPtr, false)
}
