package exetimer

import (
	"fmt"
	"time"
)

type TimeName string

type ExeTimeMap struct {
	mp map[TimeName]time.Time
}

func NewExeTimeMap() *ExeTimeMap {
	return &ExeTimeMap{
		mp: make(map[TimeName]time.Time),
	}
}

func (t *ExeTimeMap) IsExistTime(name TimeName) (bool, time.Time) {
	v, ok := t.mp[name]
	if ok {
		return true, v
	}

	return false, time.Time{}
}

func (t *ExeTimeMap) MemoTime(name TimeName) {
	if ok, _ := t.IsExistTime(name); ok {
		panic("ExeTimeMap: MemoTime() called twice")
	}

	t.mp[name] = time.Now()
}

func (t *ExeTimeMap) MakeTimer(name1, name2 TimeName) (*ExeTimer, error) {
	ok, time1 := t.IsExistTime(name1)
	if !ok {
		return nil, fmt.Errorf("ExeTimeMap: MakeTimer() %v is not exist", name1)
	}
	ok, time2 := t.IsExistTime(name1)
	if !ok {
		return nil, fmt.Errorf("ExeTimeMap: MakeTimer() %v is not exist", name2)
	}

	if time1.After(time2) {
		time1, time2 = time2, time1
	}

	return newExeTime(time1, time2), nil
}
