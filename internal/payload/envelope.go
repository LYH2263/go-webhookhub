package payload

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"time"
	"unicode/utf8"
)

// Envelope 投递给端点的 JSON 外壳。Data 在 body 已是 JSON 时原样嵌入。
type Envelope struct {
	ID        string          `json:"id"`
	Event     string          `json:"event"`
	Timestamp string          `json:"timestamp"`
	Attempt   int             `json:"attempt"`
	Data      json.RawMessage `json:"data"`
}

// Wrap 把原始 body 封进信封。body 会先拷贝再使用。
func Wrap(id, event string, ts time.Time, attempt int, body []byte) ([]byte, error) {
	body = CloneBytes(body)
	data, err := asJSONData(body)
	if err != nil {
		return nil, err
	}
	env := Envelope{
		ID:        id,
		Event:     event,
		Timestamp: ts.UTC().Format(time.RFC3339Nano),
		Attempt:   attempt,
		Data:      data,
	}
	return json.Marshal(env)
}

func asJSONData(body []byte) (json.RawMessage, error) {
	if len(body) == 0 {
		return json.RawMessage("null"), nil
	}
	trimmed := bytes.TrimSpace(body)
	if json.Valid(trimmed) {
		return json.RawMessage(CloneBytes(trimmed)), nil
	}
	if utf8.Valid(trimmed) {
		quoted, err := json.Marshal(string(trimmed))
		if err != nil {
			return nil, err
		}
		return json.RawMessage(quoted), nil
	}
	enc := base64.StdEncoding.EncodeToString(trimmed)
	quoted, err := json.Marshal(map[string]string{"encoding": "base64", "raw": enc})
	if err != nil {
		return nil, err
	}
	return json.RawMessage(quoted), nil
}

// UnwrapData 取出信封中的 data 字段（拷贝）。
func UnwrapData(raw []byte) ([]byte, error) {
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return CloneBytes(raw), nil
	}
	return CloneBytes(env.Data), nil
}

// MustWrap 忽略错误，失败时退回原始拷贝。
func MustWrap(id, event string, ts time.Time, attempt int, body []byte) []byte {
	b, err := Wrap(id, event, ts, attempt, body)
	if err != nil {
		return CloneBytes(body)
	}
	return b
}
