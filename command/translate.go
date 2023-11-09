package command

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hios/config"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Message struct {
	ID    string
	Value string
}

func init() {
	rootCommand.AddCommand(translateCmd)
}

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "开始翻译",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("-------translate start-------\n")
		err := translate()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("-------translate end-------\n")
	},
}

func translate() error {
	var (
		chineseAllText      []string
		directoryPath       string = "./app"
		translationFilePath string = "./i18n/lang/zh.yaml"
	)
	//
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理文本文件
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			chineseList, err := extractChineseTextFromFile(path)
			if err != nil {
				fmt.Printf("Failed to extract Chinese texts from %s: %v\n", path, err)
			} else {
				for _, chineseText := range chineseList {
					exists := false
					for _, str := range chineseAllText {
						if strings.Contains(str, chineseText) {
							exists = true
							break
						}
					}
					if !exists {
						chineseAllText = append(chineseAllText, chineseText)
					}
				}
			}
		}
		return nil
	})
	//
	// 读取现有的翻译文件
	translations := make(map[string]string)
	translationFile, err := os.ReadFile(translationFilePath)
	if err == nil {
		err = yaml.Unmarshal(translationFile, &translations)
		if err != nil {
			fmt.Printf("Failed to parse translation file: %v\n", err)
		}
	}
	// 过滤已经翻译的
	var tobeTranslateTexts []string
	for _, chineseText := range chineseAllText {
		exists := false
		for _, str := range translations {
			if strings.Contains(str, chineseText) {
				exists = true
				break
			}
		}
		if !exists {
			tobeTranslateTexts = append(tobeTranslateTexts, chineseText)
		}
	}
	// 将翻译结果写入到翻译文件中
	for _, chineseText := range tobeTranslateTexts {
		for _, lang := range config.Language {
			// 读取现有的翻译文件
			filePath := "./i18n/lang/" + lang + ".yaml"
			translations := make(map[string]string)
			translationFile, err := os.ReadFile(filePath)
			if err == nil {
				err = yaml.Unmarshal(translationFile, &translations)
				if err != nil {
					fmt.Printf("Failed to parse translation file: %v\n", err)
					break
				}
			}
			// 翻译
			if lang == "zh" {
				translations[chineseText] = chineseText
			} else {
				translations[chineseText] = translateText(chineseText, lang)
			}
			// 将新的YAML数据写入文件
			newYamlData, err := yaml.Marshal(translations)
			err = os.WriteFile(filePath, newYamlData, os.ModePerm)
			if err != nil {
				fmt.Printf("Error writing YAML file: %v", err)
				break
			}
		}
	}
	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}
	return err
}

func extractChineseTextFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return extractChineseText(file)
}

func extractChineseText(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	var paragraph strings.Builder
	chineseMap := make(map[string]bool)

	chinesePattern := regexp.MustCompile(`"([\p{Han}]+)"`) // 匹配被双引号包含的中文文本

	for scanner.Scan() {
		line := scanner.Text()
		paragraph.WriteString(line)
		if strings.Contains(line, "\"") {
			chineseMatches := chinesePattern.FindAllStringSubmatch(line, -1)
			for _, match := range chineseMatches {
				chineseText := match[1] // 匹配结果中的第一个捕获组为中文文本
				chineseMap[chineseText] = true
			}
		}

		if line == "" {
			// 遇到空行，则表示一个段落结束
			paragraphString := paragraph.String()
			if strings.Contains(paragraphString, "\"") {
				chineseMatches := chinesePattern.FindAllStringSubmatch(paragraphString, -1)
				for _, match := range chineseMatches {
					chineseText := match[1] // 匹配结果中的第一个捕获组为中文文本
					chineseMap[chineseText] = true
				}
			}
			// 重置段落缓冲区
			paragraph.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 将去重后的中文文本转换为切片
	chineseList := make([]string, 0, len(chineseMap))
	for chineseText := range chineseMap {
		chineseList = append(chineseList, chineseText)
	}

	return chineseList, nil
}

// 翻译程序
func translateText(text, lang string) string {
	return googleTranslate(text, lang)
}

func logError(msg string, err error) {
	fmt.Printf("\033[41;31m%s %s\033[0m\n", msg, err)
}

// 谷歌翻译
func googleTranslate(text, tl string) string {
	q := url.QueryEscape(text)
	url := fmt.Sprintf("https://translate.google.com/translate_a/single?client=gtx&hl=zh&dt=t&sl=auto&tl=%s&q=%s&ie=UTF-8&oe=UTF-8&multires=1&otf=0&pc=1&trs=1&ssel=0&tsel=0&kc=1&tk=736127.854017", tl, q)
	res, err := http.Get(url)
	if err != nil {
		logError("谷歌翻译错误，请求失败:", err)
		return ""
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logError("谷歌翻译错误:", err)
		return ""
	}
	if res.StatusCode == 200 {
		var jsonObj []interface{}
		if err := json.Unmarshal(body, &jsonObj); err != nil {
			logError("谷歌翻译错误，解析接口返回的JSON内容失败:", err)
			return ""
		}
		if len(jsonObj) > 0 && len(jsonObj[0].([]interface{})) > 0 && len(jsonObj[0].([]interface{})[0].([]interface{})) > 0 {
			var result []string
			for _, h := range jsonObj[0].([]interface{}) {
				result = append(result, h.([]interface{})[0].(string))
			}
			uniqueResult := unique(result)
			return strings.Join(uniqueResult, "")
		}
		logError("谷歌翻译错误，请求失败", nil)
	} else {
		logError("歌翻译错误，请求失败: 请留意网络是否可以正常访问google", nil)
	}
	return ""
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
