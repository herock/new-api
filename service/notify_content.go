package service

import (
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/dto"
)

func renderNotifyContent(content string, values []interface{}) string {
	for _, value := range values {
		content = strings.Replace(content, dto.ContentValueParam, fmt.Sprintf("%v", value), 1)
	}
	return content
}
