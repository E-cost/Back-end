package types

type ContactInfo struct {
	Name         string
	Surname      string
	Email        string
	ContactPhone string
}

type ContactConfirmationInfo struct {
	ContactId  string
	EmailId    string
	SecretCode string
}

type OrderOutput struct {
	EmailId    string
	SecretCode string
	OrderId    string
}

type EnterOrderItem struct {
	ID       string
	Category string
	Quantity int
}
