package domain

import "strings"

func SpanKindValue(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "0", "unspecified":
		return "unspecified"
	case "1", "internal":
		return "internal"
	case "2", "server":
		return "server"
	case "3", "client":
		return "client"
	case "4", "producer":
		return "producer"
	case "5", "consumer":
		return "consumer"
	default:
		return strings.TrimSpace(strings.ToLower(value))
	}
}
