package domain

type ExpenseType string

const (
	ExpenseTypeOperational    ExpenseType = "operational"
	ExpenseTypeAdministrative ExpenseType = "administrative"
	ExpenseTypePersonal       ExpenseType = "personal"
	ExpenseTypeExtraordinary  ExpenseType = "extraordinary"
)

func IsValidExpenseType(t ExpenseType) bool {
	switch t {
	case ExpenseTypeOperational,
		ExpenseTypeAdministrative,
		ExpenseTypePersonal,
		ExpenseTypeExtraordinary:
		return true
	}
	return false
}
