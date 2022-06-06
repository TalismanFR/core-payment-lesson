package dto

type VendorChargeResult struct {
	vendor string
}

func (v VendorChargeResult) Vendor() string {
	return v.vendor
}

func NewVendorChargeResult(vendor string) *VendorChargeResult {
	return &VendorChargeResult{vendor: vendor}
}

func (c *VendorChargeResult) IsFailed() bool {

	return true
}
