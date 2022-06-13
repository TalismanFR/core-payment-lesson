package dto

type VendorChargeResult struct {
	vendor   string
	isFailed bool
	info     string
}

func (v VendorChargeResult) Vendor() string {
	return v.vendor
}

func NewVendorChargeResult(vendor string, status bool, info string) *VendorChargeResult {
	return &VendorChargeResult{vendor: vendor, isFailed: status, info: info}
}

func (v *VendorChargeResult) IsFailed() bool {
	return v.isFailed
}

func (v VendorChargeResult) Info() string {
	return v.info
}
