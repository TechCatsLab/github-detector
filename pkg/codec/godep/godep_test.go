/*
 * Revision History:
 *     Initial: 2018/07/27        Li Zebang
 */

package godep

import (
	"testing"
)

var (
	testgodep = `{
    "ImportPath": "github.com/tools/godep",
    "GoVersion": "go1.7",
    "GodepVersion": "v74",
    "Deps": [
        {
            "ImportPath": "github.com/kr/fs",
            "Rev": "2788f0dbd16903de03cb8186e5c7d97b69ad387b"
        },
        {
            "ImportPath": "github.com/kr/pretty",
            "Comment": "go.weekly.2011-12-22-24-gf31442d",
            "Rev": "f31442d60e51465c69811e2107ae978868dbea5c"
        },
        {
            "ImportPath": "github.com/kr/text",
            "Rev": "6807e777504f54ad073ecef66747de158294b639"
        },
        {
            "ImportPath": "github.com/pmezard/go-difflib/difflib",
            "Rev": "f78a839676152fd9f4863704f5d516195c18fc14"
        },
        {
            "ImportPath": "golang.org/x/tools/go/vcs",
            "Rev": "1f1b3322f67af76803c942fd237291538ec68262"
        }
    ]
}`

	godep = Godeps{
		ImportPath:   "github.com/tools/godep",
		GoVersion:    "go1.7",
		GodepVersion: "v74",
		Deps: []Dependency{
			Dependency{
				ImportPath: "github.com/kr/fs",
				Rev:        "2788f0dbd16903de03cb8186e5c7d97b69ad387b",
			},
			Dependency{
				ImportPath: "github.com/kr/pretty",
				Comment:    "go.weekly.2011-12-22-24-gf31442d",
				Rev:        "f31442d60e51465c69811e2107ae978868dbea5c",
			},
			Dependency{
				ImportPath: "github.com/kr/text",
				Rev:        "6807e777504f54ad073ecef66747de158294b639",
			},
			Dependency{
				ImportPath: "github.com/pmezard/go-difflib/difflib",
				Rev:        "f78a839676152fd9f4863704f5d516195c18fc14",
			},
			Dependency{
				ImportPath: "golang.org/x/tools/go/vcs",
				Rev:        "1f1b3322f67af76803c942fd237291538ec68262",
			},
		},
	}
)

func Test_godepParser_ParseGodeps(t *testing.T) {
	godeps, err := GodepParser.ParseGodeps([]byte(testgodep))
	if err != nil {
		t.Errorf("Test_godepParser_ParseGodeps() error %s\n", err)
	}
	if godeps.ImportPath != godep.ImportPath {
		t.Errorf("Test_godepParser_ParseGodeps() godeps.ImportPath = %s, not %s\n", godeps.ImportPath, godep.ImportPath)
	}
	if godeps.GoVersion != godep.GoVersion {
		t.Errorf("Test_godepParser_ParseGodeps() godeps.GoVersion = %s, not %s\n", godeps.GoVersion, godep.GoVersion)
	}
	if godeps.GodepVersion != godep.GodepVersion {
		t.Errorf("Test_godepParser_ParseGodeps() godeps.GodepVersion = %s, not %s\n", godeps.GodepVersion, godep.GodepVersion)
	}
	for index := range godeps.Packages {
		if godeps.Packages[index] != godep.Packages[index] {
			t.Errorf("Test_godepParser_ParseGodeps() godeps.Packages[%d] = %s, not %s\n", index, godeps.Packages[index], godep.Packages[index])
		}
	}
	for index := range godeps.Deps {
		if godeps.Deps[index].ImportPath != godep.Deps[index].ImportPath ||
			godeps.Deps[index].Comment != godep.Deps[index].Comment ||
			godeps.Deps[index].Rev != godep.Deps[index].Rev {
			t.Errorf("Test_godepParser_ParseGodeps() godeps.Deps[%d] = %v, not %v\n", index, godeps.Deps[index], godep.Deps[index])
		}
	}
}
