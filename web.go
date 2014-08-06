/*
	Staff directory pilot API
	sits along side neo4j server serving up requests
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/vly/unimelb_staffdir/staffdir"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"regexp"
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

// Clean phone numbers
func CleanPhone(b string) string {
	b = strings.TrimSpace(b)
	if len(b) > 1 {
		reg, err := regexp.Compile("[^0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		b = reg.ReplaceAllString(b, "")
		if len(b) == 5 {
			if b[0] == 56 {
				b = "903" + b
			} else {
				b = "834" + b
			}
		} else if len(b) == 6 {
			b = "83" + b
		} else {
			if len(b) > 8 {
				b = b[len(b)-8:]
			}
		}
		b = "03 " + b
	}
	// regex := regexp.MustCompile(".{1,3}")
	// regex.FindAllString(s, n)
	// b = regex/.{1,4}/)
	return b
}

// Clean mobile numbers
func CleanMobile(b string) string {
	if len(b) > 1 {
		b = strings.TrimSpace(b)
		reg, err := regexp.Compile("[^0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		b = reg.ReplaceAllString(b, "")
		if len(b) > 9 {
			b = "0" + b[len(b)-9:]
		}
	}
	// regex := regexp.MustCompile(".{1,3}")
	// regex.FindAllString(s, n)
	// b = regex/.{1,4}/)
	return b
}

// Clean email
func CleanEmail(b string) string {
	b = strings.Replace(b, " ", "", -1)
	b = strings.ToLower(b)
	return b
}

// Clean names
func CleanName(name string, prefName string, lastName string) string {
	nameTemp := strings.Split(name, " ")
	if len(prefName) > 1 {
		name = fmt.Sprintf("%s %s", prefName, lastName)
	} else if len(nameTemp) > 2 {
		name = fmt.Sprintf("%s %s", nameTemp[0], lastName)
	}
	return name
}

// Clean positions (incl. default)
func CleanPosition(position string) string {
	if len(position) == 0 {
		position = "Casual / Honorary"
	}
	return position
}

func CleanDetails(b interface{}) staffdir.PersonDetail {
	k := b.(staffdir.PersonDetail)
	k.Name = CleanName(k.Name, k.PrefName, k.LastName)
	k.Phone = CleanPhone(k.Phone)
	if len(k.Mobile) > 0 {
		k.Mobile = CleanMobile(k.Mobile)
	}
	k.Email = CleanEmail(k.Email)

	if len(k.Position) < 1 {
		k.Position = CleanPosition(k.Position)
	}
	return k
}

func CleanSummary(b interface{}) staffdir.PersonSummary {
	k := b.(staffdir.PersonSummary)
	k.Name = CleanName(k.Name, k.PrefName, k.LastName)
	k.Phone = CleanPhone(k.Phone)
	if len(k.Mobile) > 0 {
		k.Mobile = CleanMobile(k.Mobile)
	}
	k.Email = CleanEmail(k.Email)
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
			k := CleanSummary(b)
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
			k := CleanDetails(b)
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
	S3_LOC := os.Getenv("S3_LOC")
	LOCAL_LOC := os.Getenv("LOCAL_LOC")

	m := martini.Classic()
	fmt.Println("Initialising")
	db := new(staffdir.Database)
	db.Connect(ENDPOINT)

	// provide webhook endpoint for S3 pulldown
	m.Get("/staffdir/update", func(params martini.Params, res http.ResponseWriter, r *http.Request) (int, string) {
		status := "OK"

		if len(S3_LOC) < 5 || len(LOCAL_LOC) < 5 {
			status = "Failed, no S3 bucket or local location set"
			return 501, fmt.Sprintf("{\"status\": \"%s\"}", status)
		}

		cmd = exec.Command("/usr/local/bin/aws", "s3", "sync", S3_LOC, LOCAL_LOC)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			LogFile(fmt.Sprintf("%s", err))
			status = "Failed to exec AWS cli s3 sync"
			LogFile(fmt.Sprintf("%s", out))
		} else {
			cmd := exec.Command("service", "staffdir", "stop")
			if err := cmd.Run(); err != nil {
				status = "Failed to stop staffdir service"
			}
			cmd = exec.Command("ln", "-s", "/usr/share/nginx/www/web", "/usr/local/bin/staffdir_api")
			if err := cmd.Run(); err != nil {
				status = "Failed to replace api sym links"
			}
			cmd = exec.Command("service", "staffdir", "start")
			if err := cmd.Run(); err != nil {
				status = "Failed to start staffdir service"
			}
		}
		return 200, fmt.Sprintf("{\"status\": \"%s\"}", status)
	})

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

	// lookup person
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

		var temp []byte
		var tempOut []staffdir.PersonSummary
		if len(out) > 0 {
			for _, b := range out {
				k := CleanSummary(b)
				tempOut = append(tempOut, k)
			}
			temp, _ = json.Marshal(tempOut)
		}

		return 200, fmt.Sprintf("{\"size\": %d, \"data\": %s}", len(out), string(temp))
	})

	// get suggested peopple
	m.Get("/staffdir/suggestions", func(res http.ResponseWriter, r *http.Request) (int, string) {
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
		var tempOut []staffdir.Suggestion
		if len(out) > 0 {
			for _, b := range out {
				k := CleanSummary(b)
				z := new(staffdir.Suggestion)
				z.Name = k.Name
				z.Department = k.Department
				z.Phone = k.Phone
				z.Email = k.Email
				z.Id = k.Id

				tempOut = append(tempOut, *z)
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

		output := fmt.Sprintf("{\"size\": %d, \"data\": {\"person\": %s", 1, personOut)
		if len(managerOut) > 1 {
			output += fmt.Sprintf(", \"manager\": %s", managerOut)
		}
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
