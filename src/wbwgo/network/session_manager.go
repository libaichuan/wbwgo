package network

import (
	"log"
	"sync"
	"sync/atomic"
)

type SessionManager struct {
	sessions map[int32]*Session

	cur_session_id int32

	ses_guard sync.RWMutex
}

const (
	atomic_count = 10
)

func (self *SessionManager) AddSession(s *Session) {
	self.ses_guard.Lock()
	defer self.ses_guard.Unlock()

	var cur_atomic_count int = atomic_count

	var id int32

	for cur_atomic_count > 0 {

		id = atomic.AddInt32(&self.cur_session_id, 1)

		if _, ok := self.sessions[id]; !ok {
			break
		}

		cur_atomic_count--
	}

	if cur_atomic_count == 0 {
		log.Println("SessionManager::AddSession same key")
	}

	s.id = id

	self.sessions[id] = s

	log.Printf("add new session id:%d\n", id)
}

func (self *SessionManager) GetSessionById(id int32) *Session {
	self.ses_guard.Lock()
	defer self.ses_guard.Unlock()

	if v, ok := self.sessions[id]; ok {
		return v
	}

	return nil
}

func (self *SessionManager) RemoveSessionById(id int32) {
	self.ses_guard.Lock()
	defer self.ses_guard.Unlock()

	delete(self.sessions, id)
}

func (self *SessionManager) GetSessionCount() int {
	self.ses_guard.Lock()
	defer self.ses_guard.Unlock()

	return len(self.sessions)
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[int32]*Session),
	}
}
