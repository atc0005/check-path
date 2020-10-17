// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Stat("C:/temp/")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
		}
		log.Fatal(err)
	}

	fmt.Println(file.Size())
	fmt.Println(file.IsDir())
	fmt.Printf("%#v\n", file)

}
