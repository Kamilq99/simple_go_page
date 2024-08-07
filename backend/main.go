package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	data  = []map[string]string{}
	mutex = &sync.Mutex{}
)

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/submit", submitData)

	fmt.Println("Serwer działa na porcie 8080")
	http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
}

func submitData(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var entry map[string]string
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("POST: Odpowiedź 400 - Bad Request")
			return
		}
		mutex.Lock()
		data = append(data, entry)
		mutex.Unlock()

		// Logowanie dodanego wpisu
		fmt.Printf("POST: Dodano wpis: %+v\n", entry)
		w.WriteHeader(http.StatusNoContent) // 204 No Content
		fmt.Println("POST: Odpowiedź 204 - No Content")
	} else if r.Method == http.MethodGet {
		mutex.Lock()
		defer mutex.Unlock()

		if len(data) == 0 {
			w.WriteHeader(http.StatusNoContent) // 204 No Content
			fmt.Println("GET: Odpowiedź 204 - No Content")
			return
		}

		json.NewEncoder(w).Encode(data)
		// Logowanie, że dane zostały pobrane
		fmt.Printf("GET: Wysłano dane: %+v\n", data)
		w.WriteHeader(http.StatusOK) // 200 OK
		fmt.Println("GET: Odpowiedź 200 - OK")
	}
}
