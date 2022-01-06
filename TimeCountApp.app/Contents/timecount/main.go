package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"timecount/icon"

	"github.com/getlantern/systray"
	"gocv.io/x/gocv"
)

var (
	msgTemp = "{\"msg_type\": \"text\", \"content\": {\"text\": \"该休息了，已经连续工作 %d 分钟了\"}}"
)

const maxSecond = 3600

type TimeCount struct {
	startTime      int64
	mx             sync.Mutex
	WebhookUrl     string // webhook通知方式
	FaceDataPath   string // 人脸识别地址
	CloseFaceRecog bool
	classifier     gocv.CascadeClassifier
}

func NewTimeCount(facePath string) *TimeCount {
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load(facePath) {
		fmt.Println("load model failed")
		return nil
	}
	return &TimeCount{
		FaceDataPath: facePath,
		classifier:   classifier,
	}
}

func (t *TimeCount) ResetTime() {
	t.mx.Lock()
	t.startTime = time.Now().Unix()
	t.mx.Unlock()
}

func (t *TimeCount) StartTime() int64 {
	return t.startTime
}

func (t *TimeCount) Recognition() (bool, error) {
	webCam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer webCam.Close()

	img := gocv.NewMat()
	defer img.Close()

	ss := time.Now().Unix()
	for {
		// 监测一秒
		if time.Now().Unix()-ss >= 3 {
			break
		}
		if ok := webCam.Read(&img); !ok {
			return false, errors.New("read img from webcam failed")
		}

		if img.Empty() {
			continue
		}

		rects := t.classifier.DetectMultiScale(img)
		fmt.Println("face numbers: ", len(rects))
		if len(rects) > 0 {
			return true, nil
		}

		time.Sleep(time.Millisecond * 200)
	}
	return false, nil
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
	t.classifier.Close()
	time.Sleep(time.Second * 1)
}

var (
	webhookUrl   string
	faceDataPath string
)

func main() {
	flag.StringVar(&webhookUrl, "webhook", "", "webhook通知地址")
	flag.StringVar(&faceDataPath, "recog", "", "人脸识别特征文件地址")
	flag.Parse()
	tc := NewTimeCount(faceDataPath)
	if tc.FaceDataPath == "" {
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
