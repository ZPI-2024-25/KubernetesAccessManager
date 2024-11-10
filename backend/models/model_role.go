package models

type Role struct {
	Name string `json:"name,omitempty"`
	Permit []Operation `json:"permit,omitempty"`
	Deny []Operation `json:"deny,omitempty"`
	Subroles []string `json:"subroles,omitempty"`
}