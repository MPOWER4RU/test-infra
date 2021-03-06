/*
Copyright 2018 The Knative Authors

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

package resourcetree

import (
	"container/list"
	"reflect"
)

// ResourceTree encapsulates a tree corresponding to a resource type.
type ResourceTree struct {
	ResourceName string
	Root NodeInterface
	forest *ResourceForest
}

func (r *ResourceTree) createNode(field string, parent NodeInterface, t reflect.Type) NodeInterface {
	var n NodeInterface
	switch t.Kind() {
	case reflect.Struct:
		n = new(StructKindNode)
	case reflect.Array, reflect.Slice:
		n = new(ArrayKindNode)
	case reflect.Ptr, reflect.UnsafePointer, reflect.Uintptr:
		n = new(PtrKindNode)
	case reflect.Bool, reflect.String, reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n = new(BasicTypeKindNode)
	default:
		n = new(OtherKindNode) // Maps, interfaces, etc
	}

	n.initialize(field, parent, t, r)

	if len(t.PkgPath()) != 0 {
		typeName := t.PkgPath() + "." + t.Name()
		if _, ok := r.forest.ConnectedNodes[typeName]; !ok {
			r.forest.ConnectedNodes[typeName] = list.New()
		}
		r.forest.ConnectedNodes[typeName].PushBack(n)
	}

	return n
}

// BuildResourceTree builds a resource tree by calling into analyzeType method starting from root.
func (r *ResourceTree) BuildResourceTree(t reflect.Type) {
	r.Root = r.createNode(r.ResourceName, nil, t)
	r.Root.buildChildNodes(t)
}