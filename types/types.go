package types

type Person struct {
	Id  uint64 `json:"id"`
	Name  string `json:"name"`
	Surname string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age uint8 `json:"age"`
	Gender string `json:"gender"`
	Nationality string `json:"nationality"`
}

type Age struct {
	Count  int `json:"count"`
	Name string `json:"name"`
	Age uint8 `json:"age"`	
}

type Gender struct {
	Count  int `json:"count"`
	Name string `json:"name"`
	Gender string `json:"gender"`
	Probability float32 `json:"probability"`
}

type Country struct {
	CountryId  string `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Nationality struct {
	Count  int `json:"count"`
	Name string `json:"name"`
	Country []Country
}