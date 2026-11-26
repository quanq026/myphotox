package tensorflow

import (
	"testing"

	"github.com/photoprism/photoprism/pkg/fs"
)

var assetsPath = fs.Abs("../../../assets")

func TestLoadSavedModelDisabled(t *testing.T) {
	model, err := LoadSavedModel("", []string{"serve"})
	if err == nil {
		t.Fatal("expected error because TensorFlow is disabled in MyPhotoX")
	}
	if model != nil {
		t.Fatalf("expected nil model, got %v", model)
	}
}
