# mac os系统久坐提醒器

### 配置app服务

```shell
cd ./Contents/timecount

go build main.go

# 执行即可，webhook可填webhook通知地址（如钉钉、飞书、企业微信等webhook地址），recog可填人脸识别特征文件绝对地址，项目自带TimeCountApp.app/Contents/timecount/data/haarcascade_frontalface_default.xml，可放任意位置；
./timecount -webhook xxxx -recog /data/haarcascade_frontalface_default.xml
```