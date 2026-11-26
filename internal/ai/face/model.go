package face

import (
	"image"

	"github.com/photoprism/photoprism/internal/ai/tensorflow"
)

// Model is a stubbed face embedding model when TensorFlow is disabled in MyPhotoX.
type Model struct {
	name       ModelName
	resolution int
	dims       int
}

// NewModel returns a disabled face model stub.
func NewModel(name ModelName, _, _ string, resolution int, meta *tensorflow.ModelInfo, _ bool) *Model {
	if name == "" {
		name = ModelFaceNet
	}
	if resolution == 0 {
		resolution = CropSize.Width
	}
	if meta == nil {
		meta = new(tensorflow.ModelInfo)
	}
	return &Model{
		name:       name,
		resolution: resolution,
		dims:       modelDims(name, meta),
	}
}

func modelDims(name ModelName, meta *tensorflow.ModelInfo) int {
	if meta != nil && meta.Output != nil && meta.Output.NumOutputs > 0 {
		return int(meta.Output.NumOutputs)
	}
	if info := FindEmbeddingModel(name); info != nil {
		return info.Dims
	}
	return len(NullEmbedding)
}

func (m *Model) ModelName() ModelName              { return m.name }
func (m *Model) Dims() int                         { return m.dims }
func (m *Model) CropSize() (width, height int)     { return m.resolution, m.resolution }
func (m *Model) Aligned() bool                     { return false }
func (m *Model) Close() error                      { return nil }
func (m *Model) Init() error                       { return nil }
func (m *Model) ModelLoaded() bool                 { return false }
func (m *Model) Run(_ image.Image) Embeddings      { return nil }

func (m *Model) Detect(fileName string, minSize int, _ bool, _ int) (Faces, error) {
	return Detect(fileName, minSize)
}
