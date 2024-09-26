/*
Copyright 2024-present Volodymyr Konstanchuk and contributors

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

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/componego/componego/internal/utils"
)

func TestIsParentContext(t *testing.T) {
	parentCtx := context.Background()
	childCtx1, cancel1 := context.WithCancel(parentCtx)
	defer cancel1()
	childCtx2, cancel2 := context.WithTimeout(childCtx1, time.Minute*60)
	defer cancel2()
	childCtx3, cancel3 := context.WithDeadline(childCtx1, time.Now().Add(time.Minute*60))
	defer cancel3()
	testCases := [...]struct {
		a      context.Context
		b      context.Context
		status bool
	}{
		{parentCtx, childCtx1, true},  // 1
		{childCtx1, parentCtx, false}, // 2
		//
		{parentCtx, childCtx2, true},  // 3
		{childCtx2, parentCtx, false}, // 4
		//
		{childCtx1, childCtx2, true},  // 5
		{childCtx2, childCtx1, false}, // 6
		//
		{childCtx1, childCtx3, true},  // 7
		{childCtx3, childCtx1, false}, // 8
		//
		{parentCtx, context.Background(), true},            // 9
		{context.Background(), context.Background(), true}, // 10
		{context.Background(), parentCtx, true},            // 11
		{context.Background(), childCtx1, true},            // 12
		{context.Background(), context.TODO(), false},      // 13
		{context.TODO(), context.Background(), false},      // 14
	}
	for i, testCase := range testCases {
		if utils.IsParentContext(testCase.a, testCase.b) != testCase.status {
			message := "Failed to check parent context (#%d). Make sure the framework is compatible with your version of Go."
			t.Fatalf(message, i+1)
		}
	}
}
