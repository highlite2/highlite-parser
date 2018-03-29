package imprt

import (
	"testing"

	"highlite2-import/internal/highlite"

	"github.com/stretchr/testify/assert"
)

func TestAttributesImport_GetBrandId(t *testing.T) {
	i := &AttributesImport{
		brandMap: map[string]string{
			"name1": "id1",
			"name2": "id2",
		},
	}

	id1, ok1 := i.GetBrandID(highlite.Product{Brand: "Name1"})
	id2, ok2 := i.GetBrandID(highlite.Product{Brand: "name3"})

	assert.True(t, ok1)
	assert.Equal(t, "id1", id1)
	assert.False(t, ok2)
	assert.Equal(t, "", id2)
}
