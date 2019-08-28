package main

import (
	"os"
	"fmt"
	"io/ioutil"
)

func main() {
	name := os.Args[2]

	var outFilename string
	if len(os.Args) >= 4 {
		outFilename = os.Args[3]
	} else {
		outFilename = name+".go"
	}

	

	data,_ := ioutil.ReadFile(os.Args[1])
	for _,datum := range data {

	}
}