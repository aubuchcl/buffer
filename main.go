package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aubuchcl/buffer/capbuff"
)

func main() {
	newbuffer := capbuff.NewBuffer(50000)
	z, err := newbuffer.WriteString("thisisastring")
	if err != nil {
		os.Exit(0)
	}

	resp, err := http.Get("http://www.google.com")
	if err != nil {
		os.Exit(0)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(0)
	}

	newbuffer.Write(bs)

	//if I needed a higher capacity on newbuffer
	//newbuffer.Grow(10000)

	// fmt.Println(newbuffer.Cap()) -- returns 60k
	//fmt.Println(newbuffer)
	nextresp, err := http.Get("http://www.lipsum.com")
	if err != nil {
		os.Exit(0)
	}
	//attempts to append the information from the io.reader
	//into the available unused space byteslice and if
	//it doesnt error it writes to the byteslice
	y, err := newbuffer.ReadFrom(nextresp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(y)

	fmt.Println(z)
	fmt.Println(newbuffer)

}
