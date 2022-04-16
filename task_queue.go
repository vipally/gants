// Copyright 2022 The @gxlb Authors. All rights reserved.
// Use of this source code is governed by a MIT License.
// License that can be found in the LICENSE file.

package gants

// taskQueue is lockfree queue of task
type taskQueue struct {
}

func (q taskQueue) Push(t *task) error {
	return nil
}

func (q taskQueue) Pop() (*task, bool) {
	return nil, false
}

func (q taskQueue) Len() int {
	return 0
}
