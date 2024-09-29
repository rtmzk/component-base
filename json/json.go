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

//go:build !jsoniter

package json

import "encoding/json"

// RawMessage is exported by component-base/json package.
type RawMessage = json.RawMessage

var (
	// Marshal is exported by component-base/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by component-base/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by component-base/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by component-base/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by component-base/json package.
	NewEncoder = json.NewEncoder
)
