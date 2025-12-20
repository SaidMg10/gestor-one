package domain

type IncomeType string

const (
	IncomeTypeInvoice     IncomeType = "invoice"
	IncomeTypeReceipt     IncomeType = "receipt"
	IncomeTypeTransfer    IncomeType = "transfer"
	IncomeTypeDepositSlip IncomeType = "deposit_slip"
)

func IsValidIncomeType(t IncomeType) bool {
	switch t {
	case IncomeTypeInvoice,
		IncomeTypeReceipt,
		IncomeTypeTransfer,
		IncomeTypeDepositSlip:
		return true
	}
	return false
}
