// Copyright 2024 rtmzk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package posixsignal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rtmzk/component-base/shutdown"
)

// Name defines PosixSignalManager name
const Name = "PosixSignalManager"

// PosixSignalManager implements ShutdownManager interface.
// Initialize it with NewPosixSignalManager.
type PosixSignalManager struct {
	signal []os.Signal
}

// NewPosixSignalManager initializes the PosixSignalManager.
// You can provide os.Signal-s to listen for, as arguments. If none are given
// it will default to SIGINT and SIGTERM
func NewPosixSignalManager(sigs ...os.Signal) *PosixSignalManager {
	if len(sigs) == 0 {
		sigs = make([]os.Signal, 2)
		sigs[0] = os.Interrupt
		sigs[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{
		signal: sigs,
	}
}

// GetName return posix signal shutdown manager name
func (posixSignalManager *PosixSignalManager) GetName() string { return Name }

// Start starts listening for posix os.Signal
func (posixSignalManager *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)

		signal.Notify(c, posixSignalManager.signal...)

		// block until receive a signal
		<-c

		gs.StartShutdown(posixSignalManager)
	}()

	return nil
}

// ShutdownStart does nothing
func (posixSignalManager *PosixSignalManager) ShutdownStart() error {
	return nil
}

// ShutdownFinished exit the app with os.Exit(0)
func (posixSignalManager *PosixSignalManager) ShutdownFinished() error {
	os.Exit(0)

	return nil
}
