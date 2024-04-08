package api

import (
	"strconv"
	"strings"
	"time"
)

type Metadata map[string]any

func NewMetadata(m map[string]any) Metadata {
	if len(m) == 0 {
		return nil
	}
	md := make(map[string]any)
	for k, v := range m {
		md[strings.ToLower(k)] = v
	}
	return Metadata(md)
}

func (md Metadata) IsExists(key string) bool {
	_, ok := md[strings.ToLower(key)]
	return ok
}

func (md Metadata) Get(key string) any {
	if md != nil {
		return md[strings.ToLower(key)]
	}
	return nil
}

func (md Metadata) GetString(keys ...string) (v string) {
	if md == nil {
		return
	}

	for _, key := range keys {
		if !md.IsExists(key) {
			continue
		}

		switch vv := md.Get(key).(type) {
		case string:
			v = vv
		case int:
			v = strconv.FormatInt(int64(vv), 10)
		case int64:
			v = strconv.FormatInt(vv, 10)
		case uint:
			v = strconv.FormatUint(uint64(vv), 10)
		case uint64:
			v = strconv.FormatUint(uint64(vv), 10)
		case bool:
			v = strconv.FormatBool(vv)
		case float32:
			v = strconv.FormatFloat(float64(vv), 'f', -1, 32)
		case float64:
			v = strconv.FormatFloat(float64(vv), 'f', -1, 64)
		}
		break
	}

	return
}

func (md Metadata) GetBool(keys ...string) (v bool) {
	if md == nil {
		return false
	}

	for _, key := range keys {
		if !md.IsExists(key) {
			continue
		}
		switch vv := md.Get(key).(type) {
		case bool:
			v = vv
		case int:
			v = vv != 0
		case string:
			v, _ = strconv.ParseBool(vv)
		}
		break
	}

	return

}

func (md Metadata) GetInt(keys ...string) (v int) {
	if md == nil {
		return
	}

	for _, key := range keys {
		if !md.IsExists(key) {
			continue
		}
		switch vv := md.Get(key).(type) {
		case bool:
			if vv {
				v = 1
			}
		case int:
			v = vv
		case string:
			v, _ = strconv.Atoi(vv)
		}
		break
	}

	return
}

func (md Metadata) GetFloat(keys ...string) (v float64) {
	if md == nil {
		return
	}

	for _, key := range keys {
		if !md.IsExists(key) {
			continue
		}

		switch vv := md.Get(key).(type) {
		case float64:
			v = vv
		case int:
			v = float64(vv)
		case string:
			v, _ = strconv.ParseFloat(vv, 64)
		}
		break
	}
	return
}

func (md Metadata) GetDuration(keys ...string) (v time.Duration) {
	if md == nil {
		return
	}

	for _, key := range keys {
		if !md.IsExists(key) {
			continue
		}

		switch vv := md.Get(key).(type) {
		case int:
			v = time.Duration(vv) * time.Second
		case string:
			v, _ = time.ParseDuration(vv)
			if v == 0 {
				n, _ := strconv.Atoi(vv)
				v = time.Duration(n) * time.Second
			}
		}
		break
	}
	return
}
