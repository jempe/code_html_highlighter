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
	ID          int
	Order       int
	Indent      int
	Text        string
	Type        string
	Block       int
	ParentBlock int
}
type Block struct {
	ID      int
	LineIDs []int
}

type CodeLineList struct {
	LineList  []Line
	BlockList []Block
}

func (c *CodeLineList) AddToList(Text string, Indent int, Order int, BlockNumber int, ParentBlock int) (blockId int) {

	firstBlockNumber := 101

	// check if the BlockList is empty
	if len(c.BlockList) == 0 {
		block := Block{ID: firstBlockNumber, LineIDs: []int{}}
		c.BlockList = append(c.BlockList, block)
		BlockNumber = firstBlockNumber
	}

	if BlockNumber == 0 {
		newBlockNumber := c.GetLastBlock() + 1
		block := Block{ID: newBlockNumber, LineIDs: []int{}}
		c.BlockList = append(c.BlockList, block)

		BlockNumber = newBlockNumber
	}

	var line Line

	line.ID = len(c.LineList) + 1
	line.Indent = Indent
	line.Text = Text
	line.Order = Order
	line.Block = BlockNumber
	line.ParentBlock = ParentBlock

	c.LineList = append(c.LineList, line)

	c.AddLineToBlock(BlockNumber, line.ID)

	return BlockNumber
}

func (c *CodeLineList) AddLineToBlock(BlockNumber int, LineID int) {
	for i, block := range c.BlockList {
		if block.ID == BlockNumber {
			c.BlockList[i].LineIDs = append(c.BlockList[i].LineIDs, LineID)
		}
	}
}

func (c *CodeLineList) GetBlockOfLine(LineID int) (BlockID int) {
	for _, block := range c.BlockList {
		for _, lineID := range block.LineIDs {
			if lineID == LineID {
				return block.ID
			}
		}
	}

	return 0
}

func (c *CodeLineList) GetLastLineOfBlock(BlockID int) (line Line) {
	for _, block := range c.BlockList {
		if block.ID == BlockID {
			return c.LineList[block.LineIDs[len(block.LineIDs)-1]-1]
		}
	}

	return Line{}
}

func (c *CodeLineList) GetLastBlock() (BlockID int) {
	return c.BlockList[len(c.BlockList)-1].ID
}

func (c *CodeLineList) Print() {
	for _, line := range c.LineList {
		//fmt.Println(line.ID, line.Text, line.Order, line.Block, line.ParentBlock)
		lineOrder := line.Block*1000 + line.Order*10

		fmt.Println(fmt.Sprintf("<div class=line data-indent=%d data-order=%d>%s</div>", line.Indent, lineOrder, line.Text))
	}
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

	var codeLines CodeLineList

	var lastIndent = 0
	var blockLine = 0
	var blockNumber = 0
	var parentBlockNumber = 0

	scanner := bufio.NewScanner(strings.NewReader(string(highlighted)))
	for scanner.Scan() {
		line := scanner.Text()

		//count the tabs at the beginning of the line
		indent := indentRegex.FindString(line)

		// get length of indent
		indentLength := len(indent) / indentChars

		line = strings.TrimSuffix(line, " ")
		line = strings.Replace(line, "</code></pre></div>", "", -1)
		line = strings.Replace(line, "<div class=\"highlight\"><pre><code>", "", -1)

		var isStartOfBlock = false

		if lastIndent < indentLength || strings.HasSuffix(line, "{") {
			isStartOfBlock = true
		}

		var isEndOfBlock = false

		if lastIndent > indentLength || strings.HasSuffix(line, "}") {
			isEndOfBlock = true
		}

		if isStartOfBlock {
			fmt.Println("start of block", line)
			blockNumber = 0
			blockLine = 0
		} else if isEndOfBlock {
			lastLineOfParentBlock := codeLines.GetLastLineOfBlock(parentBlockNumber)
			blockNumber = lastLineOfParentBlock.Block
			parentBlockNumber = lastLineOfParentBlock.ParentBlock

			blockLine = lastLineOfParentBlock.Order + 1
		}

		if isEndOfBlock || isStartOfBlock {
			lastIndent = indentLength
		}

		blockNumber = codeLines.AddToList(line, indentLength, blockLine, blockNumber, parentBlockNumber)
		blockLine++

		if isStartOfBlock {
			blockLine = 0
		}

		if strings.HasSuffix(line, "{") {
			parentBlockNumber = blockNumber
		}

		if strings.HasSuffix(line, "}") {
			blockNumber = 0
		}

	}

	codeLines.Print()

}
