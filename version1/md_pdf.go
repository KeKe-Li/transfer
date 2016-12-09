package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"io/ioutil"
	"path/filepath"
	"log"
	"os/exec"
	"os/user"

	"github.com/russross/blackfriday"
	//"github.com/pleximus/go-htmltopdf"
)

const (
	htmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_COMPLETE_PAGE |
		blackfriday.HTML_FOOTNOTE_RETURN_LINKS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	extensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_HARD_LINE_BREAK
)

func main() {
	/*
	  FIXME
	  usage

	*/

	//解析输入和输出地址
	pathPtr := flag.String("input", "input/read.md", "需要转换的输入markdown文件路径")

	outputPtr := flag.String("output", "output/read1.html", "输出html路径")

	// 解析
	flag.Parse()

	//FIXME
	//判断输入错误
	if *pathPtr == "" {
		log.Fatal(errors.New("error:出错了你必须提供你的markdown文件路径"))
	}else{
		log.Fatal(pathPtr)

	}


	//请检查输入的路径和清除文件路径
	path := cleanPath(pathPtr)

	outputPath := outputFilePath(outputPtr, path)

	rawMarkdown := readInput(path)

	fmt.Println(string(rawMarkdown))

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	outputOpts := blackfriday.Options{
		Extensions:extensions,
	}


	htmlBody := blackfriday.MarkdownOptions(rawMarkdown, renderer, outputOpts, )

	//打印到标准输出中进度消息
	fmt.Printf("转换 %s. 输出位于%s\n", filepath.Base(path), outputPath)

	//0644 是标准权限配置，对 Apache 的 html 文件
	ioutil.WriteFile(outputPath, htmlBody, 0644)


	////读取html生成的pdf
	//converter := html2pdf.New()  // 创建一个新的 html2pdf对象
	//
	////ReadFile
	//file, err := ioutil.ReadFile("output/read.html")
	//if err != nil {
	//	fmt.Println(err)
	//}
	////fmt.Println(string(file))
	//
	//
	//converter.SetData(string(file))

	//generatePdf
	result := exec.Command("wkhtmltopdf", "output/read.html", "output/read.pdf").Run()
	if result == nil {

		fmt.Printf("succesfully generated")
	} else {
		fmt.Printf("Error on trying to generate PDF: %s", result)
	}

	//converter.Destroy()

}

//FIXME
func cleanPath(pathPtr *string) (cleaned string,) {
	//确保文件是一个markdown文件
	if filepath.Ext(*pathPtr) != ".md" {
		log.Fatal(errors.New("出错:你必须用一个markdown(.md)文件"))
	}


	//unix的波浪符
	if string((*pathPtr)[0]) == "~" {
		expandTilde(pathPtr)
	}

	cleaned = filepath.Clean(*pathPtr)

	return

}

func expandTilde(pathPtr *string) {

	usr, _ := user.Current()
	*pathPtr = strings.Replace(*pathPtr, "~", usr.HomeDir, 1)
}


func readInput(filename string) []byte {

	rawMarkdown, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return rawMarkdown
}

func outputFilePath(outputPtr *string, inputPath string) (output string) {
	if *outputPtr == "" {
		inputDir, inputFile := filepath.Split(inputPath)
		outputFile := strings.TrimSuffix(inputFile, ".md") + ".html"
		output = filepath.Join(inputDir, outputFile)
	} else {
		output = *outputPtr
	}


	output, _ = filepath.Abs(output)

	return
}





