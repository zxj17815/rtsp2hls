package models

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

type Camera struct {
	ID            int
	Name          string
	RtspUrl       string
	HlsFileStatic string
	HlsFileUrl    string
	Pid           int
	State         int
	CreateTime    time.Time
	UpdateTime    time.Time
}

func GetCameras() (cameras []Camera, count int) {
	// 获取全部记录
	result := DB.Find(&cameras)
	// SELECT * FROM users;

	count = int(result.RowsAffected) // 返回找到的记录数，相当于 `len(users)`
	if result.Error == nil {
		return cameras, 0
	} // returns error
	return
}

func GetCamera(id string) (cameras Camera, err error) {
	result := DB.First(&cameras, id)
	if result.Error != nil {
		return cameras, result.Error
	} // returns error
	return
}

func (c *Camera) CreateCamera() (created bool, err error) {
	c.CreateTime = time.Now()
	c.UpdateTime = time.Now()
	result := DB.Create(&c)
	if result.Error == nil {
		created = true
		err = nil
	} else {
		created = false
		err = result.Error
	}
	return created, err
}

func (c *Camera) UpdateCamera(data map[string]interface{}) (updated bool, err error) {
	data["UpdateTime"] = time.Now()
	delete(data, "CreateTime")
	delete(data, "Pid")
	delete(data, "ID")
	result := DB.Model(c).Updates(data)
	if result.Error == nil {
		updated = true
		err = nil
	} else {
		updated = false
		err = result.Error
	}
	return updated, err
}

func (c *Camera) OpenCamera() (status bool, err error) {
	if _, err = HlsStart(c); err != nil {
		status = false
	} else {
		status = true
	}
	return status, err
}

func (c *Camera) CloseCamera() (status bool, err error) {
	if p, error := os.FindProcess(c.Pid); error != nil {
		err = error
		if err.Error() == "OpenProcess: The parameter is incorrect." {
			DB.Model(&c).Select("Pid", "State").Updates(Camera{
				Pid:   0,
				State: 0,
			})
		}
	} else {
		p.Kill()
		p.Release()
		DB.Model(&c).Select("Pid", "State").Updates(Camera{
			Pid:   0,
			State: 0,
		})
		status = true
	}
	return status, err
}

func KeepAlive(cmd *exec.Cmd, camera *Camera) {
	err := cmd.Wait()
	if err != nil {
		time.Sleep(5 * time.Second)
		DB.First(&camera, camera.ID)
		if camera.State != 0 {
			HlsStart(camera)
		}
	}
}

func HlsStart(camera *Camera) (pid int, err error) {
	staticFile := camera.HlsFileStatic

	cmd := exec.Command("ffmpeg", "-i", camera.RtspUrl, "-c", "copy", "-f", "hls", "-hls_time", "2.0", "-hls_list_size", "1", "-hls_wrap", "5", staticFile)

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	err = cmd.Start()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	} else {
		pid = cmd.Process.Pid
		DB.Model(camera).Select("Pid", "State").Updates(Camera{
			Pid:   pid,
			State: 1,
		})
		go KeepAlive(cmd, camera)
	}
	return pid, err
}
