package services

type Data struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DataService struct {
	data []Data
}

func NewDataService() *DataService {
	return &DataService{
		data: []Data{
			{ID: 1, Name: "Rk"},
			{ID: 2, Name: "Data 2"},
			{ID: 3, Name: "Data 3"},
		},
	}
}

func (s *DataService) GetAllData() ([]Data, error) {
	return s.data, nil
}
