package tensorflow

import (
	"errors"
)

// SavedModel is a stubbed TensorFlow model handle when TensorFlow CGo is disabled.
type SavedModel struct{}

// SavedModel returns an error because TensorFlow is disabled in MyPhotoX.
func LoadSavedModel(_ string, _ []string) (*SavedModel, error) {
	return nil, errors.New("tensorflow: disabled in MyPhotoX")
}
