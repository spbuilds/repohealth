package report

import (
	"encoding/json"
	"io"

	"github.com/spbuilds/repohealth/internal/model"
)

// JSON renders the report as JSON.
func JSON(w io.Writer, r *model.Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
