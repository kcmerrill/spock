package channels

type template struct {
	ID         string `json:"id"`
	Properties struct {
		Attempts int    `json:"attempts"`
		Module   string `json:"module"`
		Error    string `json:"error"`
		Output   string `json:"output"`
		TaskName string `json:"name"`
	} `json:"properties"`
}
