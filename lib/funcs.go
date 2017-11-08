package lib

import "log"

func handErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
