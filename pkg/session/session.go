package session

import (
	"agent-in-go/pkg/agentcore"
	"sync"
)

type Session struct {
	ID    string
	agent *agentcore.Agent
}

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]*Session
	factory  func() *agentcore.Agent
}

func NewSessionStore(factory func() *agentcore.Agent) *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
		factory:  factory,
	}
}

func (s *SessionStore) Get(id string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[id]; ok {
		return sess
	}
	sess := &Session{ID: id, agent: s.factory()}
	s.sessions[id] = sess
	return sess
}

func (s *SessionStore) Ask(sessionID, question string) string {
	sess := s.Get(sessionID)
	return sess.agent.Run(question) // maxSteps lives on the agent
}

func (s *SessionStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
}
