package cmd

import (
	"fmt"

	"github.com/chihqiang/tlsctl/pkg/stdout"
	"github.com/chihqiang/tlsctl/pkg/structs"
)

// printTagTable 打印 map 配置中各 key 的字段标签表（用于 help:deploy / help:dns）。
func printTagTable(title string, maps structs.KeysMaps[string, structs.Tag]) error {
	tp := stdout.NewTablePrinter()
	for _, key := range maps.Keys {
		tp.SetTitle(fmt.Sprintf("`%s` %s Environment variables and other tags", key, title))
		tp.Add([]string{"Field Name", "Type", "Environment", "JSON", "Yaml"})
		for _, tag := range maps.Maps[key] {
			tp.Add([]string{tag.Field, tag.Type, tag.Env, tag.Json, tag.Yaml})
		}
		if err := tp.Print(); err != nil {
			return err
		}
		tp.Reset()
	}
	return nil
}
