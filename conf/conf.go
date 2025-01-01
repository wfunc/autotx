package conf

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xmap"
)

type JSONFile struct {
	Users       map[string]xmap.M
	Seeds       map[string]string
	SeedsRevert map[string]string
	NotDo       map[string]string
	Do          map[string]string
	Lock        sync.RWMutex
}

var Conf *JSONFile

func Bootstrap() {
	Conf = NewJSONFile()
	Conf.Load()
}

func NewJSONFile() *JSONFile {
	return &JSONFile{
		Users:       map[string]xmap.M{},
		Seeds:       map[string]string{},
		SeedsRevert: map[string]string{},
		NotDo:       map[string]string{},
		Do:          map[string]string{},
		Lock:        sync.RWMutex{},
	}
}

func (j *JSONFile) Load() (err error) {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	err = ReadJSON("conf/users.json", &j.Users)
	if err != nil {
		xlog.Infof("ReadJSON(conf/users.json) Failed with err %v", err)
	}
	err = ReadJSON("conf/seeds.json", &j.Seeds)
	if err != nil {
		xlog.Infof("ReadJSON(conf/seeds.json) Failed with err %v", err)
	}
	err = ReadJSON("conf/not_do.json", &j.NotDo)
	if err != nil {
		xlog.Infof("ReadJSON(conf/not_do.json) Failed with err %v", err)
	}
	err = ReadJSON("conf/do.json", &j.Do)
	if err != nil {
		xlog.Infof("ReadJSON(conf/do.json) Failed with err %v", err)
	}
	for k, v := range j.Seeds {
		j.SeedsRevert[v] = k
	}
	return
}

func (j *JSONFile) Save(keys ...string) (err error) {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	if len(keys) > 0 {
		for _, key := range keys {
			filename := ""
			var v interface{}
			switch key {
			case "user":
				filename = "conf/users.json"
				v = j.Users
			case "seed":
				filename = "conf/seeds.json"
				v = j.Seeds
			case "not_do":
				filename = "conf/not_do.json"
				v = j.NotDo
			case "do":
				filename = "conf/do.json"
				v = j.Do
			}
			if len(filename) > 0 {
				err = WriteJSON(filename, v)
				if err != nil {
					return
				}
			}
		}
		return
	}
	err = WriteJSON("conf/users.json", j.Users)
	if err != nil {
		return
	}
	err = WriteJSON("conf/seeds.json", j.Seeds)
	if err != nil {
		return
	}
	err = WriteJSON("conf/not_do.json", j.NotDo)
	if err != nil {
		return
	}
	err = WriteJSON("conf/do.json", j.Do)
	if err != nil {
		return
	}
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
		j.Save("user")
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
		j.Save("user")
	}()
	delete(j.Users, username)
}

func (j *JSONFile) UpdateUser(username, key, value string) {
	if j == nil || j.Users == nil {
		return
	}
	j.Lock.Lock()
	defer func() {
		j.Lock.Unlock()
		j.Save("user")
	}()
	j.Users[username][key] = value
}

func (j *JSONFile) IsAllSignIN() bool {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	isAllSignIN := true
	now := time.Now()
	layout := "2006-01-02 15:04:05"
	for _, user := range j.Users {
		signIN := user.Str("signIN")
		if len(signIN) > 0 {
			// 使用 time.Parse 将字符串解析为 time.Time
			parsedTime, err := time.ParseInLocation(layout, signIN, time.Local)
			if err == nil {
				if !(parsedTime.Year() == now.Year() && parsedTime.Month() == now.Month() && parsedTime.Day() == now.Day()) {
					isAllSignIN = false
					break
				}
			} else {
				isAllSignIN = false
				break
			}
		} else {
			isAllSignIN = false
			break
		}
	}
	return isAllSignIN
}

func (j *JSONFile) GetSeeds() map[string]string {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	return j.Seeds
}

func (j *JSONFile) SetSeeds(seeds map[string]string) {
	j.Lock.Lock()
	defer func() {
		j.Lock.Unlock()
		j.Save("seed")
	}()

	for k, v := range seeds {
		j.Seeds[k] = v
		j.SeedsRevert[v] = k
	}
}

func (j *JSONFile) GetSeedsRevert() map[string]string {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	return j.SeedsRevert
}

func (j *JSONFile) GetSeedName(seedID string) string {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	return j.SeedsRevert[seedID]
}

func (j *JSONFile) AddNotDo(key, value string) {
	j.Lock.Lock()
	defer func() {
		j.Lock.Unlock()
		j.Save("not_do")
	}()
	if _, ok := j.NotDo[key]; !ok {
		j.NotDo[key] = ""
	}
	if len(j.NotDo[key]) < 1 {
		j.NotDo[key] = value
		return
	}
	values := strings.Split(j.NotDo[key], ",")
	for _, v := range values {
		if v == value {
			return
		}
	}
	values = append(values, value)
	j.NotDo[key] = strings.Join(values, ",")
}

func (j *JSONFile) RemoveNotDo(key, value string) {
	j.Lock.Lock()
	defer func() {
		j.Lock.Unlock()
		j.Save("not_do")
	}()
	if _, ok := j.NotDo[key]; !ok {
		j.NotDo[key] = ""
	}
	values := strings.Split(j.NotDo[key], ",")
	for i, v := range values {
		if v == value {
			values = append(values[:i], values[i+1:]...)
			break
		}
	}
	j.NotDo[key] = strings.Join(values, ",")
}

func (j *JSONFile) ListNotDo(key string) []string {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	if _, ok := j.NotDo[key]; !ok {
		return []string{}
	}
	if len(j.NotDo[key]) < 1 {
		return []string{}
	}
	return strings.Split(j.NotDo[key], ",")
}

func (j *JSONFile) LoadNotDo(key string) string {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	if _, ok := j.NotDo[key]; !ok {
		return ""
	}
	return j.NotDo[key]
}

func (j *JSONFile) AddDo(key, value string) {
	j.Lock.Lock()
	defer func() {
		j.Lock.Unlock()
		j.Save("do")
	}()
	if _, ok := j.Do[key]; !ok {
		j.Do[key] = ""
	}
	if len(j.Do[key]) < 1 {
		j.Do[key] = value
		return
	}
	values := strings.Split(j.Do[key], ",")
	for _, v := range values {
		if v == value {
			return
		}
	}
	values = append(values, value)
	j.Do[key] = strings.Join(values, ",")
}

func (j *JSONFile) RemoveDo(key, value string) {
	j.Lock.Lock()
	defer func() {
		j.Lock.Unlock()
		j.Save("do")
	}()
	if _, ok := j.Do[key]; !ok {
		j.Do[key] = ""
	}
	values := strings.Split(j.Do[key], ",")
	for i, v := range values {
		if v == value {
			values = append(values[:i], values[i+1:]...)
			break
		}
	}
	j.Do[key] = strings.Join(values, ",")
}

func (j *JSONFile) ListDo(key string) []string {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	if _, ok := j.Do[key]; !ok {
		return []string{}
	}
	if len(j.Do[key]) < 1 {
		return []string{}
	}
	return strings.Split(j.Do[key], ",")
}

func (j *JSONFile) LoadDo(key string) string {
	j.Lock.RLock()
	defer j.Lock.RUnlock()
	if _, ok := j.Do[key]; !ok {
		return ""
	}
	return j.Do[key]
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
