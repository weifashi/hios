package database

import (
	"hios/core"
	"hios/database/migrations"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/go-gormigrate/gormigrate/v2"
)

// 把驼峰名转 -
func convertToKebabCase(str string) string {
	var result strings.Builder
	for i, ch := range str {
		if unicode.IsUpper(ch) {
			if i > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(ch))
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

// 初始化数据库
func Init() error {
	value := reflect.ValueOf(migrations.InitDatabase{})
	structType := reflect.TypeOf(migrations.InitDatabase{})
	// 遍历结构体的方法
	defaultMigration := []*gormigrate.Migration{}
	for i := 0; i < structType.NumMethod(); i++ {
		method := structType.Method(i)
		methodValue := value.MethodByName(method.Name)
		result := methodValue.Call(nil)
		resultValue := result[0].Interface().(*gormigrate.Migration)
		if resultValue.ID == "" {
			resultValue.ID = convertToKebabCase(method.Name)
		} else {
			resultValue.ID = resultValue.ID + "-" + convertToKebabCase(method.Name)
		}
		defaultMigration = append(defaultMigration, resultValue)
	}
	// 排序
	sort.Slice(defaultMigration, func(i, j int) bool {
		return defaultMigration[i].ID < defaultMigration[j].ID
	})
	//
	defaultOptions := &gormigrate.Options{
		TableName:                 "xw_hios_migrations",
		IDColumnName:              "id",
		IDColumnSize:              255,
		UseTransaction:            false,
		ValidateUnknownMigrations: false,
	}
	m := gormigrate.New(core.DB, defaultOptions, defaultMigration)
	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}
