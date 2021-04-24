package hasura

import (
	"reflect"
	"testing"

	"github.com/dipdup-net/go-lib/config"
	"github.com/stretchr/testify/assert"
)

type testTable struct {
	Field1 string
	Field2 int64
}

type testTable2 struct {
	Field1 string
	Field2 int64 `gorm:"primaryKey"`
}

func (testTable2) TableName() string {
	return "fake_name"
}

type testTable3 struct {
	Field1 string `gorm:"-"`
	Field2 int64  `gorm:"primaryKey"`
}

func (testTable3) TableName() (int64, error) {
	return 0, nil
}

type testTable4 struct {
	Field1 string `gorm:"-"`
	Field2 int64
	inherited
}

type inherited struct {
	Field3 string
}

func (testTable4) TableName() int64 {
	return 0
}

func Test_getTableName(t *testing.T) {
	tests := []struct {
		name  string
		model interface{}
		want  string
	}{
		{
			name:  "Test",
			model: &testTable{},
			want:  "test_table",
		}, {
			name:  "Test 2",
			model: &testTable2{},
			want:  "fake_name",
		}, {
			name:  "Test 3",
			model: &testTable3{},
			want:  "test_table_3",
		}, {
			name:  "Test 4",
			model: &testTable4{},
			want:  "test_table_4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := reflect.ValueOf(tt.model)
			typ := reflect.TypeOf(tt.model)
			for typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			got := getTableName(value, typ)
			if got != tt.want {
				t.Errorf("getTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getColumns(t *testing.T) {
	tests := []struct {
		name  string
		model interface{}
		want  []string
	}{
		{
			name:  "Test",
			model: &testTable{},
			want:  []string{"field_1", "field_2"},
		}, {
			name:  "Test 2",
			model: &testTable2{},
			want:  []string{"field_1", "field_2"},
		}, {
			name:  "Test 3",
			model: &testTable3{},
			want:  []string{"field_2"},
		}, {
			name:  "Test 4",
			model: &testTable4{},
			want:  []string{"field_2", "field_3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := reflect.TypeOf(tt.model)
			for typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			if got := getColumns(typ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	type args struct {
		cfg    config.Database
		models []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				cfg: config.Database{
					Kind: "mysql",
				},
				models: []interface{}{
					&testTable{}, &testTable2{}, &testTable3{}, &testTable4{},
				},
			},
			want: `{"tables":[{"array_relationships":[],"object_relationships":[],"select_permissions":[{"permission":{"allow_aggregations":true,"columns":["field_1","field_2"],"filter":{}},"role":"user"}],"table":{"name":"test_table","schema":"public"}},{"array_relationships":[],"object_relationships":[],"select_permissions":[{"permission":{"allow_aggregations":true,"columns":["field_1","field_2"],"filter":{}},"role":"user"}],"table":{"name":"fake_name","schema":"public"}},{"array_relationships":[],"object_relationships":[],"select_permissions":[{"permission":{"allow_aggregations":true,"columns":["field_2"],"filter":{}},"role":"user"}],"table":{"name":"test_table_3","schema":"public"}},{"array_relationships":[],"object_relationships":[],"select_permissions":[{"permission":{"allow_aggregations":true,"columns":["field_2","field_3"],"filter":{}},"role":"user"}],"table":{"name":"test_table_4","schema":"public"}}],"version":2}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.args.cfg, tt.args.models...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotStr, err := json.MarshalToString(got)
			if err != nil {
				t.Errorf("MarshalToString() error = %v", err)
				return
			}
			assert.JSONEq(t, tt.want, gotStr)
		})
	}
}
