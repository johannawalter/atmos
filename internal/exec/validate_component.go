package exec

import (
	"fmt"
	c "github.com/cloudposse/atmos/pkg/config"
	u "github.com/cloudposse/atmos/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
)

// ExecuteValidateComponentCmd executes `validate component` command
func ExecuteValidateComponentCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments. The command requires one argument 'componentName'")
	}

	componentName := args[0]

	flags := cmd.Flags()

	stack, err := flags.GetString("stack")
	if err != nil {
		return err
	}

	schemaPath, err := flags.GetString("schema-path")
	if err != nil {
		return err
	}

	schemaType, err := flags.GetString("schema-type")
	if err != nil {
		return err
	}

	_, err = ExecuteValidateComponent(componentName, stack, schemaPath, schemaType)
	if err != nil {
		return err
	}

	return nil
}

// ExecuteValidateComponent validates a component in a stack using JsonSchema, OPA or CUE schema documents
func ExecuteValidateComponent(componentName string, stack string, schemaPath string, schemaType string) (bool, error) {
	var configAndStacksInfo c.ConfigAndStacksInfo
	configAndStacksInfo.ComponentFromArg = componentName
	configAndStacksInfo.Stack = stack

	configAndStacksInfo.ComponentType = "terraform"
	configAndStacksInfo, err := ProcessStacks(configAndStacksInfo, true)
	if err != nil {
		u.PrintErrorVerbose(err)
		configAndStacksInfo.ComponentType = "helmfile"
		configAndStacksInfo, err = ProcessStacks(configAndStacksInfo, true)
		if err != nil {
			return false, err
		}
	}

	componentSection := configAndStacksInfo.ComponentSection

	return ValidateComponent(componentName, componentSection, schemaPath, schemaType)
}

// ValidateComponent validates the component config using JsonSchema, OPA or CUE schema documents
func ValidateComponent(componentName string, componentSection any, schemaPath string, schemaType string) (bool, error) {
	ok := true
	var err error

	if schemaPath != "" && schemaType != "" {
		fmt.Println()
		u.PrintInfo(fmt.Sprintf("Validating the component '%s' using '%s' file '%s'", componentName, schemaType, schemaPath))

		ok, err = validateComponentInternal(componentSection, schemaPath, schemaType)
		if err != nil {
			return false, err
		}
	} else {
		validations, err := FindValidationSection(componentSection.(map[string]any))
		if err != nil {
			return false, err
		}

		for _, v := range validations {
			schemaPath = v.SchemaPath
			schemaType = v.SchemaType

			fmt.Println()
			u.PrintInfo(fmt.Sprintf("Validating the component '%s' using '%s' file '%s'", componentName, schemaType, schemaPath))
			if v.Description != "" {
				u.PrintMessage(v.Description)
			}

			ok2, err := validateComponentInternal(componentSection, schemaPath, schemaType)
			if err != nil {
				return false, err
			}
			if !ok2 {
				ok = false
			}
		}
	}

	fmt.Println()

	return ok, nil
}

func validateComponentInternal(componentSection any, schemaPath string, schemaType string) (bool, error) {
	if schemaType != "jsonschema" && schemaType != "opa" && schemaType != "cue" {
		return false, fmt.Errorf("invalid schema type '%s'. Supported types: jsonschema, opa, cue", schemaType)
	}

	// Check if the file pointed to by 'schemaPath' exists.
	// If not, join it with the schemas `base_path` from the CLI config
	var filePath string
	if u.FileExists(schemaPath) {
		filePath = schemaPath
	} else {
		switch schemaType {
		case "jsonschema":
			{
				filePath = path.Join(c.Config.BasePath, c.Config.Schemas.JsonSchema.BasePath, schemaPath)
			}
		case "opa":
			{
				filePath = path.Join(c.Config.BasePath, c.Config.Schemas.Opa.BasePath, schemaPath)
			}
		case "cue":
			{
				filePath = path.Join(c.Config.BasePath, c.Config.Schemas.Cue.BasePath, schemaPath)
			}
		}

		if !u.FileExists(filePath) {
			return false, fmt.Errorf("the file '%s' does not exist for schema type '%s'", schemaPath, schemaType)
		}
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	schemaText := string(fileContent)
	var ok bool

	switch schemaType {
	case "jsonschema":
		{
			ok, err = ValidateWithJsonSchema(componentSection, filePath, schemaText)
			if err != nil {
				return false, err
			}
		}
	case "opa":
		{
			ok, err = ValidateWithOpa(componentSection, filePath, schemaText)
			if err != nil {
				return false, err
			}
		}
	case "cue":
		{
			ok, err = ValidateWithCue(componentSection, filePath, schemaText)
			if err != nil {
				return false, err
			}
		}
	}

	return ok, nil
}

// FindValidationSection finds 'validation' section in the component config
func FindValidationSection(componentSection map[string]any) (c.Validation, error) {
	validationSection := map[any]any{}

	if i, ok := componentSection["settings"].(map[any]any); ok {
		if i2, ok2 := i["validation"].(map[any]any); ok2 {
			validationSection = i2
		}
	}

	var result c.Validation

	err := mapstructure.Decode(validationSection, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
