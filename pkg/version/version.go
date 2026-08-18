// Package version 提供构建版本信息，供各服务在启动日志与 health 响应中上报。
package version

// 以下变量在构建时通过 -ldflags 注入，默认值用于本地开发。
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Info 汇总一次构建的版本元信息。
type Info struct {
	Version string
	Commit  string
	Date    string
}

// Get 返回当前构建的版本信息。
func Get() Info {
	return Info{Version: version, Commit: commit, Date: date}
}
