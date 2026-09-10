package fp

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONFile[T any] struct {
	IsEqual func(a, b T) bool
	Merge   func(existing *T, newItem T)
}

func (j *JSONFile[T]) Save(fileName string, newItem T) error {
	data, err := j.Load(fileName)
	if err != nil {
		return err
	}
	for i := range data {
		if j.IsEqual(data[i], newItem) {
			// 找到已存在，执行合并逻辑后写回
			j.Merge(&data[i], newItem)
			return writeJSON(fileName, data)
		}
	}
	// 不存在则追加
	data = append(data, newItem)
	return writeJSON(fileName, data)
}

func (j *JSONFile[T]) Load(fileName string) ([]T, error) {
	var data []T
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return data, nil
	}
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if len(fileBytes) == 0 {
		return data, nil
	}
	if err := json.Unmarshal(fileBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return data, nil
}

func (j *JSONFile[T]) Remove(fileName string, remove func(t T) bool) error {
	data, err := j.Load(fileName)
	if err != nil {
		return err
	}
	var newItems []T
	for _, datum := range data {
		if remove(datum) {
			continue
		}
		newItems = append(newItems, datum)
	}
	return writeJSON(fileName, newItems)
}

// writeJSON 统一 JSON 落盘：两空格缩进 + 0600 权限。
func writeJSON(fileName string, data any) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	if err := os.WriteFile(fileName, jsonBytes, 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
