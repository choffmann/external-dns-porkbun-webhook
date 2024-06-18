package porkbun

type ZoneStatus string

const (
	Active ZoneStatus = "ACTIVE"
)

type Zone struct {
	Domain string     `json:"domain"`
	Status ZoneStatus `json:"status"`
	TLD    string     `json:"tld"`
}
