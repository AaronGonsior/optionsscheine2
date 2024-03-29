package optionsscheine2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"math"
)

func check(err error){
	if err!=nil{
		fmt.Println(err)
	}
}


func WriteJson(path string, content string){
	// Open a file for writing
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// Encode the string as JSON and write it to the file
	if err := json.NewEncoder(file).Encode(content); err != nil {
		fmt.Println(err)
		return
	}
}

func LoadJson(path string) string{

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	// Decode the JSON data from the file
	var readStr string
	if err := json.NewDecoder(file).Decode(&readStr); err != nil {
		fmt.Println(err)
		return ""
	}

	// Return the decoded string
	return fmt.Sprint(readStr)

}

func JsonToOptions(path string) []Option {
	str := LoadJson(path)
	str = strings.Replace(str,"[{","",-1)
	str = strings.Replace(str,"}]","",-1)
	strAr := strings.Split(str,"} {")
	var options []Option
	var params []string
	var newOption Option
	var tmp float64
	for _,ar := range strAr {
		params = strings.Split(ar," ")

		spc,err := strconv.ParseFloat(params[5],64)
		check(err)
		sp, err := strconv.ParseFloat(params[6],64)
		check(err)

		newOption = Option{
			Cfi:                 params[0],
			Contract_type:       params[1],
			Exerciese_style:     params[2],
			Expiration_date:     params[3],
			Primaty_exchange:    params[4],
			Shares_per_contract: int(spc),
			Strike_price:        sp,
			Ticker:              params[7],
			Underlying_ticker:   params[8],
		}


		tmp, err = strconv.ParseFloat(params[9],64)
		check(err)
		newOption.Volume = int(tmp)
		newOption.Vw,err = strconv.ParseFloat(params[10],64)
		check(err)
		newOption.Open,err = strconv.ParseFloat(params[11],64)
		check(err)
		newOption.Close,err = strconv.ParseFloat(params[12],64)
		check(err)
		newOption.High,err = strconv.ParseFloat(params[13],64)
		check(err)
		newOption.Low,err = strconv.ParseFloat(params[14],64)
		check(err)
		tmp,err = strconv.ParseFloat(params[15],64)
		check(err)
		newOption.T = int(tmp)
		tmp,err = strconv.ParseFloat(params[16],64)
		check(err)
		newOption.N = int(tmp)


		options = append(options, newOption)
	}

	return options
}

type OptionURLReq struct {
	Ticker string
	Contract_type string
	ApiKey string
	StrikeRange []int
	DateRange []string
	//expDateTime, err := time.Parse(dateFormat,expDateStr)
	//check(err)
}

type Option struct {
	Cfi string
	Contract_type string
	Exerciese_style string
	Expiration_date string
	Primaty_exchange string
	Shares_per_contract int
	Strike_price float64
	Ticker string
	Underlying_ticker string
	Volume int
	Vw float64
	Open float64
	Close float64
	High float64
	Low float64
	T int
	N int
}

func (o Option) Print() string{
	readStr := fmt.Sprint(o)
	readStr = strings.Replace(readStr,"} {","\n",-1)
	readStr = strings.Replace(readStr,"}]","",-1)
	readStr = strings.Replace(readStr,"[{","",-1)
	return readStr
}

func GetOptions(optreq OptionURLReq, nMax int) ([]Option , string, error) {
	debug := false
	print := true

	if print {
		fmt.Println("Pulling options for option request:")
		fmt.Printf("ticker=%v\nstrikeRange=%v\nContract_type=%v\nDateRange=%v\n",optreq.Ticker,optreq.StrikeRange,optreq.Contract_type,optreq.DateRange)
	}

	if nMax == -1 {
		nMax = 10000
	}

	log := ""
	var msg string

	msg = fmt.Sprintln("optreq: ", optreq)
	log += msg

	// Get URL for option request
	optURL, err := URLoption(optreq)
	check(err)
	msg = fmt.Sprintln("optURL: ",optURL)
	log += msg

	var body string
	var res string
	var nextURL string = optURL


	var optionsStr []string

	requestCounter := 0
	optionCounter := 0

	var dataStr string
	var dataAr []string
	var n int = 0
	for ok := true ; ok ; ok = strings.Contains(body,optreq.Ticker) && n<=nMax {

		n++


		// Do next url request
		res, body, _ = APIRequest(nextURL,1)
		/*
		if err != nil {
			return nil, "", err
		}
		 */

		msg = fmt.Sprintln("response: ", res)
		log += msg

		// extract data
		dataStr = strings.Split(body,"\"results\":[")[1]
		dataStr = strings.Split(dataStr, "]")[0]
		dataAr = strings.Split(dataStr,"},{")
		dataAr[0] = strings.Replace(dataAr[0],"{","",-1)
		dataAr[len(dataAr)-1] = strings.Replace(dataAr[len(dataAr)-1],"}","",-1)

		// save dataAr into optionsStr
		for _,data := range dataAr {
			msg = fmt.Sprintln("Add to optionsStr: " , data)
			log += msg
			optionsStr = append(optionsStr,data)

			optionCounter++
			if print /*&& !debug*/ {
				fmt.Printf("\r %v API requests successfully made - %v available options found",requestCounter,optionCounter)
			}
		}



		// print response
		msg = fmt.Sprintln("res.Body:\n",body,"\n")
		log += msg
		if !strings.Contains(body,"next_url"){
			break
		}
		nextURL = strings.Split(body,"\"next_url\":")[1]
		nextURL = strings.Replace(nextURL,"\"","",-1)
		nextURL = strings.Replace(nextURL,"}","",-1)
		nextURL = strings.Replace(nextURL,"\n", "",-1)


		// filter out next url
		nextURL = strings.Split(body,"\"next_url\":")[1]
		nextURL = strings.Replace(nextURL,"\"","",-1)
		nextURL = strings.Replace(nextURL,"}","",-1)
		//nextURL += "&apiKey=***********" + optreq.ApiKey[len(optreq.ApiKey)-4:len(optreq.ApiKey)-1]
		urlPrint := nextURL + "&apiKey=***********" + optreq.ApiKey[len(optreq.ApiKey)-4:len(optreq.ApiKey)]
		nextURL += "&apiKey=" + optreq.ApiKey

		// filted out next url
		msg = fmt.Sprintln("nextURL:"+urlPrint)
		if debug{fmt.Println(msg)}
		log += msg

		requestCounter++
		if print /*&& !debug*/ {
			fmt.Printf("\r %v API requests successfully made - %v available options found",requestCounter,optionCounter)
		}

	}
	if print /*&& !debug*/ {
		fmt.Println("")
	}

	var options []Option
	var params []string
	//var tmp float64

	for _,opt := range optionsStr {
		opt = strings.Replace(opt,"\"","",-1)
		opt = strings.Replace(opt,"O:","",1)
		params = strings.Split(opt,",")
		for i,p := range params {
			params[i] = strings.Split(p,":")[1]
		}

		tmp,err := strconv.ParseFloat(params[5],64)
		check(err)
		spc := int(tmp)
		tmp, err = strconv.ParseFloat(params[6], 64)
		check(err)
		//sp := int(tmp)
		sp := tmp

		options = append(options, Option{
			Cfi:                 params[0],
			Contract_type:       params[1],
			Exerciese_style:     params[2],
			Expiration_date:     params[3],
			Primaty_exchange:    params[4],
			Shares_per_contract: spc,
			Strike_price:        sp,
			Ticker:              params[7],
			Underlying_ticker:   params[8],
		})

	}

	options = completeOptions(options,optreq.ApiKey)

	//filter strike_price
	var newOptions []Option

	//filter strike_range
	if debug{fmt.Println("strike_range: ",optreq.StrikeRange)}
	for _,opt := range options {
		if opt.Strike_price > float64(optreq.StrikeRange[0]) && opt.Strike_price < float64(optreq.StrikeRange[1]){
			newOptions = append(newOptions,opt)
		}
	}
	options = newOptions


	return options, log, nil

}

func completeOptions(options []Option, apiKey string) []Option {
	debug := false

	//var res, body string
	fmt.Println("There are ",len(options), " options to pull from the API.")
	if apiKey == "" {
		fmt.Println("With a free license this will take approx. ",len(options)/5, " minutes.")
	} else {
		//fmt.Println("With a premium license this will take approx. ",len(options)/5, " minutes.")
	}
	for j,opt := range options {
		url := "https://api.polygon.io/v2/aggs/ticker/O:"+opt.Ticker
		url += "/prev?adjusted=true&apiKey="+apiKey

		_, body, err := APIRequest(url,1)
		if err != nil {
			fmt.Println(err)
			if debug {fmt.Println("Removing option from list")}
			options[j] = Option{}
			continue
		}

		var dataStr string
		var dataAr []string

		// extract data
		dataStr = strings.Split(body,"\"results\":[{")[1]
		dataStr = strings.Split(dataStr, "}]")[0]
		dataStr = strings.Replace(dataStr,"O:","",1)
		dataAr = strings.Split(dataStr,",")
		for i,d := range dataAr {
			dataAr[i] = strings.Replace(strings.Split(d,":")[1],"\"","",-1)
		}

		tmp,err := strconv.ParseFloat(dataAr[1],64)
		check(err)
		options[j].Volume = int(tmp)
		options[j].Vw,err = strconv.ParseFloat(dataAr[2],64)
		check(err)
		options[j].Open,err = strconv.ParseFloat(dataAr[3],64)
		check(err)
		options[j].Close,err = strconv.ParseFloat(dataAr[4],64)
		check(err)
		options[j].High,err = strconv.ParseFloat(dataAr[5],64)
		check(err)
		options[j].Low,err = strconv.ParseFloat(dataAr[6],64)
		check(err)
		tmp,err = strconv.ParseFloat(dataAr[7],64)
		check(err)
		options[j].T = int(tmp)
		tmp,err = strconv.ParseFloat(dataAr[8],64)
		check(err)
		options[j].N = int(tmp)

		fmt.Printf("\r %.2f %% (%v out of %v) of available options pulled and converted",float64(j)/float64(len(options))*100,j,len(options))
	}
	fmt.Println("")

	//filter out empty options
	var newOptions []Option
	for _,opt := range options{
		if opt.Ticker != "" {
			newOptions = append(newOptions,opt)
		}
	}

	return newOptions

}

func APIRequest (url string, iteration int) (string,string,error) {
	debug := false
	//print := true

	//var requestCounter int = 0

	req, err := http.NewRequest("GET", url, nil)
	check(err)
	var res *http.Response
	res, err = http.DefaultClient.Do(req)
	check(err)
	retryNr := 1
	maxRetry := 12
	for res == nil {
		fmt.Println("Response is nil (possibly due to connection loss), waiting for 5 seconds and retrying (",retryNr,")")
		waitTime := 5 * time.Second
		time.Sleep( waitTime )
		res, err = http.DefaultClient.Do(req)
		check(err)
		retryNr++
		if retryNr >= maxRetry{
			fmt.Println("API Request Response is nil and kept being nil for ", int(waitTime.Seconds())*maxRetry ," seconds. Try establishing a better internet connection.")
			return "", "", err
		}
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if debug{
		fmt.Println("Request made:")
		fmt.Println(url)
		fmt.Println("Response and body:")
		fmt.Println(res,"\n", string(body))
	}

	//len(strings.Split(strings.Split(body,"\"results\":[{")[1], "}]")[0])



	//check if result in body is too short

	//check for error
	var errormsg string
	if strings.Contains(string(body),"ERROR") {
		//return "", "", fmt.Errorf("body empty - probably exceeded api request limit, please wait and try again.")
		errormsg = strings.Split(string(body),"\"error\":")[1]
		errormsg = strings.Split(errormsg,"}")[0]
		fmt.Println("An error occured: \n "+errormsg+"\nWaiting 60 seconds and retrying...")
		time.Sleep(60*time.Second)
		return APIRequest(url,1)
	}

	if len(strings.Split(string(body),"\"results\":[{"))<2 {
		if debug {fmt.Println("no result")}
		for iteration < 3 {
			if debug{fmt.Printf("ReRequesting in 500 milliseconds. That will be the %v%v reRequest.",iteration,stndrdth(iteration))}
			time.Sleep(500*time.Millisecond)
			return APIRequest(url, iteration+1)
		}
		return "", "", fmt.Errorf("no results")
	}

	//fmt.Println("\ndebug: ",strings.Split(strings.Split(string(body),"\"results\":")[1],"]")[0],len(strings.Split(strings.Split(string(body),"\"results\":")[1],"]")[0]))

	if len(strings.Split(strings.Split(string(body),"\"results\":")[1],"]")[0])<5{
		if debug {fmt.Println("no result")}
		for iteration < 3 {
			if debug {fmt.Println("ReRequesting in 500 milliseconds. That will be the ",iteration,stndrdth(iteration)," reRequest.")}
			time.Sleep(500*time.Millisecond)
			return APIRequest(url, iteration+1)
		}
		if debug {fmt.Printf("ReRequesting in 500 milliseconds. That will be the %v%v reRequest.",iteration,stndrdth(iteration))}
		time.Sleep(500*time.Millisecond)
		return APIRequest(url, iteration+1)
	}

	if debug{fmt.Println("API Request successfully made")}

	//fmt.Println("response: ", res)
	//fmt.Println("response body :", string(body))
	return fmt.Sprint(res) , string(body), nil
}

func stndrdth (n int) string {
	switch math.Mod(float64(n), 10) {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func URLoption(req OptionURLReq) (string, error) {
	if len(req.Ticker)==0 || len(req.ApiKey)==0{
		return "", fmt.Errorf("ticker or api missing")
	}

	req.Ticker = strings.ToUpper(req.Ticker)
	var url string
	url = "https://api.polygon.io/v3/reference/options/contracts"
	url += "?underlying_ticker="+req.Ticker
	url += "&apiKey="+req.ApiKey

	if len(req.StrikeRange) == 2 {
		url += "&strike-price.gte="+strconv.Itoa(req.StrikeRange[0])
		url += "&strike-price.lte="+strconv.Itoa(req.StrikeRange[1])
	} else if len (req.StrikeRange) == 1 {
		url += "&strike-price="+strconv.Itoa(req.StrikeRange[0])
	}

	if len(req.DateRange) == 2 {
		url += "&expiration_date.gte=" + req.DateRange[0]
		url += "&expiration_date.lte=" + req.DateRange[1]
	}
	if len(req.Contract_type)>0 {
		url += "&contract_type="+req.Contract_type
	}

	return url, nil
}