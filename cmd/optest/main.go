/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	var err error

	portStr := os.Getenv("COH_HEALTH_PORT")
	if portStr == "" {
		portStr = "6676"
	}

	port, err := strconv.Atoi(portStr)
	panicIfError(err)

	http.HandleFunc("/", Handle)
	http.HandleFunc("/ready", Ready)
	http.HandleFunc("/healthz", Ready)
	http.HandleFunc("/ha", Ready)

	address := fmt.Sprintf(":%d", port)
	fmt.Printf("http server listening on %s\n", address)
	err = http.ListenAndServe(address, nil)
	panicIfError(err)
}

func Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	env := make(map[string]string)

	for _, s := range os.Environ() {
		kv := strings.Split(s, "=")
		if len(kv) == 1 {
			env[kv[0]] = ""
		} else {
			env[kv[0]] = os.Getenv(kv[0])
		}
	}

	m := make(map[string]interface{})
	m["env"] = env
	m["args"] = os.Args[1:]

	j, err := json.Marshal(m)
	panicIfError(err)

	_, err = w.Write(j)
	panicIfError(err)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
