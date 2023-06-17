package optional_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/optional"
)

func TestRef(t *testing.T) {
	t.Parallel()

	stringv := optional.Ref("test")
	assert.Equal(t, "test", *stringv)

	boolv := optional.Ref(true)
	assert.Equal(t, true, *boolv)

	intv := optional.Ref(420)
	assert.Equal(t, 420, *intv)

	floatv := optional.Ref(42.0)
	assert.Equal(t, 42.0, *floatv)

	type strct struct{}
	structv := optional.Ref(strct{})
	assert.Equal(t, strct{}, *structv)
}
