package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"timecount/icon"

	"github.com/getlantern/systray"
)

var (
	msgTemp = "{\"msg_type\": \"text\", \"content\": {\"text\": \"该休息了，已经连续工作 %d 分钟了\"}}"
	tc      = TimeCount{}
)

const maxSecond = 3600

type TimeCount struct {
	startTime          int64
	mx                 sync.Mutex
	WebhookUrl         string // webhook通知方式
	FaceRecognitionUrl string // 人脸识别地址
	CloseFaceRecog     bool
}

func (t *TimeCount) ResetTime() {
	t.mx.Lock()
	t.startTime = time.Now().Unix()
	t.mx.Unlock()
}

func (t *TimeCount) StartTime() int64 {
	return t.startTime
}

func (t *TimeCount) OnReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("初始化...")
	systray.SetTooltip("总时间计算")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	resetBtn := systray.AddMenuItem("Reset Time", "reset the count time")
	closeBtn := systray.AddMenuItem("Close Recognition", "关闭人脸识别功能")
	go func() {
		for {
			select {
			case <-mQuitOrig.ClickedCh:
				fmt.Println("Requesting quit")
				systray.Quit()
				fmt.Println("Finished quitting")
			case <-resetBtn.ClickedCh:
				t.ResetTime()
				systray.SetTitle("持续工作 0分0秒")
			case <-closeBtn.ClickedCh:
				if t.CloseFaceRecog {
					t.CloseFaceRecog = false
					closeBtn.SetTitle("Close Recognition")
				} else {
					t.CloseFaceRecog = true
					closeBtn.SetTitle("Open Recognition")
				}
			}
		}
	}()
}

func (t *TimeCount) OnExit() {
	fmt.Println("exit")
	time.Sleep(time.Second * 1)
}

func main() {
	flag.StringVar(&tc.WebhookUrl, "webhook", "", "webhook通知地址")
	flag.StringVar(&tc.FaceRecognitionUrl, "recog", "", "人脸识别服务器地址")
	flag.Parse()
	if tc.FaceRecognitionUrl == "" {
		tc.CloseFaceRecog = true
	}

	go func() {
		for {
			if tc.StartTime() == 0 {
				tc.ResetTime()
			}
			now := time.Now()
			less := now.Unix() - tc.StartTime()
			msg := fmt.Sprintf("持续工作 %d分%d秒", less/60, less%60)
			systray.SetTitle(msg)
			// 一分钟发一次
			if less >= maxSecond {
				if tc.WebhookUrl != "" && less%60 == 0 {
					newMsg := fmt.Sprintf(msgTemp, less/60)
					http.Post(tc.WebhookUrl, "application/json", strings.NewReader(newMsg))
				}
				if tc.FaceRecognitionUrl != "" && tc.CloseFaceRecog == false {
					resp, _ := http.Get(tc.FaceRecognitionUrl)
					if resp.Body == nil {
						continue
					}
					rs, _ := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					if string(rs) == "ok" {
						tc.ResetTime()
						fmt.Println("reset startTime")
					}
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	systray.Run(tc.OnReady, tc.OnExit)
}
