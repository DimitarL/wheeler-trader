package common

import (
	"fmt"
	"net/url"
	"strconv"
)

func ExtractQueryParameters(query url.Values, stringParams []string, numericParams []string) (map[string]interface{}, error) {
	params := make(map[string]interface{})

	for _, queryParamName := range stringParams {
		err := addStringParam(query, params, queryParamName)
		if err != nil {
			return params, err
		}
	}

	for _, queryParamName := range numericParams {
		err := addIntParam(query, params, queryParamName)
		if err != nil {
			return params, err
		}
	}

	return params, nil
}

func addStringParam(query url.Values, params map[string]interface{}, param string) error {
	paramValue, ok, err := RetrieveStringQueryParam(query, param)
	if err != nil {
		return err
	}
	if ok {
		params[param] = paramValue
	}

	return nil
}

func addIntParam(query url.Values, params map[string]interface{}, param string) error {
	paramValue, ok, err := RetrieveIntQueryParam(query, param)
	if err != nil {
		return err
	}
	if ok {
		params[param] = paramValue
	}

	return nil
}

func RetrieveStringQueryParam(query url.Values, param string) (string, bool, error) {
	values, ok := query[param]
	if !ok {
		return "", false, nil
	}
	if len(values) > 1 {
		return "", true, fmt.Errorf("query parameter '%s' must be used only once", param)
	}

	return values[0], true, nil
}

func RetrieveIntQueryParam(query url.Values, param string) (int, bool, error) {
	values, ok := query[param]
	if !ok {
		return -1, false, nil
	}
	if len(values) > 1 {
		return -1, true, fmt.Errorf("query parameter '%s' must be used only once", param)
	}

	intValue, err := strconv.Atoi(values[0])
	if err != nil {
		return -1, true, fmt.Errorf("query parameter '%s' is not valid integer: %w", param, err)
	}

	return intValue, true, nil
}
