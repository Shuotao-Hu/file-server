//客户端，实现本地文件上传，服务器文件下载的功能
package server

import (
	"bytes"
"errors"
"fmt"
"io/ioutil"
"mime/multipart"
"net/http"
)

func main() {
	createMemoToAmazon()
}

func createMemoToAmazon() error {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	//err := w.WriteField("file", "我的世界！！！")
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("resp status:" + fmt.Sprint(resp.StatusCode))
	}

	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fw, err := w.CreateFormFile("uploadfile", "baidu.com")      //fw
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = fw.Write(bin)
	if err != nil {
		fmt.Println(err)
		return err
	}
	w.Close()

	req, err := http.NewRequest("POST", "http://localhost:8888/upload", buf)
	if err != nil {
		fmt.Println("req err: ", err)
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("resp err: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("resp status:" + fmt.Sprint(resp.StatusCode))
	}

	return nil
}
