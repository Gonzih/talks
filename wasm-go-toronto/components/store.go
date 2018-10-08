package main

type Store struct {
	subs  []chan bool
	store map[string]interface{}
}

func (s *Store) Set(k string, v interface{}) {
	s.store[k] = v
	for _, sub := range s.subs {
		select {
		case sub <- true:
		default:
		}
	}
}

func (s *Store) Get(k string) interface{} {
	return s.store[k]
}

func (s *Store) Subscribe(ch chan bool) {
	s.subs = append(s.subs, ch)
}

func NewStore() *Store {
	s := &Store{}
	s.store = make(map[string]interface{}, 0)
	return s
}
