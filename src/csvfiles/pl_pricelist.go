package csvfiles

import "log"

type PlPriceList struct {
	PricelistId    string
	PricelistType  string
	Priority       string
	CurrencyFactor string
	TregionId      string
	CurrencyId     string
	CarrierId      string
	ServiceType    string
	ValidFrom      string
}

type PlPriceListData struct {
	Log  *log.Logger
	Data []PlPriceList
}

func NewPlPriceListData(logger *log.Logger) *PlPriceListData {
	return &PlPriceListData{
		Log: logger,
	}
}

func (pl *PlPriceListData) GetData() interface{} {
	return pl.Data
}

func (pl *PlPriceListData) Write(p []byte) (n int, err error) {
	data, err := GetCSVData(p)
	if err != nil {
		pl.Log.Println(err)
		return
	}
	for _, d := range data {
		if len(d) > 8 {
			pl.Data = append(pl.Data, PlPriceList{
				PricelistId:    d[0],
				PricelistType:  d[1],
				Priority:       d[2],
				CurrencyFactor: d[3],
				TregionId:      d[4],
				CurrencyId:     d[5],
				CarrierId:      d[6],
				ServiceType:    d[7],
				ValidFrom:      d[8],
			})
		}
	}
	n = len(p)
	return
}
