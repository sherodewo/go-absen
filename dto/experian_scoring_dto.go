package dto

type ExperianScoringDto struct {
	IsIndosat  string `json:"is_indosat" form:"is_indosat" validate:"required"`
	Experian   string `json:"experian" form:"experian" validate:"required"`
	Internal   string `json:"internal" form:"internal" validate:"required"`
	ScoreLos   string `json:"score_los" form:"score_los" validate:"required"`
	FinalScore string `json:"final_score" form:"final_score" validate:"required"`
	Notes      string `json:"notes" form:"notes"`
}

type ExperianScoringUpdateDto struct {
	IsIndosat  string `json:"is_indosat" form:"is_indosat" validate:"required"`
	Experian   string `json:"experian" form:"experian" validate:"required"`
	Internal   string `json:"internal" form:"internal" validate:"required"`
	ScoreLos   string `json:"score_los" form:"score_los" validate:"required"`
	FinalScore string `json:"final_score" form:"final_score" validate:"required"`
	Notes      string `json:"notes" form:"notes"`
}
