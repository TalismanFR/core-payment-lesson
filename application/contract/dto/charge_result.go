package dto

type ChargeResult struct {
	status     int
	statusName string
	uuid       string
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

func NewChargeResult(status int, statusName string, uuid string) *ChargeResult {
	return &ChargeResult{status: status, statusName: statusName, uuid: uuid}
}
