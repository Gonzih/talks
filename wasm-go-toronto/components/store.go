package main

var notificationChan = make(chan bool)

type Store struct {
	store map[string]interface{}
}

func (s *Store) Set(k string, v interface{}) {
	s.store[k] = v
	select {
	case notificationChan <- true:
	default:
	}
}

func (s *Store) Get(k string) interface{} {
	return s.store[k]
}

func NewStore() *Store {
	s := &Store{}
	s.store = make(map[string]interface{}, 0)
	return s
}
