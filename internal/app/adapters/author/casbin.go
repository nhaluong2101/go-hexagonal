package author

import (
	"database/sql"
	"fmt"
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure/configs"
	_casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"go.uber.org/fx"
	"log"

	_ "github.com/lib/pq"
)

type CasbinConfig struct {
	DSN        string
	DriverName string
	Enforcer   *_casbin.Enforcer
}

func NewCasbinConfig(db *configs.DB) *CasbinConfig {
	casbinConfig := &CasbinConfig{
		DSN:        db.DSN,
		DriverName: db.DriverName,
	}
	casbinConfig.Enforcer = casbinConfig.NewEnforcer()
	return casbinConfig
}

func (casbinConfig *CasbinConfig) LoadModelFromDB() (string, error) {

	db, err := sql.Open(casbinConfig.DriverName, casbinConfig.DSN)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var modelText string
	query := "SELECT model_text FROM casbin_model WHERE model_name = 'rbac_model' LIMIT 1"
	err = db.QueryRow(query).Scan(&modelText)
	if err != nil {
		return "", err
	}
	return modelText, nil

}

func (casbin *CasbinConfig) NewEnforcer() *_casbin.Enforcer {

	driverName := casbin.DriverName
	dsn := casbin.DSN

	// Create XORM adapter for Casbin
	// dbSpecified = "true" is for automatically creating the 'casbin_rule' table
	adapter, err := xormadapter.NewAdapter(driverName, dsn, true)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v\n", err)
	}
	// Assume loadModelFromDB pulls your model configuration from the DB
	modelText, err := casbin.LoadModelFromDB()
	if err != nil {
		log.Fatalf("Failed to load model: %v\n", err)
	}
	// Load Casbin model from the text
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		log.Fatalf("Failed to create model from string: %v\n", err)
	}
	// Create a Casbin enforcer with the adapter and model
	e, err := _casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v\n", err)
	}

	return e

}

func (casbinConfig *CasbinConfig) LoadPolicy() {
	err := casbinConfig.Enforcer.LoadPolicy()
	if err != nil {
		fmt.Errorf("failed to load policy: %w", err)
	}
}

var CasbinModule = fx.Module(
	"casbin-module",
	fx.Provide(
		NewCasbinConfig,
	),
)
