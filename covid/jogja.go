package covid

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)


type JogjaCov struct {
	payload []struct{
	xsrftoken []byte
	session   []byte
	hiddenToken []byte
	}
	f *fasthttp.Client
	//pool *sync.Pool
	wg *sync.WaitGroup
}

//The Information Provider retired for serving the information
func Jogja() *JogjaCov {
	j := &JogjaCov{f: &fasthttp.Client{}}
	j.wg = &sync.WaitGroup{}
	j.wg.Add(1)
	/*j.pool = &sync.Pool{
		New: func() interface{} {
			var a struct{
				xsrftoken string
				session string
				hiddenToken string
			}
			return a
		},
	}*/
	return j
}

//Set fasthttp.Client if not set default will be use instead.
func (j *JogjaCov) SetClient(f *fasthttp.Client){
	j.f = f
	return 
}

//the full initiation will done after 2*GenN Minutes, but for early 10 Data will be done and returned to struct
func (j *JogjaCov) Init(xsrf, sess string, genNx10 int){
	if genNx10 == 0{
		genNx10 = 1
	}
	for i:=0;i<genNx10;i++{
		j.init(xsrf,sess)
	time.Sleep(2*time.Minute)	
	}
}
//Init used for generate xsrf and laravel session. parameter used from the lastest known xsrf and sess. it may take times for init to be done due to request limitation
func (j *JogjaCov) init(xsrf, sess string) {
	var a struct{
		xsrftoken []byte
		session []byte
		hiddenToken []byte
	}
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	c := fasthttp.AcquireCookie()
	
	req.SetRequestURI("https://sebaran-covid19.jogjaprov.go.id/kodepos")

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseCookie(c)
	
	for i := 0;i<10;i++{
	fasthttp.Do(req,res)
	c.ParseBytes(res.Header.PeekCookie("XSRF-TOKEN"))
	a.xsrftoken = c.Value()
	c.Reset()
	c.ParseBytes(res.Header.PeekCookie("laravel_session"))
	a.session = c.Value()
	c.Reset()
	token, err := peekToken(res.Body())
	if err != nil{
		log.Println(err)
		continue
	}
	a.hiddenToken = token
	j.payload = append(j.payload, a)
	}

	log.Println("done")
	j.wg.Done()
}

func peekToken(body []byte) ([]byte, error){
	b := bytes.NewReader(body)
	body = body[:0]
	r  := bufio.NewReader(b)
	for i:=0;i<63;i++{
		r.ReadLine()
	}
	l, _, _ := r.ReadLine()
	prefix := []byte(`value="`)
	prefixLen := len(prefix)
	suffix := []byte(`"`)[0]
	for i:=0;i<len(l)-prefixLen;i++{
		if bytes.Compare(l[i:i+prefixLen],prefix) == 0{
			for j :=i+prefixLen+1;j<len(l);j++{
				if l[j] == suffix{
					return l[i+prefixLen:j], nil
				}
			}
		}
	}
	return nil,errors.New("Data Not Found")
}

//Get Information By ZipCode and return it
func (j *JogjaCov) GetByZipCode(zipCode string) (header *fasthttp.ResponseHeader, body []byte, err error){
	j.wg.Wait()
	rand.Seed(time.Now().UnixNano())
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	n := rand.Intn(len(j.payload))
	req.SetRequestURI("https://sebaran-covid19.jogjaprov.go.id/result")
	req.Header.SetCookieBytesKV([]byte("XSRF-TOKEN"),j.payload[n].xsrftoken)
	req.Header.SetCookieBytesKV([]byte("laravel_session"),j.payload[n].session)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	s := fmt.Sprintf("_token=%v&kode_pos=%v", string(j.payload[n].hiddenToken), zipCode)
	req.AppendBodyString(s)
	if err :=fasthttp.Do(req,res); err !=nil{
		return nil, nil, err
}
log.Println(string(j.payload[n].xsrftoken))
log.Println(string(j.payload[n].session))
log.Println(string(j.payload[n].hiddenToken))
log.Println(string(req.Body()))
	return &res.Header, res.Body(), nil
}