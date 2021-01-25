### RTSP2HLS Server

#### Before
This program needs to use [FFmpeg](https://ffmpeg.org/download.html):   to work
The current conversion parameters in the program are fixed:
`ffmpeg -i "rtspurl" -c copy -f hls -hls_time 2.0 -hls_list_size 1 -hls_wrap 5 HlsFileStatic`

#### Run
Run app api in 8888 port(change it in conf app.ini)
```bash
#in windows cmd
rtsp2hle.exe
```

#### Set
Set Camera rtsp info with api      
GET `http://localhost:8888/cameras` all cameras info
POST `http://localhost:8888/cameras` create new camera 
```json
{
    "Name":"test",
    "RtspUrl":"rtsp://admin:n1234567@192.168.8.92:554/h264/ch1/main/av_stream",
    "HlsFileUrl": "test.m3u8",
    "HlsFileStatic":"d:/hls_data/test.m3u8"
}
```

PATCH `http://localhost:8888/camera/id` update info by id

#### START & STOP
GET `http://localhost:8888/camera/start/id` Start a ffmpeg thread to convert the corresponding Rtsp to HLS

GET `http://localhost:8888/camera/stop/id` Close the corresponding thread task

