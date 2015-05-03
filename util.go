package main

import "log"

func fatalIfErr(err error, what string) {
	if err != nil {
		log.Fatalf("%s failed with %s\n", what, err)
	}
}
