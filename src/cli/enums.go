package main

import "strings"

type Priority int

const (
	Low Priority = iota // 0
	Medium
	High
	Critical
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Mid"
	case High:
		return "High"
	case Critical:
		return "Crit"
	default:
		return "Invalid Priority"
	}
}

func PriorityStingInt(priority string) string {
	priority = strings.ToLower(priority)
	switch priority {
	case "low":
		return "0"
	case "mid":
		return "1"
	case "high":
		return "2"
	case "crit":
		return "3"
	default:
		return "4"
	}
}
