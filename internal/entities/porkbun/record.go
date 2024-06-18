package porkbun

type Record struct {
	Type    string `json:"type"`
	Id      string `json:"id"`
	Content string `json:"content"`
	TTL     string `json:"ttl"`
	Name    string `json:"name"`
}
