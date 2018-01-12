// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register font manager
	gopi.RegisterModule(gopi.Module{
		Name: "rpi/fonts",
		Type: gopi.MODULE_TYPE_FONTS,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(FontManager{}, app.Logger)
		},
	})
}
