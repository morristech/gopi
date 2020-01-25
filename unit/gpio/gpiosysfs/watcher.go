// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpiosysfs

import (
	"fmt"
	"os"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Watcher struct {
	filepoll gopi.FilePoll
	pins     map[uintptr]Pin

	sync.RWMutex
}

type Pin struct {
	edge gopi.GPIOEdge
	file *os.File
	pin  gopi.GPIOPin
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

func (this *Watcher) Init(config GPIO) error {
	this.Lock()
	defer this.Unlock()

	if config.FilePoll == nil {
		return gopi.ErrBadParameter.WithPrefix("FilePoll")
	} else {
		this.filepoll = config.FilePoll
	}
	this.pins = make(map[uintptr]Pin, 0)
	return nil
}

func (this *Watcher) Close() error {
	this.Lock()
	defer this.Unlock()

	errs := gopi.NewCompoundError()
	for fd, pin := range this.pins {
		errs.Add(this.filepoll.Unwatch(fd))
		errs.Add(pin.file.Close())
	}

	// Release resources
	this.filepoll = nil
	this.pins = nil

	// Return any errors
	return errs.ErrorOrSelf()
}

func (this *Watcher) Exists(logical gopi.GPIOPin) bool {
	this.RLock()
	defer this.RUnlock()

	return this.fileForPin(logical) != nil
}

func (this *Watcher) Watch(logical gopi.GPIOPin, file *os.File, edge gopi.GPIOEdge) error {
	this.Lock()
	defer this.Unlock()

	if this.fileForPin(logical) != nil {
		return gopi.ErrDuplicateItem.WithPrefix("Watch")
	} else if err := this.filepoll.Watch(file.Fd(), gopi.FILEPOLL_FLAG_READ, func(handle uintptr, _ gopi.FilePollFlags) {
		this.RLock()
		defer this.RUnlock()
		if pin, exists := this.pins[handle]; exists {
			this.handleEdge(pin)
		}
	}); err != nil {
		defer file.Close()
		return fmt.Errorf("Watch: Unable to watch %v: %w", logical, err)
	} else if _, exists := this.pins[file.Fd()]; exists {
		return gopi.ErrInternalAppError.WithPrefix("Watch")
	} else {
		this.pins[file.Fd()] = Pin{edge, file, logical}
	}
	// Return success
	return nil
}

func (this *Watcher) Unwatch(logical gopi.GPIOPin) error {
	this.Lock()
	defer this.Unlock()

	if file := this.fileForPin(logical); file == nil {
		return gopi.ErrNotFound.WithPrefix("Watch")
	} else if err := this.filepoll.Unwatch(file.Fd()); err != nil {
		return err
	} else {
		delete(this.pins, file.Fd())
		defer file.Close()
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Watcher) fileForPin(logical gopi.GPIOPin) *os.File {
	for _, pin := range this.pins {
		if pin.pin == logical {
			return pin.file
		}
	}
	return nil
}

func (this *Watcher) handleEdge(pin Pin) error {
	readBuffer := make([]byte, 1)
	pin.file.Seek(0, 0)
	if _, err := pin.file.Read(readBuffer); err != nil {
		return err
	} else {
		fmt.Println(readBuffer[0])
	}
	/*
			state := int(readBuffer[0] & 1)
		}
		if err != nil {
			return -1, err
		}
		return state, nil

		if _, err := pin.file.Seek(0, io.SeekStart); err != nil {
			return err
		} else if buf, err := pin.file.Read(pin.file); err != nil {
			return err
		} else {
			value := strings.TrimSpace(string(buf))
			fmt.Println(value)
			switch value {
			case "0":
				fmt.Println(pin.pin, pin.edge, gopi.GPIO_EDGE_FALLING)
			case "1":
				fmt.Println(pin.pin, pin.edge, gopi.GPIO_EDGE_RISING)
			}
		}
	*/

	// Return success
	return nil
}
