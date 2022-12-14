package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jason/instrument_trace/instrumenter"
	"github.com/jason/instrument_trace/instrumenter/ast"
)

var (
	wrote bool
)

func init() {
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
}

func usage() {
	fmt.Println("instrument [-w] xxx.go")
	flag.PrintDefaults()
}

func main() {
	fmt.Println(os.Args)
	flag.Usage = usage
	flag.Parse() // 解析命令行参数

	if len(os.Args) < 2 {
		usage()
		return
	}

	var file string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}

	if len(os.Args) == 2 {
		file = os.Args[1]
	}

	if filepath.Ext(file) != ".go" { // 对源文件扩展名进行校验
		usage()
		return
	}

	var ins instrumenter.Instrumenter // 声明 instrumenter.Instrumenter 接口类型变量

	// 创建以 ast 方式实现 Instrumenter 接口的 ast.instrumenter 实例
	ins = ast.New("github.com/jason/instrument_trace", "trace", "Trace")
	newSrc, err := ins.Instrument(file)
	if err != nil {
		panic(err)
	}

	if newSrc == nil {
		// add nothing to the source file. no change
		fmt.Printf("no trace added for %s\n", file)
		return
	}

	if !wrote {
		fmt.Println(string(newSrc))
		return
	}

	// 将生成的新代码的内容写回原 Go 源文件
	if err = ioutil.WriteFile(file, newSrc, 0666); err != nil {
		fmt.Printf("write %s error: %v\n", file, err)
		return
	}

	fmt.Printf("instrument trace for %s ok\n", file)
}
