/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runtime

// Note that the types provided in this file are not versioned and are intended to be
// safe to use from within all versions of every API object.

// TypeMeta is shared by all top level objects. The proper way to use it is to inline it in your type,
// like this:
//
//	type MyAwesomeAPIObject struct {
//	     runtime.TypeMeta    `json:",inline"`
//	     ... // other fields
//	}
//
// func (obj *MyAwesomeAPIObject) SetGroupVersionKind(gvk *metav1.GroupVersionKind) { metav1.UpdateTypeMeta(obj,gvk) }; GroupVersionKind() *GroupVersionKind
//
// TypeMeta is provided here for convenience. You may use it directly from this package or define
// your own with the same fields.
//
// +k8s:deepcopy-gen=false
// +protobuf=true
// +k8s:openapi-gen=true
type TypeMeta struct {
	// +optional
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty" protobuf:"bytes,1,opt,name=apiVersion"`
	// +optional
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty" protobuf:"bytes,2,opt,name=kind"`
}

// Unknown allows api objects with unknown types to be passed-through. This can be used
// to deal with the API objects from a plug-in. Unknown objects still have functioning
// TypeMeta features-- kind, version, etc.
// TODO: Make this object have easy access to field based accessors and settors for
// metadata and field mutatation.
//
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +protobuf=true
// +k8s:openapi-gen=true
type Unknown struct {
	TypeMeta `json:",inline" protobuf:"bytes,1,opt,name=typeMeta"`
	// Raw will hold the complete serialized object which couldn't be matched
	// with a registered type. Most likely, nothing should be done with this
	// except for passing it through the system.
	Raw []byte `json:"-" protobuf:"bytes,2,opt,name=raw"`
	// ContentEncoding is encoding used to encode 'Raw' data.
	// Unspecified means no encoding.
	ContentEncoding string `protobuf:"bytes,3,opt,name=contentEncoding"`
	// ContentType  is serialization method used to serialize 'Raw'.
	// Unspecified means ContentTypeJSON.
	ContentType string `protobuf:"bytes,4,opt,name=contentType"`
}
