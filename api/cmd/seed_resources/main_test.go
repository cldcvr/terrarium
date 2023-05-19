package main

import (
	"fmt"
	"testing"

	"github.com/cldcvr/terrarium/api/db/mocks"
	"github.com/cldcvr/terrarium/api/pkg/tf/schema"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_pushProvidersSchemaToDB(t *testing.T) {
	tests := []struct {
		name            string
		providersSchema *schema.ProvidersSchema
		mocks           func(*mocks.DB)
		panicFu         func(t assert.TestingT, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool
	}{
		{
			name: "success",
			providersSchema: &schema.ProvidersSchema{
				ProviderSchemas: map[string]schema.ProviderSchema{
					"mock_provider": {
						ResourceSchemas: map[string]schema.SchemaRepresentation{
							"mock_resource": {
								Block: schema.BlockRepresentation{
									Attributes: map[string]schema.AttributeRepresentation{
										"A": {},
									},
									BlockTypes: map[string]schema.BlockTypeRepresentation{
										"X": {
											Block: schema.BlockRepresentation{
												Attributes: map[string]schema.AttributeRepresentation{
													"Y": {},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			mocks: func(dbMocks *mocks.DB) {
				dbMocks.On("CreateTFProvider", mock.Anything).Return(uuid.New(), nil).Once()
				dbMocks.On("CreateTFResourceType", mock.Anything).Return(uuid.New(), nil).Once()
				dbMocks.On("CreateTFResourceAttribute", mock.Anything).Return(uuid.New(), nil).Twice()
			},
			panicFu: assert.NotPanics,
		},
		{
			name: "panic",
			providersSchema: &schema.ProvidersSchema{
				ProviderSchemas: map[string]schema.ProviderSchema{
					"mock_provider": {},
				},
			},
			mocks: func(dbMocks *mocks.DB) {
				dbMocks.On("CreateTFProvider", mock.Anything).Return(uuid.Nil, fmt.Errorf("mocked error")).Once()
			},
			panicFu: assert.Panics,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocks := &mocks.DB{}
			tt.mocks(dbMocks)

			tt.panicFu(t, func() {
				pushProvidersSchemaToDB(tt.providersSchema, dbMocks)
			})

			dbMocks.AssertExpectations(t)

		})
	}
}

func Test_loadProvidersSchema(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.ProvidersSchema
		wantErr bool
	}{
		{
			name: "success",
			args: args{"./example_schema.json"},
			want: &schema.ProvidersSchema{
				ProviderSchemas: map[string]schema.ProviderSchema{
					"mock_provider": {
						ResourceSchemas: map[string]schema.SchemaRepresentation{
							"mock_resource": {
								Block: schema.BlockRepresentation{
									Attributes: map[string]schema.AttributeRepresentation{
										"A": {
											Description: "a",
											Type:        "string",
											Computed:    true,
										},
									},
									BlockTypes: map[string]schema.BlockTypeRepresentation{
										"X": {
											Block: schema.BlockRepresentation{
												Attributes: map[string]schema.AttributeRepresentation{
													"Y": {
														Description: "y",
														Type:        "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "fail read",
			args:    args{"./invalid_file_path"},
			wantErr: true,
		},
		{
			name:    "fail unmarshal",
			args:    args{"./invalid_schema.txt"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadProvidersSchema(tt.args.filename)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
