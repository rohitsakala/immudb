/*
Copyright 2021 CodeNotary, Inc. All rights reserved.

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

package sessions

import (
	"github.com/codenotary/immudb/pkg/errors"
	"github.com/codenotary/immudb/pkg/logger"
	"os"
	"sync"
	"time"
)

var ErrGuardAlreadyRunning = errors.New("session guard already launched")
var ErrGuardNotRunning = errors.New("session guard not running")

var guard *manager

type manager struct {
	Running    bool
	sessionMux sync.Mutex
	sessions   map[string]*Session
	ticker     *time.Ticker
	done       chan bool
	logger     logger.Logger
	options    *Options
}

type Manager interface {
	SessionPresent(sessionID string) bool
	AddSession(sessionID string, sess *Session)
	RemoveSession(sessionID string)
	UpdateSessionActivityTime(sessionID string)
	UpdateHeartBeatTime(sessionID string)
	StartSessionsGuard() error
	StopSessionsGuard() error
	GetSession(sessionID string) *Session
	CountSession() int
}

func NewManager(options *Options) *manager {
	if options == nil {
		options = DefaultOptions()
	}
	guard = &manager{
		sessionMux: sync.Mutex{},
		sessions:   make(map[string]*Session),
		ticker:     time.NewTicker(options.SessionGuardCheckInterval),
		done:       make(chan bool),
		logger:     logger.NewSimpleLogger("immudb session guard", os.Stdout),
		options:    options,
	}
	return guard
}

func (sm *manager) SessionPresent(sessionID string) bool {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	if _, ok := sm.sessions[sessionID]; ok {
		return true
	}
	return false
}

func (sm *manager) AddSession(sessionID string, sess *Session) {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	sm.sessions[sessionID] = sess
	sm.logger.Debugf("created session %s", sessionID)
}

func (sm *manager) GetSession(sessionID string) *Session {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	return sm.sessions[sessionID]
}

func (sm *manager) RemoveSession(sessionID string) {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	delete(sm.sessions, sessionID)
}

func (sm *manager) UpdateSessionActivityTime(sessionID string) {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	if sess, ok := sm.sessions[sessionID]; ok {
		sess.lastActivityTime = time.Now()
		sm.logger.Debugf("updated last activity time for %s", sessionID)
	}
}

func (sm *manager) UpdateHeartBeatTime(sessionID string) {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	if sess, ok := sm.sessions[sessionID]; ok {
		sess.lastHeartBeat = time.Now()
		sm.logger.Debugf("updated last heart beat time for %s", sessionID)
	}
}

func (sm *manager) CountSession() int {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	return len(sm.sessions)
}

func (sm *manager) StartSessionsGuard() error {
	sm.sessionMux.Lock()
	if sm.Running == true {
		return ErrGuardAlreadyRunning
	}
	sm.Running = true
	sm.sessionMux.Unlock()
	for {
		select {
		case <-sm.done:
			return nil
		case <-sm.ticker.C:
			sm.expireSessions()
		}
	}
}

func (sm *manager) StopSessionsGuard() error {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	if sm.Running == false {
		return ErrGuardNotRunning
	}
	sm.Running = false
	sm.ticker.Stop()
	sm.done <- true
	sm.logger.Debugf("shutdown")
	return nil
}

func (sm *manager) expireSessions() {
	sm.sessionMux.Lock()
	defer sm.sessionMux.Unlock()
	if sm.Running {
		now := time.Now()
		sm.logger.Debugf("checking at %s", now.Format(time.UnixDate))
		for ID, sess := range sm.sessions {
			if sess.lastHeartBeat.Add(sm.options.MaxSessionIdle).Before(now) && sess.GetStatus() != IDLE {
				sess.SetStatus(IDLE)
				sm.logger.Debugf("session %s became IDLE, no more heartbeat received", ID)
			}
			if sess.lastActivityTime.Add(sm.options.MaxSessionIdle).Before(now) && sess.GetStatus() != IDLE {
				sess.SetStatus(IDLE)
				sm.logger.Debugf("session %s became IDLE due to max inactivity time", ID)
			}
			if sess.creationTime.Add(sm.options.MaxSessionAge).Before(now) {
				sess.SetStatus(DEAD)
				sm.logger.Debugf("session %s exceeded MaxSessionAge and became DEAD", ID)
			}
			if sess.state == IDLE {
				if sess.lastActivityTime.Add(sm.options.Timeout).Before(now) {
					sess.SetStatus(DEAD)
					sm.logger.Debugf("IDLE session %s is DEAD", ID)
				}
				if sess.lastHeartBeat.Add(sm.options.Timeout).Before(now) {
					sess.SetStatus(DEAD)
					sm.logger.Debugf("IDLE session %s is DEAD", ID)
				}
			}
			if sess.state == DEAD {
				sm.RemoveSession(ID)
				sm.logger.Debugf("removed DEAD session %s", ID)
			}
		}
	}
}
