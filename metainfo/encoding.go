package metainfo

import (
	"github.com/anacrolix/torrent/bencode"
)

func applyUTF8Field(v interface{}, data map[string]interface{}) error {
	resetUTF8Val(data)
	b, err := bencode.Marshal(data)
	if err != nil {
		return err
	}
	if err := bencode.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

func resetUTF8Val(data map[string]interface{}) {
	for k := range data {
		if v, ok := data[k+".utf-8"]; ok {
			data[k] = v
			continue
		}
		switch v := data[k].(type) {
		case []interface{}:
			for _, vv := range v {
				if m, ok := vv.(map[string]interface{}); ok {
					resetUTF8Val(m)
				}
			}
		case map[string]interface{}:
			resetUTF8Val(v)
		}
	}
	return
}
