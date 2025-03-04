package merge

import (
	u "github.com/cloudposse/atmos/pkg/utils"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// MergeWithOptions takes a list of maps of interface and options as input and returns a single map with the merged contents
func MergeWithOptions(inputs []map[any]any, appendSlice, sliceDeepCopy bool) (map[any]any, error) {
	merged := map[any]any{}

	for index := range inputs {
		current := inputs[index]

		if len(current) == 0 {
			continue
		}

		// Due to a bug in `mergo.Merge`
		// (Note: in the `for` loop, it DOES modify the source of the previous loop iteration if it's a complex map and `mergo` gets a pointer to it,
		// not only the destination of the current loop iteration),
		// we don't give it our maps directly; we convert them to YAML strings and then back to `Go` maps,
		// so `mergo` does not have access to the original pointers
		yamlCurrent, err := yaml.Marshal(current)
		if err != nil {
			u.PrintErrorToStdError(err)
			return nil, err
		}

		var dataCurrent map[any]any
		if err = yaml.Unmarshal(yamlCurrent, &dataCurrent); err != nil {
			u.PrintErrorToStdError(err)
			return nil, err
		}

		var opts []func(*mergo.Config)
		opts = append(opts, mergo.WithOverride, mergo.WithTypeCheck)

		// This was fixed/broken in https://github.com/imdario/mergo/pull/231/files
		// It was released in https://github.com/imdario/mergo/releases/tag/v0.3.14
		// It was not working before in `github.com/imdario/mergo` so we need to disable it in our code
		// opts = append(opts, mergo.WithOverwriteWithEmptyValue)

		if appendSlice {
			opts = append(opts, mergo.WithAppendSlice)
		}

		if sliceDeepCopy {
			opts = append(opts, mergo.WithSliceDeepCopy)
		}

		if err = mergo.Merge(&merged, dataCurrent, opts...); err != nil {
			u.PrintErrorToStdError(err)
			return nil, err
		}
	}

	return merged, nil
}

// Merge takes a list of maps of interface as input and returns a single map with the merged contents
func Merge(inputs []map[any]any) (map[any]any, error) {
	return MergeWithOptions(inputs, false, false)
}
