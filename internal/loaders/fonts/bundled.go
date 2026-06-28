package fonts

import (
	"embed"
	"encoding/json"
	"errors"
	"strings"
	"sync"
)

//go:embed assets/index.json
var indexStr string
var index map[string]string

//go:embed assets
var assets embed.FS

const dir = "assets"

type BundledFontRepository struct{}

func (BundledFontRepository) Get(val string) ([]byte, error) {
	fileName, ok := index[strings.ToLower(val)]
	if !ok {
		return nil, errors.New("no font found")
	}

	return assets.ReadFile(dir + "/" + fileName)
}

func (bfr *BundledFontRepository) Has(val string) bool {
	_, ok := index[strings.ToLower(val)]
	return ok
}

var (
	once sync.Once
	bfc  *BundledFontRepository
)

func GetRepository() *BundledFontRepository {
	once.Do(func() {
		err := json.Unmarshal([]byte(indexStr), &index)
		if err != nil {
			panic(err)
		}

		bfc = &BundledFontRepository{}
	})

	return bfc
}
