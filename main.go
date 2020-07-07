package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"

	iconv "github.com/djimenez/iconv-go"
	"github.com/urfave/cli"
	"github.com/wxnacy/wgo/arrays"
)

/*
 to-do:
 1. 注释和说明使用英文
 2. 发布v1.0.0版本
 3. chinese readme_zh.md file
*/

// TemplateOpeningTag 替换正文里的 {
var TemplateOpeningTag = "___TEMPLATE_OPENING_TAG___"

// TemplateClosingTag 替换正文里的 }
var TemplateClosingTag = "___TEMPLATE_CLOSING_TAG___"

// FormatArgs 命令行的参数
type FormatArgs struct {
	BlankSpace int
	Charset    string
	Backup     bool
	Verbose    bool
	Testing    bool
}

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "nginxConfigFormatterGo"
	app.Usage = "Nginx config file formatter"
	app.Author = "github.com/rwx------"
	app.Description = "Nginx config file formatter"
	app.UsageText = "./nginxConfigFormatterGo [-s 2] [-c utf-8] [-b] [-v] [-t] <filelists>"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "charset, c",
			Value:    "utf-8",
			Required: false,
			Usage:    "current supported charset: gbk, gb18030, windows-1252, utf-8",
		},
		cli.IntFlag{
			Name:  "space, s",
			Value: 4,
			Usage: "blank spaces indentation",
		},
		cli.BoolFlag{
			Name:     "backup, b",
			Required: false,
			Usage:    "backup the original config file",
		},
		cli.BoolFlag{
			Name:     "verbose, v",
			Required: false,
			Usage:    "verbose mode",
		},
		cli.BoolFlag{
			Name:     "testing, t",
			Required: false,
			Usage:    "perform a test run with no changes made." + "\n                              show the formatted contents on the standard output.",
		},
	}

	app.Action = func(c *cli.Context) error {

		// 赋值
		var f = FormatArgs{
			c.Int("space"),
			c.String("charset"),
			c.Bool("backup"),
			c.Bool("verbose"),
			c.Bool("testing"),
		}

		// 检查字符集
		if !checkCharset(f.Charset) {
			s := "不支持的字符集! 当前只支持: \"gbk\", \"gb18030\", \"windows-1252\", \"utf-8\"\n"
			s += "将终止配置文件的格式化!\n"
			errorMessage(s, true)
			return nil
		}

		if c.NArg() > 0 {
			for _, conf := range c.Args() {
				// 检查传入的文件是否存在
				if IsFile(conf) {
					// 进行格式化处理
					f.formatConfigFile(conf)
				} else {
					s := "需要传入一个文件, 文件: \"" + conf + "\" 不是一个文件!\n"
					errorMessage(s, true)
				}
			}
		} else {
			s := "需要传入文件\n"
			errorMessage(s, true)
		}
		return nil
	}
	app.Run(os.Args)
}

func (f *FormatArgs) formatConfigFile(configFilePath string) {
	/*
		1. 首先以正确的编码打开文件
		2. 然后以正确的编码读取文件
		3. 判断文件内容是否为空
		4. 判断是否需要备份, 若要备份, 则进行备份(以原有的编码进行备份).
			4.1 判断是否需要显示详细信息
		5. 以utf8格式转码, 然后进行文件格式化
			5.1 将格式化后的内容, 以原编码格式写入到文件.

	*/

	// 获取文件内容, 并转换为utf-8编码
	fc := ReadAll(configFilePath)
	if f.Charset != "utf-8" {
		// 转换为utf8字符集
		fcIconv, err := iconv.ConvertString(fc, f.Charset, "utf-8")
		if err != nil {
			s := fmt.Sprintf("You want convert the strings from %v to utf-8, but could not convert!", f.Charset)
			errorMessage(s, false)
			return
		}
		fc = fcIconv
	}

	// 判断文件是否为空
	if len(fc) == 0 {
		s := fmt.Sprintf("%v是一个空文件\n", configFilePath)
		errorMessage(s, false)
		return
	}

	// 此方法不用关心原来的字符集是什么, 备份的文件还是原来的字符集.
	if f.Backup {
		_, err := copyFile(configFilePath, configFilePath+"~")
		if err != nil {
			s := fmt.Sprintf("%v 备份失败\n, \n%v", configFilePath, err)
			errorMessage(s, false)
			// 当出现备份错误的时候, 不再进行后面的真正格式化
			return
		}
	}

	// 具体执行配置文件格式化
	fcNew := f.formatConfigContent(fc)

	if f.Testing {
		fmt.Println(fcNew)
	} else {
		// 进行编码格式转换
		if f.Charset != "utf-8" {
			fcNew, _ = iconv.ConvertString(fcNew, "utf-8", f.Charset)
		}

		// 写入新文件
		err := writeNewConfig(configFilePath, fcNew)
		if err != nil {
			s := fmt.Sprintf("%v 写入新文件失败\n, \n%v", configFilePath, err)
			errorMessage(s, false)
		}
	}
}

// copyFile 复制文件
func copyFile(dstName, srcName string) (writeen int64, err error) {
	src, err := os.Open(dstName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(srcName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

// checkCharset 检查是否为受支持的字符集
func checkCharset(s string) bool {
	charsetList := []string{"gbk", "gb18030", "windows-1252", "utf-8"}
	i := arrays.ContainsString(charsetList, s)
	if i == -1 {
		return false
	}
	return true
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !s.IsDir()
}

// ReadAll 读取到file中，再利用ioutil将file直接读取到[]byte中
func ReadAll(filePth string) string {
	f, err := os.Open(filePth)
	if err != nil {
		fmt.Println("read file fail", err)
		return ""
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return ""
	}

	return string(fd)
}

// decomposeLine 分解一行的内容
// 返回值 []string  分解后的行slice
// 返回值 bool 是否有分解发生
func decomposeLine(line string) (ls []string, mFlag bool) {
	mFlag = false
	ls = strings.Split(line, "\n")

	if len(ls) > 1 {
		mFlag = true
	}

	return
}

func cheackEveryChar(line string) string {
	var inQuotes = false
	var commentFlag = false
	var result string
	var lastC rune
	var lastQuote rune

	var c []rune
	c = []rune(line)
	cLen := len(c) - 1
	for i, k := range c {
		if commentFlag == true { // content after `#`
			result += string(k)
		} else { // content before `#`

			// whether inQuotes
			if (k == '"' || k == '\'') && lastC != '\\' {
				if k != lastQuote && lastQuote == 0 {
					inQuotes = true
					lastQuote = k
				} else if k == lastQuote && lastQuote != 0 {
					inQuotes = false
					lastQuote = 0
				}
			}

			if inQuotes == true {
				if k == '{' {
					result += TemplateOpeningTag
				} else if k == '}' {
					result += TemplateClosingTag
				} else {
					result += string(k)
				}
			} else {
				if k == '#' {
					result += "\n#"
					commentFlag = true
				} else {
					// `;`, `{`, `}` turn into a newline
					if k == ';' && i != cLen {
						result += ";\n"
					} else if k == '{' && i != cLen {
						result += " {\n"
					} else if k == '}' && i != 0 {
						result += "\n}\n"
					} else if k == '}' && i == 0 {
						result += "}\n"
					} else {
						// whitespaces are collapsed
						if unicode.IsSpace(k) && lastC != ' ' {
							lastC = ' '
							result += " "
							continue
						} else if unicode.IsSpace(k) && lastC == ' ' {
							lastC = ' '
							continue
						} else {
							result += string(k)
						}
					}
				}
			}
			lastC = k
		}
	}
	return result
}

func (f *FormatArgs) formatConfigContent(fc string) string {
	/*
		1. 将引号内的 {} 进行替换
		2. 将内容分割为行(\n)
		3. 按行给每行进行处理
		4. 处理行中的 '{' (opening bracket)
		5. 处理缩进情况
		6. 合并行
		7. 将括号的替换进行替换回来.
		8. 清理多余的空行
		9. 返回内容
	*/

	// 按行进行分割
	lines := strings.Split(fc, "\n")
	if f.Verbose {
		fmt.Printf("\n==Split:===\n%#v\n=======\n", lines)
	}

	lines = cleanLines(lines)
	if f.Verbose {
		fmt.Printf("\n==cleanLines:===\n%#v\n=======\n", lines)
	}

	lines = joinOpeningBracket(lines)
	if f.Verbose {
		fmt.Printf("\n==joinOpeningBracket:===\n%#v\n=======\n", lines)
	}

	lines = performIndentation(lines, f.BlankSpace)
	if f.Verbose {
		fmt.Printf("\n==performIndentation:===\n%#v\n=======\n", lines)
	}

	text := strings.Join(lines, "\n")
	if f.Verbose {
		fmt.Printf("\n==strings.Join:===\n%#v\n=======\n", text)
	}

	text = stripBracketTemplateTags(text)
	if f.Verbose {
		fmt.Printf("\n==stripBracketTemplateTags:===\n%#v\n=======\n", text)
	}

	return text
}

func stripBracketTemplateTags(content string) string {
	content = strings.ReplaceAll(content, TemplateOpeningTag, "{")
	content = strings.ReplaceAll(content, TemplateClosingTag, "}")
	return content
}

func performIndentation(lines []string, blankSpace int) []string {
	newLines := make([]string, 0, cap(lines))
	currentIndent := 0
	for _, line := range lines {
		if (!strings.HasPrefix(line, "#")) && strings.HasSuffix(line, "}") && currentIndent > 0 {
			currentIndent--
		}

		if line != "" {
			newLines = append(newLines, strings.Repeat(" ", blankSpace*currentIndent)+line)
		} else {
			newLines = append(newLines, "")
		}

		if !strings.HasPrefix(line, "#") && strings.HasSuffix(line, "{") {
			currentIndent++
		}
	}
	return newLines
}

// joinOpeningBracket 当 { 为单独一行的时候, 合并到上一行
func joinOpeningBracket(lines []string) []string {
	newLines := make([]string, 0, cap(lines))

	lastLine := ""
	for i, l := range lines {
		if lastLine != "{" {
			if (lastLine == "" && l == "") || i == 0 {
				lastLine = l
				continue
			} else if lastLine == "" && l == "}" {
				lastLine = "}"
				continue
			} else if strings.HasSuffix(lastLine, "{") && l == "" {
				continue
			} else if i > 0 && lastLine == "" && l == "{" {
				newLines[len(newLines)-1] += " {"
			} else if i > 0 && lastLine != "" && l == "{" {
				newLines = append(newLines, lastLine+" {")
			} else {
				newLines = append(newLines, lastLine)
			}
		} else if lastLine == "{" && l == "" {
			continue
		}

		lastLine = l
	}
	// 把最后一行加入进来
	newLines = append(newLines, lastLine)

	return newLines
}

func cleanLines(lines []string) []string {
	cleanedLines := make([]string, 0, cap(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			cleanedLines = append(cleanedLines, l)
		} else if strings.HasPrefix(l, "#") {
			cleanedLines = append(cleanedLines, l)
		} else {
			l = cheackEveryChar(l)
			newLines, ok := decomposeLine(l)

			if ok {
				nl := make([]string, 0, cap(newLines))
				nl = cleanAgain(newLines)
				cleanedLines = append(cleanedLines, nl...)
			} else {
				cleanedLines = append(cleanedLines, l)
			}
		}
	}
	return cleanedLines
}

func cleanAgain(lines []string) []string {
	cleanedLines := make([]string, 0, cap(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		cleanedLines = append(cleanedLines, l)
	}
	return cleanedLines
}

func writeNewConfig(Path string, content string) error {
	text := []byte(content)
	return ioutil.WriteFile(Path, text, 0644)
}

func errorMessage(s string, b bool) {

	if b == true {
		usageText := "./nginxConfigFormatterGo [-s 2] [-c utf-8] [-b] [-v] [-t] <filelists>"
		s += "\n[usage]:\n" + usageText + "\n"
	}

	fmt.Printf("\n[error]:\n%v", s)
}
