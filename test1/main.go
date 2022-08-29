package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

type Сonf struct {
	Update int `yaml:"update"`
	Every  int `yaml:"every"`
}

func (c *Сonf) getConf() {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println("yamlFile.Get err   %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

}

func main() {
	var c = new(Сonf)

	c.getConf()
	tm := make(chan int, 10)
	go c.timer(tm)

	go c.update(tm)

	fmt.Scanln()
}

func (c *Сonf) timer(tm chan int) {
	var i int
	for {
		i++

		tm <- i

		time.Sleep(time.Second)

		ost := i % c.Every
		if ost == 0 {
			fmt.Println("время с запуска", i)
		}
	}
}

func (c *Сonf) update(tm chan int) {

	for {
		up := <-tm
		ot := up % c.Update
		if ot == 0 {
			c.getConf()
		}

	}
}
