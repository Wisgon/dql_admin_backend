// 这个脚本测试dgraph当大量高并发时，用的upsert可不可以防止唯一字段user_name注册了多个
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	// 测试同时post相同名字和电话，测试结果：第一个会成功，其他两个失败，报错是:Transaction has been aborted. Please retry
	// for i := 0; i < 3; i++ {
	// 	go postConcurrent("xxx", "8888")
	// }

	// 测试不同名字不同电话同时， 测试结果：三个都成功注册
	go postConcurrent("yyy", "8787")
	go postConcurrent("iii", "7767")
	go postConcurrent("admin", "123456")

	time.Sleep(10 * time.Second)
}

func postConcurrent(name string, phone string) {
	url := "http://localhost:8063/user/regist"
	JSONStr := "{\"username\":\"" + name + "\", \"password\":\"123456\", \"phone\":\"" + phone + "\"}" // 里面字段必须是斜扛双引号，不能是单引号
	JSONB := []byte(JSONStr)

	resp, err := http.Post(url,
		"application/json",
		bytes.NewBuffer(JSONB))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
