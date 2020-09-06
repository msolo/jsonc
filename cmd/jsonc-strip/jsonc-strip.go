package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"io/ioutil"

	"github.com/msolo/jsonc"
)

var usage = `Simple tool to strip comments from JSONC, yielding a plain-old JSON.

  jsonc-strip < something.jsonc > something.json
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	out, err := jsonc.StripComments(in)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stdout.Write(out)
	if err != nil {
		log.Fatal(err)
	}
}