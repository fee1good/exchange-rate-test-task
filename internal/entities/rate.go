package entities

type PairsRate struct {
	Raw     map[string]map[string]RateValues[float64] `json:"RAW"`
	Display map[string]map[string]RateValues[string]  `json:"DISPLAY"`
}

type Pair struct {
	CryptoSymbol      string
	FiatSymbol        string
	RawRateValues     RateValues[float64]
	DisplayRateValues RateValues[string]
}

type RateValues[T string | float64] struct {
	Change24Hour    T `json:"CHANGE24HOUR"`
	ChangePCT24Hour T `json:"CHANGEPCT24HOUR"`
	Open24Hour      T `json:"OPEN24HOUR"`
	Volume24Hour    T `json:"VOLUME24HOUR"`
	Volume24HourTo  T `json:"VOLUME24HOURTO"`
	Low24Hour       T `json:"LOW24HOUR"`
	High24Hour      T `json:"HIGH24HOUR"`
	Price           T `json:"PRICE"`
	Supply          T `json:"SUPPLY"`
	MKTCap          T `json:"MKTCAP"`
}
