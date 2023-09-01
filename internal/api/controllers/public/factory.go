package public

import (
	"github.com/authzed/authzed-go/v1"
	"playbook-dispatcher/internal/api/connectors"

	"gorm.io/gorm"
)

func CreateController(database *gorm.DB, cloudConnectorClient connectors.CloudConnectorClient) ServerInterfaceWrapper {
	return ServerInterfaceWrapper{
		Handler: &controllers{
			database:             database,
			cloudConnectorClient: cloudConnectorClient,
		},
	}
}

func CreateSpiceDBController(spiceDBClient *authzed.Client, database *gorm.DB, cloudConnectorClient connectors.CloudConnectorClient) ServerInterfaceWrapper {
	return ServerInterfaceWrapper{
		Handler: &spiceDBControllers{
			spiceDBClient:        spiceDBClient,
			database:             database,
			cloudConnectorClient: cloudConnectorClient,
		},
	}
}

// implements api.ServerInterface
type controllers struct {
	database             *gorm.DB
	cloudConnectorClient connectors.CloudConnectorClient
}

// implements api.ServerInterface
type spiceDBControllers struct {
	spiceDBClient        *authzed.Client
	database             *gorm.DB
	cloudConnectorClient connectors.CloudConnectorClient
}
