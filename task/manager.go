package task

import "sync"

type ChromeManager struct {
	Users map[string]*sync.Mutex
	Lock  sync.RWMutex
}

var ChromeManagerInstance *ChromeManager

func BootstrapChromeManagerInstance() {
	ChromeManagerInstance = NewChromeManager()
}

func NewChromeManager() *ChromeManager {
	return &ChromeManager{
		Users: make(map[string]*sync.Mutex),
		Lock:  sync.RWMutex{},
	}
}

func (m *ChromeManager) GetUserLock(username string) (lock *sync.Mutex) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.Users[username]; !ok {
		m.Users[username] = &sync.Mutex{}
	}
	lock = m.Users[username]
	return
}
