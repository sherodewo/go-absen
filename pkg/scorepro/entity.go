package scorepro

import "time"

type PickleModelingDto struct {
	SupplierID       string      `json:"supplier_id"`
	ProspectID       string      `json:"prospect_id" validate:"required"`
	CbFound          *bool       `json:"cb_found" validate:"required"`
	StatusKonsumen   string      `json:"status_konsumen"`
	Journey          string      `json:"journey" validate:"required"`
	PhoneNumber      string      `json:"phone_number" validate:"required,number"`
	RequestorID      string      `json:"requestor_id" validate:"required"`
	ScoreGeneratorID string      `json:"score_generator_id" validate:"required"`
	TransactionType  string      `json:"transaction_type,omitempty"`
	Data             interface{} `json:"data"`
}

type PickleModelingKmbDto struct {
	SupplierID       string      `json:"supplier_id" validate:"required"`
	ProspectID       string      `json:"prospect_id" validate:"required"`
	CbFound          *bool       `json:"cb_found" validate:"required"`
	StatusKonsumen   string      `json:"status_konsumen"`
	PhoneNumber      string      `json:"phone_number" validate:"required,number"`
	RequestorID      string      `json:"requestor_id" validate:"required"`
	ScoreGeneratorID string      `json:"score_generator_id" validate:"required"`
	Data             interface{} `json:"data"`
}

type ScoreproResponse struct {
	ProspectID  string   `json:"prospect_id"`
	Score       *float64 `json:"score"`
	Result      string   `json:"result"`
	ScoreResult string   `json:"score_result"`
	Status      string   `json:"status"`
	PhoneNumber string   `json:"phone_number"`
}

type WgScoreproResponse struct {
	ProspectID  string   `json:"prospect_id"`
	Score       *float64 `json:"score"`
	Result      string   `json:"result"`
	MaxDSR      int      `json:"max_dsr"`
	ScoreBand   string   `json:"score_band"`
	ScoreResult string   `json:"score_result"`
	Status      string   `json:"status"`
	PhoneNumber string   `json:"phone_number"`
}

type PickleResponseIDX struct {
	Messages string      `json:"messages"`
	Errors   interface{} `json:"errors"`
	Data     struct {
		ProspectID string  `json:"transaction_id"`
		Score      float64 `json:"score"`
		Result     string  `json:"result"`
	} `json:"data"`
	ServerTime time.Time `json:"server_time"`
}

type PickleJJ struct {
	TransactionID   string `json:"transaction_id"  validate:"required"`
	ZipCode         int    `json:"zip_code"  validate:"required,number,min=2"`
	CategoryID      string `json:"category_id"  validate:"required"`
	LengthOfEmploy  *int   `json:"length_of_employ"  validate:"required,number"`
	OldestMobBanks  *int   `json:"oldest_mob_banks"  validate:"required,number"`
	FirstFourOfCell string `json:"first_four_of_cell" validate:"required,min=4"`
	Gender          string `json:"gender" validate:"required"`
	MaritalStatus   string `json:"marital_status" validate:"required"`
	ProfessionID    string `json:"profession_id" validate:"required"`
	MaxLimitPlAll   *int   `json:"max_limit_pl_all" validate:"required,number"`
	MonthlyIncome   *int   `json:"monthly_income" validate:"required,number"`
	Education       string `json:"education" validate:"required"`
	Nom036MthBanks  *int   `json:"nom03_6mth_banks" validate:"required,number"`
}

type PickleOT struct {
	TransactionID             string `json:"transaction_id" validate:"required"`
	ZipCode                   int    `json:"zip_code" validate:"required,number,min=2"`
	CategoryID                string `json:"category_id" validate:"required"`
	LengthOfEmploy            *int   `json:"length_of_employ" validate:"required,number"`
	Education                 string `json:"education" validate:"required"`
	FirstFourOfCell           string `json:"first_four_of_cell" validate:"required,min=4"`
	Gender                    string `json:"gender" validate:"required"`
	MaritalStatus             string `json:"marital_status" validate:"required"`
	Dependant                 *int   `json:"dependant" validate:"required,number"`
	ProfessionID              string `json:"profession_id" validate:"required"`
	TotPlafonAll              *int   `json:"tot_plafon_all" validate:"required,number"`
	OldestMobBanks            *int   `json:"oldest_mob_banks" validate:"required,number"`
	TotPlafon12MthBanksActive *int   `json:"tot_plafon_12mth_banks_active" validate:"required,number"`
	Nom036MthAll              *int   `json:"nom03_6mth_all" validate:"required,number"`
}

type PickleMM struct {
	TransactionID             string `json:"transaction_id" validate:"required"`
	CategoryID                string `json:"category_id" validate:"required"`
	ZipCode                   int    `json:"zip_code" validate:"required,number,min=2"`
	LengthOfEmploy            *int   `json:"length_of_employ" validate:"required,number"`
	TotPlafon12MthBanksActive *int   `json:"tot_plafon_12mth_banks_active" validate:"required,number"`
	Gender                    string `json:"gender" validate:"required"`
	MaritalStatus             string `json:"marital_status" validate:"required"`
	Dependant                 *int   `json:"dependant" validate:"required,number"`
	Education                 string `json:"education" validate:"required"`
	MonthlyIncome             *int   `json:"monthly_income" validate:"required,number"`
	FirstFourOfCell           string `json:"first_four_of_cell" validate:"required,min=4"`
	Worst6MthPl               *int   `json:"worst_6mth_pl" validate:"required,number"`
	MaxLimitPl                *int   `json:"max_limit_pl" validate:"required,number"`
}

type Pickle2WJabo struct {
	BpkbName                string `json:"bpkb_name" validate:"required"`
	FinalNom6012Mth         *int   `json:"final_nom60_12mth" validate:"required,number"`
	Gender                  string `json:"gender" validate:"required,max=1"`
	IndustryTypeID          *int   `json:"industry_type_id" validate:"required,number"`
	MaritalStatus           string `json:"marital_status" validate:"required,max=1"`
	MonthlyIncome           *int   `json:"monthly_income" validate:"required,number"`
	OldestMobPl             *int   `json:"oldest_mob_pl" validate:"required,number"`
	TotBakiDebetBanksActive *int   `json:"tot_baki_debet_banks_active" validate:"required,number"`
	TotBakiDebet3160Dpd     *int   `json:"tot_baki_debet_31_60dpd" validate:"required,number"`
	Worst24MthPl            *int   `json:"worst_24mth_pl" validate:"required,number"`
	MaxLimitOth             *int   `json:"max_limit_oth" validate:"required,number"`
	TransactionID           string `json:"transaction_id" validate:"required"`
	ZipCode                 int    `json:"zip_code" validate:"required,number,min=2"`
}

type Pickle2WOther struct {
	MonthlyIncome     *int   `json:"monthly_income" validate:"required,number"`
	Nom036MthAll      *int   `json:"nom03_6mth_all" validate:"required,number"`
	MaritalStatus     string `json:"marital_status" validate:"required,max=1"`
	TransactionID     string `json:"transaction_id" validate:"required"`
	ZipCode           int    `json:"zip_code" validate:"required,number,min=2"`
	MaxLimitPl        *int   `json:"max_limit_pl" validate:"required,number"`
	InstallmentAmount *int   `json:"installment_amount" validate:"required,number"`
	IndustryTypeID    *int   `json:"industry_type_id" validate:"required,number"`
	FirstFourOfCell   string `json:"first_four_of_cell" validate:"required,min=4"`
	LengthOfStay      *int   `json:"length_of_stay" validate:"required,number"`
	TotBakiDebet4     *int   `json:"tot_baki_debet_4" validate:"required,number"`
	Gender            string `json:"gender" validate:"required,max=1"`
	Worst24Mth        *int   `json:"worst_24mth" validate:"required,number"`
}

type ModelingPickleWgOther struct {
	TransactionID             string `json:"transaction_id" validate:"required"`
	ZipCode                   int    `json:"zip_code" validate:"required,number,min=2"`
	LengthOfEmploy            *int   `json:"length_of_employ" validate:"required"`
	Education                 string `json:"education" validate:"required"`
	FirstFourOfCell           string `json:"first_four_of_cell" validate:"required"`
	Gender                    string `json:"gender" validate:"required"`
	MaritalStatus             string `json:"marital_status" validate:"required"`
	Dependant                 *int   `json:"dependant" validate:"required"`
	ProfessionID              string `json:"profession_id" validate:"required"`
	TotPlafonAll              *int   `json:"tot_plafon_all" validate:"required"`
	OldestmobBanks            *int   `json:"oldestmob_banks" validate:"required"`
	TotPlafon12MthBanksActive *int   `json:"tot_plafon_12mth_banks_active" validate:"required"`
	Nom036MthAll              *int   `json:"nom03_6mth_all" validate:"required"`
}

type ModelingPickleWgJabo struct {
	TransactionID   string `json:"transaction_id" validate:"required"`
	ZipCode         int    `json:"zip_code" validate:"required,number,min=2"`
	LengthOfEmploy  *int   `json:"length_of_employ" validate:"required"`
	OldestmobBanks  *int   `json:"oldestmob_banks" validate:"required"`
	FirstFourOfCell string `json:"first_four_of_cell" validate:"required"`
	Gender          string `json:"gender" validate:"required"`
	MaritalStatus   string `json:"marital_status" validate:"required"`
	ProfessionID    string `json:"profession_id" validate:"required"`
	MaxLimitPlAll   *int   `json:"max_limit_pl_all" validate:"required"`
	MonthlyIncome   *int   `json:"monthly_income" validate:"required"`
	Education       string `json:"education" validate:"required"`
	Nom036MthBanks  *int   `json:"nom03_6mth_banks" validate:"required"`
}

type ModelingPickleWgProductLimit struct {
	TransactionID  string		`json:"transaction_id" validate:"required"`
	Tenor          int   		`json:"tenor" validate:"required"`
	ZipCode        int   		`json:"zip_code" validate:"required,number,min=2"`
	OTR            int   		`json:"otr" validate:"required"`
	Gender         string		`json:"gender" validate:"required"`
	MaritalStatus  string		`json:"marital_status" validate:"required"`
	Dependant      *int  		`json:"dependant" validate:"required"`
	Worst24MthPl   *int  		`json:"worst_24mth_pl" validate:"required"`
	Education      string		`json:"education" validate:"required"`
	TotBakiDebet2  *int  		`json:"tot_baki_debet2" validate:"required"`
	LengthOfEmploy *int  		`json:"length_of_employ" validate:"required"`
	MonthlyIncome  *int  		`json:"monthly_income" validate:"required"`
	Dsr            *float64		`json:"dsr_confins" validate:"required"`
	Age            int			`json:"age" validate:"required"`
}
