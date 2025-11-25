package tools

// LatLng represents a geographic location with latitude and longitude
type LatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ExtraComputation represents optional computation features
type ExtraComputation string

const (
	ExtraComputationLocalAQI                       ExtraComputation = "LOCAL_AQI"
	ExtraComputationHealthRecommendations          ExtraComputation = "HEALTH_RECOMMENDATIONS"
	ExtraComputationPollutantAdditionalInfo        ExtraComputation = "POLLUTANT_ADDITIONAL_INFO"
	ExtraComputationDominantPollutantConcentration ExtraComputation = "DOMINANT_POLLUTANT_CONCENTRATION"
	ExtraComputationPollutantConcentration         ExtraComputation = "POLLUTANT_CONCENTRATION"
)

// ColorPalette represents the color palette for UAQI
type ColorPalette string

const (
	ColorPaletteRedGreen      ColorPalette = "RED_GREEN"
	ColorPaletteIndigoPersian ColorPalette = "INDIGO_PERSIAN"
	ColorPaletteNumeric       ColorPalette = "NUMERIC"
)

// MapType represents the type of air quality heatmap
type MapType string

const (
	MapTypeUAQIRedGreen      MapType = "UAQI_RED_GREEN"
	MapTypeUAQIIndigoPersian MapType = "UAQI_INDIGO_PERSIAN"
	MapTypePM25IndigoPersian MapType = "PM25_INDIGO_PERSIAN"
	MapTypeGBRDefra          MapType = "GBR_DEFRA"
	MapTypeDEUUba            MapType = "DEU_UBA"
	MapTypeCANEc             MapType = "CAN_EC"
	MapTypeFRAAtmo           MapType = "FRA_ATMO"
	MapTypeUSAQI             MapType = "US_AQI"
)

// CustomLocalAqi represents a custom AQI configuration for a specific region
type CustomLocalAqi struct {
	RegionCode string `json:"regionCode"`
	Aqi        string `json:"aqi"`
}

// Interval represents a time period
type Interval struct {
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
}

// CurrentConditionsRequest represents the request for current air quality conditions
type CurrentConditionsRequest struct {
	Location          LatLng             `json:"location"`
	ExtraComputations []ExtraComputation `json:"extraComputations,omitempty"`
	UaqiColorPalette  ColorPalette       `json:"uaqiColorPalette,omitempty"`
	CustomLocalAqis   []CustomLocalAqi   `json:"customLocalAqis,omitempty"`
	UniversalAqi      *bool              `json:"universalAqi,omitempty"`
	LanguageCode      string             `json:"languageCode,omitempty"`
}

// ForecastRequest represents the request for air quality forecast
type ForecastRequest struct {
	Location          LatLng             `json:"location"`
	ExtraComputations []ExtraComputation `json:"extraComputations,omitempty"`
	UaqiColorPalette  ColorPalette       `json:"uaqiColorPalette,omitempty"`
	CustomLocalAqis   []CustomLocalAqi   `json:"customLocalAqis,omitempty"`
	PageSize          int                `json:"pageSize,omitempty"`
	PageToken         string             `json:"pageToken,omitempty"`
	DateTime          string             `json:"dateTime,omitempty"`
	Period            *Interval          `json:"period,omitempty"`
	UniversalAqi      *bool              `json:"universalAqi,omitempty"`
	LanguageCode      string             `json:"languageCode,omitempty"`
}

// HistoryRequest represents the request for historical air quality data
type HistoryRequest struct {
	PageSize          int                `json:"pageSize,omitempty"`
	PageToken         string             `json:"pageToken,omitempty"`
	Location          LatLng             `json:"location"`
	ExtraComputations []ExtraComputation `json:"extraComputations,omitempty"`
	UaqiColorPalette  ColorPalette       `json:"uaqiColorPalette,omitempty"`
	CustomLocalAqis   []CustomLocalAqi   `json:"customLocalAqis,omitempty"`
	DateTime          string             `json:"dateTime,omitempty"`
	Hours             int                `json:"hours,omitempty"`
	Period            *Interval          `json:"period,omitempty"`
	UniversalAqi      *bool              `json:"universalAqi,omitempty"`
	LanguageCode      string             `json:"languageCode,omitempty"`
}

// AQI represents an Air Quality Index
type AQI struct {
	Code        string `json:"code,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Aqi         int    `json:"aqi,omitempty"`
	AqiDisplay  string `json:"aqiDisplay,omitempty"`
	Color       *Color `json:"color,omitempty"`
	Category    string `json:"category,omitempty"`
}

// Color represents an RGB color
type Color struct {
	Red   float32 `json:"red,omitempty"`
	Green float32 `json:"green,omitempty"`
	Blue  float32 `json:"blue,omitempty"`
	Alpha float32 `json:"alpha,omitempty"`
}

// Pollutant represents a pollutant measurement
type Pollutant struct {
	Code           string         `json:"code,omitempty"`
	DisplayName    string         `json:"displayName,omitempty"`
	FullName       string         `json:"fullName,omitempty"`
	Concentration  *Concentration `json:"concentration,omitempty"`
	AdditionalInfo string         `json:"additionalInfo,omitempty"`
}

// Concentration represents pollutant concentration
type Concentration struct {
	Value float64 `json:"value,omitempty"`
	Units string  `json:"units,omitempty"`
}

// HealthRecommendations represents health recommendations
type HealthRecommendations struct {
	GeneralPopulation      string `json:"generalPopulation,omitempty"`
	Elderly                string `json:"elderly,omitempty"`
	LungDiseasePopulation  string `json:"lungDiseasePopulation,omitempty"`
	HeartDiseasePopulation string `json:"heartDiseasePopulation,omitempty"`
	Athletes               string `json:"athletes,omitempty"`
	PregnantWomen          string `json:"pregnantWomen,omitempty"`
	Children               string `json:"children,omitempty"`
}

// CurrentConditionsResponse represents the response for current air quality conditions
type CurrentConditionsResponse struct {
	DateTime              string                 `json:"dateTime,omitempty"`
	RegionCode            string                 `json:"regionCode,omitempty"`
	Indexes               []AQI                  `json:"indexes,omitempty"`
	Pollutants            []Pollutant            `json:"pollutants,omitempty"`
	HealthRecommendations *HealthRecommendations `json:"healthRecommendations,omitempty"`
}

// HourlyForecast represents an hourly forecast
type HourlyForecast struct {
	DateTime              string                 `json:"dateTime,omitempty"`
	Indexes               []AQI                  `json:"indexes,omitempty"`
	Pollutants            []Pollutant            `json:"pollutants,omitempty"`
	HealthRecommendations *HealthRecommendations `json:"healthRecommendations,omitempty"`
}

// ForecastResponse represents the response for air quality forecast
type ForecastResponse struct {
	RegionCode      string           `json:"regionCode,omitempty"`
	HourlyForecasts []HourlyForecast `json:"hourlyForecasts,omitempty"`
	NextPageToken   string           `json:"nextPageToken,omitempty"`
}

// HourInfo represents historical hourly information
type HourInfo struct {
	DateTime              string                 `json:"dateTime,omitempty"`
	Indexes               []AQI                  `json:"indexes,omitempty"`
	Pollutants            []Pollutant            `json:"pollutants,omitempty"`
	HealthRecommendations *HealthRecommendations `json:"healthRecommendations,omitempty"`
}

// HistoryResponse represents the response for historical air quality data
type HistoryResponse struct {
	RegionCode    string     `json:"regionCode,omitempty"`
	HoursInfo     []HourInfo `json:"hoursInfo,omitempty"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}
