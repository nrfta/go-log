package log

import (
	"bytes"
	"encoding/json"
	"strings"
)

type ByteLogs struct {
	Log    *bytes.Buffer
	Parsed []map[string]interface{}
}

func (b *ByteLogs) Parse(log *bytes.Buffer) {
	// the buffer will contain the latest log since it is a pointer, but this function will need to be called
	// to have the latest values in ParsedLogs for use with LogInLogs
	if log != nil {
		b.Log = log
	}
	if b.Log == nil || b.Log.Len() == 0 {
		return
	}

	logs := strings.Split(b.Log.String(), "\n")
	objLogs := make([]map[string]interface{}, 0)
	for _, v := range logs {
		var obj map[string]interface{}
		err := json.Unmarshal([]byte(v), &obj)
		if err == nil {
			objLogs = append(objLogs, obj)
		}
	}

	b.Parsed = objLogs
}

func (b *ByteLogs) LogInLogs(key string, value interface{}) bool {
	for _, log := range b.Parsed {
		entry, ok := log[key]
		if ok && entry == value {
			return true
		}
	}

	return false
}

func (b *ByteLogs) CheckForLogEndMessage(customEnd ...string) bool {
	haveStop := false
	if len(customEnd) > 0 {
		for _, msg := range customEnd {
			haveStop = b.LogInLogs("msg", msg)
			if haveStop {
				break
			}
		}
	} else {
		goodStop := b.LogInLogs("msg", "task-consumer: successfully processed message")
		badStop := b.LogInLogs("msg", "task-consumer: failed to process message")
		haveStop = goodStop || badStop
	}

	return haveStop
}
