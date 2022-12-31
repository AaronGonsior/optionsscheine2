package main

import (
	"fmt"
	//"net/http"
	//"io/ioutil"
	"encoding/json"
	"os"
	"strings"
	"strconv"
	//"time"
	opt "github.com/AaronGonsior/optionsscheine2"
	//opt "./optionsscheine2.go"
)





func check(err error){
	if err!=nil{
		fmt.Println(err)
	}
}


/*
type OptionURLReq struct {
	ticker string
	apiKey string
	strikeRange []int
	dateRange []time.Time
	//expDateTime, err := time.Parse(dateFormat,expDateStr)
	//check(err)
}
 */


func main() {

	//fmt.Println(opt.TestFunc())

	var err error


	// exe param handling
	args := os.Args
	api := 4

	var update bool = false
	var optionOptions bool = false

	expDateStr := "2023-06-15"

	//fmt.Printf("Program arguments (%v): %v\n",len(args),args)
	for _, s := range args {
		if s=="-u"{
			update = true
		}
		if s=="--update"{
			update = true
		}
		if strings.Contains(s,"api"){
			api, err = strconv.Atoi(strings.Split(s,"=")[1])
			check(err)
		}
		if s=="-o" || s=="--options" {
			optionOptions = true
		}
		if strings.Contains(s,"exp"){
			expDateStr = strings.Split(s,"=")[1]
		}
	}

	load := !update

	fmt.Print("\n")
	fmt.Printf("%v ",args[0])
	fmt.Printf("update=%t ",update)
	fmt.Printf("load=%v ",load)
	fmt.Printf("api=%v ",api)
	fmt.Printf("options=%v ",optionOptions)
	fmt.Printf("exp=%v ",expDateStr)
	fmt.Print("\n\n")

	ticker := "TSLA"
	nMax := -1
	apiKey := opt.LoadJson("apiKey.json")
	fmt.Println("apiKey loaded:", apiKey)


	var optreq opt.OptionURLReq
	var options []opt.Option

	optreq = opt.OptionURLReq{
		Ticker:      ticker,
		ApiKey:      apiKey,
		StrikeRange: []int{50,200},
		DateRange:   []string{"2025-01-01","2026-01-01"},
	}


	if update {

		log := ""
		var msg string
		//var optionsStr []string


		options, msg = opt.GetOptions(optreq,nMax)
		log += msg

		/*
		msg = fmt.Sprintln("optreq: ", optreq)
		log += msg

		// Get URL for option request
		optURL, err := opt.URLoption(optreq)
		check(err)
		msg = fmt.Sprintln("optURL: ",optURL)
		log += msg

		var body string
		var res string
		var nextURL string = optURL


		var optionsStr []string

		var dataStr string
		var dataAr []string
		var n int = 0
		nMax := 3
		for ok := true ; ok ; ok = strings.Contains(body,ticker) && n<=nMax {

			n++


			// Do next url request
			res, body = opt.APIRequest(nextURL)
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
			}


			// print response
			msg = fmt.Sprintln("res.Body:\n",body,"\n")
			log += msg
			nextURL = strings.Split(body,"\"next_url\":")[1]
			nextURL = strings.Replace(nextURL,"\"","",-1)
			nextURL = strings.Replace(nextURL,"}","",-1)

			// filter out next url
			msg = fmt.Sprintln("nextURL:"+nextURL)
			log += msg


			// filter out next url
			nextURL = strings.Split(body,"\"next_url\":")[1]
			nextURL = strings.Replace(nextURL,"\"","",-1)
			nextURL = strings.Replace(nextURL,"}","",-1)
			msg = fmt.Sprintln("nextURL:"+nextURL)
			log += msg
			nextURL += "&apiKey=" + apiKey

		}
		 */

		msg = "optionsStr\n"
		log += msg
		for _, opt := range options {
			fmt.Println(opt)
			msg = fmt.Sprintln(opt)
			log += msg
		}




		// Open a file for writing
		file, err := os.Create("log.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		// Encode the string as JSON and write it to the file
		if err := json.NewEncoder(file).Encode(log); err != nil {
			fmt.Println(err)
			return
		}



		file, err = os.Create("options.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		// Encode the string as JSON and write it to the file
		if err := json.NewEncoder(file).Encode(fmt.Sprint(options)); err != nil {
			fmt.Println(err)
			return
		}

	}


	if !update{

		readStr := opt.LoadJson("options.json")
		readStr = strings.Replace(readStr,"} {","\n",-1)
		readStr = strings.Replace(readStr,"}]","",-1)
		readStr = strings.Replace(readStr,"[{","",-1)
		fmt.Println(readStr)

		options = opt.JsonToOptions("options.json")
		fmt.Println("loaded options: \n",options)
	}


	//filter - only for update for now

	/*
	var newOptions []opt.Option


	//filter strike_range
	fmt.Println("strike_range: ",optreq.StrikeRange)
	for _,opt := range options {
		if opt.Strike_price > optreq.StrikeRange[0] && opt.Strike_price < optreq.StrikeRange[1]{
			newOptions = append(newOptions,opt)
		}
	}
	options = newOptions


	fmt.Println(strings.Replace(fmt.Sprintln(newOptions),"} {","\n",-1))

	 */







}


/*
	switch api {
	case 1:
		req.Header.Add("X-RapidAPI-Proxy-Secret", "a755b180-f5a9-11e9-9f69-7bf51e845926")
		req.Header.Add("X-RapidAPI-Key", "afc3011b42mshe0c588092139f13p1ad2cbjsn72cd7b9e9125")
		req.Header.Add("X-RapidAPI-Host", "stock-and-options-trading-data-provider.p.rapidapi.com")
	case 2:
		req.Header.Add("X-RapidAPI-Key", "afc3011b42mshe0c588092139f13p1ad2cbjsn72cd7b9e9125")
		req.Header.Add("X-RapidAPI-Host", "real-time-finance-data.p.rapidapi.com")
	case 3:
		req.Header.Add("X-RapidAPI-Key", "afc3011b42mshe0c588092139f13p1ad2cbjsn72cd7b9e9125")
		req.Header.Add("X-RapidAPI-Host", "yahoo-finance15.p.rapidapi.com")
	}
*/




/*
	switch api {
	case 1:
		apiInfo = "Stock and Options Trading Data Provider - By Sam Johnson"
		url = "https://stock-and-options-trading-data-provider.p.rapidapi.com/options/"+ticker

	case 2:
		apiInfo = "Real-Time Finance Data - By OpenWeb Ninja"
		url = "https://real-time-finance-data.p.rapidapi.com/search?query="+stockname

	case 3:
		apiInfo = "Yahoo Finance - By API Datacenter"
		if !options{
			url = "https://yahoo-finance15.p.rapidapi.com/api/yahoo/qu/quote/"+ticker
		} else {
			url = "https://yahoo-finance15.p.rapidapi.com/api/yahoo/op/option/"+ticker
			if expDateStr!=""{
				url += "?expiration="+strconv.FormatInt(expDateUnix,10)
			}
		}

	case 4:
		apiInfo = "Polygon.io API"


		//apiKey :=
		//url = "https://api.polygon.io/v3/reference"


		//url+= "/options/contracts"
		//url+= "?underlying_ticker="+ticker
		//url+="&contract_type=call"
		//url+="&expiration_date.lte="+expDateStr


		//url+="&apiKey="+apiKey

		//url="https://api.polygon.io/v3/reference/tickers?active=true&apiKey="

		url="https://api.polygon.io/v3/reference/options/contracts?underlying_ticker=TSLA&apiKey="

	}
*/


/*
	expDateTime, err := time.Parse("2006-01-02",expDateStr)
	check(err)
	expDateRFC3339 := expDateTime.Format(time.RFC3339)
	expDateUnix := expDateTime.Unix()

	//fmt.Printf("%v\n%v\n%v\n%v\n",expDateStr,expDateTime,expDateRFC3339,expDateUnix)

	var apiInfo string
	var url string

	apiInfo = "Polygon.io API"
	//url="https://api.polygon.io/v3/reference/options/contracts?underlying_ticker="+ticker+"&apiKey="
	url="https://api.polygon.io/v3/reference/options/contracts?underlying_ticker="+ticker+"&apiKey=

	fmt.Println("\napiInfo: "+apiInfo)
		fmt.Println("api call: "+url+"\n")
*/
