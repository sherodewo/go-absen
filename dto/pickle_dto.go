package dto

type PickleDto struct {
	OrderID string `json:"OrderID" validate:"required"`
	CatID   string `json:"CatID_" validate:"required"`
	DP      *int   `json:"DP_" validate:"required"`
	NTF     *int   `json:"NTF_" validate:"required"`
	ANGS    *int   `json:"ANGS_" validate:"required"`
	Tnr     *int   `json:"Tnr_" validate:"required"`
	Gnr     string `json:"Gnr_" validate:"required"`
	Age     *int   `json:"Age_" validate:"required"`
	Mrtl    string `json:"Mrtl_" validate:"required"`
	Dpnt    *int   `json:"Dpnt_" validate:"required"`
	Hmsts   string `json:"Hmsts_" validate:"required"`
	Lnstay  *int   `json:"Lnstay_" validate:"required"`
	Educ    string `json:"Educ_" validate:"required"`
	Prof    string `json:"Prof_" validate:"required"`
	Jobex   *int   `json:"Jobex_" validate:"required"`
	Inc     *int   `json:"Inc_" validate:"required"`
	VarInc  *int   `json:"VarInc_" validate:"required"`
	SpsInc  *int   `json:"SpsInc_" validate:"required"`
	Zip     string `json:"Zip_" validate:"required"`
	IIR     *int   `json:"IIR_" validate:"required"`
}

type PickleLimitDto struct {
	DepV        float64 `json:"Dep_v"`
	EducationV  string  `json:"Education_v"`
	GenderV     string  `json:"Gender_v"`
	HomestsV    string  `json:"Homests_v"`
	INCMJTV     float64 `json:"INCMJT_v"`
	JobExpV     float64 `json:"JobExp_v"`
	LenOfStayV  float64 `json:"LenOfStay_v"`
	MaritalStsV string  `json:"MaritalSts_v"`
	OrderID     string  `json:"OrderID"`
	ProfessionV string  `json:"Profession_v"`
	ZIP3V       string  `json:"ZIP3_v"`
	AgeV        float64 `json:"age_v"`
}

type PickleModelingDto struct {
	SupplierID       string      `json:"supplier_id" validate:"required"`
	ProspectID       string      `json:"prospect_id" validate:"required"`
	CbFound          bool        `json:"cb_found" validate:"required"`
	StatusKonsumen   string      `json:"status_konsumen"`
	PhoneNumber      string      `json:"phone_number" validate:"required,number"`
	RequestorID      string      `json:"requestor_id" validate:"required"`
	ScoreGeneratorID string      `json:"score_generator_id" validate:"required"`
	Data             interface{} `json:"data"`
}

type PickleJJ struct {
	TransactionID   string `json:"transaction_id"`
	HomeZip         int64  `json:"home_zip"`
	CategoryID      string `json:"category_id"`
	LengthOfEmploy  int64  `json:"length_of_employ"`
	OldestMobBank   int64  `json:"oldest_mob_bank"`
	FirstFourOfCell string `json:"first_four_of_cell"`
	Gender          string `json:"gender"`
	MaritalStatus   string `json:"marital_status"`
	ProfessionID    string `json:"profession_id"`
	MaxLimitPlAll   int64  `json:"max_limit_pl_all"`
	MonthlyIncome   int64  `json:"monthly_income"`
	Education       string `json:"education"`
	NumMthBanks     int64  `json:"nom_03_6_mth_banks"`
}

type PickleMM struct {
	TransactionID             string `json:"transaction_id"`
	HomeZip                   int64  `json:"home_zip"`
	CategoryID                string `json:"category_id"`
	LengthOfEmploy            int64  `json:"length_of_employ"`
	TotPlafon12MthBanksActive int64  `json:"tot_plafon_12_mth_banks_active"`
	Gender                    string `json:"gender"`
	MaritalStatus             string `json:"marital_status"`
	Dependant                 int64  `json:"dependant"`
	Education                 string `json:"education"`
	MonthlyIncome             int64  `json:"monthly_income"`
	FirstFourOfCell           string `json:"first_four_of_cell"`
	Worst6MthPl               int64  `json:"worst_6_mth_pl"`
	MaxLimit                  int64  `json:"max_limit"`
	Score                     string `json:"score"`
}

type PickleOther struct {
	TransactionID             string `json:"transaction_id"`
	HomeZip                   int64  `json:"home_zip"`
	CategoryID                string `json:"category_id"`
	LengthOfEmploy            int64  `json:"length_of_employ"`
	Education                 string `json:"education"`
	FirstFourOfCell           string `json:"first_four_of_cell"`
	Gender                    string `json:"gender"`
	MaritalStatus             string `json:"marital_status"`
	Dependant                 int64  `json:"dependant"`
	ProfessionID              string `json:"profession_id"`
	TotalPlafonAll            int64  `json:"total_plafon_all"`
	OldestMobBank             int64  `json:"oldest_mob_bank"`
	TotPlafon12MthBanksActive int64  `json:"tot_plafon_12_mth_banks_active"`
	NumMthBanks               int64  `json:"nom_03_6_mth_all"`
}
