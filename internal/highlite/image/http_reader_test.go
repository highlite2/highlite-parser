package image

import (
	_ "context"
	_ "io/ioutil"
	"testing"
	_ "time"

	_ "github.com/stretchr/testify/assert"
)

// TODO
func TestDownloadImages(t *testing.T) {
	// arrange
	//p := HTTPProvider{}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()
	//var images = []string{"100232.jpg", "100232_detail1.jpg", "100232_draw1.jpg"}

	// act
	//result, err := p.GetImages(ctx, images)

	// assert
	//assert.Nil(t, err)
	//for _, name := range images {
	//	reader, ok := result[name]
	//	assert.True(t, ok)
	//
	//	bytes, err := ioutil.ReadAll(reader)
	//	assert.Nil(t, err)
	//
	//	assert.Nil(t, ioutil.WriteFile(name, bytes, 0644))
	//}
}
