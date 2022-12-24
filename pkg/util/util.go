package util

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

func LoadJsonFile(filepath string) map[string]interface{} {
	jsonFile, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return make(map[string]interface{})
	}
	var loadedJson = make(map[string]interface{})
	err = json.Unmarshal(jsonFile, &loadedJson)
	if err != nil {
		log.Printf("Could not parse json from %v\nERROR: %v\n", filepath, err)
		return make(map[string]interface{})
	}
	log.Printf("Loaded %v \n", filepath)
	return loadedJson
}

func LoadMultipleJsonFiles(metadataFile string) map[string]interface{} {
	retMap := make(map[string]interface{})
	if len(metadataFile) > 0 {
		log.Println("Read metada from file(s)")
		if metadataFile != "" {
			for _, mdFile := range strings.Split(metadataFile, ",") {
				log.Printf("Try to load json file %v\n", mdFile)
				input := LoadJsonFile(mdFile)
				for k, v := range input {
					retMap[k] = v
				}
			}
		}
	}
	return retMap
}

func YesOrNoQuestion(question string) bool {
	fmt.Printf("%v (Y/n) ", question)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(strings.TrimSpace(text))
	if text != "" && text != "y" && text != "yes" {
		return false
	}
	return true
}

// Set field by json name
func SetValueInStructByJsonKey(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("tag %v provided does not define a json tag", fieldName)
	}
	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jname, _ := findJsonName(tag)
		fieldNames[jname] = i
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return fmt.Errorf("field %s does not exist within the provided item", fieldName)
	}
	fieldVal := v.Field(fieldNum)
	fieldVal.Set(reflect.ValueOf(value))
	return nil
}

// TODO improve approach
func IsPropperHash(hash string) bool {
	_, err := hex.DecodeString(hash)
	return err == nil && len(hash) == 64
}

func GetJsonFromRaw(raw []byte) (string, error) {
	loadedJson, err := GetMapFromRaw(raw)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", loadedJson), err

}
func GetMapFromRaw(raw []byte) (map[string]interface{}, error) {
	var loadedJson = make(map[string]interface{})
	err := json.Unmarshal(raw, &loadedJson)
	return loadedJson, err
}

func GetPrettyJsonFromString(ugly string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(ugly), "", "  ")
	return fmt.Sprintf("%s\n", out.Bytes()), err
}

func GetPrettyJsonFromMap(input map[string]interface{}) (string, error) {
	log.Println("Convert metada to json")
	jsonStr, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return "", err
	}
	return string(jsonStr), nil
}
