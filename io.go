package main

import (
	"encoding/json"
	"io/ioutil"
)

func GenConfigName() string {
	return "data.wp"
}

func mResponse(v interface{}) string {
	data, _ := json.Marshal(&v)
	return string(data)
}

//Config functions
func SaveConfig(v interface{}) error {

	str := mResponse(v)
	pathoffile := GenConfigName()
	strbytes := []byte(str)
	err := ioutil.WriteFile(pathoffile, strbytes, 0700)
	strbytes = nil

	return err

}
func LoadConfig(targ interface{}) error {

	pathoffile := GenConfigName()
	data, err := ioutil.ReadFile(pathoffile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, targ)

	return err
}
