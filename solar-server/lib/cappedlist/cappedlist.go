package cappedlist

import (
	"sync"
)

type Entry interface{}

type CappedList struct {
	entries         []Entry
	removedCallback func(Entry)
	cap             int
	sync.RWMutex
}

func New(cap int) *CappedList {
	return &CappedList{make([]Entry, cap), nil, cap, sync.RWMutex{}}
}

func (l *CappedList) Add(e Entry) {
	l.Lock()
	defer l.Unlock()
	l.entries = append(l.entries, e)
	if len(l.entries) > l.cap {
		r := l.entries[0]
		l.entries = l.entries[1:]
		if l.removedCallback != nil {
			l.removedCallback(r)
		}
	}
}

func (l *CappedList) At(index int) Entry {
	l.RLock()
	defer l.RUnlock()
	length := len(l.entries)
	if index >= 0 && index < length {
		return l.entries[index]
	}
	return nil
}

func (l *CappedList) Last() Entry {
	l.RLock()
	defer l.RUnlock()
	n := len(l.entries)
	if n > 0 {
		return l.entries[n-1]
	}
	return nil
}

func (l *CappedList) All() []Entry {
	l.RLock()
	defer l.RUnlock()
	cpy := make([]Entry, len(l.entries))
	copy(cpy, l.entries)
	return cpy
}

func (l *CappedList) RegisterRemovedEntryCallback(callback func(Entry)) {
	l.removedCallback = callback
}
