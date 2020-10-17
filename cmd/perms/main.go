// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/phayes/permbits"
	"github.com/rs/zerolog/log"
)

func main() {

	/*
		- require presence
		  - `group-write`
		  - `group-read`
		  - `group-execute`

		- assert not present
		  - `no-group-write`
		  - `no-group-read`
		  - `no-group-execute`
		  - `no-other-write`
		  - `no-other-read`
		  - `no-other-execute`

	*/

	// testModeValue := "-rw-rw-rw-"

	file, err := os.Stat(filepath.Join("/tmp", "go1.15.3.linux-amd64.tar.gz"))
	// file, err := os.Stat("vscode-inno-updater-1603112887.log")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal().Err(err).Msg("File does not exist.")
		}
		log.Fatal().Err(err).Msg("failed to stat file; unknown error occurred")
	}

	fmt.Println(file.Mode())
	// fmt.Println(testModeValue)

	perms := permbits.FileMode(file.Mode())

	if !perms.GroupWrite() {
		fmt.Println("group cannot write")
	}
	if !perms.OtherWrite() {
		fmt.Println("other cannot write")
	}

	if !perms.UserWrite() {
		fmt.Println("user cannot write")
	}

}
