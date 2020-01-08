// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfgen

import (
	"github.com/pkg/errors"
	csgen "github.com/pulumi/pulumi/pkg/codegen/dotnet"

	"github.com/pulumi/pulumi-terraform-bridge/pkg/tfbridge"
)

type csharpGenerator struct {
	pkg         string
	version     string
	info        tfbridge.ProviderInfo
	overlaysDir string
	outDir      string
}

// newCSharpGenerator returns a language generator that understands how to produce Go packages.
func newCSharpGenerator(pkg, version string, info tfbridge.ProviderInfo, overlaysDir, outDir string) langGenerator {
	return &csharpGenerator{
		pkg:         pkg,
		version:     version,
		info:        info,
		overlaysDir: overlaysDir,
		outDir:      outDir,
	}
}

// typeName returns a type name for a given resource type.
func (g *csharpGenerator) typeName(r *resourceType) string {
	if r.info != nil && r.info.CSharpName != "" {
		return r.info.CSharpName
	}
	return r.name
}

// emitPackage emits an entire package pack into the configured output directory with the configured settings.
func (g *csharpGenerator) emitPackage(pack *pkg) error {
	ppkg, err := genPulumiSchema(pack, g.pkg, g.version, g.info)
	if err != nil {
		return errors.Wrap(err, "generating Pulumi schema")
	}

	var extraCSharpFiles map[string][]byte
	if csi := g.info.JavaScript; csi != nil && csi.Overlay != nil {
		extraCSharpFiles, err = getOverlayFiles(csi.Overlay, ".cs", g.outDir)
		if err != nil {
			return err
		}
	}

	files, err := csgen.GeneratePackage(tfgen, ppkg, extraCSharpFiles)
	if err != nil {
		return errors.Wrap(err, "generating Pulumi package")
	}

	for f, contents := range files {
		if err := emitFile(g.outDir, f, contents); err != nil {
			return errors.Wrapf(err, "emitting file %v", f)
		}
	}
	return nil
}
