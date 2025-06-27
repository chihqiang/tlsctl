package structs

import (
	"reflect"
	"sort"
)

type Tag struct {
	Key   string
	Field string
	Type  string
	Json  string
	Xml   string
	Yaml  string
	Env   string
}

type KeysMaps[K comparable, D any] struct {
	Keys []K
	Maps map[K][]D
}

func TagsMaps[T any](mapConfig map[string]T) KeysMaps[string, Tag] {
	result := make(map[string][]Tag)
	var keys []string
	for k := range mapConfig {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, name := range keys {
		config := mapConfig[name]
		v := reflect.ValueOf(config)

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if !v.IsValid() || v.Kind() != reflect.Struct {
			continue
		}
		t := v.Type()
		var configFieldType reflect.Type
		found := false
		// 查找第一个嵌套的 struct 或 *struct 字段作为配置字段体
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			ft := f.Type
			if ft.Kind() == reflect.Struct {
				configFieldType = ft
				found = true
				break
			}
			if ft.Kind() == reflect.Ptr && ft.Elem().Kind() == reflect.Struct {
				configFieldType = ft.Elem()
				found = true
				break
			}
		}
		// 如果没有嵌套 struct，就用自身
		if !found {
			configFieldType = t
		}
		var tags []Tag
		for i := 0; i < configFieldType.NumField(); i++ {
			f := configFieldType.Field(i)
			env := f.Tag.Get("env")
			if env == "" {
				continue
			}
			tags = append(tags, Tag{
				Key:   name,
				Field: f.Name,
				Type:  f.Type.String(),
				Json:  f.Tag.Get("json"),
				Xml:   f.Tag.Get("xml"),
				Yaml:  f.Tag.Get("yaml"),
				Env:   env,
			})
		}

		if len(tags) > 0 {
			result[name] = tags
		}
	}
	return OrderMaps[string, Tag](result)
}

func OrderMaps[K comparable, D any](input map[K][]D) KeysMaps[K, D] {
	ordered := KeysMaps[K, D]{
		Keys: make([]K, 0, len(input)),
		Maps: make(map[K][]D, len(input)),
	}
	// 收集 key
	for k := range input {
		ordered.Keys = append(ordered.Keys, k)
	}
	// 排序（仅支持常见类型）
	sort.Slice(ordered.Keys, func(i, j int) bool {
		switch any(ordered.Keys[i]).(type) {
		case string:
			return any(ordered.Keys[i]).(string) < any(ordered.Keys[j]).(string)
		case int:
			return any(ordered.Keys[i]).(int) < any(ordered.Keys[j]).(int)
		case int64:
			return any(ordered.Keys[i]).(int64) < any(ordered.Keys[j]).(int64)
		default:
			// 不可比较类型，保持原顺序（不排序）
			return false
		}
	})
	// 构造有序 map
	for _, k := range ordered.Keys {
		ordered.Maps[k] = input[k]
	}
	return ordered
}
