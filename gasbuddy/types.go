package gasbuddy

type StationResponse struct {
	Data struct {
		Station struct {
			ID     string `json:"id"`
			Prices []struct {
				FuelProduct string `json:"fuelProduct"`
				LongName    string `json:"longName"`
				Credit      struct {
					Price          float32 `json:"price"`
					FormattedPrice string  `json:"formattedPrice"`
					PostedTime     string  `json:"postedTime"`
				} `json:"credit"`
			} `json:"prices"`
		} `json:"station"`
	} `json:"data"`
}
