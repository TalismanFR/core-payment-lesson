package vo

type Amount uint64

type Currency string

const (
	USD Currency = "USD"
	RUB Currency = "RUB"
	BYN Currency = "BYN"
	UAH Currency = "UAH"
)

type Status string

const (
	StatusSuccessful = "SUCCESSFUL"
	StatusNew        = "NEW"
)
