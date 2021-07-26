package main

import (
	"log"
	"net/url"
	"os"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/sarjsheff/hiklib"
)

func main() {
	if config() < 0 {
		return
	}

	log.Println(hiklib.HikVersion())
	log.Println(cfg)

	run()
}

func run() {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	ur, err := url.Parse(cfg.Url)
	if err != nil {
		log.Fatalln(err)
	}
	opt := mqtt.NewClientOptions()
	opt.ClientID = "hikmqtt" + uuid.Must(uuid.NewRandom()).String()
	opt.Servers = append(opt.Servers, ur)
	if cfg.Username != "" {
		opt.Username = cfg.Username
	}
	if cfg.Password != "" {
		opt.Password = cfg.Password
	}

	cli := mqtt.NewClient(opt)
	c := cli.Connect()
	<-c.Done()
	if c.Error() != nil {
		log.Fatalln(c.Error().Error())
	}
	// log.Println(t)
	// s := cli.Subscribe("#", 0, func(c mqtt.Client, m mqtt.Message) {
	// 	log.Println(m.Topic())
	// })
	// <-s.Done()
	// if c.Error() != nil {
	// 	log.Fatalln(c.Error().Error())
	// }

	log.Println("Connect to cameras.")
	var wg sync.WaitGroup

	for _, c := range cfg.Cams {
		wg.Add(1)
		go c.run(&cli, &wg)
	}
	wg.Wait()
	// u, dev := hiklib.HikLogin(*camipFlag, *camuserFlag, *campassFlag)
	// imgpath := "/tmp/" + *camipFlag + "image.jpg"
	// ticker := time.NewTicker(time.Millisecond * 5000)
	// for t := range ticker.C {
	// 	log.Println("Tick at", t)
	// 	if hiklib.HikCaptureImage(u, dev.ByStartChan, imgpath) == 1 {
	// 		f, err := os.Open(imgpath)
	// 		if err == nil {
	// 			image, _, err := image.Decode(f)
	// 			if err == nil {
	// 				newImage := resize.Resize(800, 0, image, resize.Lanczos3)
	// 				var b bytes.Buffer
	// 				w := bufio.NewWriter(&b)
	// 				//w, _ := os.Create("/app/image.jpg")
	// 				err = jpeg.Encode(w, newImage, nil)
	// 				//w.Close()
	// 				log.Println(base64.StdEncoding.EncodeToString(b.Bytes()))
	// 				<-cli.Publish("camsnap/"+*camipFlag, 0, true, base64.StdEncoding.EncodeToString(b.Bytes())).Done()
	// 			}
	// 			f.Close()
	// 		}
	// 	} else {
	// 		log.Fatalln("Error capture image")
	// 	}
	// }
}
