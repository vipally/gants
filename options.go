// Copyright 2022 The @gxlb Authors. All rights reserved.
// Use of this source code is governed by a MIT License.
// License that can be found in the LICENSE file.

package gants

import (
	"runtime"
)

// Option is the optional function.
type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// Options contains all options which will be applied when instantiating a gants pool.
type Options struct {
	EnableStatistic bool

	MaxWorkerCount int
}

func WithStatistic() Option {
	return func(opts *Options) {
		opts.EnableStatistic = true
	}
}

func WithMaxWorkerCount(maxWorkerCount int) Option {
	return func(opts *Options) {
		opts.MaxWorkerCount = maxWorkerCount
	}
}

func limitWorkerCount(maxWorkerCount int) int {
	const maxChannelLen = 4096
	if maxWorkerCount > 0 {
		if maxWorkerCount >= maxChannelLen {
			return maxChannelLen
		}
		return maxWorkerCount
	}
	return runtime.NumCPU() * 2
}
