package validator

import (
	log "github.com/dulisoft/spirit/core/log/zapx"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func TrimSpace(fl validator.FieldLevel) bool {
	value := fl.Field()
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			// is nil, no validate
			return true
		}

		value = value.Elem()
	}

	if value.Kind() != reflect.String {
		log.Warnf("field type not is string, kind: [%v]", value.Kind())
		return true
	}

	if !value.CanSet() {
		log.Warnf("field not can set, struct name: [%v], field name: [%v]", fl.Top().Type().Name(), fl.StructFieldName())
		return false
	}

	value.SetString(strings.TrimSpace(value.String()))

	return true
}

func verifyModelID(fl validator.FieldLevel) bool {
	value := fl.Field()
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return true
		}

		value = value.Elem()
	}

	if value.Kind() != reflect.String {
		log.Warnf("field type not is string, kind: [%v]", value.Kind())
		return false
	}

	omit := fl.Param() == "omit"

	idStr := strings.TrimSpace(value.String())
	if len(idStr) == 0 {
		if omit {
			return true
		}

		log.Errorf("id string show is empty")
		return false
	}

	ui64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Errorf("id real type is not uint64, err: %v", err)
		return false
	}

	if !omit && ui64 < 1 {
		log.Error("id lt 1")
		return false
	}

	value.SetString(idStr)

	return true
}

func VerifyDateString(fl validator.FieldLevel) bool {
	//日期格式必须符合2021-01-01
	f := fl.Field().String()
	compile := regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2}$")
	return compile.Match([]byte(f))
}

func VerifyTimeString(fl validator.FieldLevel) bool {
	//时间格式必须符合13:12
	f := fl.Field().String()
	compile := regexp.MustCompile("^\\d{2}:\\d{2}$")
	return compile.Match([]byte(f))
}

// ValidateSnowflakeID 雪花ID结构验证器
func ValidateSnowflakeID(fl validator.FieldLevel) bool {
	// 雪花ID的结构验证逻辑
	value := fl.Field().String()
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil || n <= 0 {
		return false
	}
	return true
}

func VerifyXssString(fl validator.FieldLevel) bool {
	// can be empty
	f := fl.Field().String()
	f = strings.TrimSpace(f)
	f = XssEscape(f)
	fl.Field().SetString(f)
	return true
}

func VerifyHost(fl validator.FieldLevel) bool {
	f := fl.Field().String()
	if regexp.MustCompile("^[0-9a-zA-Z\\n\\n]([-.\\\\w]*[0-9a-zA-Z])*$").Match([]byte(f)) {
		return true //url
	}
	if regexp.MustCompile("^(?:(?:1[0-9][0-9]\\\\.)|(?:2[0-4][0-9]\\\\.)|(?:25[0-5]\\\\.)|(?:[1-9][0-9]\\\\.)|(?:[0-9]\\\\.)){3}(?:(?:1[0-9][0-9])|(?:2[0-4][0-9])|(?:25[0-5])|(?:[1-9][0-9])|(?:[0-9]))").Match([]byte(f)) {
		return true //ipv4
	}
	if regexp.MustCompile("\\b(?:[a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}\\b").Match([]byte(f)) {
		return true //ipv6
	}
	return false
}

func XssEscape(values string) string {
	if values == "" {
		return values
	}
	special := strings.NewReplacer(`<`, `&lt;`, `>`, `&gt;`, `select`, `查询`, `drop`, `删除表`, `delete`, `删除数据`, `update`, `更新`, `insert`,
		`插入`, `SELECT`, `查询`, `DROP`, `删除表`, `DELETE`, `删除数据`, `UPDATE`, `更新`, `INSERT`, `插入`, `script`, `脚本`, `SCRIPT`, `脚本`, `ALTER`,
		`修改结构`, `alter`, `修改结构`, `create`, `创建`, `CREATE`, `创建`)
	return special.Replace(values)
}
