package dto

type VendorRefundResult struct {
	vendor   string
	message  string
	reason   string
	isFailed bool
}

func (v VendorRefundResult) Vendor() string {
	return v.vendor
}

func NewVendorRefundResult(vendor, message, reason string, isFailed bool) *VendorRefundResult {
	return &VendorRefundResult{
		vendor:   vendor,
		message:  message,
		reason:   reason,
		isFailed: isFailed,
	}
}

func (v VendorRefundResult) IsFailed() bool {
	return true
}
