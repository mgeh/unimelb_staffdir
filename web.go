/*
	Staff directory pilot API
	sits along side neo4j server serving up requests
*/

package main

import (
	"github.com/vly/unimelb_staffdir/staffdir"
	// "./staffdir"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
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

// Clean names
func CleanNameDetails(b interface{}) staffdir.PersonDetail {
	k := b.(staffdir.PersonDetail)
	nameTemp := strings.Split(k.Name, " ")
	if len(k.PrefName) > 1 {
		k.Name = fmt.Sprintf("%s %s", k.PrefName, nameTemp[len(nameTemp)-1])
	} else if len(nameTemp) > 2 {
		k.Name = fmt.Sprintf("%s %s", nameTemp[0], nameTemp[len(nameTemp)-1])
	}
	return k
}

func CleanNameSummary(b interface{}) staffdir.PersonSummary {
	k := b.(staffdir.PersonSummary)
	nameTemp := strings.Split(k.Name, " ")
	if len(k.PrefName) > 1 {
		k.Name = fmt.Sprintf("%s %s", k.PrefName, nameTemp[len(nameTemp)-1])
	} else if len(nameTemp) > 2 {
		k.Name = fmt.Sprintf("%s %s", nameTemp[0], nameTemp[len(nameTemp)-1])
	}
	return k
}

// Convert returned neo4j results to structs
func ProcessResults(t interface{}) []interface{} {
	z := reflect.ValueOf(t).Elem()
	s := make([]interface{}, z.Len())
	for i := 0; i < z.Len(); i++ {
		s[i] = z.Index(i).Interface()
	}
	return s
}

// Output json blob
func ProcessSummaries(t interface{}) string {
	out := ProcessResults(t)
	fmt.Println(out)
	var temp []byte
	var tempOut []staffdir.PersonSummary
	if len(out) > 0 {
		for _, b := range out {
			k := CleanNameSummary(b)
			tempOut = append(tempOut, k)
		}
		temp, _ = json.Marshal(tempOut)
	}
	return string(temp)
}

// Output json blob
func ProcessDetails(t interface{}) string {
	out := ProcessResults(t)
	fmt.Println(out)
	var temp []byte
	var tempOut []staffdir.PersonDetail
	if len(out) > 0 {
		for _, b := range out {
			k := CleanNameDetails(b)
			tempOut = append(tempOut, k)
		}
		temp, _ = json.Marshal(tempOut)
	}
	return string(temp)
}

// preflight headers

func SetHeaders(res *http.ResponseWriter) *http.ResponseWriter {
	(*res).Header().Set("Access-Control-Allow-Origin", "*")
	(*res).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	return res
}

// Main function for the API, starts up martini instance
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ENDPOINT := os.Getenv("NEO4J_URL")

	m := martini.Classic()
	fmt.Println("Initialising")
	db := new(staffdir.Database)
	db.Connect(ENDPOINT)

	m.Get("/staffdir/department/:query", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		db.Connect(ENDPOINT)
		SetHeaders(&res)
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
		}
		return 200, fmt.Sprintf("{\"size\": %d, \"data\": %s}", len(out), string(temp))
	})

	// process authentication
	m.Get("/staffdir/person", func(res http.ResponseWriter, r *http.Request) (int, string) {
		db.Connect(ENDPOINT)
		SetHeaders(&res)
		block := ""
		if r.FormValue("q") != "" {
			block = r.FormValue("q")
		} else {
			return 200, ""
		}
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
				k := CleanNameSummary(b)
				tempOut = append(tempOut, k)
			}
			temp, _ = json.Marshal(tempOut)
		}

		return 200, fmt.Sprintf("{\"size\": %d, \"data\": %s}", len(out), string(temp))
	})

	m.Get("/staffdir/manager/:email", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		db.Connect(ENDPOINT)
		SetHeaders(&res)
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
		}

		return 200, fmt.Sprintf("{\"size\": %d, \"data\": %s}", len(out), string(temp))
	})

	m.Get("/staffdir/colleagues/:email", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		db.Connect(ENDPOINT)
		SetHeaders(&res)
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
		}

		return 200, fmt.Sprintf("{\"size\": %d, \"data\": %s}", len(out), string(temp))
	})

	m.Get("/staffdir/reports/:email", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		db.Connect(ENDPOINT)
		SetHeaders(&res)
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
		}

		return 200, fmt.Sprintf("{\"size\": %d, \"data\": %s}", len(out), string(temp))
	})

	m.Get("/staffdir/details", func(res http.ResponseWriter, r *http.Request) (int, string) {
		db.Connect(ENDPOINT)
		SetHeaders(&res)
		block := ""
		if r.FormValue("id") != "" {
			block = r.FormValue("id")
		} else {
			return 200, ""
		}
		person, err := db.LookupPerson(block)
		manager, ok := db.LookupManager(block)
		colleagues, ok := db.LookupColleagues(block)
		reports, ok := db.LookupReports(block)

		if ok != nil || err != nil {
			log.Fatalln("issue with results")
		}
		personOut := ProcessDetails(person)
		managerOut := ProcessSummaries(manager)
		colleaguesOut := ProcessSummaries(colleagues)
		reportsOut := ProcessSummaries(reports)

		output := fmt.Sprintf("{\"size\": %d, \"data\": {\"person\": %s, \"manager\": %s", 1, personOut, managerOut)
		if len(colleaguesOut) > 1 {
			output += fmt.Sprintf(", \"colleagues\": %s", colleaguesOut)
		}
		if len(reportsOut) > 1 {
			output += fmt.Sprintf(", \"reports\": %s", reportsOut)
		}
		log.Println(output)

		return 200, output + "}}"
	})

	m.Options("/", func(res http.ResponseWriter) {
		SetHeaders(&res)
	})

	m.Options("/staffdir/colleagues/:val", func(res http.ResponseWriter) {
		SetHeaders(&res)
	})

	m.Options("/staffdir/person/:val", func(res http.ResponseWriter) {
		SetHeaders(&res)
	})

	m.NotFound(func() string {
		// handle 404
		return "Something went wrong."
	})

	err := http.ListenAndServe(":"+os.Getenv("PORT"), m)
	if err != nil {
		log.Fatal(err)
	}
}
