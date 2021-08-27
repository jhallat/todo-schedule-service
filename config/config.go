package config

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
)


type propDef struct {
	index int
	field string
	prop  string
	env   string
}


func Scan(config interface{}) error {
	propDefs := parsePropDefs(config)
	profile, err := loadFlags(propDefs, config)
	if err != nil {
		return err
	}
	populateEnvironment(propDefs, config)
	if profile != "" {
		profile := "app-" + profile + ".config"
		err = populateConfiguration(profile, propDefs, config)
		if err != nil {
			return err
		}
	}
	err = populateConfiguration("app.config", propDefs, config)
	return err
}

func parsePropDefs(config interface{}) []propDef {
	propDefs := make([]propDef,0)
	val := reflect.Indirect(reflect.ValueOf(config))
	for i := 0; i < val.Type().NumField(); i++ {
		propDefs = append(propDefs, propDef {
			i,
			val.Type().Field(i).Name,
			val.Type().Field(i).Tag.Get("prop"),
			val.Type().Field(i).Tag.Get("env"),
		})
	}
	return propDefs
}

func loadFlags(propDefs []propDef, config interface{}) (string, error) {
	flags := make(map[string]*string)
	profile := flag.String("profile", "", "profile")
	for _, item := range propDefs {
		value := flag.String(item.prop, "", item.field)
		flags[item.field] = value
	}
	flag.Parse()
	for _, item := range propDefs {
		if *flags[item.field] != "" {
			err := setFieldIfEmpty(config, item.field, *flags[item.field])
			if err != nil {
				return "", err
			}
		}
	}
	return *profile, nil
}

func populateEnvironment(propDefs []propDef, config interface{}) {
	for _, item := range propDefs {
		env := os.Getenv(item.env)
		if env != "" {
			setFieldIfEmpty(config, item.field, env)
		}
	}
}

func populateConfiguration(filename string, propDefs []propDef, config interface{} ) error {

	propMap := make(map[string]string)
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		file, _ := os.Open("app.config")
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var key string
			var value string
			num, err := fmt.Sscanf(scanner.Text(), "%s %s", &key, &value)
			if err != nil || num != 2 {
				return errors.New(fmt.Sprintf("invalid property file %s", filename))
			}
			propMap[key] = value
		}
	}
	for _, item := range propDefs {
		if propMap[item.prop] != "" {
			err = setFieldIfEmpty(config, item.field, propMap[item.prop])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func setFieldIfEmpty(config interface{}, field string, value string) error {
	val := reflect.ValueOf(config).Elem()
	fld := val.FieldByName(field)
	if fld.IsValid() {
		if fld.String() != "" {
			return nil
		}
		if fld.CanSet() {
			fld.SetString(value)
		} else {
			return errors.New(fmt.Sprintf("Cannot set %s", field))
		}
	} else {
		return errors.New(fmt.Sprintf("%s is not a valid field", field))
	}
	return nil
}