package main

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	msgTempFeiShu        = "{\"msg_type\": \"text\", \"content\": {\"text\": \"该休息了，已经连续工作 %d 分钟了\"}}"
	msgTempCompanyWechat = "{\"msgtype\": \"text\", \"text\": {\"content\": \"该休息了，已经连续工作 %d 分钟了\"}}"
	msgTempDingTalk      = "{\"msgtype\": \"text\", \"text\": {\"content\": \"该休息了，已经连续工作 %d 分钟了\"}}"
)

var (
	webhookUrl   string
	faceDataPath string
)

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func main() {
	flag.StringVar(&webhookUrl, "webhook", "", "webhook通知地址")
	flag.StringVar(&faceDataPath, "recog", "", "人脸识别特征文件地址")
	flag.Parse()

	if webhookUrl == "" {
		fmt.Println("Notify: webhook通知地址为空")
	}
	if faceDataPath == "" {
		faceDataPath = "./haarcascade_frontalface_default.xml"
		if !isExists(faceDataPath) {
			faceDataPath = ""
		}
	}
	tc := NewTimeCount(faceDataPath, webhookUrl)
	if tc.FaceDataPath == "" {
		tc.CloseFaceRecog = true
		fmt.Println("Notify: 人脸识别特征文件地址为空")
	}
	if strings.Contains(webhookUrl, "open.feishu.cn") {
		tc.WebhookTemp = msgTempFeiShu
	} else if strings.Contains(webhookUrl, "weixin.qq.com") {
		tc.WebhookTemp = msgTempCompanyWechat
	} else {
		tc.WebhookTemp = msgTempDingTalk
	}

	go func() {
		for {
			if tc.StartTime() == 0 {
				tc.ResetTime()
			}
			now := time.Now()
			less := now.Unix() - tc.StartTime()
			msg := fmt.Sprintf("持续 %d分%d秒", less/60, less%60)
			systray.SetTitle(msg)
			// 一分钟发一次
			if less >= maxSecond {
				// 休息时间忽略
				h := now.Hour()
				if h < 9 || (h >= 12 && h <= 13) || h >= 19 {
					tc.ResetTime()
					continue
				}
				if tc.WebhookUrl != "" && less%60 == 0 {
					newMsg := fmt.Sprintf(tc.WebhookTemp, less/60)
					http.Post(tc.WebhookUrl, "application/json", strings.NewReader(newMsg))
				}
				if tc.FaceDataPath != "" && tc.CloseFaceRecog == false && less%60 == 0 {
					result, err := tc.Recognition()
					if err != nil {
						fmt.Println(err)
						continue
					}
					if !result {
						tc.ResetTime()
					}
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	systray.Run(tc.OnReady, tc.OnExit)
}
