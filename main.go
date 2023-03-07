package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/d4l3k/go-highlight"
)

var (
	file   = flag.String("file", "", "specify the file")
	spaces = flag.Int("spaces", 0, "indentation spaces, if is 0 then use tabs")
	lang   = flag.String("lang", "go", "specify the language")
	help   = flag.Bool("help", false, "show help")
)

type Line struct {
	Indent int
	Text   string
	Order  int
}

func main() {
	flag.Parse()

	if *file == "" || *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	bytes, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	highlighted, err := highlight.HTML(*lang, bytes)
	if err != nil {
		log.Fatal(err)
	}

	indentRegex := regexp.MustCompile(`^\t+`)

	indentChars := 1

	if *spaces > 0 {
		indentRegex = regexp.MustCompile(`\s+`)
		indentChars = *spaces
	}

	var codeLineList []Line

	scanner := bufio.NewScanner(strings.NewReader(string(highlighted)))
	for scanner.Scan() {
		line := scanner.Text()

		//count the tabs at the beginning of the line
		indent := indentRegex.FindString(line)

		// get length of indent
		indentLength := len(indent) / indentChars

		codeLine := Line{Indent: indentLength}

		line = strings.Replace(line, "</code></pre></div>", "", -1)
		line = strings.Replace(line, "<div class=\"highlight\"><pre><code>", "", -1)

		codeLine.Text = fmt.Sprintf("<div class=line data-indent=%d>%s</div>", indentLength, line)

		codeLineList = append(codeLineList, codeLine)
	}

	for _, line := range codeLineList {
		fmt.Println(line.Text)
	}
}
