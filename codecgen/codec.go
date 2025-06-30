package codecgen

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run() {
	tablename := flag.String("t", "", "会根据这个表明生成对应的cache文件")
	s := flag.Int64("s", 0, "cache 的缓存过期时间，单位s")
	h := flag.Int64("h", 0, "cache 的缓存过期时间，单位小时,默认3")
	d := flag.String("d", "passive", "redis的那个模块的db,按业务区分。目前提供 story,property,block,user...")
	primaryAlisaName := flag.String("uq", "id", "默认id，但你的表如果唯一索引锁uid，这里你就可以用uid")
	mode := flag.String("m", "slp", "给个项目的go.mod的包名")

	if *tablename == "" {
		fmt.Println("必须输入-t参数，db表名的意思")
		return
	}

	if *s <= 0 && *h <= 0 {
		fmt.Println("必须输入-s或者-d参数，优先级-s > -h，指定过期时间")
		return
	}
	seconds := int64(60 * 60 * 3)
	if *s > 0 {
		seconds = int64(*s)
	} else if *h > 0 {
		seconds = int64(*h * 60 * 60)
	} else {
		fmt.Println("无法解析到过期时间参数")
		return
	}
	tableName := *tablename
	pbName := FirstUppers(tableName)

	redisDb := fmt.Sprintf("%s%s", "Redis", FirstUppers(*d))
	ttt := strings.Replace(templateV2, "%pbName", pbName, 1000)
	ttt = strings.Replace(ttt, "%lowerName", FirstLower(pbName), 1000)
	ttt = strings.Replace(ttt, "%mode", *mode, 1000)
	ttt = strings.Replace(ttt, "%tableName", pbName+"TableKey", 1000)
	ttt = strings.Replace(ttt, "%tableKey", tableName, 1000)
	ttt = strings.Replace(ttt, "%ttl", fmt.Sprintf("%d", seconds), 1000)
	ttt = strings.Replace(ttt, "%redis-db", redisDb, 100)
	ttt = strings.Replace(ttt, "%entity_id", FirstUpper(*primaryAlisaName), 100)
	ttt = strings.Replace(ttt, "%id", strings.ToLower(*primaryAlisaName), 100)
	generate(tableName, ttt)
}

// FirstUpper 字符串首字母大写

func FirstUppers(s string) string {
	rStr := ""
	strs := strings.Split(s, "_")
	for _, str := range strs {
		str = FirstUpper(str)
		rStr += str
	}
	return rStr
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}

	return strings.ToUpper(s[:1]) + s[1:]

}
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]

}

func generate(filename string, context string) {
	path := "codec"
	//创建一个新文件，写入内容 5 句 "http://c.biancheng.net/golang/"
	filePath := fmt.Sprintf("./rpc/server/internal/cache/%s/%s_codec.go", path, filename)
	b, err := PathExists(filePath)
	if err != nil {
		panic(err)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("生成文件的路径:%s", absPath))

	if b {
		if err = os.Remove(filePath); err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}

	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(context)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()

	msg, err := RunCommand("./", "gofmt", "-l", "-w", absPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("success. 缓存code路径：%s  msg=%s", filePath, msg))

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RunCommand(path, name string, arg ...string) (msg string, err error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = path
	err = cmd.Run()
	log.Println(cmd.Args)
	if err != nil {
		log.Println("err", err.Error(), "cmd", cmd.Args)
	}
	return
}
