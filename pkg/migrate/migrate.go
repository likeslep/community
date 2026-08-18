// Package migrate 提供基于 SQL 文件的版本化数据库迁移（plan.md §41）。
package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
)

// migration 是一条迁移：一个 version 对应 up / down 两个 SQL。
type migration struct {
	version int
	name    string
	up      string
	down    string
}

// Runner 执行版本化迁移，版本记录在 schema_migrations 表。
type Runner struct {
	db *sql.DB
	ms []migration
}

// New 构造迁移执行器。
func New(db *sql.DB) *Runner { return &Runner{db: db} }

// Load 从 fs.FS 加载 dir 下的迁移文件（形如 0001_name.up.sql / 0001_name.down.sql）。
func (r *Runner) Load(fsys fs.FS, dir string) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return err
	}
	byVersion := map[int]*migration{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		v, name, isUp, ok := parseName(e.Name())
		if !ok {
			continue
		}
		m := byVersion[v]
		if m == nil {
			m = &migration{version: v, name: name}
			byVersion[v] = m
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		if isUp {
			m.up = string(data)
		} else {
			m.down = string(data)
		}
	}
	for _, m := range byVersion {
		if m.up == "" {
			return fmt.Errorf("migrate: 版本 %d 缺少 up 脚本", m.version)
		}
		r.ms = append(r.ms, *m)
	}
	sort.Slice(r.ms, func(i, j int) bool { return r.ms[i].version < r.ms[j].version })
	return nil
}

// Up 应用所有未应用的迁移，每个迁移在一个事务内执行。
func (r *Runner) Up(ctx context.Context) error {
	if err := r.ensureTable(ctx); err != nil {
		return err
	}
	current, err := r.version(ctx)
	if err != nil {
		return err
	}
	for _, m := range pending(r.ms, current) {
		if err := r.apply(ctx, m, m.up); err != nil {
			return fmt.Errorf("migrate: 应用 %d_%s 失败: %w", m.version, m.name, err)
		}
	}
	return nil
}

// Down 回滚最近一个已应用的迁移。
func (r *Runner) Down(ctx context.Context) error {
	if err := r.ensureTable(ctx); err != nil {
		return err
	}
	current, err := r.version(ctx)
	if err != nil {
		return err
	}
	for i := len(r.ms) - 1; i >= 0; i-- {
		m := r.ms[i]
		if m.version == current {
			if m.down == "" {
				return fmt.Errorf("migrate: 版本 %d 缺少 down 脚本", m.version)
			}
			return r.rollback(ctx, m)
		}
	}
	return nil
}

// Version 返回当前已应用的迁移版本（未初始化时为 0）。
func (r *Runner) Version(ctx context.Context) (int, error) {
	if err := r.ensureTable(ctx); err != nil {
		return 0, err
	}
	return r.version(ctx)
}

func (r *Runner) ensureTable(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version BIGINT PRIMARY KEY,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func (r *Runner) version(ctx context.Context) (int, error) {
	var v int
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&v)
	return v, err
}

func (r *Runner) apply(ctx context.Context, m migration, sqlText string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // 提交后 Rollback 为 no-op
	if _, err := tx.ExecContext(ctx, sqlText); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES (?)`, m.version); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Runner) rollback(ctx context.Context, m migration) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // 提交后 Rollback 为 no-op
	if _, err := tx.ExecContext(ctx, m.down); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM schema_migrations WHERE version = ?`, m.version); err != nil {
		return err
	}
	return tx.Commit()
}

// pending 返回版本大于 current 的迁移。
func pending(ms []migration, current int) []migration {
	var out []migration
	for _, m := range ms {
		if m.version > current {
			out = append(out, m)
		}
	}
	return out
}

// parseName 解析迁移文件名：0001_create_users.up.sql → (1, "create_users", up, true)。
func parseName(filename string) (version int, name string, isUp bool, ok bool) {
	s := strings.TrimSuffix(filename, ".sql")
	switch {
	case strings.HasSuffix(s, ".up"):
		isUp = true
		s = strings.TrimSuffix(s, ".up")
	case strings.HasSuffix(s, ".down"):
		s = strings.TrimSuffix(s, ".down")
	default:
		return 0, "", false, false
	}
	i := strings.Index(s, "_")
	if i <= 0 {
		return 0, "", false, false
	}
	v, err := strconv.Atoi(s[:i])
	if err != nil {
		return 0, "", false, false
	}
	return v, s[i+1:], isUp, true
}
