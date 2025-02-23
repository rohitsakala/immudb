/*
Copyright 2022 CodeNotary, Inc. All rights reserved.

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
package tbtree

import (
	"os"
	"time"

	"github.com/codenotary/immudb/embedded/appendable"
	"github.com/codenotary/immudb/embedded/appendable/multiapp"
	"github.com/codenotary/immudb/pkg/logger"
)

const DefaultMaxNodeSize = 4096
const DefaultFlushThld = 100_000
const DefaultSyncThld = 1_000_000
const DefaultFlushBufferSize = 4096
const DefaultCleanUpPercentage float32 = 0
const DefaultMaxActiveSnapshots = 100
const DefaultRenewSnapRootAfter = time.Duration(1000) * time.Millisecond
const DefaultCacheSize = 100_000
const DefaultFileMode = os.FileMode(0755)
const DefaultFileSize = 1 << 26 // 64Mb
const DefaultMaxKeyLen = 1024
const DefaultCompactionThld = 2
const DefaultDelayDuringCompaction = time.Duration(10) * time.Millisecond

const DefaultNodesLogMaxOpenedFiles = 10
const DefaultHistoryLogMaxOpenedFiles = 1
const DefaultCommitLogMaxOpenedFiles = 1

const MinNodeSize = 128
const MinCacheSize = 1

type AppFactoryFunc func(
	rootPath string,
	subPath string,
	opts *multiapp.Options,
) (appendable.Appendable, error)

type Options struct {
	log logger.Logger

	flushThld          int
	syncThld           int
	flushBufferSize    int
	cleanupPercentage  float32
	maxActiveSnapshots int
	renewSnapRootAfter time.Duration
	cacheSize          int
	readOnly           bool
	fileMode           os.FileMode

	nodesLogMaxOpenedFiles   int
	historyLogMaxOpenedFiles int
	commitLogMaxOpenedFiles  int

	maxKeyLen int

	compactionThld        int
	delayDuringCompaction time.Duration

	// options below are only set during initialization and stored as metadata
	maxNodeSize int
	fileSize    int

	appFactory AppFactoryFunc
}

func DefaultOptions() *Options {
	return &Options{
		log:                   logger.NewSimpleLogger("immudb ", os.Stderr),
		flushThld:             DefaultFlushThld,
		syncThld:              DefaultSyncThld,
		flushBufferSize:       DefaultFlushBufferSize,
		cleanupPercentage:     DefaultCleanUpPercentage,
		maxActiveSnapshots:    DefaultMaxActiveSnapshots,
		renewSnapRootAfter:    DefaultRenewSnapRootAfter,
		cacheSize:             DefaultCacheSize,
		readOnly:              false,
		fileMode:              DefaultFileMode,
		maxKeyLen:             DefaultMaxKeyLen,
		compactionThld:        DefaultCompactionThld,
		delayDuringCompaction: DefaultDelayDuringCompaction,

		nodesLogMaxOpenedFiles:   DefaultNodesLogMaxOpenedFiles,
		historyLogMaxOpenedFiles: DefaultHistoryLogMaxOpenedFiles,
		commitLogMaxOpenedFiles:  DefaultCommitLogMaxOpenedFiles,

		// options below are only set during initialization and stored as metadata
		maxNodeSize: DefaultMaxNodeSize,
		fileSize:    DefaultFileSize,
	}
}

func validOptions(opts *Options) bool {
	return opts != nil &&
		opts.maxNodeSize >= MinNodeSize &&
		opts.flushThld > 0 &&
		opts.flushThld <= opts.syncThld &&
		opts.flushBufferSize > 0 &&
		opts.cleanupPercentage >= 0 && opts.cleanupPercentage <= 100 &&
		opts.nodesLogMaxOpenedFiles > 0 &&
		opts.historyLogMaxOpenedFiles > 0 &&
		opts.commitLogMaxOpenedFiles > 0 &&

		opts.maxActiveSnapshots > 0 &&
		opts.renewSnapRootAfter >= 0 &&
		opts.cacheSize >= MinCacheSize &&
		opts.maxKeyLen > 0 &&
		opts.compactionThld > 0 &&
		opts.log != nil
}

func (opts *Options) WithLog(log logger.Logger) *Options {
	opts.log = log
	return opts
}

func (opts *Options) WithAppFactory(appFactory AppFactoryFunc) *Options {
	opts.appFactory = appFactory
	return opts
}

func (opts *Options) WithFlushThld(flushThld int) *Options {
	opts.flushThld = flushThld
	return opts
}

func (opts *Options) WithSyncThld(syncThld int) *Options {
	opts.syncThld = syncThld
	return opts
}

func (opts *Options) WithFlushBufferSize(size int) *Options {
	opts.flushBufferSize = size
	return opts
}

func (opts *Options) WithCleanupPercentage(cleanupPercentage float32) *Options {
	opts.cleanupPercentage = cleanupPercentage
	return opts
}

func (opts *Options) WithMaxActiveSnapshots(maxActiveSnapshots int) *Options {
	opts.maxActiveSnapshots = maxActiveSnapshots
	return opts
}

func (opts *Options) WithRenewSnapRootAfter(renewSnapRootAfter time.Duration) *Options {
	opts.renewSnapRootAfter = renewSnapRootAfter
	return opts
}

func (opts *Options) WithCacheSize(cacheSize int) *Options {
	opts.cacheSize = cacheSize
	return opts
}

func (opts *Options) WithReadOnly(readOnly bool) *Options {
	opts.readOnly = readOnly
	return opts
}

func (opts *Options) WithFileMode(fileMode os.FileMode) *Options {
	opts.fileMode = fileMode
	return opts
}

func (opts *Options) WithNodesLogMaxOpenedFiles(nodesLogMaxOpenedFiles int) *Options {
	opts.nodesLogMaxOpenedFiles = nodesLogMaxOpenedFiles
	return opts
}

func (opts *Options) WithHistoryLogMaxOpenedFiles(historyLogMaxOpenedFiles int) *Options {
	opts.historyLogMaxOpenedFiles = historyLogMaxOpenedFiles
	return opts
}

func (opts *Options) WithCommitLogMaxOpenedFiles(commitLogMaxOpenedFiles int) *Options {
	opts.commitLogMaxOpenedFiles = commitLogMaxOpenedFiles
	return opts
}

func (opts *Options) WithMaxKeyLen(maxKeyLen int) *Options {
	opts.maxKeyLen = maxKeyLen
	return opts
}

func (opts *Options) WithMaxNodeSize(maxNodeSize int) *Options {
	opts.maxNodeSize = maxNodeSize
	return opts
}

func (opts *Options) WithFileSize(fileSize int) *Options {
	opts.fileSize = fileSize
	return opts
}

func (opts *Options) WithCompactionThld(compactionThld int) *Options {
	opts.compactionThld = compactionThld
	return opts
}

func (opts *Options) WithDelayDuringCompaction(delay time.Duration) *Options {
	opts.delayDuringCompaction = delay
	return opts
}
