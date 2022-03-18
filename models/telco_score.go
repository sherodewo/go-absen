package models

import "time"

type ExperianScore struct {
	PayLoad PayloadExperian `json:"PAYLOAD"`
}

type ExperianScoreResp struct {
	Payload PayloadResp `json:"PAYLOAD"`
}

type ReqHeader struct {
	ClientID    string  `json:"CLIENTID"`
	PartnerID   string  `json:"PARTNERID"`
	ProductID   string  `json:"PRODUCTID"`
	TransID     string  `json:"TRANSID"`
	ReqDateTime float64 `json:"REQDATETIME"` //example yyyymmddhhmmss
	MSISDN      string  `json:"MSISDN"`
	ServiceName string  `json:"SERVICENAME"`
	Model       string  `json:"MODEL"`
	Token       string  `json:"TOKEN"`
}

type ResHeader struct {
	ClientID    string  `json:"CLIENTID"`
	PartnerID   string  `json:"PARTNERID"`
	ProductID   string  `json:"PRODUCTID"`
	TransID     *string `json:"TRANSID"`
	ReqDateTime float64 `json:"REQDATETIME"` //example yyyymmddhhmmss
	MSISDN      string  `json:"MSISDN"`
	ServiceName string  `json:"SERVICENAME"`
	BetID       string  `json:"BETID"`
	Status      int64   `json:"STATUS"`
	Desc        string  `json:"DESC"`
	Model       string  `json:"MODEL"`
}

type BodyResp struct {
	CreditStatus CreditStatus `json:"CREDITSTATUS"`
}

type ReqBody struct {
	Data Data `json:"DATA"`
}

type Data struct {
	OTP string `json:"OTP"`
}

type CreditStatus struct {
	Score  *float64 `json:"Score"`
	Result *string  `json:"Result"`
}

type PayloadExperian struct {
	ReqHeader ReqHeader `json:"REQHEADER"`
	Body      ReqBody   `json:"BODY"`
}

type PayloadResp struct {
	ResHeader ResHeader `json:"RESHEADER"`
	Body      BodyResp  `json:"BODY"`
}

type TokenResponse struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	ExpiresIn   float64 `json:"expires_in"`
	Scope       string  `json:"scope"`
}

type CreditScoreResp struct {
	ExperianID  string   `json:"experian_id"`
	ProspectID  string   `json:"ProspectID"`
	Score       *float64 `json:"score"`
	Result      *string  `json:"result"`
	Status      string   `json:"status"`
	PhoneNumber string   `json:"phone_number"`
}

type CreditScoreRespIdx struct {
	TransactionID    string   `json:"transaction_id"`
	ProspectID    string   `json:"prospect_id"`
	SourceDecison string   `json:"source_decison"`
	Score         *float64 `json:"score"`
	Result        *string  `json:"result"`
	Status        string   `json:"status"`
	PhoneNumber   string   `json:"phone_number"`
}

type PickleResponse struct {
	OrderID string   `json:"OrderID"`
	Code    string   `json:"code"`
	Result  string   `json:"result"`
	Score   *float64 `json:"score"`
}


type PickleResponseIDX struct {
	Messages string      `json:"messages"`
	Errors   interface{} `json:"errors"`
	Data     struct {
		ProspectID string `json:"prospect_id"`
		Score      float64    `json:"score"`
		Result     string `json:"result"`
	} `json:"data"`
	ServerTime time.Time `json:"server_time"`
}

type PickleLimitResponse struct {
	ProspectID string   `json:"ProspectID"`
	Result     string   `json:"result"`
	Score      *float64 `json:"score"`
}
