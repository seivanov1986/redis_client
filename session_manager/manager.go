package session_manager

import (
	"sync"
	"time"
)

type sessionManager struct {
	keys  map[string]interface{}
	mutex sync.Mutex
}

func New() *sessionManager {
	return &sessionManager{
		keys: map[string]interface{}{},
	}
}

func (s *sessionManager) Set(key string, value interface{}, expiration time.Duration) error {
	s.mutex.Lock()
	s.keys[key] = value
	s.mutex.Unlock()

	return nil
}

func (s *sessionManager) Get(key string) (string, error) {
	s.mutex.Lock()
	value, ok := s.keys[key]
	s.mutex.Unlock()
	if !ok {
		return "", nil
	}

	strValue, ok := value.(string)
	if !ok {
		return "", nil
	}

	return strValue, nil
}

func (s *sessionManager) Exists(keys ...string) (bool, error) {
	for _, key := range keys {
		s.mutex.Lock()
		_, ok := s.keys[key]
		s.mutex.Unlock()

		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func (s *sessionManager) Del(keys ...string) (bool, error) {
	for _, key := range keys {
		s.mutex.Lock()
		delete(s.keys, key)
		s.mutex.Unlock()
	}

	return false, nil
}
