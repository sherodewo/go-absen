package credential

import (
	"github.com/kreditplus/scorepro/utils"
	"os"
)

var (
	AppHost               string
	AppPort               string
	DbHost                string
	DbPort                string
	DbUsername            string
	DbPassword            string
	DbName                string
	DbNameScorepro        string
	NewRelicConfigLicense string
	ExperianBaseUrl       string
	ExperianClientID      string
	ExperianPartnerID     string
	ExperianProductID     string
	ExperianUsername      string
	ExperianPassword      string
)

func CredentialsConfig() (err error) {
	AppHost, err = utils.DecryptCredential(os.Getenv("APP_HOST"))
	AppPort, err = utils.DecryptCredential(os.Getenv("APP_PORT"))
	DbHost, err = utils.DecryptCredential(os.Getenv("DB_HOST"))
	DbPort, err = utils.DecryptCredential(os.Getenv("DB_PORT"))
	DbUsername, err = utils.DecryptCredential(os.Getenv("DB_USERNAME"))
	DbPassword, err = utils.DecryptCredential(os.Getenv("DB_PASSWORD"))
	DbName, err = utils.DecryptCredential(os.Getenv("DB_DATABASE"))
	DbNameScorepro, err = utils.DecryptCredential(os.Getenv("DB_DATABASE_SCOREPRO_V2"))
	NewRelicConfigLicense, err = utils.DecryptCredential(os.Getenv("NEWRELIC_CONFIG_LICENSE"))
	ExperianBaseUrl, err = utils.DecryptCredential(os.Getenv("EXPERIAN_BASE_URL"))
	ExperianClientID, err = utils.DecryptCredential(os.Getenv("EXPERIAN_CLIENT_ID"))
	ExperianPartnerID, err = utils.DecryptCredential(os.Getenv("EXPERIAN_PARTNERID"))
	ExperianProductID, err = utils.DecryptCredential(os.Getenv("EXPERIAN_PRODUCTID"))
	ExperianUsername, err = utils.DecryptCredential(os.Getenv("EXPERIAN_USERNAME"))
	ExperianPassword, err = utils.DecryptCredential(os.Getenv("EXPERIAN_PASSWORD"))

	return err
}
