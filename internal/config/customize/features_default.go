package customize

import (
	"os"
	"reflect"
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/photoprism/photoprism/pkg/clean"
)

// DefaultFeatures holds the baseline feature flags applied to new settings instances.
// Values may be overridden at startup via PHOTOPRISM_DISABLE_FEATURES.
var DefaultFeatures FeatureSettings

// init wires DefaultFeatures from defaults and environment overrides.
func init() {
	DefaultFeatures = initDefaultFeatures()
}

// NewFeatures returns a copy of the default feature flags so callers can mutate
// the result without changing the shared defaults.
func NewFeatures() FeatureSettings {
	return DefaultFeatures
}

// liteEnabledFeatures defines the subset of UI/API features enabled in PhotoPrism Lite.
var liteEnabledFeatures = map[string]struct{}{
	"albums":    {},
	"archive":   {},
	"delete":    {},
	"download":  {},
	"export":    {},
	"edit":      {},
	"batchedit": {},
	"favorites": {},
	"folders":   {},
	"import":    {},
	"library":   {},
	"search":    {},
	"account":   {},
	"settings":  {},
	"share":     {},
	"upload":    {},
	"videos":    {},
}

// initDefaultFeatures builds the package-level defaults and applies any disable
// list supplied via PHOTOPRISM_DISABLE_FEATURES.
func initDefaultFeatures() FeatureSettings {
	features := FeatureSettings{}

	disabled := buildDisabledSet(os.Getenv("PHOTOPRISM_DISABLE_FEATURES"))

	val := reflect.ValueOf(&features).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldKey := clean.FieldNameLower(field.Name)
		if _, enabled := liteEnabledFeatures[fieldKey]; !enabled {
			continue
		}

		if len(disabled) > 0 {
			candidates := []string{
				clean.FieldNameLower(field.Tag.Get("json")),
				fieldKey,
			}

			if isDisabled(disabled, candidates) {
				continue
			}
		}

		val.Field(i).SetBool(true)
	}

	return features
}

// buildDisabledSet tokenizes the disable list into normalized feature names.
func buildDisabledSet(disable string) map[string]struct{} {
	if disable == "" {
		return nil
	}

	parts := strings.FieldsFunc(disable, func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\t' || r == '\n'
	})

	disabled := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		name := clean.FieldNameLower(part)
		if name == "" {
			continue
		}

		disabled[name] = struct{}{}
		disabled[clean.FieldNameLower(inflection.Singular(name))] = struct{}{}
		disabled[clean.FieldNameLower(inflection.Plural(name))] = struct{}{}
	}

	return disabled
}

// isDisabled checks whether any of the candidate field names are present in the disabled set.
func isDisabled(disabled map[string]struct{}, candidates []string) bool {
	if len(disabled) == 0 {
		return false
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}

		if _, ok := disabled[candidate]; ok {
			return true
		}
	}

	return false
}
