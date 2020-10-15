package operationRepository

type RefillForm struct {
	AccountId string
	Sum       float64
	Currency  string
	Details   struct {
		Source string
	} // Details of operation
}

type TransferForm struct {
	From     string
	To       string
	Sum      float64
	Currency string
	Details  struct {
		Comment string
	} // Details of operation
}
