package models

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/yuchanns/kong-exercise-microservices/utils/helpers"
	"gorm.io/gorm"
)

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(helpers.GetHostID())
	if err != nil {
		panic(err)
	}
}

type BaseModel struct{}

const (
	DefaultFormat = "2006-01-02 15:04:05"
)

func (t BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now().Format(DefaultFormat)
	if field := tx.Statement.Schema.LookUpField("created_at"); field != nil {
		tx.Statement.SetColumn("created_at", now)
	}
	if field := tx.Statement.Schema.LookUpField("modified_at"); field != nil {
		tx.Statement.SetColumn("modified_at", now)
	}
	tx.Statement.SetColumn("id", node.Generate().Int64())
	return
}

func (t BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now().Format(DefaultFormat)
	if field := tx.Statement.Schema.LookUpField("modified_at"); field != nil {
		tx.Statement.SetColumn("modified_at", now)
	}
	return
}
