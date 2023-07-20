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

func GetISO3166ByCountryCode(countryDialingCode int64) ISO3166 {

	iso3166 := ISO3166{}

	for _, i := range GetISO3166() {

		dialingCode, _ := strconv.ParseInt(i.CountryCode,10,64)
		if dialingCode == countryDialingCode {
			iso3166 = i
			break
		}
	}
	return iso3166
}

func IsE164compliant(msisdn string) bool {

	if len(msisdn) < 11 {

		return false
	}

	if !strings.HasPrefix(msisdn, "+") {

		msisdn = fmt.Sprintf("+%s", msisdn)
	}

	match, _ := regexp.MatchString(`^\+[1-9]\d{1,14}$`, msisdn)
	return match

}

func GetCountry(CountryCode int64) (ISO3166, error) {

	// get country code Alpha3
	iso := GetISO3166ByCountryCode(CountryCode)
	return iso, nil
}

func GetCountryFromMsisdn(msisdn string) (ISO3166, error) {

	match, err := regexp.MatchString(`^\+[1-9]\d{1,14}$`, msisdn)

	if err != nil {

		if IsE164compliant(msisdn) {

			if !strings.HasPrefix(msisdn, "+") {

				msisdn = fmt.Sprintf("+%s", msisdn)
			}

		} else {

			fmt.Printf("got error doing regex %s  \n", err.Error())
			return ISO3166{
				Alpha2:             "KE",
				Alpha3:             "KEN",
				CountryCode:        "254",
				CountryName:        "KENYA",
				MobileBeginWith:    nil,
				PhoneNumberLengths: nil,
			}, err
		}
	}

	if !match {

		if IsE164compliant(msisdn) {

			if !strings.HasPrefix(msisdn, "+") {

				msisdn = fmt.Sprintf("+%s", msisdn)
			}

		} else {

			log.Printf("%s is invalid E164 number", msisdn)

			return ISO3166{
				Alpha2:             "KE",
				Alpha3:             "KEN",
				CountryCode:        "254",
				CountryName:        "KENYA",
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
	dialingCode, _ := strconv.ParseInt(CountryCode,10,64)
	iso := GetISO3166ByCountryCode(dialingCode)
	return iso, nil

}

func GetCountryISO(msisdn string) string {

	match, err := regexp.MatchString(`^\+[1-9]\d{1,14}$`, msisdn)

	if err != nil {

		if IsE164compliant(msisdn) {

			if !strings.HasPrefix(msisdn, "+") {

				msisdn = fmt.Sprintf("+%s", msisdn)
			}

		} else {

			fmt.Printf("got error doing regex %s  \n", err.Error())
			return ""
		}
	}

	if !match {

		if IsE164compliant(msisdn) {

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
	dialingCode, _ := strconv.ParseInt(CountryCode,10,64)
	CountryCode = GetISO3166ByCountryCode(dialingCode).Alpha2
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

	// harmonize phone numbers

	// get country code
	Country, err := GetCountryFromMsisdn(msisdn)
	if err != nil {

		return MCCMNCData{}
	}

	phonenumber, err := Parse(msisdn,strings.ToUpper(Country.Alpha2))
	if err != nil {

		return MCCMNCData{}
	}

	// get carrier
	carrier, err := GetCarrierForNumber(phonenumber,"en")
	if err != nil {

		return MCCMNCData{}
	}

	mcc, mnc,network, err := GetMNCMCCFromIsoAndCarrier(Country.Alpha2, carrier)

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

func GetMSISDNWithCountryCode(msisdn, countryCode string) MCCMNCData {

	msisdn = ParseOld(msisdn, countryCode)

	// get country code
	Country, err := GetCountryFromMsisdn(msisdn)
	if err != nil {

		return MCCMNCData{}
	}

	phonenumber, err := Parse(msisdn,strings.ToUpper(Country.Alpha2))
	if err != nil {

		return MCCMNCData{}
	}

	// get carrier
	carrier, err := GetCarrierForNumber(phonenumber,"en")
	if err != nil {

		return MCCMNCData{}
	}

	mcc, mnc,network, err := GetMNCMCCFromIsoAndCarrier(Country.Alpha2, carrier)

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

func GetMsisdnWithDialingCode(dialingCode, msisdn int64) MCCMNCData {

	// get country code
	Country, err := GetCountry(dialingCode)
	if err != nil {

		return MCCMNCData{}
	}

	phonenumber, err := Parse(fmt.Sprintf("%d",msisdn),strings.ToUpper(Country.Alpha2))
	if err != nil {

		return MCCMNCData{}
	}

	// get carrier
	carrier, err := GetCarrierForNumber(phonenumber,"en")
	if err != nil {

		return MCCMNCData{}
	}

	mcc, mnc,network, err := GetMNCMCCFromIsoAndCarrier(Country.Alpha2, carrier)

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
