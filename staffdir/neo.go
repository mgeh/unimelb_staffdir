package staffdir

import (
	"errors"
	"github.com/jmcvetta/neoism"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Database struct {
	db       *neoism.Database
	endpoint string
}

type PersonSummary struct {
	Name       string `json:"a.name"`
	Position   string `json:"a.position"`
	Department string `json:"a.department"`
	Phone      string `json:"a.phone"`
	Mobile     string `json:"a.mobile"`
	Email      string `json:"a.email"`
	LastName   string `json:"a.l_name"`
	PrefName   string `json:"a.pref_name"`
	Gender     string `json:"a.gender"`
	Id         int    `json:"pid"`
}

type PersonDetail struct {
	Name             string `json:"a.name"`
	Position         string `json:"a.position"`
	PositionGroup    string `json:"a.position_group"`
	Department       string `json:"a.department"`
	DepartmentNumber string `json:"a.department_number"`
	LocCampus        string `json:"a.loc_campus"`
	LocBuilding      string `json:"a.loc_building"`
	LocFloor         string `json:"a.loc_floor"`
	LocRoom          string `json:"a.loc_room"`
	Phone            string `json:"a.phone"`
	Mobile           string `json:"a.mobile"`
	Email            string `json:"a.email"`
	Gender           string `json:"a.gender"`
	LastName         string `json:"a.l_name"`
	PrefName         string `json:"a.pref_name"`
}

func (db *Database) Connect(hostname string) (database *neoism.Database, ok bool) {
	hostname = hostname + "/db/data/"
	database, err := neoism.Connect(hostname)
	if err != nil {
		log.Fatal(err)
	}
	db.db = database
	ok = true
	return
}

func (db *Database) ProcessName(in string) (out string) {
	parts := strings.Split(in, " ")

	if len(parts) > 1 {
		out = strings.Join(parts, ".* ")
		out += ".*"
	} else {
		out = "(^|.* )" + in + ".*"
	}
	return
}

func (db *Database) ProcessQuery(in string) (qtype string, query string) {
	in = strings.ToLower(in)
	isNum, _ := regexp.MatchString("^[0-9 ()+]{4,15}$", in)
	if strings.Contains(in, "@") {
		qtype = "email"
	} else if isNum {
		qtype = "phone"
		query = ".*" + in + ".*"
	} else {
		qtype = "name"
		query = db.ProcessName(in)
	}
	if len(query) == 0 {
		query = in + ".*"
	}
	log.Println(query)
	return
}

func (db *Database) SearchPeople(query string) (results interface{}, err error) {
	var qtype string
	qtype, query = db.ProcessQuery(query)
	cq := neoism.CypherQuery{
		Statement:  "MATCH (a:Person) WHERE a." + qtype + " =~{name} RETURN a.name, a.position, a.department, a.phone, a.mobile, a.email, a.l_name, a.pref_name, a.gender, id(a) AS pid ORDER BY a.position DESC, a.name LIMIT 100",
		Parameters: neoism.Props{"name": "(?i)" + query},
		Result:     &[]PersonSummary{},
	}
	// db.Session.Log = true
	db.db.Cypher(&cq)
	results = cq.Result
	log.Println(results)
	if results == nil {
		err = errors.New("No results returned")
	}
	return
}

// department search results
func (db *Database) SearchDepartment(query string) (results interface{}, err error) {
	departmentList := map[string]string{
		"its": "Information Technology Services",
		"eng": "Engineering",
	}
	searchTerm := strings.ToLower(query)
	if _, ok := departmentList[searchTerm]; ok {
		searchTerm = departmentList[searchTerm]
	}
	cq := neoism.CypherQuery{
		Statement:  "MATCH (a:Person) WHERE a.department =~ {department} RETURN a.name, a.position, a.department, a.phone, a.mobile, a.email ORDER BY a.position DESC, a.name LIMIT 1000",
		Parameters: neoism.Props{"department": "(?i)" + searchTerm},
		Result:     &[]PersonSummary{},
	}
	// db.Session.Log = true
	db.db.Cypher(&cq)
	results = cq.Result
	if results == nil {
		err = errors.New("No results returned")
	}
	return
}

// individual person lookup
func (db *Database) LookupPerson(query string) (results interface{}, err error) {
	val, _ := strconv.Atoi(query)
	cq := neoism.CypherQuery{
		Statement:  "MATCH (a:Person) WHERE id(a) = {id} RETURN a.name, a.position, a.position_group, a.department, a.department_number, a.loc_campus, a.loc_building, a.loc_floor, a.loc_room, a.phone, a.mobile, a.email, a.gender, a.l_name, a.pref_name",
		Parameters: neoism.Props{"id": val},
		Result:     &[]PersonDetail{},
	}
	// db.Session.Log = true
	db.db.Cypher(&cq)
	results = cq.Result
	if results == nil {
		err = errors.New("No results returned")
	}
	return
}

// lookup of a person's manager
func (db *Database) LookupManager(query string) (results interface{}, err error) {
	val, _ := strconv.Atoi(query)
	cq := neoism.CypherQuery{
		Statement:  "MATCH (a:Person)-[:MANAGES]->(b:Person) WHERE id(b) = {id} RETURN a.name, a.position, a.department, a.phone, a.mobile, a.email, a.l_name, a.pref_name, a.gender, id(a) AS pid LIMIT 1",
		Parameters: neoism.Props{"id": val},
		Result:     &[]PersonSummary{},
	}
	// db.Session.Log = true
	db.db.Cypher(&cq)
	results = cq.Result
	if results == nil {
		err = errors.New("No results returned")
	}
	return
}

// lookup of a person's colleagues
func (db *Database) LookupColleagues(query string) (results interface{}, err error) {
	val, _ := strconv.Atoi(query)
	cq := neoism.CypherQuery{
		Statement:  "MATCH (b:Person)<-[:MANAGES]-(c:Person)-[:MANAGES]->(a:Person) WHERE id(b) = {id} RETURN a.name, a.position, a.department, a.phone, a.mobile, a.email, a.l_name, a.pref_name, a.gender, id(a) AS pid ORDER BY a.position DESC, a.name LIMIT 100",
		Parameters: neoism.Props{"id": val},
		Result:     &[]PersonSummary{},
	}
	// db.Session.Log = true
	db.db.Cypher(&cq)
	results = cq.Result
	if results == nil {
		err = errors.New("No results returned")
	}
	return
}

// lookup for a person's direct reports
func (db *Database) LookupReports(query string) (results interface{}, err error) {
	val, _ := strconv.Atoi(query)
	cq := neoism.CypherQuery{
		Statement:  "MATCH (b:Person)-[:MANAGES]->(a:Person) WHERE id(b) = {id} RETURN a.name, a.position, a.department, a.phone, a.mobile, a.email, a.l_name, a.pref_name, a.gender, id(a) AS pid ORDER BY a.position DESC, a.name LIMIT 100",
		Parameters: neoism.Props{"id": val},
		Result:     &[]PersonSummary{},
	}
	// db.Session.Log = true
	db.db.Cypher(&cq)
	results = cq.Result
	if results == nil {
		err = errors.New("No results returned")
	}
	return
}
