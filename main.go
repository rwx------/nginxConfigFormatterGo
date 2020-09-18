package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"

	iconvgo "github.com/djimenez/iconv-go"
	"github.com/urfave/cli"
	"github.com/wxnacy/wgo/arrays"
)

// TemplateOpeningTag  change { to this const
const TemplateOpeningTag = "___TEMPLATE_OPENING_TAG___"

// TemplateClosingTag  change } to this const
const TemplateClosingTag = "___TEMPLATE_CLOSING_TAG___"

// FormatArgs args
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
			Usage:    "perform a test run with no changes made",
		},
	}

	app.Action = func(c *cli.Context) error {
		var f = FormatArgs{
			c.Int("space"),
			c.String("charset"),
			c.Bool("backup"),
			c.Bool("verbose"),
			c.Bool("testing"),
		}

		// check charset
		if !checkCharset(f.Charset) {
			s := `Do not support the charst!` + "\n"
			s += `We now support this charsets:"gbk", "gb18030", "windows-1252", "utf-8"` + "\n"
			errorMessage(s, true)
			return nil
		}

		if c.NArg() > 0 {
			for _, conf := range c.Args() {
				if isFile(conf) {
					f.formatConfigFile(conf)
				} else {
					s := "You should give a filename to the programe to format.\n"
					errorMessage(s, true)
				}
			}
		} else {
			s := "You should give a filename to the programe to format.\n"
			errorMessage(s, true)
		}
		return nil
	}
	app.Run(os.Args)
}

func (f *FormatArgs) formatConfigFile(configFilePath string) {

	fc := readAll(configFilePath)
	if f.Charset != "utf-8" {
		// convert the content to utf-8
		fcIconv, err := iconvgo.ConvertString(fc, f.Charset, "utf-8")
		if err != nil {
			s := fmt.Sprintf("You want convert the strings from %v to utf-8, but could not convert!", f.Charset)
			errorMessage(s, false)
			return
		}
		fc = fcIconv
	}

	// if the file is empty
	if len(fc) == 0 {
		s := fmt.Sprintf("%v is an empty file.\n", configFilePath)
		errorMessage(s, false)
		return
	}

	if f.Backup {
		_, err := copyFile(configFilePath, configFilePath+"~")
		if err != nil {
			s := fmt.Sprintf("%v backup failed\n, \n%v", configFilePath, err)
			errorMessage(s, false)
			// if backup failed, then no further format would do.
			return
		}
	}

	// formating
	fcNew := f.formatConfigContent(fc)

	if f.Testing {
		fmt.Println(fcNew)
	} else {
		// change the charset back
		if f.Charset != "utf-8" {
			fcNew, _ = iconvgo.ConvertString(fcNew, "utf-8", f.Charset)
		}

		// formated content write into the file
		err := writeNewConfig(configFilePath, fcNew)
		if err != nil {
			s := fmt.Sprintf("%v formated content wrote error\n, \n%v", configFilePath, err)
			errorMessage(s, false)
		}
	}
}

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

// checkCharset check the charset parameter
func checkCharset(s string) bool {
	charsetList := []string{"gbk", "gb18030", "windows-1252", "utf-8"}
	i := arrays.ContainsString(charsetList, s)
	if i == -1 {
		return false
	}
	return true
}

func isFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !s.IsDir()
}

func readAll(filePth string) string {
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

// decomposeLine decompose one line
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

// joinOpeningBracket join brackets and collapse multi blank lines
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
	// the last line
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
