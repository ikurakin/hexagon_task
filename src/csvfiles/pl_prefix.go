package csvfiles

import "log"

type PlPrefix struct {
	RegionId string
	PrefixId string
	Prefix   string
}

type PlPrefixData struct {
	Log  *log.Logger
	Data []PlPrefix
}

func NewPlPrefixData(logger *log.Logger) *PlPrefixData {
	return &PlPrefixData{
		Log: logger,
	}
}

func (pd *PlPrefixData) GetData() interface{} {
	return pd.Data
}

func (pd *PlPrefixData) Write(p []byte) (n int, err error) {
	data, err := GetCSVData(p)
	if err != nil {
		pd.Log.Println(err)
		return
	}
	for _, d := range data {
		if len(d) > 2 {
			pd.Data = append(pd.Data, PlPrefix{
				RegionId: d[0],
				PrefixId: d[1],
				Prefix:   d[2],
			})
		}
	}
	n = len(p)
	return
}
