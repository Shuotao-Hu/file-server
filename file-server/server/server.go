//实现展示文件，上传文件，下载文件的功能,支持并发、断点续传

package server

import (
	"fmt"
	//"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

//var err = yaml.Unmarshal()

const path = "file/files/"
const pathBre = "file/bre_pos/"

var ServeHandle = http.FileServer(http.Dir(path))

func Download(w http.ResponseWriter, r *http.Request) {
	if r.Header["Range"] == nil{
		downloadNormal(w, r)
	}else{
		downloadContinue(w, r)
	}
}

func downloadNormal(w http.ResponseWriter, r *http.Request){
	fileName := filepath.Base(r.URL.Path) // filename
	pathName := path+fileName
	file, err := os.OpenFile(pathName,os.O_RDWR,0666)
	if err != nil{
		w.Write([]byte("open file failed. "))
		return
	}
	defer file.Close()
	fileHeader := make([]byte,512)
	fileStat, _ := file.Stat()
	w.Header().Set("Content-Disposition", "attachment; filename=" + fileName)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(),10))
	_, err = io.Copy(w, file)
	if err != nil{
		w.Write([]byte("copy file failed. "))
		return
	}
	return
}

func downloadContinue(w http.ResponseWriter, r *http.Request){
	fileName := filepath.Base(r.URL.Path)
	rangeData := r.Header["Range"]
	re := regexp.MustCompile("[0-9]+")
	rangeStr := re.FindAllString(rangeData[0], -1)
	rangeLow, _ := strconv.Atoi(rangeStr[0])
	rangeHigh, _ := strconv.Atoi(rangeStr[1])
	rangeSize := rangeHigh - rangeLow + 1

	pathName := path+fileName
	file, err := os.OpenFile(pathName,os.O_RDWR,0666)
	if err != nil{
		w.Write([]byte("open file failed. "))
		return
	}
	defer file.Close()
	fileHeader := make([]byte,512)    //设置响应头部信息
	fileStat, _ := file.Stat()
	w.Header().Set("Content-Disposition", "attachment; filename=" + fileName)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Range", "bytes "+rangeStr[0]+"-"+rangeStr[1]+"/"+strconv.FormatInt(fileStat.Size(),10))
	file.Seek(int64(rangeLow), 0)
	fileStream := make([]byte, 512)
	sizeTemp := rangeSize
	for {
		if sizeTemp <= 512{
			_, err := file.Read(fileStream[:sizeTemp])
			if err != nil{
				w.Write([]byte("read file failed. "))
				return
			}
			_,err = w.Write(fileStream)
			if err != nil {
				w.Write([]byte("write file failed. "))
				return
			}
			break
		}else{
			_, err := file.Read(fileStream)
			if err != nil{
				w.Write([]byte("read file failed. "))
				return
			}
			sizeTemp -= 512
			_, err = w.Write(fileStream)
			if err != nil {
				w.Write([]byte("write file failed. "))
				return
			}
		}
	}
	return
}


func Upload(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := r.MultipartForm //记录和解析申请包内容
		files := m.File["filename"] //filename,上传的文件
		for i, _ := range files {
			file, _ := files[i].Open() //打开file
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()
			f := path + files[i].Filename              //在服务器./files/下创建同名文件
			dst, _ := os.OpenFile(f, os.O_CREATE, 0666) //打开dst
			defer dst.Close()
			var pos int = 0                               // 读取断点位置
			b := pathBre + files[i].Filename + "_bre"
			fileBre, _ := os.OpenFile(b,os.O_CREATE,0666)
			posS := make([]byte, 128)
			posN, err := fileBre.Read(posS)
			if err != nil{
				fmt.Println(err)
				return
			}
			pos, _ = strconv.Atoi(string(posS[:posN]))
			fileBre.Close()
			file.Seek(int64(pos), io.SeekStart)      //设置光标位置
			dst.Seek(int64(pos), io.SeekStart)
			fmt.Printf("当前读取位置为：%d\n", pos)//    注意字节数和读取字符不同

			var buf []byte = make([]byte, 128) // 复制文件内容

			for {
				count, err := file.Read(buf)
				if err == io.EOF {
					fmt.Println("读取完成")
					break
				}
				dst.Write(buf[:count])
				pos += count
				fileBre, err := os.OpenFile(b,os.O_RDWR,0766)
				if err != nil{
					fmt.Println(err)
				}
				fileBre.Write([]byte(fmt.Sprintf("%d", pos)))
				fileBre.Close()
				fmt.Printf("已上传字节数：%d\n", pos)

			}

		}
		w.Write([]byte("file upload success. "))

}

