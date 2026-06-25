// Copyright 2019 Google LLC All Rights Reserved.
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

package v1

// Priority is a sorting option for GameServers with Counters or Lists based on the available capacity,
// i.e. the current Capacity value, minus either the Count value or List length.
type Priority struct {
	// Type: Sort by a "Counter" or a "List".
	Type string `json:"type"`
	// Key: The name of the Counter or List. If not found on the GameServer, has no impact.
	Key string `json:"key"`
	// Order: Sort by "Ascending" or "Descending". "Descending" a bigger available capacity is preferred.
	// "Ascending" would be smaller available capacity is preferred.
	// The default sort order is "Ascending"
	Order string `json:"order"`
}
