package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Config struct {
	Url      string
	Username string
	Password string
	Cams     []Cam
}

type Cam struct {
	Ip       string
	Username string
	Password string
	Name     string
}

var configFlag = flag.String("c", "hikmqtt.json", "config file")
var urlFlag = flag.String("s", "", "MQTT server url")
var usernameFlag = flag.String("u", "", "MQTT username")
var passwordFlag = flag.String("p", "", "MQTT password")
var camipFlag = flag.String("ci", "", "Camera ip")
var camuserFlag = flag.String("cu", "", "Camera username")
var campassFlag = flag.String("cp", "", "Camera password")
var helpFlag = flag.Bool("h", false, "Help")
var cfg Config

func config() int {
	flag.Parse()
	if *helpFlag {
		flag.PrintDefaults()
		return -1
	}
	bt, err := os.ReadFile(*configFlag)
	if err == nil {
		err = json.Unmarshal(bt, &cfg)
		if err != nil {
			log.Printf("Error parsing config: %s\n", err)
			return -2
		}
	} else {
		log.Println("Config not found.")
	}
	// cfg = Config{
	// 	Url:      *urlFlag,
	// 	Username: *usernameFlag,
	// 	Password: *passwordFlag,
	// }
	if *urlFlag != "" {
		cfg.Url = *urlFlag
	}
	if *usernameFlag != "" {
		cfg.Username = *usernameFlag
	}
	if *passwordFlag != "" {
		cfg.Password = *passwordFlag
	}

	if *camipFlag != "" {
		cfg.Cams = append(cfg.Cams, Cam{
			Ip:       *camipFlag,
			Username: *camuserFlag,
			Password: *campassFlag,
		})
	}

	if len(cfg.Cams) == 0 {
		log.Println("Cam not defined.")
		return -3
	}

	return 0
}
