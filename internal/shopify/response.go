package shopify

type AddressValidationResponse struct {
	Result struct {
		Verdict struct {
			InputGranularity      string `json:"inputGranularity"`
			ValidationGranularity string `json:"validationGranularity"`
			GeocodeGranularity    string `json:"geocodeGranularity"`
			AddressComplete       bool   `json:"addressComplete"`
			HasInferredComponents bool   `json:"hasInferredComponents"`
		} `json:"verdict"`
		Address struct {
			FormattedAddress string `json:"formattedAddress"`
			PostalAddress    struct {
				RegionCode         string   `json:"regionCode"`
				LanguageCode       string   `json:"languageCode"`
				PostalCode         string   `json:"postalCode"`
				AdministrativeArea string   `json:"administrativeArea"`
				Locality           string   `json:"locality"`
				AddressLines       []string `json:"addressLines"`
			} `json:"postalAddress"`
			AddressComponents []struct {
				ComponentName struct {
					Text         string `json:"text"`
					LanguageCode string `json:"languageCode,omitempty"`
				} `json:"componentName"`
				ComponentType     string `json:"componentType"`
				ConfirmationLevel string `json:"confirmationLevel"`
				Inferred          bool   `json:"inferred,omitempty"`
			} `json:"addressComponents"`
		} `json:"address"`
		Geocode struct {
			Location struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"location"`
			PlusCode struct {
				GlobalCode string `json:"globalCode"`
			} `json:"plusCode"`
			Bounds struct {
				Low struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"low"`
				High struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"high"`
			} `json:"bounds"`
			PlaceId    string   `json:"placeId"`
			PlaceTypes []string `json:"placeTypes"`
		} `json:"geocode"`
		Metadata struct {
			Business    bool `json:"business"`
			Residential bool `json:"residential"`
		} `json:"metadata"`
		UspsData struct {
			StandardizedAddress struct {
				FirstAddressLine        string `json:"firstAddressLine"`
				CityStateZipAddressLine string `json:"cityStateZipAddressLine"`
				City                    string `json:"city"`
				State                   string `json:"state"`
				ZipCode                 string `json:"zipCode"`
				ZipCodeExtension        string `json:"zipCodeExtension"`
			} `json:"standardizedAddress"`
			DeliveryPointCode       string `json:"deliveryPointCode"`
			DeliveryPointCheckDigit string `json:"deliveryPointCheckDigit"`
			DpvConfirmation         string `json:"dpvConfirmation"`
			DpvFootnote             string `json:"dpvFootnote"`
			DpvCmra                 string `json:"dpvCmra"`
			DpvVacant               string `json:"dpvVacant"`
			DpvNoStat               string `json:"dpvNoStat"`
			CarrierRouteIndicator   string `json:"carrierRouteIndicator"`
			PostOfficeCity          string `json:"postOfficeCity"`
			PostOfficeState         string `json:"postOfficeState"`
			FipsCountyCode          string `json:"fipsCountyCode"`
			County                  string `json:"county"`
			ElotNumber              string `json:"elotNumber"`
			ElotFlag                string `json:"elotFlag"`
			AddressRecordType       string `json:"addressRecordType"`
			DpvNoStatReasonCode     int    `json:"dpvNoStatReasonCode"`
			DpvDrop                 string `json:"dpvDrop"`
			DpvThrowback            string `json:"dpvThrowback"`
			DpvNonDeliveryDays      string `json:"dpvNonDeliveryDays"`
			DpvNoSecureLocation     string `json:"dpvNoSecureLocation"`
			DpvPbsa                 string `json:"dpvPbsa"`
			DpvDoorNotAccessible    string `json:"dpvDoorNotAccessible"`
			DpvEnhancedDeliveryCode string `json:"dpvEnhancedDeliveryCode"`
		} `json:"uspsData"`
	} `json:"result"`
	ResponseId string `json:"responseId"`
}
