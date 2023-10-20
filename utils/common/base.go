package common

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// PrintError 错误输出
func PrintError(msg string) {
	fmt.Printf("\033[1;31m" + msg + " \033[0m\n")
}

// PrintSuccess 正确输出
func PrintSuccess(msg string) {
	fmt.Printf("\033[1;32m" + msg + " \033[0m\n")
}

// CheckOs 判断系统类型
func CheckOs() bool {
	return runtime.GOOS == "darwin" || runtime.GOOS == "linux"
}

// Test 正则判断
func Test(str, pattern string) bool {
	re := regexp.MustCompile(pattern)
	if re.MatchString(str) {
		return true
	} else {
		return false
	}
}

// RunDir 前面加上绝对路径
func RunDir(path string, a ...any) string {
	wd, _ := os.Getwd()
	if len(a) > 0 {
		path = fmt.Sprintf(path, a...)
	}
	return fmt.Sprintf("%s%s", wd, path)
}

// GeneratePassword 生成随机密码
// 密码字符集，可任意添加你需要的字符，1数字、2大小写字母、21小写字母、22大写字母、默认全部
func GeneratePassword(length int, t string) string {
	var chars string
	switch t {
	case "1":
		chars = "0123456789"
	case "2":
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "21":
		chars = "abcdefghijklmnopqrstuvwxyz"
	case "22":
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	default:
		if t == "" {
			chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		} else {
			chars = t
		}
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	passwordStr := make([]byte, length)
	max := len(chars) - 1
	for i := 0; i < length; i++ {
		passwordStr[i] = chars[r.Intn(max)]
	}
	return string(passwordStr)
}

// IsKind 判断一个值是否为指定的类型
func IsKind(v interface{}, kind reflect.Kind) bool {
	return reflect.TypeOf(v).Kind() == kind
}

// InSlice 判断一个值是否在切片中
func InSlice(v interface{}, s []interface{}) bool {
	for _, item := range s {
		if item == v {
			return true
		}
	}
	return false
}

// Unique 去重
func Unique(s []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	result := []interface{}{}
	for _, item := range s {
		if !m[item] {
			m[item] = true
			result = append(result, item)
		}
	}
	return result
}

// UniqueInt int去重
func UniqueInt(ints []int) []int {
	uniqueMap := make(map[int]bool)
	uniqueInts := []int{}
	for _, i := range ints {
		if _, ok := uniqueMap[i]; !ok {
			uniqueMap[i] = true
			uniqueInts = append(uniqueInts, i)
		}
	}
	return uniqueInts
}

// UniqueInt int去重
func UniqueStr(ints []string) []string {
	uniqueMap := make(map[string]bool)
	uniqueInts := []string{}
	for _, i := range ints {
		if _, ok := uniqueMap[i]; !ok {
			uniqueMap[i] = true
			uniqueInts = append(uniqueInts, i)
		}
	}
	return uniqueInts
}

// 获取后缀名图标相对地址
func ExtIcon(ext string) string {
	switch ext {
	case "docx":
		return "images/ext/doc.png"
	case "xlsx":
		return "images/ext/xls.png"
	case "pptx":
		return "images/ext/ppt.png"
	case "ai", "avi", "bmp", "cdr", "doc", "eps", "gif", "mov", "mp3", "mp4", "pdf", "ppt", "pr", "psd", "rar", "svg", "tif", "txt", "xls", "zip":
		return "images/ext/" + ext + ".png"
	default:
		return "images/ext/file.png"
	}
}

// 相对路径补全
func FillUrl(path string) string {
	if path == "" {
		return path
	}
	if strings.HasPrefix(path, "//") ||
		strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "ftp://") ||
		strings.HasPrefix(path, "/") ||
		strings.HasPrefix(strings.ReplaceAll(path, " ", ""), "data:image/") {
		return path
	}
	// TODO
	// if ($_A['__fill_url_remote_url'] === true) {
	// 	return "{{RemoteURL}}" . $str;
	// }
	// try {
	// 	return url($str);
	// } catch (\Throwable) {
	// 	return self::getSchemeAndHost() . "/" . $str;
	// }

	return path
}

// UnFillUrl 相对路径去掉
func UnFillUrl(path string) string {
	return path
}

// BoolToInt bool转int
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IsUrl 判断是否是url
func IsHttpsUrl(str string) bool {
	if len(str) < 8 {
		return false
	}
	matched, _ := regexp.MatchString("^https*://", str)
	return matched
}
