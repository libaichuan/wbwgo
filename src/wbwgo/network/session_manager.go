package network

import (
	"sync"
)

type SessionManager struct {
	sessions map[int]*Session

	ses_guard sync.RWMutex
}

func (self *SessionManager) AddSession(s *Session) {

}
