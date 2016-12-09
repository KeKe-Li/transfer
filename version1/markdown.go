package main

import (
	"os"
	"fmt"
	"sort"
	"io/ioutil"
	"path/filepath"
)


func main() {

	//多个markdown文件转换为一个markdown文件
	strs := make([]string, 0)

	//遍历文件夹并把文件或文件夹名称加入相应的slice
	filepath.Walk("input", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			//println(info.Name())
			strs = append(strs, info.Name())
			sort.Strings(strs)
		}
		return nil
	})

	fmt.Println(strs)
	//fmt.Println("Strings:", strs)


	//OpenFile
	file, _ := os.OpenFile("output/all.md", os.O_APPEND |os.O_RDWR, 0644)
	defer file.Close()


	//range
	for _,value := range strs{
		data ,err :=ioutil.ReadFile("input/" + value)
		if err != nil{
			fmt.Println(err.Error())
		}

		file.Write(data)

	}




}

