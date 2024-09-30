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

package shutdown

import "sync"

type ShutdownCallback interface {
	OnShutdown(managerName string) error
}

type ShutdownFunc func(string) error

func (f ShutdownFunc) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

type ShutdownManager interface {
	GetName() string
	Start(gs GSInterface) error
	ShutdownStart() error
	ShutdownFinished() error
}

type ErrorHandler interface {
	HandleError(err error)
}

type ErrorFunc func(err error)

func (e ErrorFunc) HandleError(err error) {
	e(err)
}

type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallback ShutdownCallback)
}

type GracefulShutdown struct {
	callbacks    []ShutdownCallback
	managers     []ShutdownManager
	errorHandler ErrorHandler
}

func (gs *GracefulShutdown) New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]ShutdownCallback, 10),
		managers:  make([]ShutdownManager, 3),
	}
}

func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}

	return nil
}

func (gs *GracefulShutdown) AddShutdownManager(manager ShutdownManager) {
	gs.managers = append(gs.managers, manager)
}

func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

func (gs *GracefulShutdown) SetErrorHandler(handler ErrorHandler) {
	gs.errorHandler = handler
}

func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	gs.ReportError(sm.ShutdownStart())

	var wg sync.WaitGroup

	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallback ShutdownCallback) {
			defer wg.Done()

			gs.ReportError(shutdownCallback.OnShutdown(sm.GetName()))
		}(shutdownCallback)
	}

	wg.Wait()
	gs.ReportError(sm.ShutdownFinished())
}

func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.HandleError(err)
	}
}
