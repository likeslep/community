// Package config 提供基于环境变量的配置加载（plan.md §33 Phase 0）。
//
// 使用 struct tag 声明配置项：
//
//	type Config struct {
//	    Addr string        `env:"HTTP_ADDR" default:":8080"`
//	    DSN  string        `env:"DB_DSN" required:"true"`
//	    Num  int           `env:"WORKERS" default:"4"`
//	}
package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Load 将环境变量绑定到 cfg（必须为指向 struct 的非空指针）。
// 支持 tag：env（变量名）、default（默认值）、required（必填）。
func Load(cfg any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("config: Load 需要指向 struct 的非空指针")
	}
	if v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config: Load 需要 struct，实际 %s", v.Elem().Kind())
	}
	return loadStruct(v.Elem())
}

func loadStruct(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		sf := t.Field(i)
		if sf.PkgPath != "" { // 跳过未导出字段
			continue
		}
		key := sf.Tag.Get("env")
		if key == "" {
			continue
		}

		raw, ok := os.LookupEnv(key)
		switch {
		case ok && raw != "":
			// 使用环境变量值
		case sf.Tag.Get("default") != "":
			raw = sf.Tag.Get("default")
		case sf.Tag.Get("required") == "true":
			return fmt.Errorf("config: 缺少必填环境变量 %s", key)
		default:
			continue // 非必填且无默认值，保持零值
		}

		if err := setField(f, raw); err != nil {
			return fmt.Errorf("config: %s: %w", key, err)
		}
	}
	return nil
}

func setField(f reflect.Value, raw string) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(raw)
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		f.SetBool(b)
	case reflect.Int, reflect.Int64:
		if f.Type() == reflect.TypeFor[time.Duration]() {
			d, err := time.ParseDuration(raw)
			if err != nil {
				return err
			}
			f.SetInt(int64(d))
			return nil
		}
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		f.SetInt(n)
	default:
		return fmt.Errorf("不支持的类型 %s", f.Type())
	}
	return nil
}
