package exec

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"

	cfg "github.com/cloudposse/atmos/pkg/config"
	u "github.com/cloudposse/atmos/pkg/utils"
)

var (
	localRepoIsNotGitRepoError  = errors.New("the local repo is not a Git repository. Check that it was initialized and has '.git' folder")
	remoteRepoIsNotGitRepoError = errors.New("the target remote repo is not a Git repository. Check that it was initialized and has '.git' folder")
)

// ExecuteDescribeAffectedWithTargetRepoClone clones the remote repo using `ref` or `sha`, processes stack configs
// and returns a list of the affected Atmos components and stacks given two Git commits
func ExecuteDescribeAffectedWithTargetRepoClone(
	cliConfig cfg.CliConfiguration,
	ref string,
	sha string,
	sshKeyPath string,
	sshKeyPassword string,
	verbose bool,
) ([]cfg.Affected, error) {

	localRepo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: false,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "%v", localRepoIsNotGitRepoError)
	}

	// Get the Git config of the local repo
	localRepoConfig, err := localRepo.Config()
	if err != nil {
		return nil, errors.Wrapf(err, "%v", localRepoIsNotGitRepoError)
	}

	localRepoWorktree, err := localRepo.Worktree()
	if err != nil {
		return nil, errors.Wrapf(err, "%v", localRepoIsNotGitRepoError)
	}

	localRepoPath := localRepoWorktree.Filesystem.Root()

	// Get the remotes of the local repo
	keys := []string{}
	for k := range localRepoConfig.Remotes {
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		return nil, localRepoIsNotGitRepoError
	}

	// Get the origin URL of the current remoteRepo
	remoteUrls := localRepoConfig.Remotes[keys[0]].URLs
	if len(remoteUrls) == 0 {
		return nil, localRepoIsNotGitRepoError
	}

	repoUrl := remoteUrls[0]
	if repoUrl == "" {
		return nil, localRepoIsNotGitRepoError
	}

	// Clone the remote repo
	// https://git-scm.com/book/en/v2/Git-Internals-Git-References
	// https://git-scm.com/docs/git-show-ref
	// https://github.com/go-git/go-git/tree/master/_examples
	// https://stackoverflow.com/questions/56810719/how-to-checkout-a-specific-sha-in-a-git-repo-using-golang
	// https://golang.hotexamples.com/examples/gopkg.in.src-d.go-git.v4.plumbing/-/ReferenceName/golang-referencename-function-examples.html
	// https://stackoverflow.com/questions/58905690/how-to-identify-which-files-have-changed-between-git-commits
	// https://github.com/src-d/go-git/issues/604

	// Create a temp dir to clone the remote repo to
	tempDir, err := os.MkdirTemp("", strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return nil, err
	}

	defer removeTempDir(tempDir)

	u.PrintInfoVerbose(verbose, fmt.Sprintf("\nCloning repo '%s' into the temp dir '%s'", repoUrl, tempDir))

	cloneOptions := git.CloneOptions{
		URL:          repoUrl,
		NoCheckout:   false,
		SingleBranch: false,
	}

	// If `ref` flag is not provided, it will clone the HEAD of the default branch
	if ref != "" {
		cloneOptions.ReferenceName = plumbing.ReferenceName(ref)
		u.PrintInfoVerbose(verbose, fmt.Sprintf("\nChecking out Git ref '%s' ...\n", ref))
	} else {
		u.PrintInfoVerbose(verbose, "\nChecking out the HEAD of the default branch ...\n")
	}

	if verbose {
		cloneOptions.Progress = os.Stdout
	}

	// Clone private repos using SSH
	// https://gist.github.com/efontan/e8e8818dc0845d3bd7bf1343c984ae7b
	// https://github.com/src-d/go-git/issues/550
	if sshKeyPath != "" {
		sshKeyContent, err := os.ReadFile(sshKeyPath)
		if err != nil {
			return nil, err
		}

		sshPublicKey, err := ssh.NewPublicKeys("git", sshKeyContent, sshKeyPassword)
		if err != nil {
			return nil, err
		}

		// Use the SSH key to clone the repo
		cloneOptions.Auth = sshPublicKey

		// Update the repo URL to SSH format
		// https://mirrors.edge.kernel.org/pub/software/scm/git/docs/git-clone.html
		cloneOptions.URL = strings.Replace(repoUrl, "https://", "ssh://", 1)
	}

	remoteRepo, err := git.PlainClone(tempDir, false, &cloneOptions)
	if err != nil {
		return nil, err
	}

	remoteRepoHead, err := remoteRepo.Head()
	if err != nil {
		return nil, err
	}

	if ref != "" {
		u.PrintInfoVerbose(verbose, fmt.Sprintf("\nChecked out Git ref '%s'\n", ref))
	} else {
		u.PrintInfoVerbose(verbose, fmt.Sprintf("\nChecked out Git ref '%s'\n", remoteRepoHead.Name()))
	}

	// Check if a commit SHA was provided and checkout the repo at that commit SHA
	if sha != "" {
		u.PrintInfoVerbose(verbose, fmt.Sprintf("\nChecking out commit SHA '%s' ...\n", sha))

		w, err := remoteRepo.Worktree()
		if err != nil {
			return nil, err
		}

		checkoutOptions := git.CheckoutOptions{
			Hash:   plumbing.NewHash(sha),
			Create: false,
			Force:  true,
			Keep:   false,
		}

		err = w.Checkout(&checkoutOptions)
		if err != nil {
			return nil, err
		}

		u.PrintInfoVerbose(verbose, fmt.Sprintf("\nChecked out commit SHA '%s'\n", sha))
	}

	affected, err := executeDescribeAffected(cliConfig, localRepoPath, tempDir, localRepo, remoteRepo, verbose)
	if err != nil {
		return nil, err
	}

	return affected, nil
}

// ExecuteDescribeAffectedWithTargetRepoPath uses `repo-path` to access the target repo, processes stack configs
// and returns a list of the affected Atmos components and stacks given two Git commits
func ExecuteDescribeAffectedWithTargetRepoPath(
	cliConfig cfg.CliConfiguration,
	repoPath string,
	verbose bool,
) ([]cfg.Affected, error) {

	localRepo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: false,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "%v", localRepoIsNotGitRepoError)
	}

	// Check the Git config of the local repo
	_, err = localRepo.Config()
	if err != nil {
		return nil, errors.Wrapf(err, "%v", localRepoIsNotGitRepoError)
	}

	localRepoWorktree, err := localRepo.Worktree()
	if err != nil {
		return nil, errors.Wrapf(err, "%v", localRepoIsNotGitRepoError)
	}

	localRepoPath := localRepoWorktree.Filesystem.Root()

	remoteRepo, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{
		DetectDotGit:          false,
		EnableDotGitCommonDir: false,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "%v", remoteRepoIsNotGitRepoError)
	}

	// Check the Git config of the remote target repo
	_, err = remoteRepo.Config()
	if err != nil {
		return nil, errors.Wrapf(err, "%v", remoteRepoIsNotGitRepoError)
	}

	affected, err := executeDescribeAffected(cliConfig, localRepoPath, repoPath, localRepo, remoteRepo, verbose)
	if err != nil {
		return nil, err
	}

	return affected, nil
}

func executeDescribeAffected(
	cliConfig cfg.CliConfiguration,
	localRepoFileSystemPath string,
	remoteRepoFileSystemPath string,
	localRepo *git.Repository,
	remoteRepo *git.Repository,
	verbose bool,
) ([]cfg.Affected, error) {

	localRepoHead, err := localRepo.Head()
	if err != nil {
		return nil, err
	}

	remoteRepoHead, err := remoteRepo.Head()
	if err != nil {
		return nil, err
	}

	if verbose {
		u.PrintInfo(fmt.Sprintf("Current working repo HEAD: %s", localRepoHead))
		u.PrintInfo(fmt.Sprintf("Remote repo HEAD: %s", remoteRepoHead))
	}

	currentStacks, err := ExecuteDescribeStacks(cliConfig, "", nil, nil, nil, false)
	if err != nil {
		return nil, err
	}

	localRepoFileSystemPathAbs, err := filepath.Abs(localRepoFileSystemPath)
	if err != nil {
		return nil, err
	}

	basePath := cliConfig.BasePath

	// Handle `atmos` absolute base path.
	// Absolute base path can be set in the `base_path` attribute in `atmos.yaml`, or using the ENV var `ATMOS_BASE_PATH` (as it's done in `geodesic`)
	// If the `atmos` base path is absolute, find the relative path between the local repo path and the `atmos` base path.
	// This relative path (the difference) is then used below to join with the remote (cloned) repo path.
	if path.IsAbs(basePath) {
		basePath, err = filepath.Rel(localRepoFileSystemPathAbs, basePath)
		if err != nil {
			return nil, err
		}
	}

	// Update paths to point to the cloned remote repo dir
	cliConfig.StacksBaseAbsolutePath = path.Join(remoteRepoFileSystemPath, basePath, cliConfig.Stacks.BasePath)
	cliConfig.TerraformDirAbsolutePath = path.Join(remoteRepoFileSystemPath, basePath, cliConfig.Components.Terraform.BasePath)
	cliConfig.HelmfileDirAbsolutePath = path.Join(remoteRepoFileSystemPath, basePath, cliConfig.Components.Helmfile.BasePath)

	cliConfig.StackConfigFilesAbsolutePaths, err = u.JoinAbsolutePathWithPaths(
		path.Join(remoteRepoFileSystemPath, basePath, cliConfig.Stacks.BasePath),
		cliConfig.StackConfigFilesRelativePaths,
	)
	if err != nil {
		return nil, err
	}

	remoteStacks, err := ExecuteDescribeStacks(cliConfig, "", nil, nil, nil, true)
	if err != nil {
		return nil, err
	}

	u.PrintInfoVerbose(verbose, fmt.Sprintf("\nGetting current working repo commit object..."))

	localCommit, err := localRepo.CommitObject(localRepoHead.Hash())
	if err != nil {
		return nil, err
	}

	u.PrintInfoVerbose(verbose, fmt.Sprintf("Got current working repo commit object"))
	u.PrintInfoVerbose(verbose, fmt.Sprintf("Getting current working repo commit tree..."))

	localTree, err := localCommit.Tree()
	if err != nil {
		return nil, err
	}

	u.PrintInfoVerbose(verbose, fmt.Sprintf("Got current working repo commit tree"))
	u.PrintInfoVerbose(verbose, fmt.Sprintf("Getting remote repo commit object..."))

	remoteCommit, err := remoteRepo.CommitObject(remoteRepoHead.Hash())
	if err != nil {
		return nil, err
	}

	u.PrintInfoVerbose(verbose, fmt.Sprintf("Got remote repo commit object"))
	u.PrintInfoVerbose(verbose, fmt.Sprintf("Getting remote repo commit tree..."))

	remoteTree, err := remoteCommit.Tree()
	if err != nil {
		return nil, err
	}

	u.PrintInfoVerbose(verbose, fmt.Sprintf("Got remote repo commit tree"))
	u.PrintInfoVerbose(verbose, fmt.Sprintf("Finding diff between the current working branch and remote target branch ..."))

	// Find a slice of Patch objects with all the changes between the current working and remote trees
	patch, err := localTree.Patch(remoteTree)
	if err != nil {
		return nil, err
	}

	u.PrintInfoVerbose(verbose, fmt.Sprintf("Found diff between the current working branch and remote target branch"))
	u.PrintInfoVerbose(verbose, "\nChanged files:")

	var changedFiles []string
	for _, fileStat := range patch.Stats() {
		u.PrintMessageVerbose(verbose && fileStat.Name != "", fileStat.Name)
		changedFiles = append(changedFiles, fileStat.Name)
	}

	u.PrintMessageVerbose(verbose, "")

	affected, err := findAffected(currentStacks, remoteStacks, cliConfig, changedFiles)
	if err != nil {
		return nil, err
	}

	return affected, nil
}

// findAffected returns a list of all affected components in all stacks
func findAffected(
	currentStacks map[string]any,
	remoteStacks map[string]any,
	cliConfig cfg.CliConfiguration,
	changedFiles []string,
) ([]cfg.Affected, error) {

	res := []cfg.Affected{}
	var err error

	for stackName, stackSection := range currentStacks {
		if stackSectionMap, ok := stackSection.(map[string]any); ok {
			if componentsSection, ok := stackSectionMap["components"].(map[string]any); ok {
				if terraformSection, ok := componentsSection["terraform"].(map[string]any); ok {
					for componentName, compSection := range terraformSection {
						if componentSection, ok := compSection.(map[string]any); ok {
							if metadataSection, ok := componentSection["metadata"].(map[any]any); ok {
								// Skip abstract components
								if metadataType, ok := metadataSection["type"].(string); ok {
									if metadataType == "abstract" {
										continue
									}
								}
								// Check `metadata` section
								if !isEqual(remoteStacks, stackName, "terraform", componentName, metadataSection, "metadata") {
									affected := cfg.Affected{
										ComponentType: "terraform",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "stack.metadata",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
							// Check if any files in the component's folder have changed
							if component, ok := componentSection["component"].(string); ok && component != "" {
								if isComponentFolderChanged(component, "terraform", cliConfig, changedFiles) {
									affected := cfg.Affected{
										ComponentType: "terraform",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "component",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
							// Check `vars` section
							if varSection, ok := componentSection["vars"].(map[any]any); ok {
								if !isEqual(remoteStacks, stackName, "terraform", componentName, varSection, "vars") {
									affected := cfg.Affected{
										ComponentType: "terraform",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "stack.vars",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
							// Check `env` section
							if envSection, ok := componentSection["env"].(map[any]any); ok {
								if !isEqual(remoteStacks, stackName, "terraform", componentName, envSection, "env") {
									affected := cfg.Affected{
										ComponentType: "terraform",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "stack.env",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
							// Check `settings` section
							if settingsSection, ok := componentSection["settings"].(map[any]any); ok {
								if !isEqual(remoteStacks, stackName, "terraform", componentName, settingsSection, "settings") {
									affected := cfg.Affected{
										ComponentType: "terraform",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "stack.settings",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
						}
					}
				}

				if helmfileSection, ok := componentsSection["helmfile"].(map[string]any); ok {
					for componentName, compSection := range helmfileSection {
						if componentSection, ok := compSection.(map[string]any); ok {
							if metadataSection, ok := componentSection["metadata"].(map[any]any); ok {
								// Skip abstract components
								if metadataType, ok := metadataSection["type"].(string); ok {
									if metadataType == "abstract" {
										continue
									}
								}
								// Check `metadata` section
								if !isEqual(remoteStacks, stackName, "helmfile", componentName, metadataSection, "metadata") {
									affected := cfg.Affected{
										ComponentType: "helmfile",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "stack.metadata",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
							// Check if any files in the component's folder have changed
							if component, ok := componentSection["component"].(string); ok && component != "" {
								if isComponentFolderChanged(component, "helmfile", cliConfig, changedFiles) {
									affected := cfg.Affected{
										ComponentType: "helmfile",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "component",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
							// Check `vars` section
							if varSection, ok := componentSection["vars"].(map[any]any); ok {
								if !isEqual(remoteStacks, stackName, "helmfile", componentName, varSection, "vars") {
									affected := cfg.Affected{
										ComponentType: "helmfile",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "stack.vars",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
							// Check `env` section
							if envSection, ok := componentSection["env"].(map[any]any); ok {
								if !isEqual(remoteStacks, stackName, "helmfile", componentName, envSection, "env") {
									affected := cfg.Affected{
										ComponentType: "helmfile",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "stack.env",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
							// Check `settings` section
							if settingsSection, ok := componentSection["settings"].(map[any]any); ok {
								if !isEqual(remoteStacks, stackName, "helmfile", componentName, settingsSection, "settings") {
									affected := cfg.Affected{
										ComponentType: "helmfile",
										Component:     componentName,
										Stack:         stackName,
										Affected:      "stack.settings",
									}
									res, err = appendToAffected(cliConfig, componentName, stackName, componentSection, res, affected)
									if err != nil {
										return nil, err
									}
									continue
								}
							}
						}
					}
				}
			}
		}
	}

	return res, nil
}

// appendToAffected adds an item to the affected list, and adds the Spacelift stack and Atlantis project (if configured)
func appendToAffected(
	cliConfig cfg.CliConfiguration,
	componentName string,
	stackName string,
	componentSection map[string]any,
	affectedList []cfg.Affected,
	affected cfg.Affected,
) ([]cfg.Affected, error) {

	if affected.ComponentType == "terraform" {
		varSection := map[any]any{}
		settingsSection := map[any]any{}

		if i, ok2 := componentSection["vars"]; ok2 {
			varSection = i.(map[any]any)
		}

		if i, ok2 := componentSection["settings"]; ok2 {
			settingsSection = i.(map[any]any)
		}

		// Affected Spacelift stack
		spaceliftStackName, err := BuildSpaceliftStackNameFromComponentConfig(
			cliConfig,
			componentName,
			stackName,
			settingsSection,
			varSection,
		)

		if err != nil {
			return nil, err
		}

		affected.SpaceliftStack = spaceliftStackName

		// Affected Atlantis project
		atlantisProjectName, err := BuildAtlantisProjectNameFromComponentConfig(
			cliConfig,
			componentName,
			settingsSection,
			varSection,
		)

		if err != nil {
			return nil, err
		}

		affected.AtlantisProject = atlantisProjectName
	}

	// Check `component` section and add `ComponentPath` to the output
	affected.ComponentPath = BuildComponentPath(cliConfig, componentSection, affected.ComponentType)
	affected.StackSlug = fmt.Sprintf("%s-%s", stackName, strings.Replace(componentName, "/", "-", -1))

	return append(affectedList, affected), nil
}

// isEqual compares a section of a component from the remote stacks with a section of a local component
func isEqual(
	remoteStacks map[string]any,
	localStackName string,
	componentType string,
	localComponentName string,
	localSection map[any]any,
	sectionName string,
) bool {

	if remoteStackSection, ok := remoteStacks[localStackName].(map[string]any); ok {
		if remoteComponentsSection, ok := remoteStackSection["components"].(map[string]any); ok {
			if remoteComponentTypeSection, ok := remoteComponentsSection[componentType].(map[string]any); ok {
				if remoteComponentSection, ok := remoteComponentTypeSection[localComponentName].(map[string]any); ok {
					if remoteSection, ok := remoteComponentSection[sectionName].(map[any]any); ok {
						if reflect.DeepEqual(localSection, remoteSection) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// isComponentFolderChanged checks if a component folder changed (has changed files in it)
func isComponentFolderChanged(
	component string,
	componentType string,
	cliConfig cfg.CliConfiguration,
	changedFiles []string,
) bool {

	var pathPrefix string

	switch componentType {
	case "terraform":
		pathPrefix = path.Join(cliConfig.BasePath, cliConfig.Components.Terraform.BasePath, component)
	case "helmfile":
		pathPrefix = path.Join(cliConfig.BasePath, cliConfig.Components.Helmfile.BasePath, component)
	}

	if u.SliceOfPathsContainsPath(changedFiles, pathPrefix) {
		return true
	}
	return false
}
