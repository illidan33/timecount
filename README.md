# 系统久坐提醒器
- 每隔一个小时提醒站立或休息，如果人脸特征识别功能是打开的，便识别计算机前的人脸，如果存在人脸，便会每隔一分钟提醒一次，直到没有检测到人脸特征或者手动重置时间。由于使用的opencv，故人脸特征识别率不是特别高，用于粗略识别，还是够用了。
- 由于使用的systray，理论上可以跨平台编译执行。

### 配置app服务

```shell
# 如果需要使用人脸特征识别，需要安装opencv
## ubuntu
sudo apt install opencv
## mac os
brew install opencv

# 编译服务
cd ./Contents/timecount

go build main.go

# 执行即可，webhook可填webhook通知地址（如钉钉、飞书、企业微信等webhook地址），recog可填人脸识别特征文件绝对地址，项目自带TimeCountApp.app/Contents/timecount/data/haarcascade_frontalface_default.xml，可放任意位置；
./timecount -webhook xxxx -recog /data/haarcascade_frontalface_default.xml
```