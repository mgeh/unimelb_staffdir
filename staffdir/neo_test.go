package staffdir

import (
	//"encoding/json"
	// "log"
	"os"
	"reflect"
	"testing"
)

const (
	//ENDPOINT   = "http://weapon.its.unimelb.edu.au/db/data"
	PERSON     = "Val"
	EMAIL      = "tania.elliott@unimelb.edu.au"
	DEPARTMENT = "ITS"
	PHONE      = "7966"
)

// TestConnectNeo checks the Neo4j DB connection functionality
func TestConnectNeo(t *testing.T) {
	ENDPOINT := os.Getenv("NEO4J_URL")
	db := new(Database)
	_, ok := db.Connect(ENDPOINT)
	//verify there are no errors
	if ok == false {
		t.Fail()
	}
	// verify DB object is there
	if db == nil {
		t.Fail()
	}
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

// Test main staff search functionality
func TestSearchPeople(t *testing.T) {
	ENDPOINT := os.Getenv("NEO4J_URL")
	db := new(Database)
	db.Connect(ENDPOINT)

	// Test results
	tests := map[string]string{
		"Val L": "Val Lyashov",
		"8632":  "Val Lyashov",
		"val.lyashov@unimelb.edu": "Val Lyashov",
	}
	for in, out := range tests {
		results, err := db.SearchPeople(in)
		record := ""
		if err != nil {
			t.Fail()
		}
		for _, b := range ProcessResults(results) {
			if b.(PersonSummary).Name == out {
				record = b.(PersonSummary).Name
				break
			}

		}
		// check if expected result has been found
		if out != record {
			t.Fail()
		}

	}
}

// Test name processing functionality
func TestProcessName(t *testing.T) {
	db := new(Database)

	// Test results
	tests := map[string]string{
		"val l":  "val.* l.*",
		"val":    "(^|.* )val.*",
		"lyasho": "(^|.* )lyasho.*",
	}

	for a, b := range tests {
		result := db.ProcessName(a)
		if b != result {
			t.Fail()
		}
	}
}

// Test department search results
// func TestSearchDepartment(t *testing.T) {
// 	tests := map[string]string{
// 		"itS val":     "Val Lyashov",
// 		"engineering": "Tony Zara",
// 		"marketing":   "Neil Ang",
// 	}
// 	db := new(Database)
// 	db.Connect(ENDPOINT)
// 	for a, b := range tests {
// 		results, err := db.SearchDepartment(a)
// 		if err != nil {
// 			t.Fail()
// 		}

// 		temp := ProcessResults(results)
// 		out := ""
// 		if len(a) > 1 {
// 			for _, z := range temp {
// 				log.Println(z.(PersonSummary).Name)
// 				if z.(PersonSummary).Name == b {
// 					out = z.(PersonSummary).Name
// 					break
// 				}
// 			}
// 			if out != b {
// 				log.Printf("%s \n", out)
// 				t.Fail()
// 			}

// 		} else {
// 			if temp[0] == nil {
// 				t.Fail()
// 			}
// 		}
// 	}
// }

// Test individual person lookup
func TestLookupPerson(t *testing.T) {
	ENDPOINT := os.Getenv("NEO4J_URL")
	db := new(Database)
	db.Connect(ENDPOINT)
	tests := map[string]string{
		"15504": "Tania Elliott",
		"546":   "Val Lyashov",
	}
	for a, b := range tests {
		results, err := db.LookupPerson(a)
		if err != nil {
			t.Fail()
		}
		temp := ProcessResults(results)
		if len(temp) == 0 {
			t.Fail()
		} else if b != temp[0].(PersonDetail).Name {
			t.Fail()
		}

	}
}

// Test lookup of a person's manager
func TestLookupManager(t *testing.T) {
	ENDPOINT := os.Getenv("NEO4J_URL")
	db := new(Database)
	db.Connect(ENDPOINT)
	tests := map[string]string{
		"15504": "Sendur Kathirgamanathan",
		"546":   "Tania Elliott",
	}
	for a, b := range tests {
		results, err := db.LookupManager(a)
		if err != nil {
			t.Fail()
		}
		temp := ProcessResults(results)[0]
		if b != temp.(PersonSummary).Name {
			t.Fail()
		}
	}
}

// // Test lookup of a person's colleagues
// func TestLookupColleagues(t *testing.T) {
// 	db := new(Database)
// 	db.Connect(ENDPOINT)
// 	tests := map[string]string{
// 		"tania.elliott@unimelb.edu.au": "Steven Wojnarowski",
// 		"val.lyashov@unimelb.edu.au":   "Greg Shea",
// 	}
// 	for a, b := range tests {
// 		results, err := db.LookupColleagues(a)
// 		if err != nil {
// 			t.Fail()
// 		}
// 		temp := ProcessResults(results)[0]
// 		if b != temp.(PersonSummary).Name {
// 			log.Printf("%s %s\n", b, temp.(PersonSummary).Name)
// 			t.Fail()
// 		}
// 	}

// }

// Test lookup for a person's direct reports
func TestLookupReports(t *testing.T) {
	ENDPOINT := os.Getenv("NEO4J_URL")
	db := new(Database)
	db.Connect(ENDPOINT)
	tests := map[string]string{
		"15504": "Val Lyashov",
	}
	for a, b := range tests {
		results, err := db.LookupReports(a)
		if err != nil {
			t.Fail()
		}
		temp := ProcessResults(results)
		out := ""
		if len(a) > 1 {
			for _, z := range temp {
				if z.(PersonSummary).Name == b {
					out = z.(PersonSummary).Name
					break
				}
			}
			if b != out {
				t.Fail()
			}
		} else {
			if temp[0] == nil {
				t.Fail()
			}
		}
	}
}

// Test lookup of a person's functional area and employment category (academic or professional)
// is this even necessary?!
func TestLookupPersonShort(t *testing.T) {

}
