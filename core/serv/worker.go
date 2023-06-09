package serv

import (
	"fmt"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/util"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

type Task struct {
	w string     //工作目录
	d *gorm.DB   //数据库链接
	v data.Video //视频
}

func NewTask(workspace string, db *gorm.DB, video data.Video) *Task {
	return &Task{workspace, db, video}
}

func (my *Task) process(worker string) {
	var wg sync.WaitGroup
	wg.Add(2)

	id := strconv.Itoa(int(my.v.Id))
	timestamp := strconv.FormatInt(my.v.UploadAt.UnixNano()/1e6, 10)

	//标题文件
	log.Println("生成标题文件:" + my.v.Vid)
	titleFile := path.Join(my.w, id, fmt.Sprintf("t0-%s.txt", my.v.Vid))
	err := util.WriteFile(strings.NewReader(my.v.Title), titleFile)
	if err != nil {
		return
	}

	//索引文件
	log.Println("生成索引文件:" + my.v.Vid)
	txtFile := path.Join(my.w, id, fmt.Sprintf("%s.txt", id))
	filepath := fmt.Sprintf("daren/2215630453359/zip/%s___%s.zip", timestamp, my.v.Vid)
	content := []string{my.v.Aid, filepath, timestamp, my.v.Vid, my.v.Fid}
	err = util.WriteFile(strings.NewReader(strings.Join(content, "\n")), txtFile)
	if err != nil {
		return
	}

	//下载封面
	log.Println("开始下载封面:" + my.v.Vid)
	coverFile := path.Join(my.w, id, fmt.Sprintf("c1-%s.jpg", my.v.Vid))
	go func() {
		defer wg.Done()
		err = util.DownloadFile(my.v.Cover, coverFile)
		if err != nil {
			panic(err)
			return
		}
	}()

	//下载视频
	log.Println("开始下载视频:" + my.v.Vid)
	videoFile := path.Join(my.w, id, fmt.Sprintf("v2-%s.mp4", my.v.Vid))
	go func() {
		defer wg.Done()
		err = util.DownloadFile(my.v.Url, videoFile)
		if err != nil {
			panic(err)
			return
		}
	}()

	wg.Wait()

	//打包文件
	log.Println("开始打包文件:" + my.v.Vid)
	zipFile := path.Join(my.w, id, fmt.Sprintf("%s.zip", my.v.Vid))
	err = util.Compress(zipFile, titleFile, videoFile, coverFile)
	if err != nil {
		return
	}

	//删除文件
	_ = os.Remove(titleFile)
	_ = os.Remove(coverFile)
	_ = os.Remove(videoFile)

	//上传文件
	log.Println("开始上传文件:" + my.v.Vid)
	err = util.UploadFile(zipFile, fmt.Sprintf("%s___", timestamp))
	if err != nil {
		return
	}
	//上传索引
	log.Println("开始上传索引:" + my.v.Vid)
	err = util.UploadFile(txtFile, "index/")
	if err != nil {
		return
	}

	//更新状态
	log.Println("开始更新状态:" + my.v.Vid)
	my.d.Model(my.v).Update("state", 1)
}
