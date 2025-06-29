package main

import (
	"flag"
	"fmt"
	"github.com/olaola-chat/slpctl/codecgen"
	"github.com/olaola-chat/slpctl/codegen"
	"log"
	"os"
)

func main() {
	opType := flag.String("op", "state", "默认游戏状态机模板生成")
	switch *opType {
	case "state":
		stateExec()
	case "codec":
		codecgen.CodecExec()
	}
	flag.Parse()
}

// 状态机
func stateExec() {
	jsonFolder := flag.String("p", "./rpc/server/internal/room_game/state/json", "游戏状态机的默认目录")
	jsonFile := flag.String("f", "", "游戏状态机的默认配置文件名称")
	outputDir := flag.String("o", "./rpc/server/internal/room_game", "输出目录")

	if *jsonFile == "" {
		flag.Usage()
		fmt.Printf("-f 用户指定配置json的文件名")
		os.Exit(1)
	}

	jPath := fmt.Sprintf("%s/%s", *jsonFolder, *jsonFile)
	generator, nErr := codegen.NewGameGenerator(jPath, *outputDir)
	if nErr != nil {
		log.Fatalf("生成失败了: %v", nErr)
		return
	}
	if err := generator.Generate(); err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("游戏代码已成功生成到目录: %s\n", *outputDir)
}
