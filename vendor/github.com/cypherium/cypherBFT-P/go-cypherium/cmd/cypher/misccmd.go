// Copyright 2016 The cypherBFT Authors
// This file is part of cypherBFT.
//
// cypherBFT is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// cypherBFT is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with cypherBFT. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/cypherium/cypherBFT-P/go-cypherium/cmd/utils"
	"github.com/cypherium/cypherBFT-P/go-cypherium/cph"
	"github.com/cypherium/cypherBFT-P/go-cypherium/params"
	"gopkg.in/urfave/cli.v1"
)

var (
	//	makecacheCommand = cli.Command{
	//		Action:    utils.MigrateFlags(makecache),
	//		Name:      "makecache",
	//		Usage:     "Generate cphash verification cache (for testing)",
	//		ArgsUsage: "<blockNum> <outputDir>",
	//		Category:  "MISCELLANEOUS COMMANDS",
	//		Description: `
	//The makecache command generates an cphash cache in <outputDir>.
	//
	//This command exists to support the system testing project.
	//Regular users do not need to execute it.
	//`,
	//	}
	//	makedagCommand = cli.Command{
	//		Action:    utils.MigrateFlags(makedag),
	//		Name:      "makedag",
	//		Usage:     "Generate cphash mining DAG (for testing)",
	//		ArgsUsage: "<blockNum> <outputDir>",
	//		Category:  "MISCELLANEOUS COMMANDS",
	//		Description: `
	//The makedag command generates an cphash DAG in <outputDir>.
	//
	//This command exists to support the system testing project.
	//Regular users do not need to execute it.
	//`,
	//	}
	versionCommand = cli.Command{
		Action:    utils.MigrateFlags(version),
		Name:      "version",
		Usage:     "Print version numbers",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is supposed to be machine-readable.
`,
	}
	licenseCommand = cli.Command{
		Action:    utils.MigrateFlags(license),
		Name:      "license",
		Usage:     "Display license information",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
	}
	setupCommand = cli.Command{
		Action:    utils.MigrateFlags(setup),
		Name:      "setup",
		Usage:     "Generate cphash verification cache (for testing)",
		ArgsUsage: "<blockNum> <outputDir>",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The setup command generates an cphash cache in <outputDir>.

This command exists to support the system testing project.
Regular users do not need to execute it.
`,
	}
)

// makecache generates an cphash verification cache into the provided folder.
//func makecache(ctx *cli.Context) error {
//	args := ctx.Args()
//	if len(args) != 2 {
//		utils.Fatalf(`Usage: cypher makecache <block number> <outputdir>`)
//	}
//	block, err := strconv.ParseUint(args[0], 0, 64)
//	if err != nil {
//		utils.Fatalf("Invalid block number: %v", err)
//	}
//	cphash.MakeCache(block, args[1])
//
//	return nil
//}

// makedag generates an cphash mining DAG into the provided folder.
//func makedag(ctx *cli.Context) error {
//	args := ctx.Args()
//	if len(args) != 2 {
//		utils.Fatalf(`Usage: cypher makedag <block number> <outputdir>`)
//	}
//	block, err := strconv.ParseUint(args[0], 0, 64)
//	if err != nil {
//		utils.Fatalf("Invalid block number: %v", err)
//	}
//	cphash.MakeDataset(block, args[1])
//
//	return nil
//}

func version(ctx *cli.Context) error {
	fmt.Println(strings.Title(clientIdentifier))
	fmt.Println("Version:", params.Version)
	if gitCommit != "" {
		fmt.Println("Git Commit:", gitCommit)
	}
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Protocol Versions:", cph.ProtocolVersions)
	fmt.Println("Network Id:", cph.DefaultConfig.NetworkId)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}

func license(_ *cli.Context) error {
	fmt.Println(`Cypher is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cypher is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with cypher. If not, see <http://www.gnu.org/licenses/>.`)
	return nil
}

func setup(c *cli.Context) error {
	//napp.InteractiveConfig(suites.MustFind("Ed25519"), "cypher")
	return nil
}
