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

import (
	"agones.dev/agones/pkg/apis"
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// GameServerAllocationAllocated is allocation successful
	GameServerAllocationAllocated GameServerAllocationState = "Allocated"
	// GameServerAllocationUnAllocated when the allocation is unsuccessful
	GameServerAllocationUnAllocated GameServerAllocationState = "UnAllocated"
	// GameServerAllocationContention when the allocation is unsuccessful
	// because of contention
	GameServerAllocationContention GameServerAllocationState = "Contention"
)

// GameServerAllocationState is the Allocation state
type GameServerAllocationState string

// +genclient
// +genclient:onlyVerbs=create
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GameServerAllocation is the data structure for allocating against a set of
// GameServers, defined `selectors` selectors
type GameServerAllocation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GameServerAllocationSpec   `json:"spec"`
	Status            GameServerAllocationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GameServerAllocationList is a list of GameServer Allocation resources
type GameServerAllocationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []GameServerAllocation `json:"items"`
}

// GameServerAllocationSpec is the spec for a GameServerAllocation
type GameServerAllocationSpec struct {
	// MultiClusterPolicySelector if specified, multi-cluster policies are applied.
	// Otherwise, allocation will happen locally.
	MultiClusterSetting MultiClusterSetting `json:"multiClusterSetting,omitempty" hash:"ignore"`

	// Deprecated: use field Selectors instead. If Selectors is set, this field is ignored.
	// Required is the GameServer selector from which to choose GameServers from.
	// Defaults to all GameServers.
	Required GameServerSelector `json:"required,omitempty" hash:"ignore"`

	// Deprecated: use field Selectors instead. If Selectors is set, this field is ignored.
	// Preferred is an ordered list of preferred GameServer selectors
	// that are optional to be fulfilled, but will be searched before the `required` selector.
	// If the first selector is not matched, the selection attempts the second selector, and so on.
	// If any of the preferred selectors are matched, the required selector is not considered.
	// This is useful for things like smoke testing of new game servers.
	Preferred []GameServerSelector `json:"preferred,omitempty" hash:"ignore"`

	// [Stage: Beta]
	// [FeatureFlag:CountsAndLists]
	// `Priorities` configuration alters the order in which `GameServers` are searched for matches to the configured `selectors`.
	//
	// Priority of sorting is in descending importance. I.e. The position 0 `priority` entry is checked first.
	//
	// For `Packed` strategy sorting, this priority list will be the tie-breaker within the least utilised infrastructure, to ensure optimal
	// infrastructure usage while also allowing some custom prioritisation of `GameServers`.
	//
	// For `Distributed` strategy sorting, the entire selection of `GameServers` will be sorted by this priority list to provide the
	// order that `GameServers` will be allocated by.
	// +optional
	Priorities []agonesv1.Priority `json:"priorities,omitempty"`

	// Ordered list of GameServer label selectors.
	// If the first selector is not matched, the selection attempts the second selector, and so on.
	// This is useful for things like smoke testing of new game servers.
	// Note: This field can only be set if neither Required or Preferred is set.
	Selectors []GameServerSelector `json:"selectors,omitempty" hash:"ignore"`

	// Scheduling strategy. Defaults to "Packed".
	// +optional
	Scheduling apis.SchedulingStrategy `json:"scheduling"`

	// MetaPatch is optional custom metadata that is added to the game server at allocation
	// You can use this to tell the server necessary session data
	MetaPatch MetaPatch `json:"metadata,omitempty" hash:"ignore"`

	// [Stage: Beta]
	// [FeatureFlag:CountsAndLists]
	// Counter actions to perform during allocation.
	// +optional
	Counters map[string]CounterAction `json:"counters,omitempty" hash:"ignore"`
	// [Stage: Beta]
	// [FeatureFlag:CountsAndLists]
	// List actions to perform during allocation.
	// +optional
	Lists map[string]ListAction `json:"lists,omitempty" hash:"ignore"`
}

// GameServerSelector contains all the filter options for selecting
// a GameServer for allocation.
type GameServerSelector struct {
	// See: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
	metav1.LabelSelector `json:",inline"`
	// GameServerState specifies which State is the filter to be used when attempting to retrieve a GameServer
	// via Allocation. Defaults to "Ready". The only other option is "Allocated", which can be used in conjunction with
	// label/annotation/player selectors to retrieve an already Allocated GameServer.
	GameServerState *agonesv1.GameServerState `json:"gameServerState,omitempty"`
	// [Stage:Alpha]
	// [FeatureFlag:PlayerAllocationFilter]
	// +optional
	// Players provides a filter on minimum and maximum values for player capacity when retrieving a GameServer
	// through Allocation. Defaults to no limits.
	Players *PlayerSelector `json:"players,omitempty"`
	// [Stage: Beta]
	// [FeatureFlag:CountsAndLists]
	// Counters provides filters on minimum and maximum values
	// for a Counter's count and available capacity when retrieving a GameServer through Allocation.
	// Defaults to no limits.
	// +optional
	Counters map[string]CounterSelector `json:"counters,omitempty"`
	// [Stage: Beta]
	// [FeatureFlag:CountsAndLists]
	// Lists provides filters on minimum and maximum values
	// for List capacity, and for the existence of a value in a List, when retrieving a GameServer
	// through Allocation. Defaults to no limits.
	// +optional
	Lists map[string]ListSelector `json:"lists,omitempty"`
}

// PlayerSelector is the filter options for a GameServer based on player counts
type PlayerSelector struct {
	MinAvailable int64 `json:"minAvailable,omitempty"`
	MaxAvailable int64 `json:"maxAvailable,omitempty"`
}

// CounterSelector is the filter options for a GameServer based on the count and/or available capacity.
type CounterSelector struct {
	// MinCount is the minimum current value. Defaults to 0.
	// +optional
	MinCount int64 `json:"minCount"`
	// MaxCount is the maximum current value. Defaults to 0, which translates as max(in64).
	// +optional
	MaxCount int64 `json:"maxCount"`
	// MinAvailable specifies the minimum capacity (current capacity - current count) available on a GameServer. Defaults to 0.
	// +optional
	MinAvailable int64 `json:"minAvailable"`
	// MaxAvailable specifies the maximum capacity (current capacity - current count) available on a GameServer. Defaults to 0, which translates to max(int64).
	// +optional
	MaxAvailable int64 `json:"maxAvailable"`
}

// ListSelector is the filter options for a GameServer based on List available capacity and/or the
// existence of a value in a List.
type ListSelector struct {
	// ContainsValue says to only match GameServers who has this value in the list. Defaults to "", which is all.
	// +optional
	ContainsValue string `json:"containsValue"`
	// MinAvailable specifies the minimum capacity (current capacity - current count) available on a GameServer. Defaults to 0.
	// +optional
	MinAvailable int64 `json:"minAvailable"`
	// MaxAvailable specifies the maximum capacity (current capacity - current count) available on a GameServer. Defaults to 0, which is translated as max(int64).
	// +optional
	MaxAvailable int64 `json:"maxAvailable"`
}

// CounterAction is an optional action that can be performed on a Counter at allocation.
type CounterAction struct {
	// Action must to either "Increment" or "Decrement" the Counter's Count. Must also define the Amount.
	// +optional
	Action *string `json:"action,omitempty"`
	// Amount is the amount to increment or decrement the Count. Must be a positive integer.
	// +optional
	Amount *int64 `json:"amount,omitempty"`
	// Capacity is the amount to update the maximum capacity of the Counter to this number. Min 0, Max int64.
	// +optional
	Capacity *int64 `json:"capacity,omitempty"`
}

// ListAction is an optional action that can be performed on a List at allocation.
type ListAction struct {
	// AddValues appends values to a List's Values array. Any duplicate values will be ignored.
	// +optional
	AddValues []string `json:"addValues,omitempty"`
	// DeleteValues removes values from a List's Values array. Any nonexistant values will be ignored.
	// +optional
	DeleteValues []string `json:"deleteValues,omitempty"`
	// Capacity updates the maximum capacity of the Counter to this number. Min 0, Max 1000.
	// +optional
	Capacity *int64 `json:"capacity,omitempty"`
}

// MultiClusterSetting specifies settings for multi-cluster allocation.
type MultiClusterSetting struct {
	Enabled        bool                 `json:"enabled,omitempty"`
	PolicySelector metav1.LabelSelector `json:"policySelector,omitempty"`
}

// MetaPatch is the metadata used to patch the GameServer metadata on allocation
type MetaPatch struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// GameServerAllocationStatus is the status for an GameServerAllocation resource
type GameServerAllocationStatus struct {
	// GameServerState is the current state of an GameServerAllocation, e.g. Allocated, or UnAllocated
	State          GameServerAllocationState       `json:"state"`
	GameServerName string                          `json:"gameServerName"`
	Ports          []agonesv1.GameServerStatusPort `json:"ports,omitempty"`
	Address        string                          `json:"address,omitempty"`
	Addresses      []corev1.NodeAddress            `json:"addresses,omitempty"`
	NodeName       string                          `json:"nodeName,omitempty"`
	// If the allocation is from a remote cluster, Source is the endpoint of the remote agones-allocator.
	// Otherwise, Source is "local"
	Source   string                            `json:"source"`
	Metadata *GameServerMetadata               `json:"metadata,omitempty"`
	Counters map[string]agonesv1.CounterStatus `json:"counters,omitempty"`
	Lists    map[string]agonesv1.ListStatus    `json:"lists,omitempty"`
}

// GameServerMetadata is the metadata from the allocated game server at allocation time
type GameServerMetadata struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
