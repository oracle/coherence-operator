/*
 * Copyright (c) 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oracle/coherence-go-client/v2/coherence"
	"log"
	"net/http"
	"strconv"
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	cache coherence.NamedMap[int, Person]
	ctx   = context.TODO()
)

func main() {
	session, err := coherence.NewSession(ctx, coherence.WithPlainText())
	if err != nil {
		log.Println("unable to connect to Coherence", err)
		return
	}
	defer session.Close()

	cache, err = coherence.GetNamedMap[int, Person](session, "people")
	if err != nil {
		log.Println("unable to create namedMap 'people'", err)
		return
	}

	http.HandleFunc("/api/people", personHandler)
	http.HandleFunc("/api/people/", personByIDHandler)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func personHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var p Person
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		_, err := cache.Put(ctx, p.ID, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(p)

	case http.MethodGet:
		var people = make([]Person, 0)

		for ch := range cache.Values(ctx) {
			if ch.Err != nil {
				http.Error(w, ch.Err.Error(), http.StatusInternalServerError)
				return
			}
			people = append(people, ch.Value)
		}
		_ = json.NewEncoder(w).Encode(people)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func personByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/people/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, err1 := cache.Get(ctx, id)

		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		if p == nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err1 = json.NewEncoder(w).Encode(&p)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		old, err2 := cache.Remove(ctx, id)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		if old == nil {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
