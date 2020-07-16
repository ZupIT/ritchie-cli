package credsingle

// Todo fix setter test
//
// func TestSet(t *testing.T) {
// 	tmp := os.TempDir()
// 	setter := NewSetter(tmp, ctxFinder, sessManager)
//
// 	tests := []struct {
// 		name string
// 		in   credential.Detail
// 		out  error
// 	}{
// 		{
// 			name: "github credential",
// 			in:   githubCred,
// 			out:  nil,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := setter.Set(tt.in)
// 			if got != tt.out {
// 				t.Errorf("Set(%s) got %v, want %v", tt.name, got, tt.out)
// 			}
// 		})
// 	}
// }
