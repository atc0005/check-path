//go:build !windows

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
	"os/user"
	"strconv"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	// file, err := os.Stat(filepath.Join("/tmp", "vscode-inno-updater-1603112887.log"))
	file, err := os.Stat("/tmp/go1.15.3.linux-amd64.tar.gz")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal().Err(err).Msg("file does not exist.")
		}
		log.Fatal().Err(err).Msg("failed to stat file; unknown error occurred")
	}

	// fmt.Println(file.Mode())
	// fmt.Printf("%#v", file.Sys())

	stat, statOK := file.Sys().(*syscall.Stat_t)
	if !statOK {
		log.Fatal().Msg("failed to access syscall.Stat_t")
	}

	UID := int(stat.Uid)
	GID := int(stat.Gid)

	fmt.Println("UID:", UID)
	fmt.Println("GID:", GID)

	gidNum := strconv.Itoa(GID)
	group, err := user.LookupGroupId(gidNum)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to lookup group for gid")
	}
	fmt.Println("group name:", group.Name)

	uidNum := strconv.Itoa(UID)
	user, err := user.LookupId(uidNum)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to lookup user for uid")
	}
	fmt.Println("user name:", user.Username)

}
