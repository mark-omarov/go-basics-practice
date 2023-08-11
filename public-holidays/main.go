package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

type Country struct {
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type Holiday struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

const COUNTRIES_URL = "https://date.nager.at/api/v3/AvailableCountries"
const HOLIDAYS_URL = "https://date.nager.at/api/v3/PublicHolidays"

func main() {
	var countries []Country
	if err := fetchData(COUNTRIES_URL, &countries); err != nil {
		fmt.Println("Error fetching countries:", err)
		return
	}

	countryOptions := make([]string, len(countries))

	for i, country := range countries {
		countryOptions[i] = fmt.Sprintf("%s (%s)", country.Name, country.CountryCode)
	}

	prompt := &survey.Select{
		Message: "Choose a country:",
		Options: countryOptions,
	}
	var selectedIndex int
	survey.AskOne(prompt, &selectedIndex)

	selectedCountry := countries[selectedIndex]
	fmt.Printf("You selected: %s (%s)\n", selectedCountry.Name, selectedCountry.CountryCode)

	year := time.Now().Year()
	var holidays []Holiday
	if err := fetchData(fmt.Sprintf("%s/%d/%s", HOLIDAYS_URL, year, selectedCountry.CountryCode), &holidays); err != nil {
		fmt.Println("Error fetching public holidays:", err)
		return
	}

	fmt.Printf("Public holidays in %s (%s) in %d:\n", selectedCountry.Name, selectedCountry.CountryCode, year)
	for _, holiday := range holidays {
		fmt.Printf("  %s (%s)\n", holiday.Name, holiday.Date)
	}
}

func fetchData(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}
