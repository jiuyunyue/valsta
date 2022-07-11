package vals

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/jiuyunyue/valsta/src/types"
)

var (
	//go:embed all.json
	UptimeJson []byte // nolint: golint

	Uptime types.Uptime
)

func init() {
	if err := json.Unmarshal(UptimeJson, &Uptime); err != nil {
		panic(err)
	}
}
