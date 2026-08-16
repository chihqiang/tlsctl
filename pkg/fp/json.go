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
	var data []T
	if _, err := os.Stat(fileName); err == nil {
		fileBytes, err := os.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		_ = json.Unmarshal(fileBytes, &data)
		for i := range data {
			if j.IsEqual(data[i], newItem) {
				// 找到已存在，执行合并逻辑
				j.Merge(&data[i], newItem)
				// 写回文件并返回
				jsonBytes, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to encode JSON: %w", err)
				}
				if err := os.WriteFile(fileName, jsonBytes, 0644); err != nil {
					return fmt.Errorf("failed to write file: %w", err)
				}
				return nil
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking file status failed: %w", err)
	}
	// 不存在则追加
	data = append(data, newItem)
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	if err := os.WriteFile(fileName, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (j *JSONFile[T]) Load(fileName string) ([]T, error) {
	var data []T
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	err = json.Unmarshal(file, &data)
	return data, err
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
	jsonBytes, _ := json.MarshalIndent(newItems, "", "  ")
	return os.WriteFile(fileName, jsonBytes, 0600)
}
