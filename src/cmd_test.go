package main

import (
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

func CheckPid(pid int, ch chan interface{}) {
	if p, err := os.FindProcess(pid); err != nil {
		ch <- err
	} else {
		ch <- p.Pid
		for i := 0; i < 5; i++ {
			time.Sleep(6 * time.Second)
			ch <- p
		}
		killError := p.Kill()
		ch <- "kill"
		ch <- killError
	}
	if p2, e := os.FindProcess(pid); e != nil {
		ch <- e
	} else {
		ch <- p2
	}
}

func KeepAlive(cmd *exec.Cmd, cid int, ch chan interface{}) {
	err := cmd.Wait()
	if err != nil {
		ch <- cid
	}
}

func TestName(t *testing.T) {
	ch := make(chan interface{})
	//ffmpeg -i "rtsp://admin:Zyadmin888@192.168.8.2:554/h264/ch1/main/av_stream" -c copy -f hls -hls_time 2.0 -hls_list_size 1 -hls_wrap 5 test.m3u8
	cmd := exec.Command("ffmpeg", "-i", "rtsp://admin:Zyadmin888@192.168.8.2:554/h264/ch1/main/av_stream", "-c", "copy", "-f", "hls", "-hls_time", "2.0", "-hls_list_size", "1", "-hls_wrap", "5", "D:\\hls\\test.m3u8")
	err := cmd.Start()
	//pid := cmd.Process.Pid
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	} else {
		go KeepAlive(cmd, 1, ch)
	}
	for {
		t.Log(<-ch)
	}
}

func TestPid(t *testing.T) {
	newP, err := os.FindProcess(13784)
	if err != nil {
		t.Logf("pid err: %s\n", err)
	} else {
		t.Log(newP)
	}
}
