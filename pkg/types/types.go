package pkg

type Graph struct {
	Nodes *[]Node `json:"nodes"`
	Links *[]Link `json:"links"`
}

type Node struct {
	Id      string `json:"id"`
	Project string `json:string`
	Seed    bool   `json:boolean`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}
