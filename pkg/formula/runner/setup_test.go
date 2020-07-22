package runner

import (
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestSetup(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	ritHome := fmt.Sprintf("%s/.rit-setup", tmpDir)
	repoPath := fmt.Sprintf("%s/repos/commons", ritHome)

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	_ = streams.Unzip("../../../testdata/ritchie-formulas-test.zip", repoPath)

	type in struct {
		def         formula.Definition
		makeBuild   formula.MakeBuilder
		batBuild    formula.BatBuilder
		dockerBuild formula.DockerBuilder
		localFlag   bool
	}

	type out struct {
		want formula.Setup
		err  error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "docker success",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				dockerBuild: dockerBuildMock{},
				localFlag:   false,
			},
			out: out{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			formulaSetup := NewSetup(ritHome, in.makeBuild, in.dockerBuild, in.batBuild, dirManager, fileManager)

			got, err := formulaSetup.Setup(in.def, in.localFlag)

			fmt.Println(got)
			fmt.Println(err)
		})
	}
}

type makeBuildMock struct {
	err error
}

func (ma makeBuildMock) Build(string) error {
	return ma.err
}

type batBuildMock struct {
	err error
}

func (ba batBuildMock) Build(string) error {
	return ba.err
}

type dockerBuildMock struct {
	err error
}

func (do dockerBuildMock) Build(formulaPath, dockerImg string) error {
	return do.err
}
