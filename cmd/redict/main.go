package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/ellezio/redict/internal/redict"
)

var db *redict.Database

func main() {
	db = redict.NewDatabase()

	http.HandleFunc("PUT /strings/{storeKey}", setStrings)
	http.HandleFunc("GET /strings/{storeKey}", getStrings)

	fmt.Println("Listening on :3000")
	http.ListenAndServe("localhost:3000", nil)
}

func setStrings(w http.ResponseWriter, r *http.Request) {
	storeKey := r.PathValue("storeKey")
	if storeKey == "" {
		http.Error(w, "storeKey not provided", http.StatusBadRequest)
		return
	}

	var b bytes.Buffer
	_, err := b.ReadFrom(r.Body)
	if err != nil && err != io.EOF {
		fmt.Println("Failed to read request body", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	if err = db.Set(storeKey, b.Bytes()); err != nil {
		http.Error(w, fmt.Sprintf("cannot set value for store: %s", err), http.StatusBadRequest)
		return
	}
}

func getStrings(w http.ResponseWriter, r *http.Request) {
	storeKey := r.PathValue("storeKey")
	if storeKey == "" {
		http.Error(w, "storeKey not provided", http.StatusBadRequest)
		return
	}

	data, err := db.Get(storeKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot get data: %s", err), http.StatusBadRequest)
		return
	}

	w.Write(data)
}
