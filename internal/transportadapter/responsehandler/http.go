package responsehandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type httpResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func HTTPHandleResponseError(rawRes []byte) error {
	if len(rawRes) == 0 {
		return nil
	}

	if isJSONArray(rawRes) {
		return nil
	}

	httpRes, err := formatHTTPResponse(rawRes)
	if httpRes == nil && err == nil {
		return nil
	} else if err != nil {
		return err
	}

	if httpRes.Status >= 400 {
		errMsg := strings.ToLower(httpRes.Message)
		errCode := httpRes.Status
		r := fmt.Sprintf("%d %s", errCode, errMsg)
		return errors.New(r)
	}

	return nil
}

func isJSONArray(b []byte) bool {
	return json.Unmarshal(b, &[]interface{}{}) == nil
}

func formatHTTPResponse(rawRes []byte) (*httpResponse, error) {
	res := make(map[string]interface{})
	err := json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, err
	}

	if _, ok := res["message"]; !ok {
		return nil, nil
	}

	httpRes := &httpResponse{
		Message: res["message"].(string),
	}
	if e, ok := res["error"].(float64); ok {
		httpRes.Status = int(e)
	} else {
		httpRes.Status = int(res["status"].(float64))
	}

	return httpRes, nil
}
