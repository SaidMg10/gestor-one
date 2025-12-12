package domain

type ReceiptType string

const (
	ReceiptTypeInvoice     ReceiptType = "invoice"
	ReceiptTypeReceipt     ReceiptType = "receipt"
	ReceiptTypeTransfer    ReceiptType = "transfer"
	ReceiptTypeDepositSlip ReceiptType = "deposit_slip"
)

func IsValidReceiptType(t ReceiptType) bool {
	switch t {
	case ReceiptTypeInvoice,
		ReceiptTypeReceipt,
		ReceiptTypeTransfer,
		ReceiptTypeDepositSlip:
		return true
	}
	return false
}
