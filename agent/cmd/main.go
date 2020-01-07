package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"git.2dfire.net/zerodb/agent/glog"
	"gopkg.in/yaml.v3"

	"git.2dfire.net/zerodb/agent/server"
)

const (
	version     = "v1.1.0"
	defaultYaml = "./etc/app.yaml"
)

var (
	configFile string
	v          bool
)

func init() {
	flag.BoolVar(&v, "v", false, "show Agent version")
	flag.StringVar(&configFile, "config", defaultYaml, "Agent config file")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			glog.Glog.Errorf("panic %v", err)
		}
	}()

	flag.Parse()

	if v {
		fmt.Println(version)
		return
	}

	strDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return
	}

	file := defaultYaml

	if len(configFile) != 0 {
		file = configFile
	}

	data, err := ioutil.ReadFile(strDir + string(filepath.Separator) + file)
	if err != nil {
		fmt.Printf("agent read app.yaml failed: %v", err)
		return
	}

	conf := new(server.Config)
	if err = yaml.Unmarshal(data, conf); err != nil {
		fmt.Printf("agent parsel config failed: %v", err)
		return
	}

	if err = glog.CreateLogs(conf.LogConf.Path, conf.LogConf.Level); err != nil {
		fmt.Printf("agent create logs failed: %v", err)
		return
	}

	var s *server.Server
	s, err = server.NewServer(conf)
	if err != nil {
		fmt.Printf("agent new server failed: %v", err)
		return
	}
	/*
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		go func() {
			sig := <-sc
			if sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				fmt.Println("agent exit with 0")
				os.Exit(0)
			}
		}()
	*/
	if err = s.Run(); err != nil {
		fmt.Printf("agent run failed: %v", err)
		return
	}
}
