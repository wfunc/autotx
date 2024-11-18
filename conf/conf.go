package conf

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/wfunc/util/xmap"
)

type JSONFile struct {
	Users map[string]xmap.M
	Lock  sync.RWMutex
}

var Conf *JSONFile

func Bootstrap() {
	Conf = NewJSONFile()
	Conf.Load()
}

func NewJSONFile() *JSONFile {
	return &JSONFile{Users: map[string]xmap.M{}, Lock: sync.RWMutex{}}
}

func (j *JSONFile) Load() (err error) {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	err = ReadJSON("conf/users.json", &j.Users)
	return
}

func (j *JSONFile) Save() (err error) {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	err = WriteJSON("conf/users.json", j.Users)
	return
}

func (j *JSONFile) GetUser(username string) xmap.M {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	return j.Users[username]
}

func (j *JSONFile) GetUsers() map[string]xmap.M {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	return j.Users
}

func (j *JSONFile) AddUser(username, password string) {
	j.Lock.Lock()
	defer func() {
		j.Lock.Unlock()
		j.Save()
	}()
	if _, ok := j.Users[username]; !ok {
		j.Users[username] = xmap.M{}
	}
	j.Users[username]["password"] = password
}

func (j *JSONFile) RemoveUser(username string) {
	j.Lock.Lock()
	defer func() {
		j.Lock.Unlock()
		j.Save()
	}()
	delete(j.Users, username)
}

// ReadJSON will read file and unmarshal to value
func ReadJSON(filename string, v interface{}) (err error) {
	data, err := os.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, v)
	}
	return
}

// WriteJSON will marshal value to json and write to file
func WriteJSON(filename string, v interface{}) (err error) {
	data, err := json.MarshalIndent(v, "", "    ")
	if err == nil {
		err = os.WriteFile(filename, data, os.ModePerm)
	}
	return
}
