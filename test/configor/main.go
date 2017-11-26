package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string `required:"true"`
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}{}

var Config2 = struct {
	APPName string `yaml:"appname"`

	DB struct {
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Port     uint   `yaml:"port"`
	}

	Contacts []struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	}
}{}

func main() {
	//configor.New(&configor.Config{Debug: true}).Load(&Config, "config.yml")
	//fmt.Printf("config: %+v", Config)

	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, &Config2)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Printf("config: %+v", Config2)
}