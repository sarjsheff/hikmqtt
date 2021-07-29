package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nfnt/resize"
	"github.com/sarjsheff/hiklib"
)

type SnapPayload struct {
	Image string
	Ip    string
	Name  string
}

func (cam Cam) width() uint {
	if cam.W != 0 {
		return cam.W
	} else if cam.H != 0 {
		return 0
	} else {
		return 800
	}
}

func (cam Cam) height() uint {
	if cam.H != 0 {
		return cam.H
	} else if cam.W != 0 {
		return 0
	} else {
		return 0
	}
}

func (cam Cam) interval() uint {
	if cam.Interval != 0 {
		return cam.Interval
	} else {
		return 5000
	}
}

func (cam Cam) run(cli *mqtt.Client, wg *sync.WaitGroup) {
	log.Printf("Login to cam %s [%s]\n", cam.Name, cam.Ip)
	u, dev := hiklib.HikLogin(cam.Ip, cam.Username, cam.Password)
	imgpath := "/tmp/" + cam.Ip + "_image.jpg"
	ticker := time.NewTicker(time.Millisecond * time.Duration(cam.interval()))
	for t := range ticker.C {
		log.Println("Tick at", t)
		if hiklib.HikCaptureImage(u, dev.ByStartChan, imgpath) == 1 {
			f, err := os.Open(imgpath)
			if err == nil {
				image, _, err := image.Decode(f)
				if err == nil {
					newImage := resize.Resize(cam.width(), cam.height(), image, resize.Lanczos3)
					var b bytes.Buffer
					w := bufio.NewWriter(&b)
					err = jpeg.Encode(w, newImage, nil)
					if err == nil {
						bt, err := json.Marshal(SnapPayload{
							Image: base64.StdEncoding.EncodeToString(b.Bytes()),
							Name:  cam.Name,
							Ip:    cam.Ip,
						})
						if err == nil {
							<-(*cli).Publish("hikmqtt/"+cam.Name+"/snap", 0, false, bt).Done()
						} else {
							log.Println(err)
						}
					} else {
						log.Println(err)
					}
				}
				f.Close()
			}
		} else {
			log.Println("Error capture image")
		}
	}
	wg.Done()
}
