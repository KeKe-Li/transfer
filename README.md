

## markdown 转换 pdf

```
使用markdown转换Pdf的时候 需要下载 wkhtmltopdf
input是放置markdown的目录
output是输出pdf的目录
images是markdown文件中的图片的目录
```


```
cd version2

```

```
go build

```

```
 ./version2 -input=../input/ -output=../output/

  wkhtmltopdf output.html output.pdf

```
