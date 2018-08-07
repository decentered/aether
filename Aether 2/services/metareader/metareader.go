// Services > MetaReader
// This package reads the META fields of entities we receive. Effectively a JSON parser that looks for a specific field.

package metareader

import (
	"github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

/*
  This should be properly delienated, in the way that we should allow for scanning of multiple values, as well as 'return me the whole thing parsed' kind of thing. TODO.
*/
func ReadSingleField(metaField, requestedKey string) interface{} {
	mfbyte := []byte(metaField)
	v := json.Get(mfbyte, requestedKey).ToString()
	i := interface{}(v)
	return i
}
