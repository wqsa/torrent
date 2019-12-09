package metainfo

import (
	"reflect"
	"testing"
	"unicode/utf8"

	"github.com/anacrolix/torrent/bencode"
)

func hasGarbled(v interface{}) bool {
	if v == nil {
		return false
	}
	val := reflect.ValueOf(v)
	switch val.Type().Kind() {
	case reflect.String:
		return !utf8.ValidString(val.String())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			if hasGarbled(f.Interface()) {
				return true
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if hasGarbled(elem.Interface()) {
				return true
			}
		}
	case reflect.Ptr:
		if !val.IsZero() {
			return hasGarbled(val.Elem().Interface())
		}
	}
	return false
}

func TestEncoding(t *testing.T) {
	metainfo, err := LoadFromFile("testdata/field_utf8.torrent")
	if err != nil {
		t.Error("load test torrent fail", err)
	}
	var info Info
	if err := bencode.Unmarshal(metainfo.InfoBytes, &info); err != nil {
		t.Error("unmarshal infohash failed", err)
	}
	if hasGarbled(metainfo) {
		t.Error("metainfo hasGarbled")
	}
	if hasGarbled(info) {
		t.Error("info hasGarbled")
	}
}
