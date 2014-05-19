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
	ENDPOINT   = "http://localhost:7474/db/data"
	PERSON     = "Val"
	EMAIL      = "tania.elliott@unimelb.edu.au"
	DEPARTMENT = "ITS"
	PHONE      = "7966"
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
		return string(temp)
		// var data []string
		// if strings.Contains(block, "%7C") {
		// 	data = strings.Split(strings.SplitAfter(r.RequestURI, "/?")[1], "|")
		// } else {
		// 	data = strings.Split(strings.SplitAfter(r.RequestURI, "/?")[1], "|")
		// }
		// for a, b := range data {
		// 	log.Println("d: " + b)
		// }
		// output := EncodeJWT(GenDict(data))
		// ENDPOINT := "https://unimelbit.zendesk.com/access/jwt?jwt="
		// res.Header().Set("Location", (ENDPOINT + output))
	})

	// process authentication
	m.Get("/staffdir/person", func(res http.ResponseWriter, r *http.Request) (int, string) {
		block, _ := url.QueryUnescape(strings.SplitAfter(r.RequestURI, "/?")[1])
		fmt.Println(block)
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
		return string(temp)
	})

	m.Get("/staffdir/manager", func(res http.ResponseWriter, r *http.Request) (int, string) {
		block := strings.SplitAfter(r.RequestURI, "/?")[1]

		return 301, block
	})

	m.Get("/staffdir/colleagues", func(res http.ResponseWriter, r *http.Request) (int, string) {
		block := strings.SplitAfter(r.RequestURI, "/?")[1]
		return 301, block
	})

	m.Get("/staffdir/reports", func(res http.ResponseWriter, r *http.Request) (int, string) {
		block := strings.SplitAfter(r.RequestURI, "/?")[1]

		return 301, block
	})

	m.Patch("/", func() {
		// update something
	})

	// m.Post("/staffdir/ticket", func(r *http.Request, res http.ResponseWriter) (int, string) {
	// 	ENDPOINT := "https://staff.unimelb.edu.au/its/requests/index#/thankyou/"
	// 	if err := r.ParseForm(); err != nil {
	// 		log.Printf("%s", "nothing posted")
	// 	}
	// 	data_values := make(map[string]string)
	// 	for a, b := range r.Form {
	// 		data_values[a] = b[0]
	// 	}
	// 	output := SubmitTicket(data_values)
	// 	//res.Header().Set("Content-Type", "application/json")
	// 	// SubmitTicket(data_values)
	// 	res.Header().Set("Location", (ENDPOINT + output))
	// 	//return "<html><head><meta http-equiv=\"refresh\" content=\"0;URL=" + (ENDPOINT + output) + "\"></head></html>"
	// 	return 301, ""
	// })

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
	//err := http.ListenAndServeTLS(":8080", "/home/wep/tls/certs/star.its.unimelb.edu.au.chain.crt", "/home/wep/tls/private/star.its.unimelb.edu.au.key", m)
	if err != nil {
		log.Fatal(err)
	}
}
