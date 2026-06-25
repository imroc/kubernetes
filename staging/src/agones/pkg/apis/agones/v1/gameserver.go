// Copyright 2017 Google LLC All Rights Reserved.
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
	"agones.dev/agones/pkg/apis/agones"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GameServerState is the state for the GameServer
type GameServerState string

const (
	// GameServerStatePortAllocation is for when a dynamically allocating GameServer
	// is being created, an open port needs to be allocated
	GameServerStatePortAllocation GameServerState = "PortAllocation"
	// GameServerStateCreating is before the Pod for the GameServer is being created
	GameServerStateCreating GameServerState = "Creating"
	// GameServerStateStarting is for when the Pods for the GameServer are being
	// created but are not yet Scheduled
	GameServerStateStarting GameServerState = "Starting"
	// GameServerStateScheduled is for when we have determined that the Pod has been
	// scheduled in the cluster -- basically, we have a NodeName
	GameServerStateScheduled GameServerState = "Scheduled"
	// GameServerStateRequestReady is when the GameServer has declared that it is ready
	GameServerStateRequestReady GameServerState = "RequestReady"
	// GameServerStateReady is when a GameServer is ready to take connections
	// from Game clients
	GameServerStateReady GameServerState = "Ready"
	// GameServerStateShutdown is when the GameServer has shutdown and everything needs to be
	// deleted from the cluster
	GameServerStateShutdown GameServerState = "Shutdown"
	// GameServerStateError is when something has gone wrong with the Gameserver and
	// it cannot be resolved
	GameServerStateError GameServerState = "Error"
	// GameServerStateUnhealthy is when the GameServer has failed its health checks
	GameServerStateUnhealthy GameServerState = "Unhealthy"
	// GameServerStateReserved is for when a GameServer is reserved and therefore can be allocated but not removed
	GameServerStateReserved GameServerState = "Reserved"
	// GameServerStateAllocated is when the GameServer has been allocated to a session
	GameServerStateAllocated GameServerState = "Allocated"
)

// PortPolicy is the port policy for the GameServer
type PortPolicy string

const (
	// Static PortPolicy means that the user defines the hostPort to be used
	// in the configuration.
	Static PortPolicy = "Static"
	// Dynamic PortPolicy means that the system will choose an open
	// port for the GameServer in question
	Dynamic PortPolicy = "Dynamic"
	// Passthrough dynamically sets the `containerPort` to the same value as the dynamically selected hostPort.
	// This will mean that users will need to lookup what port has been opened through the server side SDK.
	Passthrough PortPolicy = "Passthrough"
	// None means the `hostPort` is ignored and if defined, the `containerPort` (optional) is used to set the port on the GameServer instance.
	None PortPolicy = "None"
)

// EvictionSafe specified whether the game server supports termination via SIGTERM
type EvictionSafe string

const (
	// EvictionSafeAlways means the game server supports termination via SIGTERM, and wants eviction signals
	// from Cluster Autoscaler scaledown and node upgrades.
	EvictionSafeAlways EvictionSafe = "Always"
	// EvictionSafeOnUpgrade means the game server supports termination via SIGTERM, and wants eviction signals
	// from node upgrades, but not Cluster Autoscaler scaledown.
	EvictionSafeOnUpgrade EvictionSafe = "OnUpgrade"
	// EvictionSafeNever means the game server should run to completion and may not understand SIGTERM. Eviction
	// from ClusterAutoscaler and upgrades should both be blocked.
	EvictionSafeNever EvictionSafe = "Never"
)

// SdkServerLogLevel is the log level for SDK server (sidecar) logs
type SdkServerLogLevel string

const (
	// SdkServerLogLevelInfo will cause the SDK server to output all messages except for debug messages.
	SdkServerLogLevelInfo SdkServerLogLevel = "Info"
	// SdkServerLogLevelDebug will cause the SDK server to output all messages including debug messages.
	SdkServerLogLevelDebug SdkServerLogLevel = "Debug"
	// SdkServerLogLevelError will cause the SDK server to only output error messages.
	SdkServerLogLevelError SdkServerLogLevel = "Error"
	// SdkServerLogLevelTrace will cause the SDK server to output all messages, including detailed tracing information.
	SdkServerLogLevelTrace SdkServerLogLevel = "Trace"
)

const (
	// ProtocolTCPUDP Protocol exposes the hostPort allocated for this container for both TCP and UDP.
	ProtocolTCPUDP corev1.Protocol = "TCPUDP"

	// DefaultPortRange is the name of the default port range.
	DefaultPortRange = "default"

	// RoleLabel is the label in which the Agones role is specified.
	// Pods from a GameServer will have the value "gameserver"
	RoleLabel = agones.GroupName + "/role"
	// GameServerLabelRole is the GameServer label value for RoleLabel
	GameServerLabelRole = "gameserver"
	// GameServerPodLabel is the label that the name of the GameServer
	// is set on the Pod the GameServer controls
	GameServerPodLabel = agones.GroupName + "/gameserver"
	// GameServerPortPolicyPodLabel is the label to identify the port policy
	// of the pod
	GameServerPortPolicyPodLabel = agones.GroupName + "/port"
	// GameServerContainerAnnotation is the annotation that stores
	// which container is the container that runs the dedicated game server
	GameServerContainerAnnotation = agones.GroupName + "/container"
	// DevAddressAnnotation is an annotation to indicate that a GameServer hosted outside of Agones.
	// A locally hosted GameServer is not managed by Agones it is just simply registered.
	DevAddressAnnotation = "agones.dev/dev-address"
	// GameServerReadyContainerIDAnnotation is an annotation that is set on the GameServer
	// becomes ready, so we can track when restarts should occur and when a GameServer
	// should be moved to Unhealthy.
	GameServerReadyContainerIDAnnotation = agones.GroupName + "/ready-container-id"
	// PodSafeToEvictAnnotation is an annotation that the Kubernetes cluster autoscaler uses to
	// determine if a pod can safely be evicted to compact a cluster by moving pods between nodes
	// and scaling down nodes.
	PodSafeToEvictAnnotation = "cluster-autoscaler.kubernetes.io/safe-to-evict"
	// SafeToEvictLabel is a label that, when "false", matches the restrictive PDB agones-gameserver-safe-to-evict-false.
	SafeToEvictLabel = agones.GroupName + "/safe-to-evict"
	// GameServerErroredAtAnnotation is an annotation that records the timestamp the GameServer entered the
	// error state. The timestamp is encoded in RFC3339 format.
	GameServerErroredAtAnnotation = agones.GroupName + "/errored-at"
	// FinalizerName is the domain name and finalizer path used to manage garbage collection of the GameServer.
	FinalizerName = agones.GroupName + "/controller"

	// NodePodIP identifies an IP address from a pod.
	NodePodIP corev1.NodeAddressType = "PodIP"

	// PassthroughPortAssignmentAnnotation is an annotation to keep track of game server container and its Passthrough ports indices
	PassthroughPortAssignmentAnnotation = "agones.dev/container-passthrough-port-assignment"

	// True is the string "true" to appease the goconst lint.
	True = "true"
	// False is the string "false" to appease the goconst lint.
	False = "false"
)

// PlayersSpec tracks the initial player capacity
type PlayersSpec struct {
	InitialCapacity int64 `json:"initialCapacity,omitempty"`
}

// Eviction specifies the eviction tolerance of the GameServer
type Eviction struct {
	// Game server supports termination via SIGTERM:
	// - Always: Allow eviction for both Cluster Autoscaler and node drain for upgrades
	// - OnUpgrade: Allow eviction for upgrades alone
	// - Never (default): Pod should run to completion
	Safe EvictionSafe `json:"safe,omitempty"`
}

// Health configures health checking on the GameServer
type Health struct {
	// Disabled is whether health checking is disabled or not
	Disabled bool `json:"disabled,omitempty"`
	// PeriodSeconds is the number of seconds each health ping has to occur in
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
	// FailureThreshold how many failures in a row constitutes unhealthy
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
	// InitialDelaySeconds initial delay before checking health
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`
}

// GameServerPort defines a set of Ports that
// are to be exposed via the GameServer
type GameServerPort struct {
	// Name is the descriptive name of the port
	Name string `json:"name,omitempty"`
	// (Alpha, PortRanges feature flag) Range is the port range name from which to select a port when using a
	// 'Dynamic' or 'Passthrough' port policy.
	// +optional
	Range string `json:"range,omitempty"`
	// PortPolicy defines the policy for how the HostPort is populated.
	// Dynamic port will allocate a HostPort within the selected MIN_PORT and MAX_PORT range passed to the controller
	// at installation time.
	// When `Static` portPolicy is specified, `HostPort` is required, to specify the port that game clients will
	// connect to
	// `Passthrough` dynamically sets the `containerPort` to the same value as the dynamically selected hostPort.
	// `None` portPolicy ignores `HostPort` and the `containerPort` (optional) is used to set the port on the GameServer instance.
	PortPolicy PortPolicy `json:"portPolicy,omitempty"`
	// Container is the name of the container on which to open the port. Defaults to the game server container.
	// +optional
	Container *string `json:"container,omitempty"`
	// ContainerPort is the port that is being opened on the specified container's process
	ContainerPort int32 `json:"containerPort,omitempty"`
	// HostPort the port exposed on the host for clients to connect to
	HostPort int32 `json:"hostPort,omitempty"`
	// Protocol is the network protocol being used. Defaults to UDP. TCP and TCPUDP are other options.
	Protocol corev1.Protocol `json:"protocol,omitempty"`
}

// SdkServer specifies parameters for the Agones SDK Server sidecar container
type SdkServer struct {
	// LogLevel for SDK server (sidecar) logs. Defaults to "Info"
	LogLevel SdkServerLogLevel `json:"logLevel,omitempty"`
	// GRPCPort is the port on which the SDK Server binds the gRPC server to accept incoming connections
	GRPCPort int32 `json:"grpcPort,omitempty"`
	// HTTPPort is the port on which the SDK Server binds the HTTP gRPC gateway server to accept incoming connections
	HTTPPort int32 `json:"httpPort,omitempty"`
}

// GameServerStatus is the status for a GameServer resource
type GameServerStatus struct {
	// GameServerState is the current state of a GameServer, e.g. Creating, Starting, Ready, etc
	State   GameServerState        `json:"state"`
	Ports   []GameServerStatusPort `json:"ports"`
	Address string                 `json:"address"`
	// Addresses is the array of addresses at which the GameServer can be reached; copy of Node.Status.addresses.
	// +optional
	Addresses     []corev1.NodeAddress `json:"addresses"`
	NodeName      string               `json:"nodeName"`
	ReservedUntil *metav1.Time         `json:"reservedUntil"`
	// [Stage:Alpha]
	// [FeatureFlag:PlayerTracking]
	// +optional
	Players *PlayerStatus `json:"players"`
	// (Beta, CountsAndLists feature flag) Counters and Lists provides the configuration for generic tracking features.
	// +optional
	Counters map[string]CounterStatus `json:"counters,omitempty"`
	// +optional
	Lists map[string]ListStatus `json:"lists,omitempty"`
	// Eviction specifies the eviction tolerance of the GameServer.
	// +optional
	Eviction *Eviction `json:"eviction,omitempty"`
	// immutableReplicas is present in gameservers.agones.dev but omitted here (it's always 1).
}

// GameServerStatusPort shows the port that was allocated to a
// GameServer.
type GameServerStatusPort struct {
	Name string `json:"name,omitempty"`
	Port int32  `json:"port"`
}

// PlayerStatus stores the current player capacity values
type PlayerStatus struct {
	Count    int64    `json:"count"`
	Capacity int64    `json:"capacity"`
	IDs      []string `json:"ids"`
}

// CounterStatus stores the current counter values and maximum capacity
type CounterStatus struct {
	Count    int64 `json:"count"`
	Capacity int64 `json:"capacity"`
}

// ListStatus stores the current list values and maximum capacity
type ListStatus struct {
	Capacity int64    `json:"capacity"`
	Values   []string `json:"values"`
}
