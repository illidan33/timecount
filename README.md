# mac os系统久坐提醒器

### 配置人脸识别服务

- 人脸识别用于检测是否还在电脑前工作，如果还在则发送提醒，不在则重置时间;
- 原始图片放recognition，修改main.py中对应路径;

```shell
cd ./Contents/recognition
# 激活python虚拟环境
source ./venv/bin/activate

# 启动python服务
python3 main.py
```

### 配置app服务

```shell
cd ./Contents/timecount

go build main.go

# 执行即可，xxxx可填webhook通知地址（如钉钉、飞书、企业微信等webhook地址），yyyy可填上面配置的python地址；
./timecount -webhook xxxx -recog yyyy
```