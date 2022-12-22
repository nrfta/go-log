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
	// logs are parsed here for convenience of using logs in a gomege.Eventually block to catch new logs
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
