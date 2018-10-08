package models

import (
	"demo/conf"
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Deal request body
func Deal(pkgName string) string {
	var x string
	x = strings.Replace(pkgName, "\n", "", -1)
	x = strings.Replace(x, "\r", "", -1)
	if strings.HasSuffix(x, "/") {
		x = x[:len(x)-1]
	}
	x = strings.Replace(x, "/", "-", -1)
	return x
}

// Find data via pkg name
func Find(pkgName string) (map[string]interface{}, error) {
	var (
		jsonData interface{}
		needList []string
	)
	finalData := make(map[string]interface{})
	needList = conf.Config.Show
	fileData, err := ioutil.ReadFile(conf.Config.Path + pkgName)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(fileData, &jsonData); err != nil {
		return nil, err
	}
	for _, each := range needList {
		finalData[each] = jsonData.(map[string]interface{})[each]
	}
	return finalData, nil
}
