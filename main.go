package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func determinEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("fetch error:%v", err)

		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e

}

func parseContent(content []byte) {
	//<a href="/tag/管理" class="tag">管理</a>
	re := regexp.MustCompile(`<a href="([^"]+)" class="tag">([^"]+)</a>`)
	match := re.FindAllSubmatch(content, -1)
	for _, m := range match {
		fmt.Printf("m[0]:%s m[1]:%s m[2]:%s\n", m[0], m[1], m[2])
	}

}

func main() {
	client := &http.Client{}
	request, _ := http.NewRequest(`GET`, `https://book.douban.com/`, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error statis code: %d", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determinEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	result, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(nil)
	}
	parseContent(result)

}
