/*
Copyright 2018 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package latest

import (
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/util"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const Version string = "skaffold/v1alpha5"

// NewSkaffoldPipeline creates a SkaffoldPipeline
func NewSkaffoldPipeline() util.VersionedConfig {
	return new(SkaffoldPipeline)
}

type SkaffoldPipeline struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`

	Build    BuildConfig  `yaml:"build,omitempty"`
	Test     []TestCase   `yaml:"test,omitempty"`
	Deploy   DeployConfig `yaml:"deploy,omitempty"`
	Profiles []Profile    `yaml:"profiles,omitempty"`
}

func (c *SkaffoldPipeline) GetVersion() string {
	return c.APIVersion
}

// BuildConfig contains all the configuration for the build steps
type BuildConfig struct {
	Artifacts []*Artifact `yaml:"artifacts,omitempty"`
	TagPolicy TagPolicy   `yaml:"tagPolicy,omitempty"`
	BuildType `yaml:",inline"`
}

// TagPolicy contains all the configuration for the tagging step
type TagPolicy struct {
	GitTagger         *GitTagger         `yaml:"gitCommit,omitempty" yamltags:"oneOf=tag"`
	ShaTagger         *ShaTagger         `yaml:"sha256,omitempty" yamltags:"oneOf=tag"`
	EnvTemplateTagger *EnvTemplateTagger `yaml:"envTemplate,omitempty" yamltags:"oneOf=tag"`
	DateTimeTagger    *DateTimeTagger    `yaml:"dateTime,omitempty" yamltags:"oneOf=tag"`
}

// ShaTagger contains the configuration for the SHA tagger.
type ShaTagger struct{}

// GitTagger contains the configuration for the git tagger.
type GitTagger struct{}

// EnvTemplateTagger contains the configuration for the envTemplate tagger.
type EnvTemplateTagger struct {
	Template string `yaml:"template"`
}

// DateTimeTagger contains the configuration for the DateTime tagger.
type DateTimeTagger struct {
	Format   string `yaml:"format,omitempty"`
	TimeZone string `yaml:"timezone,omitempty"`
}

// BuildType contains the specific implementation and parameters needed
// for the build step. Only one field should be populated.
type BuildType struct {
	LocalBuild          *LocalBuild          `yaml:"local,omitempty" yamltags:"oneOf=build"`
	GoogleCloudBuild    *GoogleCloudBuild    `yaml:"googleCloudBuild,omitempty" yamltags:"oneOf=build"`
	KanikoBuild         *KanikoBuild         `yaml:"kaniko,omitempty" yamltags:"oneOf=build"`
	AzureContainerBuild *AzureContainerBuild `yaml:"acr,omitempty" yamltags:"oneOf=build"`
}

// LocalBuild contains the fields needed to do a build on the local docker daemon
// and optionally push to a repository.
type LocalBuild struct {
	SkipPush     *bool `yaml:"skipPush,omitempty"`
	UseDockerCLI bool  `yaml:"useDockerCLI,omitempty"`
	UseBuildkit  bool  `yaml:"useBuildkit,omitempty"`
}

// GoogleCloudBuild contains the fields needed to do a remote build on
// Google Cloud Build.
type GoogleCloudBuild struct {
	ProjectID   string `yaml:"projectId"`
	DiskSizeGb  int64  `yaml:"diskSizeGb,omitempty"`
	MachineType string `yaml:"machineType,omitempty"`
	Timeout     string `yaml:"timeout,omitempty"`
	DockerImage string `yaml:"dockerImage,omitempty"`
}

// KanikoBuildContext contains the different fields available to specify
// a kaniko build context
type KanikoBuildContext struct {
	GCSBucket string `yaml:"gcsBucket,omitempty" yamltags:"oneOf=buildContext"`
}

// KanikoBuild contains the fields needed to do a on-cluster build using
// the kaniko image
type KanikoBuild struct {
	BuildContext   KanikoBuildContext `yaml:"buildContext,omitempty"`
	PullSecret     string             `yaml:"pullSecret,omitempty"`
	PullSecretName string             `yaml:"pullSecretName,omitempty"`
	Namespace      string             `yaml:"namespace,omitempty"`
	Timeout        string             `yaml:"timeout,omitempty"`
	Image          string             `yaml:"image,omitempty"`
}

// AzureContainerBuild contains the fields needed to do a build
// on Azure Container Registry
type AzureContainerBuild struct {
	Credentials       AzureContainerBuildCredentials `yaml:"credentials,omitempty"`
	ContainerRegistry string                         `yaml:"containerRegistry,omitempty"`
	ResourceGroup     string                         `yaml:"resourceGroup,omitempty"`
}

// AzureContainerBuildCredentials contains the fields needed to
// to authorize the azure sdk
type AzureContainerBuildCredentials struct {
	SubscriptionId string `yaml:"subscriptionId,omitempty"`
	ClientId       string `yaml:"clientId,omitempty"`
	ClientSecret   string `yaml:"clientSecret,omitempty"`
	TenantId       string `yaml:"tenantId,omitempty"`
}

// TestCase is a struct containing all the specified test
// configuration for an image.
type TestCase struct {
	ImageName      string   `yaml:"image"`
	StructureTests []string `yaml:"structureTests,omitempty"`
}

// DeployConfig contains all the configuration needed by the deploy steps
type DeployConfig struct {
	DeployType `yaml:",inline"`
}

// DeployType contains the specific implementation and parameters needed
// for the deploy step. Only one field should be populated.
type DeployType struct {
	HelmDeploy      *HelmDeploy      `yaml:"helm,omitempty" yamltags:"oneOf=deploy"`
	KubectlDeploy   *KubectlDeploy   `yaml:"kubectl,omitempty" yamltags:"oneOf=deploy"`
	KustomizeDeploy *KustomizeDeploy `yaml:"kustomize,omitempty" yamltags:"oneOf=deploy"`
}

// KubectlDeploy contains the configuration needed for deploying with `kubectl apply`
type KubectlDeploy struct {
	Manifests       []string     `yaml:"manifests,omitempty"`
	RemoteManifests []string     `yaml:"remoteManifests,omitempty"`
	Flags           KubectlFlags `yaml:"flags,omitempty"`
}

// KubectlFlags describes additional options flags that are passed on the command
// line to kubectl either on every command (Global), on creations (Apply)
// or deletions (Delete).
type KubectlFlags struct {
	Global []string `yaml:"global,omitempty"`
	Apply  []string `yaml:"apply,omitempty"`
	Delete []string `yaml:"delete,omitempty"`
}

// HelmDeploy contains the configuration needed for deploying with helm
type HelmDeploy struct {
	Releases []HelmRelease `yaml:"releases,omitempty"`
}

// KustomizeDeploy contains the configuration needed for deploying with kustomize.
type KustomizeDeploy struct {
	KustomizePath string       `yaml:"path,omitempty"`
	Flags         KubectlFlags `yaml:"flags,omitempty"`
}

type HelmRelease struct {
	Name              string                 `yaml:"name"`
	ChartPath         string                 `yaml:"chartPath"`
	ValuesFiles       []string               `yaml:"valuesFiles"`
	Values            map[string]string      `yaml:"values,omitempty"`
	Namespace         string                 `yaml:"namespace"`
	Version           string                 `yaml:"version"`
	SetValues         map[string]string      `yaml:"setValues"`
	SetValueTemplates map[string]string      `yaml:"setValueTemplates"`
	Wait              bool                   `yaml:"wait"`
	RecreatePods      bool                   `yaml:"recreatePods"`
	Overrides         map[string]interface{} `yaml:"overrides"`
	Packaged          *HelmPackaged          `yaml:"packaged"`
	ImageStrategy     HelmImageStrategy      `yaml:"imageStrategy"`
}

// HelmPackaged represents parameters for packaging helm chart.
type HelmPackaged struct {
	// Version sets the version on the chart to this semver version.
	Version string `yaml:"version"`

	// AppVersion set the appVersion on the chart to this version
	AppVersion string `yaml:"appVersion"`
}

type HelmImageStrategy struct {
	HelmImageConfig `yaml:",inline"`
}

type HelmImageConfig struct {
	HelmFQNConfig        *HelmFQNConfig        `yaml:"fqn"`
	HelmConventionConfig *HelmConventionConfig `yaml:"helm"`
}

// HelmFQNConfig represents image config to use the FullyQualifiedImageName as param to set
type HelmFQNConfig struct {
	Property string `yaml:"property"`
}

// HelmConventionConfig represents image config in the syntax of image.repository and image.tag
type HelmConventionConfig struct {
}

// Artifact represents items that need to be built, along with the context in which
// they should be built.
type Artifact struct {
	ImageName    string            `yaml:"image"`
	Workspace    string            `yaml:"context,omitempty"`
	Sync         map[string]string `yaml:"sync,omitempty"`
	ArtifactType `yaml:",inline"`
}

// Profile is additional configuration that overrides default
// configuration when it is activated.
type Profile struct {
	Name   string       `yaml:"name"`
	Build  BuildConfig  `yaml:"build,omitempty"`
	Test   []TestCase   `yaml:"test,omitempty"`
	Deploy DeployConfig `yaml:"deploy,omitempty"`
}

type ArtifactType struct {
	DockerArtifact *DockerArtifact `yaml:"docker,omitempty" yamltags:"oneOf=artifact"`
	BazelArtifact  *BazelArtifact  `yaml:"bazel,omitempty" yamltags:"oneOf=artifact"`
}

// DockerArtifact describes an artifact built from a Dockerfile,
// usually using `docker build`.
type DockerArtifact struct {
	DockerfilePath string             `yaml:"dockerfile,omitempty"`
	BuildArgs      map[string]*string `yaml:"buildArgs,omitempty"`
	CacheFrom      []string           `yaml:"cacheFrom,omitempty"`
	Target         string             `yaml:"target,omitempty"`
}

// BazelArtifact describes an artifact built with Bazel.
type BazelArtifact struct {
	BuildTarget string `yaml:"target"`
}

// Parse reads a SkaffoldPipeline from yaml.
func (c *SkaffoldPipeline) Parse(contents []byte, useDefaults bool) error {
	if err := yaml.UnmarshalStrict(contents, c); err != nil {
		return err
	}

	if useDefaults {
		if err := c.SetDefaultValues(); err != nil {
			return errors.Wrap(err, "applying default values")
		}
	}

	return nil
}
