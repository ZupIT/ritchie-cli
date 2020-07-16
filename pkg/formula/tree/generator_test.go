package tree

// Todo fix
// func TestGenerate(t *testing.T) {
// 	fileManager := stream.NewFileManager()
// 	dirManager := stream.NewDirManager(fileManager)
// 	generator := NewGenerator(dirManager, fileManager)
//
// 	adder := repo.NewAdder(os.TempDir(), http.DefaultClient, generator, dirManager, fileManager)
// 	_ = adder.Add(formula.Repo{
// 		Name:     "commons",
// 		ZipUrl:   "http://localhost:8882/repos/ZupIT/ritchie-formulas/zipball/v2.0.0",
// 		Version:  "v2.0.0",
// 		Priority: 0,
// 	})
//
// 	resultDir := os.TempDir() + "/commons"
// 	_ = dirManager.Create(resultDir)
// 	defer func() {
// 		_ = dirManager.Remove(resultDir)
// 	}()
// 	tree, err := generator.Generate(resultDir)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	bytes, _ := json.MarshalIndent(tree, "", "\t")
//
// 	fmt.Println(string(bytes))
// }
