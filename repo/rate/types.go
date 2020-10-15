package rateRepository

type CreateForm struct {
	Symbol string
	Value  float64
}

type UploadRatesForm struct {
	Rates []CreateForm
}

//type ConvertRate struct {
//	RateFrom  entity.Rate
//	RateTo    entity.Rate
//	RateValue decimal.Decimal
//}
