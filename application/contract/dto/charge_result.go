package dto

type ChargeResult struct {
	status     int
	statusName string
	uuid       string
	threeDs    *ThreeDs
}

func NewChargeResult(status int, statusName string, uuid string) *ChargeResult {
	return &ChargeResult{status: status, statusName: statusName, uuid: uuid}
}

func (c ChargeResult) Status() int {
	return c.status
}

func (c ChargeResult) StatusName() string {
	return c.statusName
}

func (c ChargeResult) Uuid() string {
	return c.uuid
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
	ThreeDsStatusIncomplete ThreeDsStatus = "incomplete"
)
