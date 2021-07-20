package converterandformatter_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/savannahghi/converterandformatter"
	"github.com/stretchr/testify/assert"
)

type SampleStruct struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type EmbededStruct struct {
	FieldStruct `json:"field"`
	Hello       string `json:"hello"`
}

type FieldStruct struct {
	OnePoint string        `json:"one_point"`
	Sample   *SampleStruct `json:"sample"`
}

func TestStructToMap_Normal(t *testing.T) {
	sample := SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}

	res, err := converterandformatter.StructToMap(sample)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	fmt.Printf("%+v \n", res)
	// Output: map[name:John Doe id:12121]
	jbyt, err := json.Marshal(res)
	assert.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"id":"12121","name":"John Doe"}
}

func TestStructToMap_FieldStruct(t *testing.T) {

	sample := &SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}
	field := FieldStruct{
		Sample:   sample,
		OnePoint: "yuhuhuu",
	}

	res, err := converterandformatter.StructToMap(field)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Printf("%+v \n", res)
	// Output: map[sample:0xc4200f04a0 one_point:yuhuhuu]
	jbyt, err := json.Marshal(res)
	assert.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"one_point":"yuhuhuu","sample":{"name":"John Doe","id":"12121"}}

}

func TestStructToMap_EmbeddedStruct(t *testing.T) {

	sample := &SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}
	field := FieldStruct{
		Sample:   sample,
		OnePoint: "yuhuhuu",
	}

	embed := EmbededStruct{
		FieldStruct: field,
		Hello:       "WORLD!!!!",
	}

	res, err := converterandformatter.StructToMap(embed)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Printf("%+v \n", res)
	//Output: map[field:map[one_point:yuhuhuu sample:0xc420106420] hello:WORLD!!!!]

	jbyt, err := json.Marshal(res)
	assert.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"field":{"one_point":"yuhuhuu","sample":{"name":"John Doe","id":"12121"}},"hello":"WORLD!!!!"}
}

func TestStructToMap(t *testing.T) {
	type testStruct struct {
		FirstField  string `json:"firstField,omitempty"`
		SecondField int    `json:"secondField,omitempty"`
	}

	type args struct {
		item interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid struct",
			args: args{
				item: testStruct{
					FirstField:  "A",
					SecondField: 1.0,
				},
			},
			want: map[string]interface{}{
				"firstField":  "A",
				"secondField": 1.0,
			},
			wantErr: false,
		},
		{
			name: "invalid input - not a struct and won't marshal to JSON",
			args: args{
				item: make(chan string),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converterandformatter.StructToMap(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for k, v := range tt.want {
				valK, present := got[k]
				assert.True(t, present)
				assert.Equal(t, valK, v)
			}
		})
	}
}

func TestGenerateRandomWithNDigits(t *testing.T) {
	result, err := converterandformatter.GenerateRandomWithNDigits(5)
	if result == "" {
		t.Errorf("can't generate random with n digits")
		return
	}
	if err != nil {
		t.Errorf("can't generate random with n digits: %v", err)
		return
	}
}

func TestGenerateRandomEmail(t *testing.T) {
	email := converterandformatter.GenerateRandomEmail()
	if email == "" {
		t.Errorf("can't generate a unique email")
	}
}
