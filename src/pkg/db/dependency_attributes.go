// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
)

type DependencyAttribute struct {
	Model

	DependencyID uuid.UUID        `gorm:"uniqueIndex:dependency_attribute_unique"`
	Name         string           `gorm:"uniqueIndex:dependency_attribute_unique"`
	Schema       *jsonschema.Node `gorm:"type:jsonb"`
	Computed     bool             `gorm:"uniqueIndex:dependency_attribute_unique"` // true means output, false means input

	Dependency *Dependency `gorm:"foreignKey:DependencyID"`
}

type DependencyAttributes []*DependencyAttribute

func (dbAttr DependencyAttribute) ToProto() *terrariumpb.DependencyInputsAndOutputs {
	resp := &terrariumpb.DependencyInputsAndOutputs{}

	// Only set Title if it's not empty
	if dbAttr.Name != "" {
		resp.Title = dbAttr.Name
	}

	// Only set values from the Schema if Schema is not nil
	if dbAttr.Schema != nil {
		if dbAttr.Schema.Description != "" {
			resp.Description = dbAttr.Schema.Description
		}
		if dbAttr.Schema.Type != "" {
			resp.Type = dbAttr.Schema.Type
		}
	}

	return resp
}

func (dbAttrs DependencyAttributes) ToProto() []*terrariumpb.DependencyInputsAndOutputs {
	resp := make([]*terrariumpb.DependencyInputsAndOutputs, 0, len(dbAttrs))

	for _, dbAttr := range dbAttrs {
		protoAttr := dbAttr.ToProto()
		// Check if protoAttr is not an entirely empty object
		if protoAttr.Title != "" || protoAttr.Description != "" || protoAttr.Type != "" {
			resp = append(resp, protoAttr)
		}
	}

	return resp
}

func (a *DependencyAttribute) GetCondition() entity {
	return &DependencyAttribute{
		DependencyID: a.DependencyID,
		Name:         a.Name,
		Computed:     a.Computed,
	}
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateDependencyAttribute(e *DependencyAttribute) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"dependency_id", "name"})
}
