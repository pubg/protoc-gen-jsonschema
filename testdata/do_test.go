package testdata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/lyft/protoc-gen-star/v2/testutils"
	"github.com/pubg/protoc-gen-jsonschema/pkg/modules"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type Testcase struct {
	// Testcase Identifier
	Name string `json:"name"`

	// Testcase Description
	Description                 string `json:"description"`
	ExpectedBehaviorDescription string `json:"expectedBehaviorDescription"`

	// Testcase Input
	InputFiles      []string          `json:"inputFiles"`
	InputParameters map[string]string `json:"inputParameters"`

	// Expected result
	ExpectResultFiles  []string `json:"expectResultFiles"`
	ExpectResultIsNull bool     `json:"expectResultIsNull"`
	ExpectFail         bool     `json:"expectFail"`
}

const testdataDir = "./cases"
const optionsDir = "../"

func TestPlugin(t *testing.T) {
	dirs, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatal(err)
	}
	optionsImportPath, _ := filepath.Abs(optionsDir)

	for _, dir := range dirs {
		tc, err := loadTestcase(testdataDir, dir)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(dir.Name(), func(t *testing.T) {
			// Load testcase
			dirPath := filepath.Join(testdataDir, dir.Name())
			module, debugger := newModule(tc, dirPath)
			ast := loadInput(t, tc, dirPath, optionsImportPath)

			// Execute testcase
			artifacts := module.Execute(ast.Targets(), ast.Packages())

			// Check testcase result
			if tc.ExpectFail {
				if !debugger.Failed() {
					t.Errorf("Expect failed, but got success")
				}
				return
			}

			if tc.ExpectResultIsNull {
				if len(artifacts) != 0 {
					t.Errorf("Expect result is null, but got %d of results", len(artifacts))
				}
				return
			}

			expectedResult, err := loadExpectedResult(dirPath, tc.ExpectResultFiles)
			if err != nil {
				t.Errorf("failed to load expected result: %s", err.Error())
				return
			}
			if len(artifacts) != len(tc.ExpectResultFiles) {
				t.Errorf("Expect %d of results, but got %d of results", len(tc.ExpectResultFiles), len(artifacts))
				return
			}
			for _, artifact := range artifacts {
				fArtifact, ok := artifact.(pgs.GeneratorFile)
				if !ok {
					t.Errorf("artifact is not GeneratorFile, but %T", artifact)
					return
				}

				var actual any
				err := yaml.Unmarshal([]byte(fArtifact.Contents), &actual)
				if err != nil {
					t.Errorf("failed to convert to comparable component: %s", err.Error())
					return
				}
				expected := expectedResult.results[fArtifact.Name]

				require.Equal(t, expected, actual, "not equals at index %s: %s", tc.Name, tc.Description)
			}
		})
	}
}

func loadTestcase(parentDir string, dir os.DirEntry) (*Testcase, error) {
	if !dir.IsDir() {
		return nil, fmt.Errorf("dir %s is not directory", dir.Name())
	}

	dirPath := filepath.Join(parentDir, dir.Name())

	buf, err := os.ReadFile(filepath.Join(dirPath, "test.yaml"))
	if err != nil {
		return nil, err
	}
	testcase := &Testcase{}
	if err := yaml.Unmarshal(buf, testcase); err != nil {
		return nil, err
	}
	return testcase, nil
}

func loadInput(t *testing.T, tc *Testcase, dirPath string, optionsImportPath string) pgs.AST {
	loader := &testutils.Loader{ImportPaths: []string{dirPath, optionsImportPath}}
	inputFiles := lo.Map[string, string](tc.InputFiles, func(item string, _ int) string { return filepath.Join(dirPath, item) })
	ast := loader.LoadProtos(t, inputFiles...)

	protoFiles := map[string]pgs.File{}
	for _, pkg := range ast.Packages() {
		for _, file := range pkg.Files() {
			protoFiles[file.InputPath().String()] = file
		}
	}
	for _, inputFile := range tc.InputFiles {
		ast.Targets()[inputFile] = protoFiles[inputFile]
	}
	return ast
}

func newModule(tc *Testcase, dirPath string) (*modules.Module, pgs.MockDebugger) {
	debugger := pgs.InitMockDebugger()
	module := modules.NewModule()
	module.InitContext(pgs.Context(debugger, tc.InputParameters, dirPath))
	return module, debugger
}

type ExpectedResult struct {
	results map[string]any
}

func loadExpectedResult(dirPath string, files []string) (*ExpectedResult, error) {
	result := &ExpectedResult{results: map[string]any{}}
	for _, file := range files {
		buf, err := os.ReadFile(filepath.Join(dirPath, file))
		if err != nil {
			return nil, err
		}

		var a any
		err = yaml.Unmarshal(buf, &a)
		if err != nil {
			return nil, err
		}

		result.results[file] = a
	}
	return result, nil
}

func convertToComparableComponent(buf []byte) ([]any, error) {
	containers := [][]any{}
	err := json.Unmarshal(buf, &containers)
	if err != nil {
		return nil, err
	}

	return lo.FlatMap[[]any, any](containers, func(container []any, _ int) []any { return container }), nil
}
