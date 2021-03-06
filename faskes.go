package health

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

//Prov on siranap
var ProvSir = map[string]string{
	"ACEH": "11prop",
	"SUMATERA UTARA":"12prop",
	"SUMATERA BARAT":"13prop",
	"RIAU":"14prop",
	"JAMBI": "15prop",
	"SUMATERA SELATAN": "16prop",
	"BENGKULU":"17prop",
	"LAMPUNG":"18prop",
	"KEPULAUAN BANGKA BELITUNG": "19prop",
	"KEPULAUAN RIAU":"20prop",
	"DKI JAKARTA": "31prop",
	"JAWA BARAT": "32prop",
	"JAWA TENGAH": "33prop",
	"YOGYAKARTA": "34prop",
	"JAWA TIMUR": "35prop",
	"BANTEN" : "36prop",
	"BALI" : "51prop",
	"NUSA TENGGARA BARAT": "52prop",
	"NUSA TENGGARA TIMUR": "53prop",
	"KALIMANTAN BARAT":"61prop",
	"KALIMANTAN TENGAH":"62prop",
	"KALIMANTAN SELATAN": "63prop",
	"KALIMANTAN TIMUR":"64prop",
	"KALIMANTAN UTARA":"65prop",
	"SULAWESI UTARA":"71prop",
	"SULAWESI TENGAH":"72prop",
	"SULAWESI SELATAN":"73prop",
	"SULAWESI TENGGARA":"74prop",
	"GORONTALO":"75prop",
	"SULAWESI BARAT":"76prop",
	"MALUKU":"81prop",
	"MALUKU UTARA":"82prop",
	"PAPUA BARAT":"91prop",
	"PAPUA":"92prop",

}

const (
	COVID = iota+1
	NONCOVID
)

type Faskes struct {
	f *fasthttp.Client
	pool *sync.Pool
}

type RumahSakitCovid struct{
	Nama string `json:"nama"`
	Alamat string `json:"alamat"`
	Informasi string `json:"informasi,omitempty"`
	UpdateTerakhir string `json:"updateTerakhir,omitempty"`
	Nomor string `json:"nomor,omitempty"`
	Link string `json:"url,omitempty"`
}

type RumahSakitNonCovid struct{
	Nama string `json:"nama"`
	Alamat string `json:"alamat"`
	Bed []struct{
		Jumlah string `json:"jumlahBed"`
		Kelas string `json:"kelas"`
		Ruang string `json:"ruang"`
		UpdateTerakhir string `json:"updateTerakhir,omitempty"`
	} `json:"informasi,omitempty"`
	Nomor string `json:"nomor,omitempty"`
	Link string `json:"url,omitempty"`
}


func AcquireFaskes() *Faskes{
	fas := &Faskes{
		f: &fasthttp.Client{},
		pool: &sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{})
			},
		},
	}
	return fas
}

//data will return array of map or array of RumahSakitCovid or RumahSakitNonCovid
func (fas *Faskes) InfoRawatInap(prop, kab, jenis string) (data []byte, err error){
req := fasthttp.AcquireRequest()
res := fasthttp.AcquireResponse()
defer fasthttp.ReleaseRequest(req)
defer fasthttp.ReleaseResponse(res)

v, ok :=ProvSir[strings.ToUpper(prop)]
if !ok{
	return nil, errors.New("Nama Provinsi tidak valid")
}
d, err := fas.AmbilDataKotaRawatInap(prop)
if err != nil{
	return nil, err
	}
tipe := 0
switch strings.ToUpper(jenis){
	case "COVID":
		tipe = COVID
	case "NONCOVID":
		tipe = NONCOVID
	default:
		return nil, errors.New("Jenis Rawat Inap Tidak Valid. Pilih COVID atau NONCOVID")
}
	uri := fmt.Sprintf("https://yankes.kemkes.go.id/app/siranap/rumah_sakit?jenis=%v&propinsi=%v&kabkota=%v", tipe,v,d[strings.ToUpper(kab)])
req.SetRequestURI(uri)
if err := fasthttp.Do(req,res); err != nil{
	return nil, err
}
switch tipe{
case COVID:
res := parseCovidHospitalData(res.Body())
data, _ :=json.Marshal(&res)
return data,nil
case NONCOVID:
res := parseNonCovidHospitalData(res.Body())	
data, _ :=json.Marshal(&res)
return data, nil
default:
	return nil, errors.New("Jenis Rawat Inap Tidak Valid. Pilih COVID atau NONCOVID")
}
}

/*AmbilDataKotaRawatInap digunakan untuk mengambil list data kota menggunakan parameter provinsi

format map[string]string
data [Nama_KabKota] = KODEKOTA

contoh
data["KULON PROGO"]="3401"

*/
func (fas *Faskes) AmbilDataKotaRawatInap(prop string) (data map[string]string, err error){
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	
	v, ok :=ProvSir[strings.ToUpper(prop)]
	if !ok{
		return nil, errors.New("Nama Provinsi tidak valid")
	}
	
	req.SetRequestURI("https://yankes.kemkes.go.id/app/siranap/Kabkota?kode_propinsi="+v)
	
	fasthttp.Do(req, res)
	if res.StatusCode() != 200{
		return nil, errors.New("Terjadi kesalahan ketika memproses data")
	}
	i := fas.pool.Get().(map[string]interface{})
	data = make(map[string]string)
	if err := json.Unmarshal(res.Body(),&i); err != nil{
		return nil, err
	}
	if !i["success"].(bool){
		return nil, errors.New("Nama Provinsi tidak valid")
	}
	for _,v:= range i["data"].([]interface{}){
		nama := v.(map[string]interface{})["nama_kabkota"].(string)
		
		kode := v.(map[string]interface{})["kode_kabkota"].(string)
		data[strings.ToUpper(nama)] = kode
	}
	return data, nil
}
	
	
func parseCovidHospitalData(body []byte) (rs []RumahSakitCovid){
	br :=bytes.NewReader(body)
	doc, _ :=goquery.NewDocumentFromReader(br)
	
	sel := doc.Find("div[data-string] > .card")
	sel.Each(func(i int, s *goquery.Selection) {
		var rumahdata RumahSakitCovid
		s1 := s.Find(".card-body > .row")
		rumahdata.Nama = s1.Find(".col-md-7 > h5 ").Text()
		rumahdata.Alamat = s1.Find(".col-md-7 > p ").Text()
		s1in := s1.Find(".col-md-5 > p")
		slot := ""
		antrian := ""
		s1in.Each(func(i int, s *goquery.Selection) {
			
			switch i{
			case 1:
				
				
				slot = strings.Join(strings.Split(strings.Trim(s.Find("b").Text(), "\n"), " "), " ")
				return
			case 2:
				antrian = strings.TrimSpace(s.Text())
				return
			case 3:
				rumahdata.UpdateTerakhir=strings.TrimSpace(s.Text())
			default:
				return
			}	
		})
		
		s2 := s.Find(".card-footer > div ")
		rumahdata.Nomor = s2.Find("span").Text()
		info := fmt.Sprintf("Tersedia %v Kasur kosong IGD\n", slot)
		info += fmt.Sprintf("%v\n", antrian)
		rumahdata.Informasi = info
		rs = append(rs, rumahdata)
	})
return rs
}


func parseNonCovidHospitalData(body []byte) (rs []RumahSakitNonCovid){
	br :=bytes.NewReader(body)
	doc, _ :=goquery.NewDocumentFromReader(br)
	
	sel := doc.Find("div[data-string] > .card")
	sel.Each(func(i int, s *goquery.Selection) {
		var rumahdata RumahSakitNonCovid
		s1 := s.Find(".card-body > .row")
		rumahdata.Nama = s1.Find(".col-md-5 > h5 ").Text()
		rumahdata.Alamat = s1.Find(".col-md-5 > p ").Text()
		s1in := s1.Find(".col-md-7 > div > div")
		s1in.Each(func(i int, s *goquery.Selection) {
			var bed struct{
				Jumlah string `json:"jumlahBed"`
				Kelas string `json:"kelas"`
				Ruang string `json:"ruang"`
				UpdateTerakhir string `json:"updateTerakhir,omitempty"`
				} 

				s1 := s.Find("div > .card")
				s1.Find(".card-body > div").Each(func(i int, s *goquery.Selection) {
					
					switch i+1{
					case 1:
					bed.Jumlah = strings.TrimSpace(s.Text())
					case 2:
					bed.Kelas = strings.TrimSpace(s.Text())
					case 3:
					bed.Ruang = strings.TrimSpace(s.Text())
				}
				})
				bed.UpdateTerakhir= strings.TrimSpace(s1.Find(".card-footer > div").Text())
				rumahdata.Bed = append(rumahdata.Bed, bed)
		})
		
		s2 := s.Find(".card-footer > div ")
		rumahdata.Nomor = s2.Find("span").Text()
		rs = append(rs, rumahdata)
	})
return rs
}

func stupid(){
	
}