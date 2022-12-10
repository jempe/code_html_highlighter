package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/d4l3k/go-highlight"
)

var (
	lang = flag.String("lang", "go", "specify the language")
	file = flag.String("file", "", "specify the file")
)

func main() {
	flag.Parse()

	if *file == "" {
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

	var codeLines string

	scanner := bufio.NewScanner(strings.NewReader(string(highlighted)))
	for scanner.Scan() {
		line := scanner.Text()

		lineWithoutTabs := strings.TrimLeft(line, "\t")

		totalTabs := len(line) - len(lineWithoutTabs)

		codeLines += fmt.Sprintf("<div class=line data-indent=%d>%s</div>\n", totalTabs, line)
	}

	fmt.Println(codeLines)
}
