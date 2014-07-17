package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"regexp"
	"strings"
)

func genList(file string) *map[string]string {
	output := new(map[string]string)
	if data, err := ioutil.ReadFile(file); err == nil {
		lines := strings.Split(string(data), "\r")
		for _, b := range lines {
			tmp := strings.Split(b, ",")
			log.Println(tmp[1], tmp[3])
		}
	}
	return output
}

func main() {
	ENDPOINT := os.Getenv("PROF_PICS")
	if resp, err := http.Get(ENDPOINT); err == nil {
		defer resp.Body.Close()
		regex := regexp.MustCompile(`href="thumbnail(.*).*">t`)

		if res, err := ioutil.ReadAll(resp.Body); err == nil {
			for _, b := range strings.Split(string(res), "\n") {
				if match := regex.FindStringSubmatch(b); len(match) > 1 {
					log.Println(match[1])
				}
			}
		}
	}
}
