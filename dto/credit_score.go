package dto

type CreditScoreDto struct {
	OrderID     string `json:"OrderID" validate:"required"`
	CatID       string `json:"CatID_" validate:"required"`
	DP          int    `json:"DP_" validate:"required"`
	NTF         int    `json:"NTF_" validate:"required"`
	ANGS        int    `json:"ANGS_" validate:"required"`
	Tnr         int    `json:"Tnr_" validate:"required"`
	Gnr         string `json:"Gnr_" validate:"required"`
	Age         int    `json:"Age_" validate:"required"`
	Mrtl        string `json:"Mrtl_" validate:"required"`
	Dpnt        int    `json:"Dpnt_" validate:"required"`
	Hmsts       string `json:"Hmsts_" validate:"required"`
	Lnstay      int    `json:"Lnstay_" validate:"required"`
	Educ        string `json:"Educ_" validate:"required"`
	Prof        string `json:"Prof_" validate:"required"`
	Jobex       int    `json:"Jobex_" validate:"required"`
	Inc         int    `json:"Inc_" validate:"required"`
	VarInc      int    `json:"VarInc_" validate:"required"`
	SpsInc      int    `json:"SpsInc_" validate:"required"`
	Zip         string `json:"Zip_" validate:"required"`
	IIR         int    `json:"IIR_" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,indosat"`
}
