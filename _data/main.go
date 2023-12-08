package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Microsoft/go-winio"
	"github.com/gocarina/gocsv"
)

const (
	PIPENAME = `\\.\pipe\USDJPY`
)

type OHLCv struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func main() {
	t := 30 * time.Second
	conn, err := winio.DialPipe(`\\.\pipe\USDJPY`, &t)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var his []OHLCv
	f, _ := os.Open("./USDJPY.csv")
	gocsv.Unmarshal(f, &his)
	fmt.Println(len(his))

	s := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(s)

	for i := 0; i < len(his); i++ {
		conn.Write([]byte(fmt.Sprintf("%f,%f,%f,%f,%d\r\n", his[i].Open, his[i].High, his[i].Low, his[i].Volume, time.Now().UnixMilli())))
		time.Sleep(time.Duration(r.Intn(10)) * time.Millisecond)

		if his[i].Open == 0 {
			fmt.Println("return")
			i = 0
		}
	}
}
