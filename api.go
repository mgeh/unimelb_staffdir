/*
	Staff directory pilot API
	sits along side neo4j server serving up requests
*/

package main

import (
	"./staffdir"
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	// "time"
)

const (
	ENDPOINT = "http://uom-staffdir-neo4j.elasticbeanstalk.com:7474/db/data"
)

// log to file
func LogFile(message string) {
	f, err := os.OpenFile("staffdir_access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(message)
}

// Convert returned neo4j results to structs
func ProcessResults(t interface{}) []interface{} {
	z := reflect.ValueOf(t).Elem()
	s := make([]interface{}, z.Len())
	for i := 0; i < z.Len(); i++ {
		s[i] = z.Index(i).Interface()
	}
	// for _, n := range s {
	//  temp, _ := json.Marshal(n.(PersonSummary))
	//  fmt.Println(string(temp))
	// }
	return s
}

// Main function for the API, starts up martini instance
func main() {
	m := martini.Classic()
	fmt.Println("Initialising")
	db := new(staffdir.Database)
	db.Connect(ENDPOINT)

	m.Get("/staffdir/department", func(res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		block, _ := url.QueryUnescape(strings.SplitAfter(r.RequestURI, "/?")[1])
		fmt.Println(block)
		results, ok := db.SearchDepartment(block)
		if ok != nil {
			log.Fatalln("issue with results")
		}
		out := ProcessResults(results)
		fmt.Println(out)
		var temp []byte
		var tempOut []staffdir.PersonSummary
		if len(out) > 0 {
			for _, b := range out {
				tempOut = append(tempOut, b.(staffdir.PersonSummary))
			}
			temp, _ = json.Marshal(tempOut)
			//temp, _ := json.Marshal(out[0].(staffdir.PersonSummary))
		}
		return 200, string(temp)
	})

	// process authentication
	m.Get("/staffdir/person", func(res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		block, _ := url.QueryUnescape(strings.SplitAfter(r.RequestURI, "?")[1])
		results, ok := db.SearchPeople(block)
		if ok != nil {
			log.Fatalln("issue with results")
		}
		out := ProcessResults(results)
		fmt.Println(out)
		var temp []byte
		var tempOut []staffdir.PersonSummary
		if len(out) > 0 {
			for _, b := range out {
				tempOut = append(tempOut, b.(staffdir.PersonSummary))
			}
			temp, _ = json.Marshal(tempOut)
			//temp, _ := json.Marshal(out[0].(staffdir.PersonSummary))
		}

		return 200, string(temp)
	})

	m.Get("/staffdir/manager", func(res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		block, _ := url.QueryUnescape(strings.SplitAfter(r.RequestURI, "?")[1])
		results, ok := db.LookupManager(block)
		if ok != nil {
			log.Fatalln("issue with results")
		}
		out := ProcessResults(results)
		fmt.Println(out)
		var temp []byte
		var tempOut []staffdir.PersonSummary
		if len(out) > 0 {
			for _, b := range out {
				tempOut = append(tempOut, b.(staffdir.PersonSummary))
			}
			temp, _ = json.Marshal(tempOut)
			//temp, _ := json.Marshal(out[0].(staffdir.PersonSummary))
		}

		return 200, string(temp)
	})

	m.Get("/staffdir/colleagues", func(res http.ResponseWriter, r *http.Request) (int, string) {
		block, _ := url.QueryUnescape(strings.SplitAfter(r.RequestURI, "?")[1])
		results, ok := db.LookupColleagues(block)
		if ok != nil {
			log.Fatalln("issue with results")
		}
		out := ProcessResults(results)
		fmt.Println(out)
		var temp []byte
		var tempOut []staffdir.PersonSummary
		if len(out) > 0 {
			for _, b := range out {
				tempOut = append(tempOut, b.(staffdir.PersonSummary))
			}
			temp, _ = json.Marshal(tempOut)
			//temp, _ := json.Marshal(out[0].(staffdir.PersonSummary))
		}

		return 200, fmt.Sprintf("{\"size\": %d,\"data\": %s}", len(out), temp)
	})

	m.Get("/staffdir/reports", func(res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		block, _ := url.QueryUnescape(strings.SplitAfter(r.RequestURI, "?")[1])
		results, ok := db.LookupReports(block)
		if ok != nil {
			log.Fatalln("issue with results")
		}
		out := ProcessResults(results)
		fmt.Println(out)
		var temp []byte
		var tempOut []staffdir.PersonSummary
		if len(out) > 0 {
			for _, b := range out {
				tempOut = append(tempOut, b.(staffdir.PersonSummary))
			}
			temp, _ = json.Marshal(tempOut)
			//temp, _ := json.Marshal(out[0].(staffdir.PersonSummary))
		}

		return 200, string(temp)
	})

	m.Patch("/", func() {
		// update something
	})

	m.Put("/", func() {
		// replace something
	})

	m.Delete("/", func() {
		// destroy something
	})

	m.Options("/", func() {
		// http options
	})

	m.NotFound(func() string {
		// handle 404
		log.Printf("%s\n", "Yep...")
		return "Something went wrong."
	})

	//m.Run()
	err := http.ListenAndServe(":5001", m)
	if err != nil {
		log.Fatal(err)
	}
}
