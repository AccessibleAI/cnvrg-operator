package v1

type PriorityClass struct {
	Name        string `json:"name"`
	Value       int32  `json:"value"`
	Description string `json:"description"`
}
