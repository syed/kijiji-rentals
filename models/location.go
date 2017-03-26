package models

type LocationResult struct {
	Results []struct {
		AddressComponents []struct {
			LongName string `json:"long_name"`
			ShortName string `json:"short_name"`
			Types []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string `json:"place_id"`
		Types []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}
