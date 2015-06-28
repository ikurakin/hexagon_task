package csvfiles

import "log"

type Vprefix struct {
	CalculationId string
	VregionId     string
	VprefixId     string
	Vprefix       string
}

type VprefixData struct {
	Log  *log.Logger
	Data []Vprefix
}

func NewVprefixData(logger *log.Logger) *VprefixData {
	return &VprefixData{
		Log: logger,
	}
}

func (vd *VprefixData) GetData() interface{} {
	return vd.Data
}

func (vd *VprefixData) Write(p []byte) (n int, err error) {
	data, err := GetCSVData(p)
	if err != nil {
		vd.Log.Println(err)
		return
	}
	for _, d := range data {
		if len(d) > 3 {
			vd.Data = append(vd.Data, Vprefix{
				CalculationId: d[0],
				VregionId:     d[1],
				VprefixId:     d[2],
				Vprefix:       d[3],
			})
		}
	}
	n = len(p)
	return
}
