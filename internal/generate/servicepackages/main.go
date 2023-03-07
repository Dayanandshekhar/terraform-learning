//go:build generate
// +build generate

package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-provider-aws/internal/generate/common"
	"github.com/hashicorp/terraform-provider-aws/names"
	"golang.org/x/exp/slices"
)

func main() {
	const (
		spdFile       = `service_package_gen.go`
		spsFile       = `../../provider/service_packages_gen.go`
		namesDataFile = `../../../names/names_data.csv`
	)
	g := common.NewGenerator()

	data, err := common.ReadAllCSVData(namesDataFile)

	if err != nil {
		g.Fatalf("error reading %s: %s", namesDataFile, err)
	}

	g.Infof("Generating per-service %s", filepath.Base(spdFile))

	td := TemplateData{}

	for i, l := range data {
		if i < 1 { // no header
			continue
		}

		// Don't skip excluded packages, instead handle missing values in the template.
		// if l[names.ColExclude] != "" {
		// 	continue
		// }

		if l[names.ColProviderPackageActual] == "" && l[names.ColProviderPackageCorrect] == "" {
			continue
		}

		// See internal/generate/namesconsts/main.go.
		p := l[names.ColProviderPackageCorrect]

		if l[names.ColProviderPackageActual] != "" {
			p = l[names.ColProviderPackageActual]
		}

		dir := fmt.Sprintf("../../service/%s", p)

		if _, err := os.Stat(dir); err != nil {
			continue
		}

		// Look for Terraform Plugin Framework and SDK resource and data source annotations.
		// These annotations are implemented as comments on factory functions.
		v := &visitor{
			g: g,

			frameworkDataSources: make([]string, 0),
			frameworkResources:   make([]string, 0),
			sdkDataSources:       make(map[string]string),
			sdkResources:         make(map[string]string),
		}

		v.processDir(dir)

		if err := v.err.ErrorOrNil(); err != nil {
			g.Fatalf("%s", err.Error())
		}

		s := ServiceDatum{
			ProviderPackage:      p,
			ProviderNameUpper:    l[names.ColProviderNameUpper],
			FrameworkDataSources: v.frameworkDataSources,
			FrameworkResources:   v.frameworkResources,
			SDKDataSources:       v.sdkDataSources,
			SDKResources:         v.sdkResources,
		}

		sort.SliceStable(s.FrameworkDataSources, func(i, j int) bool {
			return s.FrameworkDataSources[i] < s.FrameworkDataSources[j]
		})
		sort.SliceStable(s.FrameworkResources, func(i, j int) bool {
			return s.FrameworkResources[i] < s.FrameworkResources[j]
		})

		filename := fmt.Sprintf("../../service/%s/%s", p, spdFile)
		d := g.NewGoFileDestination(filename)

		if err := d.WriteTemplate("servicepackagedata", spdTmpl, s); err != nil {
			g.Fatalf("error generating %s service package data: %s", p, err)
		}

		if err := d.Write(); err != nil {
			g.Fatalf("generating file (%s): %s", filename, err)
		}

		td.Services = append(td.Services, s)
	}

	sort.SliceStable(td.Services, func(i, j int) bool {
		return td.Services[i].ProviderPackage < td.Services[j].ProviderPackage
	})

	g.Infof("Generating %s", filepath.Base(spsFile))

	d := g.NewGoFileDestination(spsFile)

	if err := d.WriteTemplate("servicepackages", spsTmpl, td); err != nil {
		g.Fatalf("error generating service packages list: %s", err)
	}

	if err := d.Write(); err != nil {
		g.Fatalf("generating file (%s): %s", spsFile, err)
	}
}

type ServiceDatum struct {
	ProviderPackage      string
	ProviderNameUpper    string
	FrameworkDataSources []string
	FrameworkResources   []string
	SDKDataSources       map[string]string
	SDKResources         map[string]string
}

type TemplateData struct {
	Services []ServiceDatum
}

//go:embed spd.tmpl
var spdTmpl string

//go:embed sps.tmpl
var spsTmpl string

// Annotation processing.
var (
	frameworkDataSourceAnnotation = regexp.MustCompile(`^//\s*@FrameworkDataSource\s*$`)
	frameworkResourceAnnotation   = regexp.MustCompile(`^//\s*@FrameworkResource\s*$`)
	sdkDataSourceAnnotation       = regexp.MustCompile(`^//\s*@SDKDataSource\(\s*"([a-z0-9_]+)"\s*\)\s*$`)
	sdkResourceAnnotation         = regexp.MustCompile(`^//\s*@SDKResource\(\s*"([a-z0-9_]+)"\s*\)\s*$`)
)

type visitor struct {
	err *multierror.Error
	g   *common.Generator

	fileName     string
	functionName string
	packageName  string

	frameworkDataSources []string
	frameworkResources   []string
	sdkDataSources       map[string]string
	sdkResources         map[string]string
}

// processDir scans a single service package directory and processes contained Go sources files.
func (v *visitor) processDir(path string) {
	fileSet := token.NewFileSet()
	packageMap, err := parser.ParseDir(fileSet, path, func(fi os.FileInfo) bool {
		// Skip tests.
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)

	if err != nil {
		v.err = multierror.Append(v.err, fmt.Errorf("parsing (%s): %w", path, err))

		return
	}

	for name, pkg := range packageMap {
		v.packageName = name

		for name, file := range pkg.Files {
			v.fileName = name

			v.processFile(file)

			v.fileName = ""
		}

		v.packageName = ""
	}
}

// processFile processes a single Go source file.
func (v *visitor) processFile(file *ast.File) {
	ast.Walk(v, file)
}

// processFuncDecl processes a single Go function.
// The function's comments are scanned for annotations indicating a Plugin Framework or SDK resource or data source.
func (v *visitor) processFuncDecl(funcDecl *ast.FuncDecl) {
	v.functionName = funcDecl.Name.Name

	for _, line := range funcDecl.Doc.List {
		line := line.Text

		if m := frameworkDataSourceAnnotation.FindStringSubmatch(line); len(m) > 0 {
			if slices.Contains(v.frameworkDataSources, v.functionName) {
				v.err = multierror.Append(v.err, fmt.Errorf("duplicate Framework Data Source: %s", fmt.Sprintf("%s.%s", v.packageName, v.functionName)))
			} else {
				v.frameworkDataSources = append(v.frameworkDataSources, v.functionName)
			}
		} else if m := frameworkResourceAnnotation.FindStringSubmatch(line); len(m) > 0 {
			if slices.Contains(v.frameworkResources, v.functionName) {
				v.err = multierror.Append(v.err, fmt.Errorf("duplicate Framework Resource: %s", fmt.Sprintf("%s.%s", v.packageName, v.functionName)))
			} else {
				v.frameworkResources = append(v.frameworkResources, v.functionName)
			}
		} else if m := sdkDataSourceAnnotation.FindStringSubmatch(line); len(m) > 0 {
			name := m[1]

			if _, ok := v.sdkDataSources[name]; ok {
				v.err = multierror.Append(v.err, fmt.Errorf("duplicate SDK Data Source (%s): %s", name, fmt.Sprintf("%s.%s", v.packageName, v.functionName)))
			} else {
				v.sdkDataSources[name] = v.functionName
			}
		} else if m := sdkResourceAnnotation.FindStringSubmatch(line); len(m) > 0 {
			name := m[1]

			if _, ok := v.sdkResources[name]; ok {
				v.err = multierror.Append(v.err, fmt.Errorf("duplicate SDK Resource (%s): %s", name, fmt.Sprintf("%s.%s", v.packageName, v.functionName)))
			} else {
				v.sdkResources[name] = v.functionName
			}
		}
	}

	v.functionName = ""
}

// Visit is called for each node visited by ast.Walk.
func (v *visitor) Visit(node ast.Node) ast.Visitor {
	// Look at functions (not methods) with comments.
	if funcDecl, ok := node.(*ast.FuncDecl); ok && funcDecl.Recv == nil && funcDecl.Doc != nil {
		v.processFuncDecl(funcDecl)
	}

	return v
}
