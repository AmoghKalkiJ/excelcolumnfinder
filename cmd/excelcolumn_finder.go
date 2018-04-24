package main

import (
	find "excelcolumnfinder/pkg/finder"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/table", find.Portal)
	r.HandleFunc("/api/columnfinder", find.Columnfinder)
	fmt.Println("Started server on localhost:8089")
	if err := http.ListenAndServe(":8089", r); err != nil {
		panic(err)

	}

}
