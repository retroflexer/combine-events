package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func transformData(pIn *interface{}) (err error) {
	switch in := (*pIn).(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(in))
		for k, v := range in {
			if err = transformData(&v); err != nil {
				return err
			}
			var sk string
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			case bool:
				sk = strconv.FormatBool(k.(bool))
			case nil:
				sk = "null"
			case float64:
				sk = strconv.FormatFloat(k.(float64), 'f', -1, 64)
			default:
				return fmt.Errorf("type mismatch: expect map key string or int; got: %T", k)
			}
			m[sk] = v
		}
		*pIn = m
	case []interface{}:
		for i := len(in) - 1; i >= 0; i-- {
			if err = transformData(&in[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	searchDir := os.Args[1]

	var mergedEvents []interface{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.Name() == "events.yaml" {
			body := make(map[interface{}]interface{})
			yamlFile, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(yamlFile, body)
			if err != nil {
				return err
			}
			events := (body["items"].([]interface{}))
			for i := len(events) - 1; i >= 0; i-- {
				if err = transformData(&events[i]); err != nil {
					return err
				}
			}
			mergedEvents = append(mergedEvents, events...)
		}
		return nil
	})

	mergedMap := make(map[string]interface{})
	mergedMap["items"] = mergedEvents
	if mergedJson, err := json.Marshal(mergedMap); err != nil {
		panic(err)
	} else {
		fmt.Printf("%s\n", mergedJson)
	}
}
