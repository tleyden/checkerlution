package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/logg"
	"io"
	"time"
)

func handleChange(reader io.Reader) string {
	changes := make(map[string]interface{})
	decoder := json.NewDecoder(reader)
	decoder.Decode(&changes)
	// logg.LogTo("DEBUG", "changes: %v", changes)
	lastSeq := changes["last_seq"]
	// logg.LogTo("DEBUG", "lastSeq: %v", lastSeq)
	lastSeqAsString := lastSeq.(string)
	logg.LogTo("DEBUG", "return lastSeq: %v", lastSeq)
	time.Sleep(time.Second * 30)
	return lastSeqAsString
}
