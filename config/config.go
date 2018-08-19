package config

import "gopkg.in/yaml.v2"
import "io/ioutil"

type Conf struct {
	FailToBanLoc string `yaml:"FailToBanLoc"`
	DeviceName   string `yaml:"DeviceName"`
}

func (c *Conf) Parse(conf_location string) *Conf {
	yamlFile, err := ioutil.ReadFile(conf_location)
	if err != nil {
		panic("Conf file not found")
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic("Unable to unmarshal ")
	}
	return c
}
