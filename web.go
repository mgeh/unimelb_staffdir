/*
	Staff directory pilot API
	sits along side neo4j server serving up requests
*/

package main

import (
	"github.com/vly/unimelb_staffdir/staffdir"
	// "./staffdir"
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"os"
	"reflect"
	// "time"
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
	ENDPOINT := os.Getenv("GRAPHENEDB_URL")

	m := martini.Classic()
	fmt.Println("Initialising")
	db := new(staffdir.Database)
	db.Connect(ENDPOINT)

	m.Get("/staffdir/department/:query", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		if params["query"] == "" {
			return 200, ""
		}
		block := params["query"]
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
	m.Get("/staffdir/person/:name", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		if params["name"] == "" {
			return 200, ""
		}
		block := params["name"]
		results, err := db.SearchPeople(block)

		if err != nil {
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

	m.Get("/staffdir/manager/:email", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		if params["email"] == "" {
			return 200, ""
		}
		block := params["email"]
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

	m.Get("/staffdir/colleagues/:email", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		if params["email"] == "" {
			return 200, ""
		}
		block := params["email"]
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

	m.Get("/staffdir/reports/:email", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		if params["email"] == "" {
			return 200, ""
		}
		block := params["email"]
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
	err := http.ListenAndServe(":"+os.Getenv("PORT"), m)
	if err != nil {
		log.Fatal(err)
	}
}
