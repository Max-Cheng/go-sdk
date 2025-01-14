package thinkingdata

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"time"
)

const (
	DATE_FORMAT = "2006-01-02 15:04:05.000"
	KEY_PATTERN = "^[a-zA-Z#][A-Za-z0-9_]{0,49}$"
)

var keyPattern, _ = regexp.Compile(KEY_PATTERN)

// mergeProperties 对source的Value做了一份Copy
// Copy 至 target
// 本质上对不定Target的内存增长是不可控的(我们应该声明时注明cap?)

func mergeProperties(target, source map[string]interface{}) {
	// Better way to implement this function?
	for k, v := range source {
		target[k] = v
	}
}

// extractTime formate time string to DATE_FORMAT
func extractTime(p map[string]interface{}) string {
	if t, ok := p["#time"]; ok {
		delete(p, "#time")
		switch v := t.(type) {
		case string:
			return v
		case time.Time:
			return v.Format(DATE_FORMAT)
		default:
			return time.Now().Format(DATE_FORMAT)
		}
	}

	return time.Now().Format(DATE_FORMAT)
}

func extractStringProperty(p map[string]interface{}, key string) string {
	if t, ok := p[key]; ok {
		delete(p, key)
		v, ok := t.(string)
		if !ok {
			fmt.Fprintln(os.Stderr, "Invalid data type for "+key)
		}
		return v
	}
	return ""
}

func isNotNumber(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64:
	default:
		return true
	}
	return false
}

func formatProperties(d *Data) error {

	if d.EventName != "" {
		matched := checkPattern([]byte(d.EventName))
		if !matched {
			return errors.New("Invalid event name: " + d.EventName)
		}
	}

	if d.Properties != nil {
		for k, v := range d.Properties {
			isMatch := checkPattern([]byte(k))
			if !isMatch {
				return errors.New("Invalid property key: " + k)
			}

			if d.Type == UserAdd && isNotNumber(v) {
				return errors.New("Invalid property value: only numbers is supported by UserAdd")
			}

			//check value
			switch v.(type) {
			case bool:
			case string:
			case []string:
			case []interface{}:
			case time.Time:
				d.Properties[k] = v.(time.Time).Format(DATE_FORMAT)
			case *time.Time:
				d.Properties[k] = v.(*time.Time).Format(DATE_FORMAT)
			default:
				if isNotNumber(v) && isNotArrayOrSlice(v) {
					errorMsg := fmt.Sprintf("%v Invalid property value type. Supported types: numbers, string, time.Time, bool, array, slice", v)
					return errors.New(errorMsg)
				}
			}
		}
	}

	return nil
}

func isNotArrayOrSlice(v interface{}) bool {
	typeOf := reflect.TypeOf(v)
	switch typeOf.Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return true
	}
	return false
}

func checkPattern(name []byte) bool {
	return keyPattern.Match(name)
}
