package main

import (
    "fmt"
    "net/http"
	"log"
	"encoding/json"

	"io/ioutil"
)
var results []string
var names []string
var userSet = make(map[string]bool)

func addUser(userName string){
	_, exists := userSet[userName]
	if !exists {
		names = append(names, userName)
		userSet[userName] = true
		f, _ := json.Marshal(names)
		ioutil.WriteFile("names.json", f, 0x644);
	}
	log.Println(names)
}

//process user data
func processUserData(userJSON string) { // the data recieved from server
	// get that data out of those JSON
	var userData interface{}
	err := json.Unmarshal([]byte(userJSON), &userData)
	if err != nil {
		fmt.Println("error:", err)
	}
	//try to assert the interface into a map
	userMap, ok := userData.(map[string]interface{})
	// conditionally add to set-like array of usernames
	if ok {
	addUser(userMap["name"].(string))
	}
}

// PostHandler converts post request body to string
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		results = append(results, string(body))

		processUserData(string(body))
		fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	f, _ := ioutil.ReadFile("names.json");
	err := json.Unmarshal(f, &names)

	for _, name := range names{
		userSet[name] = true
	}

	if err != nil{
		log.Println(f)
	}

	mux := http.NewServeMux()

    mux.HandleFunc("/datasend", PostHandler)

    fs := http.FileServer(http.Dir("./"))
    mux.Handle("/", http.StripPrefix("/", fs))

	err = http.ListenAndServeTLS(":443", "dummycert.pem", "dummykey.pem", mux)
	log.Println("serving")
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
