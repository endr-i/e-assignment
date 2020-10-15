package reportRepository

import (
	"assignment/entity"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type AccountTransactionsReportForm struct {
	Date      string
	AccountId string
}

type AccountTransactionsReportData struct {
	Account      *entity.Account
	Transactions []AccountTransactionsReportTransaction
	Since        time.Time
	Till         time.Time
	DateTime     time.Time
}

func (data *AccountTransactionsReportData) GetCSVData() [][]string {
	resultLength := len(data.Transactions) + 3
	result := make([][]string, resultLength)
	if data.Account != nil && data.Account.User != nil {
		result[0] = []string{
			fmt.Sprintf("User: %s", data.Account.User.Name),
			fmt.Sprintf("Report date: %s", data.DateTime.Format(time.RFC3339)),
			fmt.Sprintf(
				"Period: %s - %s",
				data.Since.Format(time.RFC3339),
				data.Till.Format(time.RFC3339),
			),
		}
	}
	result[1] = []string{
		"Type",
		"Value",
		"Date time",
		"Currency",
		"Currency rate",
		"Account rate",
	}
	resumeValue := decimal.NewFromInt(0)
	resumeUSDValue := decimal.NewFromInt(0)
	if data.Transactions != nil {
		for i, transaction := range data.Transactions {
			result[i+2] = []string{
				transaction.Type,
				transaction.DateTime.Format(time.RFC3339),
				transaction.Currency,
				transaction.CurrencyRate.Round(2).String(),
				transaction.AccountRate.Round(2).String(),
			}
			usdValue := transaction.Value.Mul(transaction.CurrencyRate)
			resumeValue.Add(usdValue.Div(transaction.AccountRate))
			resumeUSDValue.Add(usdValue)
		}
	}
	result[resultLength-1] = []string{
		fmt.Sprintf("Summary in %s", data.Account.Currency),
		resumeValue.Round(2).String(),
		"Summary in USD",
		resumeUSDValue.Round(2).String(),
	}
	return result
}

type AccountTransactionsReportTransaction struct {
	Type         string
	Value        decimal.Decimal
	DateTime     time.Time
	Currency     string
	CurrencyRate decimal.Decimal
	AccountRate  decimal.Decimal
}

func (reportTransaction AccountTransactionsReportTransaction) MarshalJSON() ([]byte, error) {
	value, _ := reportTransaction.Value.Round(2).Float64()
	currencyRate, _ := reportTransaction.CurrencyRate.Round(2).Float64()
	accountRate, _ := reportTransaction.AccountRate.Round(2).Float64()
	return json.Marshal(struct {
		Type         string
		Value        float64
		DateTime     time.Time
		Currency     string
		CurrencyRate float64
		AccountRate  float64
	}{
		Type:         reportTransaction.Type,
		Value:        value,
		DateTime:     reportTransaction.DateTime,
		Currency:     reportTransaction.Currency,
		CurrencyRate: currencyRate,
		AccountRate:  accountRate,
	})
}

func GetAccountTransactionsReportTransaction(t *entity.Transaction) AccountTransactionsReportTransaction {
	operationType := "Unknown"
	if t.Operation != nil {
		operationType = t.Operation.GetOperationType()
	}
	return AccountTransactionsReportTransaction{
		Type:         operationType,
		Value:        t.Value,
		DateTime:     t.DateTime,
		Currency:     t.Currency.Symbol,
		CurrencyRate: t.CurrencyRateValue,
		AccountRate:  t.AccountRateValue,
	}
}
