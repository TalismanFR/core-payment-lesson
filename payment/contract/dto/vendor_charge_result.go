package dto

import "fmt"

type VendorChargeResult struct {
	vendor        string
	transactionId string
	message       string
	status        VendorChargeStatus
	receiptUrl    string
	threeDs       *VendorThreeDs
}

func (v VendorChargeResult) ThreeDs() *VendorThreeDs {
	return v.threeDs
}

func (v VendorChargeResult) Message() string {
	return v.message
}

func (v VendorChargeResult) ReceiptUrl() string {
	return v.receiptUrl
}

func NewVendorChargeResult(vendor string, transactionId string, message string, status VendorChargeStatus, receiptUrl string, threeDs *VendorThreeDs) *VendorChargeResult {
	return &VendorChargeResult{vendor: vendor, transactionId: transactionId, message: message, status: status, receiptUrl: receiptUrl, threeDs: threeDs}
}

var (
	UnknownVendorChargeStatus    = VendorChargeStatus{""}
	SuccessfulVendorChargeStatus = VendorChargeStatus{"SUCCESSFUL"}
	FailedVendorChargeStatus     = VendorChargeStatus{"FAILED"}
	Need3DSVendorChargeStatus    = VendorChargeStatus{"NEED3DS"}
)

type VendorChargeStatus struct {
	s string
}

func (cur VendorChargeStatus) String() string {
	return cur.s
}

func FromString(s string) (VendorChargeStatus, error) {
	switch s {
	case SuccessfulVendorChargeStatus.String():
		return SuccessfulVendorChargeStatus, nil
	case FailedVendorChargeStatus.String():
		return FailedVendorChargeStatus, nil
	case Need3DSVendorChargeStatus.String():
		return Need3DSVendorChargeStatus, nil

	}

	return UnknownVendorChargeStatus, fmt.Errorf("unknown vendor status: %s", s)
}

func (v VendorChargeResult) Vendor() string {
	return v.vendor
}

func (c *VendorChargeResult) IsFailed() bool {

	return true
}

type VendorThreeDs struct {
	Status      ThreeDsVendorChargeStatus
	RedirectUrl string
}

type ThreeDsVendorChargeStatus string

const (
	UnknownThreeDsVendorStatus    ThreeDsVendorChargeStatus = "unknown"
	IncompleteThreeDsVendorStatus ThreeDsVendorChargeStatus = "incomplete"
)
