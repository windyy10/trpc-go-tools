package mysql

import (
	"fmt"
	"reflect"
	"strings"
)

func GormFields(dest interface{}) []string {
	fields, _ := GetGormFields(dest)
	return fields
}

func GetGormFields(dest interface{}) ([]string, error) {
	var fields []string
	pointType := reflect.TypeOf(dest)
	if pointType.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("dest not pointer type")
	}
	elemType := pointType.Elem()
	if elemType.Kind() == reflect.Slice {
		elemType = elemType.Elem()
	}

	// 支持指针类型
	if elemType.Kind() == reflect.Pointer {
		elemType = elemType.Elem()
	}

	// 确保元素类型是结构体
	if elemType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("elemType are not structs")
	}

	// 遍历结构体字段
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		// 获取字段的gorm标签
		gormTag := field.Tag.Get("gorm")
		columnName := extractColumnName(gormTag)
		if columnName == "" {
			// 没有column直接拼字段名称, 这样方便定位问题原因
			fields = append(fields, field.Name)
		} else {
			fields = append(fields, columnName)
		}
	}
	return fields, nil
}

func extractColumnName(gormTag string) string {
	parts := strings.Split(gormTag, ":")
	if len(parts) >= 2 {
		if parts[0] == "column" {
			return parts[1]
		}
	}
	return ""
}
