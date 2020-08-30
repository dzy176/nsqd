package simple

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"testing"
)

func TestA(t *testing.T) {
	// 根据Hostname随机生成了一个数字
	hostname, _ := os.Hostname()
	h := md5.New()
	io.WriteString(h, hostname)
	defaultID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024)
	fmt.Println(defaultID)
}

func TestA1(t *testing.T) {
	//定义零值Buffer类型变量b
	var b bytes.Buffer
	//使用Write方法为写入字符串
	b.Write([]byte("你好"))
	//这个是把一个字符串拼接到Buffer里
	fmt.Fprint(&b, ",", "http://www.flysnow.org")
	//把Buffer里的内容打印到终端控制台
	b.WriteTo(os.Stdout)
}
