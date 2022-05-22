package utils

import "log"

func CheckErr(err any) {
	if err != nil {
		log.Fatal(err)
	}
}
