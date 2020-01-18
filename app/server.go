/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type server struct {
	base.App
}

////////////////////////////////////////////////////////////////////////////////
// gopi.App implementation for command-line tool

func NewServer(units ...string) (gopi.App, error) {
	this := new(server)

	// Name of the server
	name := filepath.Base(os.Args[0])

	// Append required units
	units = append(units, "server")

	// Check parameters
	if err := this.App.Init(name, units); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

func (this *server) Run() int {
	if err := this.App.Start(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) == false {
			fmt.Fprintln(os.Stderr, this.App.Flags().Name()+":", err)
			return -1
		} else {

		}
	}

	// Defer closing of instances to exit
	defer func() {
		if err := this.App.Close(); err != nil {
			fmt.Fprintln(os.Stderr, this.App.Flags().Name()+":", err)
		}
	}()

	// Start server and block until done
	if server := this.UnitInstance("server").(gopi.RPCServer); server == nil {
		fmt.Fprintln(os.Stderr, this.App.Flags().Name()+":", gopi.ErrInternalAppError.WithPrefix("server"))
		return -1
	} else {
		fmt.Println("STARTING SERVER")
		if err := server.Start(); err != nil {
			fmt.Fprintln(os.Stderr, this.App.Flags().Name()+":", err)
			return -1
		}
		fmt.Println("STOPPING SERVER")
	}

	// Success
	return 0
}
