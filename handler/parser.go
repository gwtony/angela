package handler

import (
	//"fmt"
	"strings"
	"encoding/json"
)

func DecodeCronMessage(m []byte) (*CronMessage, error){
	cm := &CronMessage{}
	//TODO: parse json
	err := json.Unmarshal(m, &cm)
	if err != nil {
		return nil, err
	}

	return cm, nil
}

func EncodeCronMessage(cm *CronMessage) ([]byte, error) {
	return json.Marshal(cm)
}

func DecodeWatchMessage(otype string, key, value []byte) (*WatchMessage, error) {
	wm := &WatchMessage{}

	//fmt.Println("Decode watch: key(%s), value{%s} ", string(key), string(value))
	if value != nil {
		err := json.Unmarshal(value, &wm)
		if err != nil {
			return nil, err
		}
	}
	wm.Type = otype
	wm.Jobid = strings.TrimPrefix(string(key), JOB_META_LOC)

	return wm, nil
}

func EncodeWatchMessage(wm *WatchMessage) ([]byte, error) {
	return json.Marshal(wm)
}

//func DecodeWatchStartMessage(otype string, key, value []byte) (*WatchStartMessage, error) {
//	wm := &WatchStartMessage{}
//
//	err := json.Unmarshal(value, &wm)
//	if err != nil {
//		return nil, err
//	}
//	wm.Type = otype
//	wm.Jobid = strings.TrimPrefix(string(key), JOB_START_LOC)
//
//	return wm, nil
//}
