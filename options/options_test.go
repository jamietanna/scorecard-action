// Copyright OpenSSF Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//nolint
package options

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ossf/scorecard/v4/options"
)

const (
	testRepo        = "good/repo"
	testResultsFile = "results.sarif"
	testToken       = "test-token"

	githubEventPathNonFork   = "testdata/non-fork.json"
	githubEventPathFork      = "testdata/fork.json"
	githubEventPathIncorrect = "testdata/incorrect.json"
	githubEventPathBadPath   = "testdata/bad-path.json"
	githubEventPathBadData   = "testdata/bad-data.json"
)

func TestNew(t *testing.T) {
	type fields struct {
		EnableSarif bool
		Format      string
		PolicyFile  string
		ResultsFile string
		Commit      string
		LogLevel    string
	}
	tests := []struct {
		name            string
		githubEventPath string
		repo            string
		resultsFile     string
		resultsFormat   string
		publishResults  string
		want            fields
		unsetToken      bool
		wantErr         bool
	}{
		{
			name:            "SuccessFormatSARIF",
			githubEventPath: githubEventPathNonFork,
			repo:            testRepo,
			resultsFormat:   "sarif",
			resultsFile:     testResultsFile,
			want: fields{
				EnableSarif: true,
				Format:      formatSarif,
				PolicyFile:  defaultScorecardPolicyFile,
				ResultsFile: testResultsFile,
				Commit:      options.DefaultCommit,
				LogLevel:    options.DefaultLogLevel,
			},
			wantErr: false,
		},
		{
			name:            "SuccessFormatJSON",
			githubEventPath: githubEventPathNonFork,
			repo:            testRepo,
			resultsFormat:   "json",
			resultsFile:     testResultsFile,
			want: fields{
				EnableSarif: true,
				Format:      options.FormatJSON,
				ResultsFile: testResultsFile,
				Commit:      options.DefaultCommit,
				LogLevel:    options.DefaultLogLevel,
			},
			wantErr: false,
		},
		{
			name:            "FailureTokenIsNotSet",
			githubEventPath: githubEventPathNonFork,
			repo:            testRepo,
			resultsFormat:   "sarif",
			resultsFile:     testResultsFile,
			want: fields{
				EnableSarif: true,
				Format:      formatSarif,
				PolicyFile:  defaultScorecardPolicyFile,
				ResultsFile: testResultsFile,
				Commit:      options.DefaultCommit,
				LogLevel:    options.DefaultLogLevel,
			},
			unsetToken: true,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, tokenEnvExists := os.LookupEnv(EnvGithubAuthToken)
			if !tokenEnvExists {
				os.Setenv(EnvGithubAuthToken, testToken)
			}
			if tt.unsetToken {
				os.Unsetenv(EnvGithubAuthToken)
			}

			_, pathEnvExists := os.LookupEnv(EnvGithubEventPath)
			if !pathEnvExists {
				if tt.githubEventPath != "" {
					os.Setenv(EnvGithubEventPath, tt.githubEventPath)
				}
			}

			_, repoEnvExists := os.LookupEnv(EnvGithubRepository)
			if !repoEnvExists {
				if tt.repo != "" {
					os.Setenv(EnvGithubRepository, tt.repo)
				}
			}

			if tt.resultsFile != "" {
				os.Setenv("SCORECARD_RESULTS_FILE", tt.resultsFile)
			}
			if tt.resultsFormat != "" {
				os.Setenv("SCORECARD_RESULTS_FORMAT", tt.resultsFormat)
			}

			opts, err := New()
			scOpts := *opts.ScorecardOpts
			got := fields{
				EnableSarif: scOpts.EnableSarif,
				Format:      scOpts.Format,
				PolicyFile:  scOpts.PolicyFile,
				ResultsFile: scOpts.ResultsFile,
				Commit:      scOpts.Commit,
				LogLevel:    scOpts.LogLevel,
			}

			if (err != nil) != tt.wantErr {
				for _, e := range os.Environ() {
					t.Logf(e)
				}
				t.Errorf("New() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(tt.want, got) {
				t.Errorf("New(): -want, +got:\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestInitialize(t *testing.T) {
	type fields struct {
		ScorecardOpts           *options.Options
		EnabledChecks           string
		EnableLicense           string
		EnableDangerousWorkflow string
		GithubEventName         string
		GithubEventPath         string
		GithubRef               string
		GithubRepository        string
		GithubWorkspace         string
		DefaultBranch           string
		IsForkStr               string
		PrivateRepoStr          string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				GithubEventPath: githubEventPathNonFork,
			},
			wantErr: false,
		},
		{
			name:    "FailureNoFieldsSet",
			wantErr: true,
		},
		{
			name: "FailureBadEventPath",
			fields: fields{
				GithubEventPath: githubEventPathBadPath,
			},
			wantErr: true,
		},
		{
			name: "FailureBadEventData",
			fields: fields{
				GithubEventPath: githubEventPathBadData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				ScorecardOpts:           tt.fields.ScorecardOpts,
				EnabledChecks:           tt.fields.EnabledChecks,
				EnableLicense:           tt.fields.EnableLicense,
				EnableDangerousWorkflow: tt.fields.EnableDangerousWorkflow,
				GithubEventName:         tt.fields.GithubEventName,
				GithubEventPath:         tt.fields.GithubEventPath,
				GithubRef:               tt.fields.GithubRef,
				GithubRepository:        tt.fields.GithubRepository,
				GithubWorkspace:         tt.fields.GithubWorkspace,
				DefaultBranch:           tt.fields.DefaultBranch,
				IsForkStr:               tt.fields.IsForkStr,
				PrivateRepoStr:          tt.fields.PrivateRepoStr,
			}
			if err := o.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("Options.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	type fields struct {
		ScorecardOpts           *options.Options
		EnabledChecks           string
		EnableLicense           string
		EnableDangerousWorkflow string
		GithubEventName         string
		GithubEventPath         string
		GithubRef               string
		GithubRepository        string
		GithubWorkspace         string
		DefaultBranch           string
		IsForkStr               string
		PrivateRepoStr          string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Success",
			fields: fields{
				ScorecardOpts: options.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				ScorecardOpts:           tt.fields.ScorecardOpts,
				EnabledChecks:           tt.fields.EnabledChecks,
				EnableLicense:           tt.fields.EnableLicense,
				EnableDangerousWorkflow: tt.fields.EnableDangerousWorkflow,
				GithubEventName:         tt.fields.GithubEventName,
				GithubEventPath:         tt.fields.GithubEventPath,
				GithubRef:               tt.fields.GithubRef,
				GithubRepository:        tt.fields.GithubRepository,
				GithubWorkspace:         tt.fields.GithubWorkspace,
				DefaultBranch:           tt.fields.DefaultBranch,
				IsForkStr:               tt.fields.IsForkStr,
				PrivateRepoStr:          tt.fields.PrivateRepoStr,
			}
			o.Print()
		})
	}
}
