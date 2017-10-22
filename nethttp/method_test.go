package nethttp_test

import (
	"github.com/cpapidas/pegasus/nethttp"

	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMethod_String(t *testing.T) {
	// Should the right const variables
	assert.Equal(t, http.MethodGet, nethttp.Get.String())
	assert.Equal(t, http.MethodHead, nethttp.Head.String(), "Should be equals to methods type")
	assert.Equal(t, http.MethodPost, nethttp.Post.String(), "Should be equals to methods type")
	assert.Equal(t, http.MethodPut, nethttp.Put.String(), "Should be equals to methods type")
	assert.Equal(t, http.MethodPatch, nethttp.Patch.String(), "Should be equals to methods type")
	assert.Equal(t, http.MethodDelete, nethttp.Delete.String(), "Should be equals to methods type")
	assert.Equal(t, http.MethodConnect, nethttp.Connect.String(), "Should be equals to methods type")
	assert.Equal(t, http.MethodOptions, nethttp.Options.String(), "Should be equals to methods type")
	assert.Equal(t, http.MethodTrace, nethttp.Trace.String(), "Should be equals to methods type")
}
