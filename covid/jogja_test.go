package covid_test

import (
	"os"
	"testing"

	"github.com/akbarfa49/health/covid"
)


func TestByZip(t *testing.T) {
	t.Run("fetching html", func(t *testing.T) {
		j := covid.Jogja()
		go j.Init("eyJpdiI6Iit3Y2l2ZlltdHdJTTBlRDJiaGNHUlE9PSIsInZhbHVlIjoiT1BEY29QajJiMzArVk9vbFpIY2Q0Z0dwRis3Z0NjSGhoaGVSU0UxY3B6YlJXWjgyZWE5UHhuUVdEeTRaT0txbiIsIm1hYyI6ImM1YTU1NmUwZDhkMzhjMmNiMGQ0NTM3ZWY1N2QxYmI5Y2EyMGRkNjk1ZmM5MTljNjQ5NmM3NmIyNjFiOWQzMTAifQ%3D%3D","eyJpdiI6IngxVU1aMGtJaDRxQk1kalIwdFdzVUE9PSIsInZhbHVlIjoiQTVxbXlmYXJleVpLMGVsL2tPUGJmVGpLR2Y2eFJ1TUVGd3FBT2FhbTBwUmM3THo3Ny9aS0hyWTFxNURzL1VIUCIsIm1hYyI6Ijg4NDg5NWE4OTRjNDlkZGE2ZTgyYzMxMjM3NGZiNjA5MWNjOTlkNzgwYzg0YjVjOTBjNTRjMjliNjA2OWFhYzEifQ%3D%3D",1)
		_, body, _ :=j.GetByZipCode("55295")
		os.WriteFile("out.html", body, 0664)
	})
}