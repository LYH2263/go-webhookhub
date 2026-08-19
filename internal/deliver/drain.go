package deliver

import (
	"io"
	"net/http"
)

const maxDrain = 1 << 20

// DrainAndClose 读取并关闭响应体。若 rc 为 nil 则直接返回。
// 必须调用 Close，否则底层 TCP 连接无法归还连接池，长跑会耗尽连接池。
func DrainAndClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrain))
}
