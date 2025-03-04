package config

// CliConfiguration structure represents schema for `atmos.yaml` CLI config
type CliConfiguration struct {
	BasePath                      string       `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
	Components                    Components   `yaml:"components" json:"components" mapstructure:"components"`
	Stacks                        Stacks       `yaml:"stacks" json:"stacks" mapstructure:"stacks"`
	Workflows                     Workflows    `yaml:"workflows" json:"workflows" mapstructure:"workflows"`
	Logs                          Logs         `yaml:"logs" json:"logs" mapstructure:"logs"`
	Commands                      []Command    `yaml:"commands" json:"commands" mapstructure:"commands"`
	Integrations                  Integrations `yaml:"integrations" json:"integrations" mapstructure:"integrations"`
	Schemas                       Schemas      `yaml:"schemas" json:"schemas" mapstructure:"schemas"`
	Initialized                   bool         `yaml:"initialized" json:"initialized" mapstructure:"initialized"`
	StacksBaseAbsolutePath        string       `yaml:"stacksBaseAbsolutePath" json:"stacksBaseAbsolutePath"`
	IncludeStackAbsolutePaths     []string     `yaml:"includeStackAbsolutePaths" json:"includeStackAbsolutePaths"`
	ExcludeStackAbsolutePaths     []string     `yaml:"excludeStackAbsolutePaths" json:"excludeStackAbsolutePaths"`
	TerraformDirAbsolutePath      string       `yaml:"terraformDirAbsolutePath" json:"terraformDirAbsolutePath"`
	HelmfileDirAbsolutePath       string       `yaml:"helmfileDirAbsolutePath" json:"helmfileDirAbsolutePath"`
	StackConfigFilesRelativePaths []string     `yaml:"stackConfigFilesRelativePaths" json:"stackConfigFilesRelativePaths"`
	StackConfigFilesAbsolutePaths []string     `yaml:"stackConfigFilesAbsolutePaths" json:"stackConfigFilesAbsolutePaths"`
	StackType                     string       `yaml:"stackType" json:"StackType"`
}

type Terraform struct {
	BasePath                string `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
	ApplyAutoApprove        bool   `yaml:"apply_auto_approve" json:"apply_auto_approve" mapstructure:"apply_auto_approve"`
	DeployRunInit           bool   `yaml:"deploy_run_init" json:"deploy_run_init" mapstructure:"deploy_run_init"`
	InitRunReconfigure      bool   `yaml:"init_run_reconfigure" json:"init_run_reconfigure" mapstructure:"init_run_reconfigure"`
	AutoGenerateBackendFile bool   `yaml:"auto_generate_backend_file" json:"auto_generate_backend_file" mapstructure:"auto_generate_backend_file"`
}

type Helmfile struct {
	BasePath              string `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
	UseEKS                bool   `yaml:"use_eks" json:"use_eks" mapstructure:"use_eks"`
	KubeconfigPath        string `yaml:"kubeconfig_path" json:"kubeconfig_path" mapstructure:"kubeconfig_path"`
	HelmAwsProfilePattern string `yaml:"helm_aws_profile_pattern" json:"helm_aws_profile_pattern" mapstructure:"helm_aws_profile_pattern"`
	ClusterNamePattern    string `yaml:"cluster_name_pattern" json:"cluster_name_pattern" mapstructure:"cluster_name_pattern"`
}

type Components struct {
	Terraform Terraform `yaml:"terraform" json:"terraform" mapstructure:"terraform"`
	Helmfile  Helmfile  `yaml:"helmfile" json:"helmfile" mapstructure:"helmfile"`
}

type Stacks struct {
	BasePath      string   `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
	IncludedPaths []string `yaml:"included_paths" json:"included_paths" mapstructure:"included_paths"`
	ExcludedPaths []string `yaml:"excluded_paths" json:"excluded_paths" mapstructure:"excluded_paths"`
	NamePattern   string   `yaml:"name_pattern" json:"name_pattern" mapstructure:"name_pattern"`
}

type Workflows struct {
	BasePath string `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
}

type Logs struct {
	Verbose bool `yaml:"verbose" json:"verbose" mapstructure:"verbose"`
	Colors  bool `yaml:"colors" json:"colors" mapstructure:"colors"`
}

type Context struct {
	Namespace     string   `yaml:"namespace" json:"namespace" mapstructure:"namespace"`
	Tenant        string   `yaml:"tenant" json:"tenant" mapstructure:"tenant"`
	Environment   string   `yaml:"environment" json:"environment" mapstructure:"environment"`
	Stage         string   `yaml:"stage" json:"stage" mapstructure:"stage"`
	Region        string   `yaml:"region" json:"region" mapstructure:"region"`
	Component     string   `yaml:"component" json:"component" mapstructure:"component"`
	BaseComponent string   `yaml:"base_component" json:"base_component" mapstructure:"base_component"`
	ComponentPath string   `yaml:"component_path" json:"component_path" mapstructure:"component_path"`
	Workspace     string   `yaml:"workspace" json:"workspace" mapstructure:"workspace"`
	Attributes    []string `yaml:"attributes" json:"attributes" mapstructure:"attributes"`
}

type ArgsAndFlagsInfo struct {
	AdditionalArgsAndFlags  []string
	SubCommand              string
	SubCommand2             string
	ComponentFromArg        string
	GlobalOptions           []string
	TerraformDir            string
	HelmfileDir             string
	ConfigDir               string
	StacksDir               string
	WorkflowsDir            string
	BasePath                string
	DeployRunInit           string
	InitRunReconfigure      string
	AutoGenerateBackendFile string
	UseTerraformPlan        bool
	PlanFile                string
	DryRun                  bool
	SkipInit                bool
	NeedHelp                bool
	JsonSchemaDir           string
	OpaDir                  string
	CueDir                  string
	RedirectStdErr          string
}

type ConfigAndStacksInfo struct {
	StackFromArg                  string
	Stack                         string
	StackFile                     string
	ComponentType                 string
	ComponentFromArg              string
	Component                     string
	ComponentFolderPrefix         string
	ComponentFolderPrefixReplaced string
	BaseComponentPath             string
	BaseComponent                 string
	FinalComponent                string
	Command                       string
	SubCommand                    string
	SubCommand2                   string
	ComponentSection              map[string]any
	ComponentVarsSection          map[any]any
	ComponentEnvSection           map[any]any
	ComponentEnvList              []string
	ComponentBackendSection       map[any]any
	ComponentBackendType          string
	AdditionalArgsAndFlags        []string
	GlobalOptions                 []string
	BasePath                      string
	TerraformDir                  string
	HelmfileDir                   string
	ConfigDir                     string
	StacksDir                     string
	WorkflowsDir                  string
	Context                       Context
	ContextPrefix                 string
	DeployRunInit                 string
	InitRunReconfigure            string
	AutoGenerateBackendFile       string
	UseTerraformPlan              bool
	PlanFile                      string
	DryRun                        bool
	SkipInit                      bool
	ComponentInheritanceChain     []string
	NeedHelp                      bool
	ComponentIsAbstract           bool
	ComponentMetadataSection      map[any]any
	TerraformWorkspace            string
	JsonSchemaDir                 string
	OpaDir                        string
	CueDir                        string
	AtmosCliConfigPath            string
	AtmosBasePath                 string
	RedirectStdErr                string
}

// Workflows

type WorkflowStep struct {
	Name    string `yaml:"name" json:"name" mapstructure:"name"`
	Command string `yaml:"command" json:"command" mapstructure:"command"`
	Stack   string `yaml:"stack,omitempty" json:"stack,omitempty" mapstructure:"stack"`
	Type    string `yaml:"type,omitempty" json:"type,omitempty" mapstructure:"type"`
}

type WorkflowDefinition struct {
	Description string         `yaml:"description,omitempty" json:"description,omitempty" mapstructure:"description"`
	Steps       []WorkflowStep `yaml:"steps" json:"steps" mapstructure:"steps"`
	Stack       string         `yaml:"stack,omitempty" json:"stack,omitempty" mapstructure:"stack"`
}

type WorkflowConfig map[string]WorkflowDefinition

type WorkflowFile map[string]WorkflowConfig

// EKS update-kubeconfig

type AwsEksUpdateKubeconfigContext struct {
	Component   string
	Stack       string
	Profile     string
	ClusterName string
	Kubeconfig  string
	RoleArn     string
	DryRun      bool
	Verbose     bool
	Alias       string
	Namespace   string
	Tenant      string
	Environment string
	Stage       string
	Region      string
}

// Component vendoring (`component.yaml` file)

type VendorComponentSource struct {
	Type          string   `yaml:"type" json:"type" mapstructure:"type"`
	Uri           string   `yaml:"uri" json:"uri" mapstructure:"uri"`
	Version       string   `yaml:"version" json:"version" mapstructure:"version"`
	IncludedPaths []string `yaml:"included_paths" json:"included_paths" mapstructure:"included_paths"`
	ExcludedPaths []string `yaml:"excluded_paths" json:"excluded_paths" mapstructure:"excluded_paths"`
}

type VendorComponentMixins struct {
	Type     string `yaml:"type" json:"type" mapstructure:"type"`
	Uri      string `yaml:"uri" json:"uri" mapstructure:"uri"`
	Version  string `yaml:"version" json:"version" mapstructure:"version"`
	Filename string `yaml:"filename" json:"filename" mapstructure:"filename"`
}

type VendorComponentSpec struct {
	Source VendorComponentSource
	Mixins []VendorComponentMixins
}

type VendorComponentMetadata struct {
	Name        string `yaml:"name" json:"name" mapstructure:"name"`
	Description string `yaml:"description" json:"description" mapstructure:"description"`
}

type VendorComponentConfig struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion" mapstructure:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind" mapstructure:"kind"`
	Metadata   VendorComponentMetadata
	Spec       VendorComponentSpec `yaml:"spec" json:"spec" mapstructure:"spec"`
}

// Custom CLI commands

type Command struct {
	Name            string                 `yaml:"name" json:"name" mapstructure:"name"`
	Description     string                 `yaml:"description" json:"description" mapstructure:"description"`
	Env             []CommandEnv           `yaml:"env" json:"env" mapstructure:"env"`
	Arguments       []CommandArgument      `yaml:"arguments" json:"arguments" mapstructure:"arguments"`
	Flags           []CommandFlag          `yaml:"flags" json:"flags" mapstructure:"flags"`
	ComponentConfig CommandComponentConfig `yaml:"component_config" json:"component_config" mapstructure:"component_config"`
	Steps           []string               `yaml:"steps" json:"steps" mapstructure:"steps"`
	Commands        []Command              `yaml:"commands" json:"commands" mapstructure:"commands"`
	Verbose         bool                   `yaml:"verbose" json:"verbose" mapstructure:"verbose"`
}

type CommandArgument struct {
	Name        string `yaml:"name" json:"name" mapstructure:"name"`
	Description string `yaml:"description" json:"description" mapstructure:"description"`
}

type CommandFlag struct {
	Name        string `yaml:"name" json:"name" mapstructure:"name"`
	Shorthand   string `yaml:"shorthand" json:"shorthand" mapstructure:"shorthand"`
	Type        string `yaml:"type" json:"type" mapstructure:"type"`
	Description string `yaml:"description" json:"description" mapstructure:"description"`
	Usage       string `yaml:"usage" json:"usage" mapstructure:"usage"`
	Required    bool   `yaml:"required" json:"required" mapstructure:"required"`
}

type CommandEnv struct {
	Key          string `yaml:"key" json:"key" mapstructure:"key"`
	Value        string `yaml:"value" json:"value" mapstructure:"value"`
	ValueCommand string `yaml:"valueCommand" json:"valueCommand" mapstructure:"valueCommand"`
}

type CommandComponentConfig struct {
	Component string `yaml:"component" json:"component" mapstructure:"component"`
	Stack     string `yaml:"stack" json:"stack" mapstructure:"stack"`
}

// Integrations

type Integrations struct {
	Atlantis Atlantis `yaml:"atlantis" json:"atlantis" mapstructure:"atlantis"`
}

// Atlantis integration

type Atlantis struct {
	Path              string                           `yaml:"path" json:"path" mapstructure:"path"`
	ConfigTemplates   map[string]AtlantisRepoConfig    `yaml:"config_templates" json:"config_templates" mapstructure:"config_templates"`
	ProjectTemplates  map[string]AtlantisProjectConfig `yaml:"project_templates" json:"project_templates" mapstructure:"project_templates"`
	WorkflowTemplates map[string]any                   `yaml:"workflow_templates" json:"workflow_templates" mapstructure:"workflow_templates"`
}

type AtlantisRepoConfig struct {
	Version                   int      `yaml:"version" json:"version" mapstructure:"version"`
	Automerge                 bool     `yaml:"automerge" json:"automerge" mapstructure:"automerge"`
	DeleteSourceBranchOnMerge bool     `yaml:"delete_source_branch_on_merge" json:"delete_source_branch_on_merge" mapstructure:"delete_source_branch_on_merge"`
	ParallelPlan              bool     `yaml:"parallel_plan" json:"parallel_plan" mapstructure:"parallel_plan"`
	ParallelApply             bool     `yaml:"parallel_apply" json:"parallel_apply" mapstructure:"parallel_apply"`
	AllowedRegexpPrefixes     []string `yaml:"allowed_regexp_prefixes" json:"allowed_regexp_prefixes" mapstructure:"allowed_regexp_prefixes"`
}

type AtlantisProjectConfig struct {
	Name                      string                        `yaml:"name" json:"name" mapstructure:"name"`
	Workspace                 string                        `yaml:"workspace" json:"workspace" mapstructure:"workspace"`
	Workflow                  string                        `yaml:"workflow,omitempty" json:"workflow,omitempty" mapstructure:"workflow"`
	Dir                       string                        `yaml:"dir" json:"dir" mapstructure:"dir"`
	TerraformVersion          string                        `yaml:"terraform_version" json:"terraform_version" mapstructure:"terraform_version"`
	DeleteSourceBranchOnMerge bool                          `yaml:"delete_source_branch_on_merge" json:"delete_source_branch_on_merge" mapstructure:"delete_source_branch_on_merge"`
	Autoplan                  AtlantisProjectAutoplanConfig `yaml:"autoplan" json:"autoplan" mapstructure:"autoplan"`
	ApplyRequirements         []string                      `yaml:"apply_requirements" json:"apply_requirements" mapstructure:"apply_requirements"`
}

type AtlantisProjectAutoplanConfig struct {
	Enabled      bool     `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	WhenModified []string `yaml:"when_modified" json:"when_modified" mapstructure:"when_modified"`
}

type AtlantisConfigOutput struct {
	Version                   int                     `yaml:"version" json:"version" mapstructure:"version"`
	Automerge                 bool                    `yaml:"automerge" json:"automerge" mapstructure:"automerge"`
	DeleteSourceBranchOnMerge bool                    `yaml:"delete_source_branch_on_merge" json:"delete_source_branch_on_merge" mapstructure:"delete_source_branch_on_merge"`
	ParallelPlan              bool                    `yaml:"parallel_plan" json:"parallel_plan" mapstructure:"parallel_plan"`
	ParallelApply             bool                    `yaml:"parallel_apply" json:"parallel_apply" mapstructure:"parallel_apply"`
	AllowedRegexpPrefixes     []string                `yaml:"allowed_regexp_prefixes" json:"allowed_regexp_prefixes" mapstructure:"allowed_regexp_prefixes"`
	Projects                  []AtlantisProjectConfig `yaml:"projects" json:"projects" mapstructure:"projects"`
	Workflows                 map[string]any          `yaml:"workflows,omitempty" json:"workflows,omitempty" mapstructure:"workflows"`
}

// Validation schemas

type JsonSchema struct {
	BasePath string `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
}

type Cue struct {
	BasePath string `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
}

type Opa struct {
	BasePath string `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
}

type Schemas struct {
	JsonSchema JsonSchema `yaml:"jsonschema" json:"jsonschema" mapstructure:"jsonschema"`
	Cue        Cue        `yaml:"cue" json:"cue" mapstructure:"cue"`
	Opa        Opa        `yaml:"opa" json:"opa" mapstructure:"opa"`
}

type ValidationItem struct {
	SchemaType  string `yaml:"schema_type" json:"schema_type" mapstructure:"schema_type"`
	SchemaPath  string `yaml:"schema_path" json:"schema_path" mapstructure:"schema_path"`
	Description string `yaml:"description" json:"description" mapstructure:"description"`
	Disabled    bool   `yaml:"disabled" json:"disabled" mapstructure:"disabled"`
}

type Validation map[string]ValidationItem

// Affected Atmos components and stacks given two Git commits

type Affected struct {
	Component       string `yaml:"component" json:"component" mapstructure:"component"`
	ComponentType   string `yaml:"component_type" json:"component_type" mapstructure:"component_type"`
	ComponentPath   string `yaml:"component_path" json:"component_path" mapstructure:"component_path"`
	Namespace       string `yaml:"namespace,omitempty" json:"namespace,omitempty" mapstructure:"namespace"`
	Tenant          string `yaml:"tenant,omitempty" json:"tenant,omitempty" mapstructure:"tenant"`
	Environment     string `yaml:"environment,omitempty" json:"environment,omitempty" mapstructure:"environment"`
	Stage           string `yaml:"stage,omitempty" json:"stage,omitempty" mapstructure:"stage"`
	Stack           string `yaml:"stack" json:"stack" mapstructure:"stack"`
	StackSlug       string `yaml:"stack_slug" json:"stack_slug" mapstructure:"stack_slug"`
	SpaceliftStack  string `yaml:"spacelift_stack,omitempty" json:"spacelift_stack,omitempty" mapstructure:"spacelift_stack"`
	AtlantisProject string `yaml:"atlantis_project,omitempty" json:"atlantis_project,omitempty" mapstructure:"atlantis_project"`
	Affected        string `yaml:"affected" json:"affected" mapstructure:"affected"`
}

type BaseComponentConfig struct {
	BaseComponentVars                      map[any]any
	BaseComponentSettings                  map[any]any
	BaseComponentEnv                       map[any]any
	FinalBaseComponentName                 string
	BaseComponentCommand                   string
	BaseComponentBackendType               string
	BaseComponentBackendSection            map[any]any
	BaseComponentRemoteStateBackendType    string
	BaseComponentRemoteStateBackendSection map[any]any
	ComponentInheritanceChain              []string
}

// Stack imports (`import` section)

type StackImport struct {
	Path    string         `yaml:"path" json:"path" mapstructure:"path"`
	Context map[string]any `yaml:"context" json:"context" mapstructure:"context"`
}

// Dependencies

type DependsOn map[any]Context

type Dependant struct {
	Component       string `yaml:"component" json:"component" mapstructure:"component"`
	ComponentType   string `yaml:"component_type" json:"component_type" mapstructure:"component_type"`
	ComponentPath   string `yaml:"component_path" json:"component_path" mapstructure:"component_path"`
	Namespace       string `yaml:"namespace,omitempty" json:"namespace,omitempty" mapstructure:"namespace"`
	Tenant          string `yaml:"tenant,omitempty" json:"tenant,omitempty" mapstructure:"tenant"`
	Environment     string `yaml:"environment,omitempty" json:"environment,omitempty" mapstructure:"environment"`
	Stage           string `yaml:"stage,omitempty" json:"stage,omitempty" mapstructure:"stage"`
	Stack           string `yaml:"stack" json:"stack" mapstructure:"stack"`
	StackSlug       string `yaml:"stack_slug" json:"stack_slug" mapstructure:"stack_slug"`
	SpaceliftStack  string `yaml:"spacelift_stack,omitempty" json:"spacelift_stack,omitempty" mapstructure:"spacelift_stack"`
	AtlantisProject string `yaml:"atlantis_project,omitempty" json:"atlantis_project,omitempty" mapstructure:"atlantis_project"`
}

// Settings

type Settings struct {
	DependsOn DependsOn `yaml:"depends_on" json:"depends_on" mapstructure:"depends_on"`
}
