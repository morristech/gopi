// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package i2c_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/i2c"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/platform"
)

func Test_I2C_000(t *testing.T) {
	t.Log("Test_I2C_000")
}

func Test_I2C_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_I2C_001, []string{"-debug"}, "platform", "i2c"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_I2C_001(app gopi.App, t *testing.T) {
	// Don't test unless on Linux
	if platform := app.Platform(); platform.Type()&gopi.PLATFORM_LINUX == 0 {
		t.Log("Skipping testing of I2C on", platform.Type())
	} else {
		i2c := app.UnitInstance("i2c").(gopi.I2C)
		if i2c == nil {
			t.Fatal(gopi.ErrInternalAppError.WithPrefix("UnitInstance() failed"))
		} else {
			t.Log(i2c)
		}

	}
}