/*
Copyright 2021 The Kubernetes Authors.

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

// This fake package is created as golang.org/x/tools/go/analysis/analysistest
// expects it to be here for loading. This package is used to test parameter
// checking.

package parameters

import (
	"context"
	"fmt"

	klog "k8s.io/klog/v2"
)

func parameterCalls() {
	// format strings (incomplete list...)
	klog.Infof("%d", 1)
	klog.InfoS("%d", 1)  // want `logging function "InfoS" should not use format specifier "%d"`
	klog.Info("%d", 1)   // want `logging function "Info" should not use format specifier "%d"`
	klog.Infoln("%d", 1) // want `logging function "Infoln" should not use format specifier "%d"`
	klog.V(1).Infof("%d", 1)
	klog.V(1).InfoS("%d", 1)  // want `logging function "InfoS" should not use format specifier "%d"`
	klog.V(1).Info("%d", 1)   // want `logging function "Info" should not use format specifier "%d"`
	klog.V(1).Infoln("%d", 1) // want `logging function "Infoln" should not use format specifier "%d"`
	klog.Errorf("%d", 1)
	klog.ErrorS(nil, "%d", 1) // want `logging function "ErrorS" should not use format specifier "%d"`
	klog.Error("%d", 1)       // want `logging function "Error" should not use format specifier "%d"`
	klog.Errorln("%d", 1)     // want `logging function "Errorln" should not use format specifier "%d"`

	klog.InfoS("hello", "value", fmt.Sprintf("%d", 1))

	// odd number of parameters
	klog.InfoS("hello", "key")       // want `Additional arguments to InfoS should always be Key Value pairs. Please check if there is any key or value missing.`
	klog.ErrorS(nil, "hello", "key") // want `Additional arguments to ErrorS should always be Key Value pairs. Please check if there is any key or value missing.`

	// non-string keys
	klog.InfoS("hello", "1", 2)
	klog.InfoS("hello", 1, 2) // want `Key positional arguments are expected to be inlined constant strings. Please replace 1 provided with string value`
	klog.ErrorS(nil, "hello", "1", 2)
	klog.ErrorS(nil, "hello", 1, 2) // want `Key positional arguments are expected to be inlined constant strings. Please replace 1 provided with string value`

	// variadic input to klog.Info*, klog.Error*, klog.LoggerWithValues functions
	kvs := []interface{}{"key1", "value1"}
	klog.InfoS("foo message", kvs) // want `Additional arguments to InfoS should always be Key Value pairs. Please check if there is any key or value missing.`
	klog.ErrorS(nil, "foo error message", kvs...)
	klog.Info(kvs...)
	klog.Error(kvs...)

	logger := klog.FromContext(context.Background())
	klog.LoggerWithValues(logger, kvs...)
	klog.LoggerWithValues(logger, kvs) // want `Additional arguments to LoggerWithValues should always be Key Value pairs. Please check if there is any key or value missing.`
}
