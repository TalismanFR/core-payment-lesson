package dto

type ChargeResult struct {
	statusCode int
	statusName string
	uuid       string
	receiptUrl string
	message    string
	threeDs    *ThreeDs
}

func NewChargeResult(status int, statusName string, uuid string, receiptUrl string, message string, threeDs *ThreeDs) *ChargeResult {
	return &ChargeResult{statusCode: status, statusName: statusName, uuid: uuid, receiptUrl: receiptUrl, message: message, threeDs: threeDs}
}

func (c ChargeResult) Status() int {
	return c.statusCode
}

func (c ChargeResult) StatusName() string {
	return c.statusName
}

func (c ChargeResult) Uuid() string {
	return c.uuid
}

func (c ChargeResult) ReceiptUrl() string {
	return c.receiptUrl
}

func (c ChargeResult) Message() string {
	return c.message
}

func (c ChargeResult) ThreeDS() *ThreeDs {
	return c.threeDs
}

func (c ChargeResult) Need3Ds() bool {
	return c.threeDs != nil
}

type ThreeDs struct {
	Status      ThreeDsStatus
	RedirectUrl string
}

type ThreeDsStatus string

const (
	UnknownThreeDsStatus    ThreeDsStatus = "UNKNOWN"
	IncompleteThreeDsStatus ThreeDsStatus = "INCOMPLETE"
	SuccessfulThreeDsStatus ThreeDsStatus = "SUCCESSFUL"
	FailedThreeDsStatus     ThreeDsStatus = "FAILED"
)
