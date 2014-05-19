package staffdir

import (
	//"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const (
	//ENDPOINT   = "http://weapon.its.unimelb.edu.au/db/data"
	ENDPOINT   = "http://localhost:7474/db/data"
	PERSON     = "Val"
	EMAIL      = "tania.elliott@unimelb.edu.au"
	DEPARTMENT = "ITS"
	PHONE      = "7966"
)

// TestConnectNeo checks the Neo4j DB connection functionality
func TestConnectNeo(t *testing.T) {
	db := new(Database)
	_, ok := db.Connect(ENDPOINT)
	//verify there are no errors
	assert.NotNil(t, ok)
	// verify DB object is there
	assert.NotNil(t, db.db)
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
	db := new(Database)
	db.Connect(ENDPOINT)

	// Test results
	tests := map[string]string{
		"Val L": "Val Lyashov",
		"8632":  "Val Lyashov",
		"val.lyashov@unimelb.edu": "Val Lyashov",
	}
	for in, out := range tests {
		results, ok := db.SearchPeople(in)
		record := ""
		assert.Nil(t, ok)
		for _, b := range ProcessResults(results) {
			if b.(PersonSummary).Name == out {
				record = b.(PersonSummary).Name
				break
			}
			// check if expected result has been found
			assert.Equal(t, out, record, "Couldn't find result")
		}

	}
}

// Test name processing functionality
func TestProcessName(t *testing.T) {
	db := new(Database)

	// Test results
	tests := map[string]string{
		"val l":  "val.* l.*",
		"val":    "(^| )val.*",
		"lyasho": "(^| )lyasho.*",
	}

	for a, b := range tests {
		result := db.ProcessName(a)
		assert.Equal(t, b, result, "Results did not match expected output")
	}
}

// Test department search results
func TestSearchDepartment(t *testing.T) {
	tests := map[string]string{
		"itS":         "Val Lyashov",
		"engineering": "Tony Zara",
		"marketing":   "Neil Ang",
	}
	db := new(Database)
	db.Connect(ENDPOINT)
	for a, b := range tests {
		results, ok := db.SearchDepartment(a)
		assert.Nil(t, ok)
		temp := ProcessResults(results)
		out := ""
		if len(a) > 1 {
			for _, z := range temp {
				if z.(PersonSummary).Name == b {
					out = z.(PersonSummary).Name
					break
				}
			}
			assert.Equal(t, b, out, "Name doesn't match.")
		} else {
			assert.Nil(t, temp[0], "Response not empty when it should be.")
		}
	}
}

// Test individual person lookup
func TestLookupPerson(t *testing.T) {
	db := new(Database)
	db.Connect(ENDPOINT)
	tests := map[string]string{
		"tania.elliott@unimelb.edu.au": "Tania Elliott",
		"val.lyashov@unimelb.edu.au":   "Val Lyashov",
	}
	for a, b := range tests {
		results, ok := db.LookupPerson(a)
		assert.Nil(t, ok)
		temp := ProcessResults(results)[0]
		assert.Equal(t, b, temp.(PersonDetail).Name, "Name doesn't match")

	}
}

// Test lookup of a person's manager
func TestLookupManager(t *testing.T) {
	db := new(Database)
	db.Connect(ENDPOINT)
	tests := map[string]string{
		"tania.elliott@unimelb.edu.au": "Sendur Kathirgamanathan",
		"val.lyashov@unimelb.edu.au":   "Tania Elliott",
	}
	for a, b := range tests {
		results, ok := db.LookupManager(a)
		assert.Nil(t, ok)
		temp := ProcessResults(results)[0]
		assert.Equal(t, b, temp.(PersonSummary).Name, "Name doesn't match")
	}
}

// Test lookup of a person's colleagues
func TestLookupColleagues(t *testing.T) {
	db := new(Database)
	db.Connect(ENDPOINT)
	tests := map[string]string{
		"tania.elliott@unimelb.edu.au": "Steven Wojnarowski",
		"val.lyashov@unimelb.edu.au":   "Greg Shea",
	}
	for a, b := range tests {
		results, ok := db.LookupColleagues(a)
		assert.Nil(t, ok)
		temp := ProcessResults(results)[0]
		assert.Equal(t, b, temp.(PersonSummary).Name, "Name doesn't match")
	}

}

// Test lookup for a person's direct reports
func TestLookupReports(t *testing.T) {
	db := new(Database)
	db.Connect(ENDPOINT)
	tests := map[string]string{
		"tania.elliott@unimelb.edu.au": "Val Lyashov",
		"val.lyashov@unimelb.edu.au":   "",
	}
	for a, b := range tests {
		results, ok := db.LookupReports(a)
		assert.Nil(t, ok)
		temp := ProcessResults(results)
		out := ""
		if len(a) > 1 {
			for _, z := range temp {
				if z.(PersonSummary).Name == b {
					out = z.(PersonSummary).Name
					break
				}
			}
			assert.Equal(t, b, out, "Name doesn't match.")
		} else {
			assert.Nil(t, temp[0], "Response not empty when it should be.")
		}
	}
}

// Test lookup of a person's functional area and employment category (academic or professional)
// is this even necessary?!
func TestLookupPersonShort(t *testing.T) {

}
