package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/domtriola/slime-mold-sim/simulation"
)

func main() {
	http.HandleFunc("/gen", genHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func genHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	options := map[string]interface{}{}

	intParams := []string{
		"width",
		"height",
		"nFrames",
		"loopCount",
		"delay",
		"sensorDegree",
		"sensorDistance",
		"scentSpreadFactor",
	}
	for _, param := range intParams {
		assignIntParam(options, query, param)
	}

	imgName := simulation.Build(options)
	http.ServeFile(w, r, imgName)
}

func assignIntParam(options map[string]interface{}, query url.Values, key string) {
	value := query.Get(key)
	if value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println(err)
		} else {
			options[key] = intValue
		}
	}
}
