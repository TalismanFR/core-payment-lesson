package dto

type VendorRefundResult struct {
	vendor string
}

func (v VendorRefundResult) Vendor() string {
	return v.vendor
}

func NewVendorRefundResult(vendor string) *VendorRefundResult {
	return &VendorRefundResult{vendor: vendor}
}

func (v VendorRefundResult) IsFailed() bool {
	return true
}
