// Package migrations 内嵌 notification-service 的数据库迁移脚本。
package migrations

import "embed"

// FS 包含 *.sql 迁移文件（up/down 成对）。
//
//go:embed *.sql
var FS embed.FS
