package deliver

import (
	"io"
	"net/http"
)

const maxDrain = 1 << 20

// DrainAndClose 读取并关闭响应体。若 rc 为 nil 则直接返回。
func DrainAndClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrain))
	_ = resp.Body.Close()
}
