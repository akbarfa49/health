package health_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/akbarfa49/health"
)


func TestRawatInap(t *testing.T) {
	t.Run("fetching RawatInap Data", func(t *testing.T) {
		j := health.AcquireFaskes()
		data, err :=j.InfoRawatInap("YOGYAKARTA", "SLEMAN","NONCOVID")
		if err != nil{
			log.Panic(err)
			return
		}
		v := health.RumahSakitCovid{}
		err = json.Unmarshal(data, &v)
		log.Println(err)
		os.WriteFile("out.json", data, 0664)
	})
}