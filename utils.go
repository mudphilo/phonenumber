package phonenumbers

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type MNCMCC struct {
	MNC int
	MCC int
}

type MCCMNCData struct {
	Mcc         string `json:"mcc"`
	Mnc         string `json:"mnc"`
	Iso         string `json:"iso"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Network     string `json:"network"`
	Msisdn     int64 `json:"msisdn"`
}

func GetISO3166ByCountryCode(country string) ISO3166 {

	iso3166 := ISO3166{}
	uppperCaseCountry := strings.ToUpper(country)

	for _, i := range GetISO3166() {

		if i.CountryCode == uppperCaseCountry {
			iso3166 = i
			break
		}
	}
	return iso3166
}

func IS_E164Compliant(msisdn string) bool {

	if !strings.HasPrefix(msisdn, "+") {

		msisdn = fmt.Sprintf("+%s", msisdn)
	}

	match, _ := regexp.MatchString(`^\+[1-9]\d{1,14}$`, msisdn)
	return match

}

func GetCountry(msisdn string) (ISO3166, error) {

	match, err := regexp.MatchString(`^\+[1-9]\d{1,14}$`, msisdn)

	if err != nil {

		if IS_E164Compliant(msisdn) {

			if !strings.HasPrefix(msisdn, "+") {

				msisdn = fmt.Sprintf("+%s", msisdn)
			}

		} else {

			fmt.Printf("got error doing regex %s  \n", err.Error())
			return ISO3166{
				Alpha2:             "",
				Alpha3:             "",
				CountryCode:        "",
				CountryName:        "",
				MobileBeginWith:    nil,
				PhoneNumberLengths: nil,
			}, err
		}
	}

	if !match {

		if IS_E164Compliant(msisdn) {

			if !strings.HasPrefix(msisdn, "+") {

				msisdn = fmt.Sprintf("+%s", msisdn)
			}

		} else {

			log.Printf("%s is invalid E164 number", msisdn)
			return ISO3166{
				Alpha2:             "",
				Alpha3:             "",
				CountryCode:        "",
				CountryName:        "",
				MobileBeginWith:    nil,
				PhoneNumberLengths: nil,
			}, err
		}
	}

	// check lengh
	//40721234567
	//+254726120256

	CountryCode := msisdn[:4]

	if len(msisdn) == 13 {

		CountryCode = msisdn[:4]

	} else if len(msisdn) == 12 {

		CountryCode = msisdn[:2]

	} else if len(msisdn) == 11 {

		CountryCode = msisdn[:2]

	}

	// get country code Alpha3
	CountryCode = strings.Replace(CountryCode, "+", "", -1)
	iso := GetISO3166ByCountryCode(CountryCode)
	return iso, nil

}

func GetCountryISO(msisdn string) string {

	match, err := regexp.MatchString(`^\+[1-9]\d{1,14}$`, msisdn)

	if err != nil {

		if IS_E164Compliant(msisdn) {

			if !strings.HasPrefix(msisdn, "+") {

				msisdn = fmt.Sprintf("+%s", msisdn)
			}

		} else {

			fmt.Printf("got error doing regex %s  \n", err.Error())
			return ""
		}
	}

	if !match {

		if IS_E164Compliant(msisdn) {

			if !strings.HasPrefix(msisdn, "+") {

				msisdn = fmt.Sprintf("+%s", msisdn)
			}

		} else {

			log.Printf("%s is invalid E164 number", msisdn)
			return ""
		}
	}

	// get country code Alpha3
	CountryCode := msisdn[:4]
	CountryCode = strings.Replace(CountryCode, "+", "", -1)
	CountryCode = GetISO3166ByCountryCode(CountryCode).Alpha2
	CountryCode = strings.ToUpper(CountryCode)
	return CountryCode

}

func GetMNCMCCFromIsoAndCarrier(iso, carrier string) (mcc int, mnc int,mycarrier string, err error)  {

	var payload []MCCMNCData
	err = json.Unmarshal([]byte(MNC_MCC_DATA),&payload)
	if err != nil {

		return 0,0,carrier, err
	}

	var mncmcc MCCMNCData

	for _, j := range payload {

		if strings.ToLower(j.Iso) == strings.ToLower(iso) {

			if strings.Contains(strings.ToLower(j.Network),strings.ToLower(carrier)) {

				mncmcc = j
				break
			}

			if strings.Contains(strings.ToLower(carrier),strings.ToLower(j.Network)) {

				mncmcc = j
				break
			}
		}
	}

	if len(carrier) == 0 {

		carrier = mncmcc.Network
	}

		mnc, _ = strconv.Atoi(mncmcc.Mnc)
	mcc, _ = strconv.Atoi(mncmcc.Mcc)
	return mcc,mnc,carrier, nil
}

func GetMNCMCCFromCountryCodeAndCarrier(countryCode, carrier string) (mcc int, mnc int, err error)  {

	var payload []MCCMNCData
	err = json.Unmarshal([]byte(MNC_MCC_DATA),&payload)
	if err != nil {

		return 0,0, err
	}

	var mncmcc MCCMNCData

	for _, j := range payload {

		if j.CountryCode == countryCode && ( strings.Contains(j.Network,carrier) || strings.Contains(carrier,j.Network)) {

			mncmcc = j
			break
		}
	}

	mnc, _ = strconv.Atoi(mncmcc.Mnc)
	mcc, _ = strconv.Atoi(mncmcc.Mcc)
	return mcc,mnc, nil
}

func GetMSISDN(msisdn string) MCCMNCData {

	// get country code
	Country, err := GetCountry(msisdn)
	if err != nil {

		return MCCMNCData{}
	}

	log.Printf("got CountryName %s ",Country.CountryName)
	log.Printf("got CountryCode %s ",Country.CountryCode)
	log.Printf("got Alpha2 %s ",Country.Alpha2)

	phonenumber, err := Parse(msisdn,strings.ToUpper(Country.Alpha2))
	if err != nil {

		return MCCMNCData{}
	}

	// get carrier
	carrier, err := GetCarrierForNumber(phonenumber,"en")
	if err != nil {

		return MCCMNCData{}
	}

	log.Printf("got carrier %s ",carrier)
	log.Printf("got carrier %s ",carrier)

	mcc, mnc,network, err := GetMNCMCCFromIsoAndCarrier(Country.Alpha2, carrier)
	log.Printf("got network %s ",network)

	international_format := Format(phonenumber, INTERNATIONAL)
	national_format := Format(phonenumber, NATIONAL)
	international_format = strings.ReplaceAll(international_format," ","")
	national_format = strings.ReplaceAll(national_format," ","")

	msisdns := strings.ReplaceAll(international_format,"+","")
	msisdns = strings.ReplaceAll(msisdns,")","")
	msisdns = strings.ReplaceAll(msisdns,"(","")
	msisdns = strings.ReplaceAll(msisdns,"-","")
	msisdns = strings.ReplaceAll(msisdns,".","")

	msisdnsFormat, _ := strconv.ParseInt(msisdns,10,64)
	cc, _ := strconv.ParseInt(Country.CountryCode,10,64)

	return MCCMNCData{
		Mcc:         fmt.Sprintf("%d",mcc),
		Mnc:         fmt.Sprintf("%d",mnc),
		Iso:         Country.Alpha2,
		Country:     Country.CountryName,
		CountryCode: fmt.Sprintf("%d",cc),
		Network:     network,
		Msisdn: msisdnsFormat,
	}
}