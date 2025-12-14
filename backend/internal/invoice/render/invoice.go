package render

type Invoice struct {
	Invoice    InvoiceMeta `json:"invoice"`
	Period     Period      `json:"period"`
	Consultant Consultant  `json:"consultant"`
	Company    Company     `json:"company"`
	Contract   Contract    `json:"contract"`
	Groups     []Group     `json:"groups"`
	Totals     Totals      `json:"totals"`
}
