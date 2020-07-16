package runner

// Todo fix setup test
// import (
// 	"net/http"
// 	"os"
// 	"testing"
//
// 	"github.com/ZupIT/ritchie-cli/pkg/formula"
//
// 	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
// )
//
// func TestDefaultSetup_Setup(t *testing.T) {
// 	def := formula.Definition{
// 		Path:    "mock/test",
// 	}
//
// 	home := os.TempDir()
//
// 	type in struct {
// 		config string
// 		bundle string
// 	}
//
// 	tests := []struct {
// 		name string
// 		in   in
// 		want error
// 	}{
// 		{
// 			name: "success",
// 			in: in{
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "config not found",
// 			in: in{
// 				config: "config-not-found",
// 			},
// 			want: ErrConfigFileNotFound,
// 		},
// 		{
// 			name: "bundle not found",
// 			in: in{
// 				bundle: "bundle-not-found",
// 			},
// 			want: ErrFormulaBinNotFound,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			in := tt.in
// 			defTest := def
// 			if in.config != "" {
// 				defTest.Config = in.config
// 			}
// 			if in.bundle != "" {
// 				defTest.Bundle = in.bundle
// 			}
//
// 			_ = fileutil.RemoveDir(home + "/formulas")
// 			setup := NewDefaultTeamSetup(home, http.DefaultClient, in.sess)
// 			_, got := setup.Setup(defTest)
//
// 			if got != nil && got.Error() != tt.want.Error() {
// 				t.Errorf("Setup(%s) got %v, want %v", tt.name, got, tt.want)
// 			}
//
// 		})
// 	}
// }
