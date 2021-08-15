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
		data, err :=j.InfoRawatInap("YOGYAKARTA", "SLEMAN","COVID")
		if err != nil{
			log.Panic(err)
			return
		}
		body,_ := json.MarshalIndent(&data, "", "   ")
		os.WriteFile("out.json", body, 0664)
	})
}