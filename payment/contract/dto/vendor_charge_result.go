package dto

type VendorChargeResult struct {
	vendor string
	status bool
	info   string
}

func (v VendorChargeResult) Vendor() string {
	return v.vendor
}

func NewVendorChargeResult(vendor string, status bool, info string) *VendorChargeResult {
	return &VendorChargeResult{vendor: vendor, status: status, info: info}
}

func (v *VendorChargeResult) IsFailed() bool {
	return v.status
}

func (v VendorChargeResult) Info() string {
	return v.info
}
