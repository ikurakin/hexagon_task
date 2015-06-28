package csvfiles

import "log"

type PlRegPrice struct {
	PricelistId string
	RegionId    string
	Price       string
	ValidFrom   string
	ValidTo     string
	RegionName  string
	TariffId    string
	TariffType  string
	Parm1       string
}

type PlRegPriceData struct {
	Log  *log.Logger
	Data []PlRegPrice
}

func NewPlRegPriceData(logger *log.Logger) *PlRegPriceData {
	return &PlRegPriceData{
		Log: logger,
	}
}

func (pr *PlRegPriceData) GetData() interface{} {
	return pr.Data
}

func (pr *PlRegPriceData) Write(p []byte) (n int, err error) {
	data, err := GetCSVData(p)
	if err != nil {
		pr.Log.Println(err)
		return
	}
	for _, d := range data {
		if len(d) > 8 {
			pr.Data = append(pr.Data, PlRegPrice{
				PricelistId: d[0],
				RegionId:    d[1],
				Price:       d[2],
				ValidFrom:   d[3],
				ValidTo:     d[4],
				RegionName:  d[5],
				TariffId:    d[6],
				TariffType:  d[7],
				Parm1:       d[8],
			})
		}
	}
	n = len(p)
	return
}
