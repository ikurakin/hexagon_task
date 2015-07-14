package csvfiles

import (
	"archive/zip"
	"database/sql"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// pl_prefix.csv(region_id) == pl_regprice.csv(region_id)
// pl_regprice.csv(pricelist_id) == pl_pricelist.csv(pricelist_id)
type ResultPrice struct {
	//vprefix fields
	VregionId string
	VprefixId string
	Vprefix   string
	//pl_pricelist fields
	PricelistId    string
	CurrencyFactor string
	CurrencyId     string
	CarrierId      string
	//pr_regprice fields
	RegionId   string
	Price      string
	RegionName string
	TariffId   string
	TariffType string
	//pl_prefix fields
	PrefixId string
	Prefix   string
}

func NewResultPrice() *ResultPrice {
	return &ResultPrice{}
}

type ResultData struct {
	sync.Mutex
	wg          sync.WaitGroup
	ParsedFiles map[string]HexWriter
	Data        []ResultPrice
	DB          *sql.DB
	Log         *log.Logger
}

func NewResultData(dbconn string) *ResultData {
	log_file, _ := os.OpenFile("process.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	logger := log.New(log_file, "margin_logger: ", log.Llongfile)
	db, err := sql.Open("sqlite3", dbconn)
	if err != nil {
		logger.Println(err)
	}
	return &ResultData{
		ParsedFiles: make(map[string]HexWriter),
		DB:          db,
		Log:         logger,
	}
}

// func (rd *ResultData) ProcessData() {
// 	r, ok := rd.ParsedFiles["pl_pricelist.csv"]
// 	if ok {
// 		data := r.(*PlPriceListData).GetData()
// 		plPrice := data.([]PlPriceList)
// 		for _, p := range plPrice {
// 			rp := NewResultPrice()
// 			rp.PricelistId = p.PricelistId
// 			rp.CurrencyFactor = p.CurrencyFactor
// 			rp.CurrencyId = p.CurrencyId
// 			rp.CarrierId = p.CarrierId
// 			rd.wg.Add(1)
// 			go rd.getRegPriceData(rp)
// 		}

// 		rd.wg.Wait()
// 		log.Println(rd.Data)
// 	}
// }

// func (rd *ResultData) getRegPriceData(resPrice *ResultPrice) {
// 	r, ok := rd.ParsedFiles["pl_regprice.csv"]
// 	if ok {
// 		regPrice := r.(*PlRegPriceData).GetData()
// 		for _, rp := range regPrice.([]PlRegPrice) {
// 			if rp.PricelistId == resPrice.PricelistId {
// 				resPrice.RegionId = rp.RegionId
// 				resPrice.Price = rp.Price
// 				resPrice.RegionName = rp.RegionName
// 				resPrice.TariffId = rp.TariffId
// 				resPrice.TariffType = rp.TariffType
// 				go rd.getPrefixData(resPrice)
// 				return
// 			}
// 		}
// 	}
// 	rd.wg.Done()
// }

func (rd *ResultData) ProcessData() {
	r, ok := rd.ParsedFiles["pl_regprice.csv"]
	if ok {
		data := r.(*PlRegPriceData).GetData()
		regPrice := data.([]PlRegPrice)
		for _, p := range regPrice {
			rp := NewResultPrice()
			rp.PricelistId = p.PricelistId
			rp.RegionId = p.RegionId
			rp.Price = p.Price
			rp.RegionName = p.RegionName
			rp.TariffId = p.TariffId
			rp.TariffType = p.TariffType
			rd.wg.Add(1)
			go rd.getPriceData(rp)
		}

		rd.wg.Wait()
		log.Println(rd.Data)
	}
}

func (rd *ResultData) getPriceData(resPrice *ResultPrice) {
	r, ok := rd.ParsedFiles["pl_pricelist.csv"]
	if ok {
		plPrice := r.(*PlPriceListData).GetData()
		for _, p := range plPrice.([]PlPriceList) {
			if p.PricelistId == resPrice.PricelistId {
				// resPrice.PricelistId = p.PricelistId
				resPrice.CurrencyFactor = p.CurrencyFactor
				resPrice.CurrencyId = p.CurrencyId
				resPrice.CarrierId = p.CarrierId
				go rd.getPrefixData(resPrice)
				return
			}
		}
	}
	rd.wg.Done()
}

func (rd *ResultData) getPrefixData(resPrice *ResultPrice) {
	r, ok := rd.ParsedFiles["pl_prefix.csv"]
	if ok {
		prefData := r.(*PlPrefixData).GetData()
		for _, p := range prefData.([]PlPrefix) {
			if p.RegionId == resPrice.RegionId {
				resPrice.PrefixId = p.PrefixId
				resPrice.Prefix = p.Prefix
				go rd.getVPrefixData(resPrice)
				return
			}
		}
	}
	rd.wg.Done()
}

func (rd *ResultData) getVPrefixData(resPrice *ResultPrice) {
	r, ok := rd.ParsedFiles["vprefix.csv"]
	if ok {
		vprefData := r.(*VprefixData).GetData()
		for _, vp := range vprefData.([]Vprefix) {
			pref, err := strconv.ParseInt(resPrice.Prefix, 10, 64)
			if err != nil {
				rd.Log.Println(err)
				continue
			}
			if vp.Vprefix == strconv.FormatInt(pref, 10) {
				resPrice.VregionId = vp.VregionId
				resPrice.VprefixId = vp.VprefixId
				resPrice.Vprefix = vp.Vprefix

				rd.Lock()
				rd.Data = append(rd.Data, *resPrice)
				rd.Unlock()
				rd.wg.Done()
				return
			}
		}
	}
	rd.wg.Done()
}

//SetParsedFiles parses csv files data from archive to a map
//filepath - path to archive
func (rd *ResultData) SetParsedFiles(filepath string) {
	r, err := zip.OpenReader(filepath)
	if err != nil {
		rd.Log.Println(err)
		return
	}
	defer r.Close()

	// Iterate through the files in the archive,
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			rd.Log.Println(err)
			return
		}
		defer rc.Close()
		var w HexWriter
		switch f.Name {
		case "vprefix.csv":
			w = NewVprefixData(rd.Log)
		case "pl_pricelist.csv":
			w = NewPlPriceListData(rd.Log)
		case "pl_prefix.csv":
			w = NewPlPrefixData(rd.Log)
		case "pl_regprice.csv":
			w = NewPlRegPriceData(rd.Log)
		}
		_, err = io.Copy(w, rc)
		if err != nil {
			rd.Log.Println(err)
			return
		}
		rd.ParsedFiles[f.Name] = w
	}
}
