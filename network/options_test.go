package network_test

import (
	"github.com/cpapidas/pegasus/network"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOptions(t *testing.T) {
	// Should panic if field mapper is not stetted
	assert.Panics(t, func() {
		options := network.NewOptions()
		options.Fields["foo"]["bar"] = "4"
	}, "Should not panics")
}

func TestBuildOptions(t *testing.T) {
	// Should build options successful
	opts := network.NewOptions()
	opts.SetField("foo", "bar", "baz")
	built := network.BuildOptions(opts.Marshal())
	assert.Equal(t, "baz", built.GetField("foo", "bar"),
		"Should build options successfully")
}

func TestOptions_SetParams(t *testing.T) {

	// Should sets the parameter
	options := network.Options{}
	options.SetParams(map[string]string{"foo": "fa"})

	assert.Equal(t, map[string]string{"foo": "fa"}, options.Fields["PARAMS"],
		`Should be equals to mapper "foo": "fa"}`)

	assert.Equal(t, map[string]string{"foo": "fa"}, options.GetParams(),
		`Should be equals to mapper map[string]string{"foo": "fa"}`)

	options.SetParam("baz", "ba")
	assert.Equal(t, map[string]string{"foo": "fa", "baz": "ba"}, options.GetParams(),
		`Should be equals to mapper map[string]string{"foo": "fa", "baz": "ba"}`)

	assert.Equal(t, "fa", options.GetParam("foo"),
		`Should be equals to string value "fa"`)
}

func TestOptions_GetParams(t *testing.T) {
	// Should gets and returns the fa value
	options := network.Options{}
	options.SetParams(map[string]string{"foo": "fa"})
	assert.Equal(t, map[string]string{"foo": "fa"}, options.GetParams(),
		`Should be equals to map[string]string{"foo": "fa"}`)
}

func TestOptions_SetParam(t *testing.T) {
	// Should set the following params
	options := &network.Options{}
	options.SetParam("baz", "ba")
	assert.Equal(t, map[string]string{"baz": "ba"}, options.Fields["PARAMS"], `Should set the param bar`)
}

func TestOptions_GetParam(t *testing.T) {
	// Should gets and returns the ba value
	options := network.Options{}
	options.SetParam("baz", "ba")
	assert.Equal(t, map[string]string{"baz": "ba"}, options.GetParams(),
		`Should gets and returns the "ba" value`)
}

func TestOptions_SetHeaders(t *testing.T) {
	// Should sets the headers
	options := network.Options{}
	options.SetHeaders(map[string]string{"foo": "fa"})
	assert.Equal(t, map[string]string{"foo": "fa"}, options.Fields["HEADERS"],
		`Should be equal to mapper map[string]string{"foo": "fa"}"`)
}

func TestOptions_GetHeaders(t *testing.T) {
	// Should sets the headers
	options := network.Options{}
	options.SetHeaders(map[string]string{"foo": "fa"})
	assert.Equal(t, map[string]string{"foo": "fa"}, options.GetHeaders(),
		`Should be equal to mapper map[string]string{"foo": "fa"}"`)
}

func TestOptions_SetHeader(t *testing.T) {
	// Should return the value fa for foo header
	options := network.Options{}
	options.SetHeader("baz", "ba")
	assert.Equal(t, map[string]string{"baz": "ba"}, options.GetHeaders(),
		`Should be equal to map[string]string{"foo":"faa"}`)
}

func TestOptions_GetHeader(t *testing.T) {
	// Should return the value fa for foo header
	options := network.Options{}
	options.SetHeaders(map[string]string{"foo": "fa"})
	assert.Equal(t, "fa", options.GetHeader("foo"),
		`Should return the value fa for foo header`)
}

func TestOptions_Marshal(t *testing.T) {
	// Should not return nil object for nil data
	data := network.NewOptions().Marshal()
	assert.NotNil(t, data, "Should not return nil object for nil data")
}

func TestOptions_Unmarshal(t *testing.T) {
	// Should not return nil object for nil data
	data := network.NewOptions().Unmarshal(nil)
	assert.NotNil(t, data, "Should not return nil object for nil data")

	// Should return nil for invalid data
	invalidData := network.NewOptions().Unmarshal([]byte("whatever"))
	assert.Nil(t, invalidData, "Should return nil for invalid data")
}

func TestOptions_SetField(t *testing.T) {
	// Should sets a field
	options := &network.Options{}
	options.SetField("foo", "faa", "Ga")
	assert.Equal(t, map[string]map[string]string{"foo":{"faa":"Ga"}},  options.Fields,
		`Should set the foo -> faa field with value Ga`)
}

func TestOptions_GetField(t *testing.T) {
	// Should gets the field
	options := &network.Options{}
	options.Fields = make(map[string]map[string]string)
	options.Fields["foo"] =  make(map[string]string)
	options.Fields["foo"]["faa"] = "Ga"
	assert.Equal(t, "Ga",  options.GetField("foo", "faa"),
		`Should set the foo -> faa field with value Ga`)

	// Should get a nil field for uninitialized fields
	options = &network.Options{}
	assert.Empty(t, options.GetField("foo", "faa"),
		`Should get a nil field for uninitialized fields`)

	// Should get a nil field for uninitialized group
	options = &network.Options{}
	options.Fields = make(map[string]map[string]string)
	assert.Empty(t, options.GetField("foo", "faa"),
		`Should get a nil field for uninitialized group`)
}
