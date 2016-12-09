package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/russross/blackfriday"
)

const (
	htmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_SAFELINK |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_COMPLETE_PAGE |
		blackfriday.HTML_HREF_TARGET_BLANK |
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

type Config struct {
	GenerateMd     string
	GeneratePdf    string
	GenerateHtml   string
	OriginPath     string
	ReadFileSuffix string
}


//newDefaultConfig
func newDefaultConfig() *Config {
	return &Config{
		ReadFileSuffix: "md",
	}
}



func main() {
	config := newDefaultConfig()

	//init config
	err := initConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//generate md
	err = generateMd(config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//generate html
	err = generateHtml(config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//generate pdf
	err = html2Pdf(config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}


//initConfig
func initConfig(config *Config) error {
	var flagOut, flagInput string

	// 设置flag参数 (变量指针，参数名，默认值，帮助信息)
	flag.StringVar(&flagInput, "input", "", "origin path is not exist")
	flag.StringVar(&flagOut, "output", "", "output path is invalid")

	flag.Usage = func() {
		fmt.Println("Usage :" + os.Args[0] + "[options]\n\nOptions:")
		fmt.Println("the output path.")
	}

	// 解析
	flag.Parse()

	err := os.MkdirAll(flagOut, 0755)
	if err != nil {
		return err
	}

	config.OriginPath = flagInput
	config.GenerateMd = flagOut + "output.md"
	config.GenerateHtml = flagOut + "output.html"
	config.GeneratePdf = flagOut + "output.pdf"

	return nil
}


//htmlToPdf
func html2Pdf(config *Config) error {
	err := exec.Command("wkhtmltopdf", config.GenerateHtml, config.GeneratePdf).Run()
	if err != nil {
		return err
	}
	return nil
}


//generateHtml
func generateHtml(config *Config) error {
	data, err := ioutil.ReadFile(config.GenerateMd)
	if err != nil {
		return err
	}

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "../output/github.css")

	outputOpts := blackfriday.Options{
		Extensions: extensions,
	}

	htmlBody := blackfriday.MarkdownOptions(data, renderer, outputOpts)
	return ioutil.WriteFile(config.GenerateHtml, htmlBody, 0644)
}



//generateMd
func generateMd(config *Config) error {
	fileInfo, err := ioutil.ReadDir(config.OriginPath)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(config.GenerateMd, os.O_APPEND | os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, value := range fileInfo {

		if checkExtension(value, config.ReadFileSuffix) {
			data, _ := ioutil.ReadFile(config.OriginPath + value.Name())

			file.Write(data)
		}
	}

	return nil
}



//checkExtension
func checkExtension(fileInfo os.FileInfo, ext string) bool {
	if fileInfo.IsDir() || !strings.HasSuffix(fileInfo.Name(), ext) {
		return false
	}
	return true
}
