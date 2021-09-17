package httphelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jhallat/todo-schedule-service/logger"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func GetRequest(writer http.ResponseWriter,
	            request *http.Request,
	            pattern string,
	            function func (map[string]string) (interface{}, error))  {
	paramMap, err := parseUrl(request.URL.Path, pattern)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	value, err := function(paramMap)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	response, err := json.Marshal(value)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
}

func GetQueryRequest(writer http.ResponseWriter,
	request *http.Request,
	function func (values url.Values) (interface{}, error))  {
	value, err := function(request.URL.Query())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	response, err := json.Marshal(value)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
}

func PostRequest(writer http.ResponseWriter,
	             request *http.Request,
	             idField string,
	             value interface{},
	             function func () (int64, error)) {

	logger.Debug("Post request called")
	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logger.Warn("Bad Request - error converting body to bytes")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &value)
	if err != nil {
		logger.Warn(fmt.Sprintf("Bad Request - error unmarshalling json - %s", bodyBytes))
		logger.Warn(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	reflected := reflect.ValueOf(value).Elem()
	id := reflected.FieldByName(idField)
	if id.Int() != 0 {
		logger.Warn(fmt.Sprintf("Bad Request - id does not equal 0 on POST, provided data = %+v", value))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	newId, err := function()
	id.SetInt(newId)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	valueJson, err := json.Marshal(value)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(valueJson)
	writer.WriteHeader(http.StatusCreated)
}

func parseUrl(url string, pattern string) (map[string]string, error) {
	paramMap := make(map[string]string)
	urlItems := strings.Split(url, "/")
	patternItems := strings.Split(pattern, "/")
	if len(urlItems) > len(patternItems) {
		urlItems = urlItems[len(urlItems)-len(patternItems):]
	}
	for index, patternItem := range patternItems {
		if strings.HasPrefix(patternItem, ":") {
			label := patternItem[1:]
			if len(urlItems) > index {
				paramMap[label] = urlItems[index]
			} else {
				return nil, errors.New("pattern does not match url")
			}
		}
	}
	return paramMap, nil
}

func PutRequest(writer http.ResponseWriter,
	            request *http.Request,
	            pattern string,
	            idField string,
	            value interface{},
	            getFunction func (map[string]string) (interface{}, error),
	            saveFunction func() error) {
	paramMap, err := parseUrl(request.URL.Path, pattern)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	value, err = getFunction(paramMap)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if value == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &value)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	reflected := reflect.ValueOf(value).Elem()
	id := reflected.FieldByName(idField)
	if id.Int() != 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = saveFunction()
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
	return
}

func DeleteRequest(writer http.ResponseWriter,
	request *http.Request,
	pattern string,
	function func(map[string]string) error)  {

	paramMap, err := parseUrl(request.URL.Path, pattern)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
    err = function(paramMap)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
}