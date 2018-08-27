/*
 * Revision History:
 *     Initial: 2018/07/27        Li Zebang
 */

package govendor

import (
	"testing"
)

var (
	testgovendor = `{
	"comment": "",
	"ignore": "test",
	"package": [
		{
			"checksumSHA1": "ZQJSiJXqf+iSGZaOZT/rxEG/Wf8=",
			"path": "github.com/Bowery/prompt",
			"revision": "94b504f42fdd503acc3b3c79ec2b517d90e0de8a",
			"revisionTime": "2018-02-09T14:12:55Z"
		},
		{
			"checksumSHA1": "6VGFARaK8zd23IAiDf7a+gglC8k=",
			"path": "github.com/dchest/safefile",
			"revision": "855e8d98f1852d48dde521e0522408d1fe7e836a",
			"revisionTime": "2015-10-22T12:31:44+02:00"
		},
		{
			"checksumSHA1": "3VJcSYFds0zeIO5opOs0AoKm3Mw=",
			"path": "github.com/google/shlex",
			"revision": "6f45313302b9c56850fc17f99e40caebce98c716",
			"revisionTime": "2015-01-27T13:39:51Z"
		},
		{
			"checksumSHA1": "GcaTbmmzSGqTb2X6qnNtmDyew1Q=",
			"path": "github.com/pkg/errors",
			"revision": "a2d6902c6d2a2f194eb3fb474981ab7867c81505",
			"revisionTime": "2016-06-27T22:23:52Z"
		},
		{
			"checksumSHA1": "CNHEeGnucEUlTHJrLS2kHtfNbws=",
			"path": "golang.org/x/sys/unix",
			"revision": "37707fdb30a5b38865cfb95e5aab41707daec7fd",
			"revisionTime": "2018-02-02T13:35:31Z"
		},
		{
			"checksumSHA1": "ULu8vRll6GP40qFH7OrHVcDSUOc=",
			"path": "golang.org/x/tools/go/vcs",
			"revision": "70252dea49c01dbfc3f800978304eb00eb15b698",
			"revisionTime": "2018-02-09T16:58:40Z"
		},
		{
			"checksumSHA1": "fALlQNY1fM99NesfLJ50KguWsio=",
			"path": "gopkg.in/yaml.v2",
			"revision": "cd8b52f8269e0feb286dfeef29f8fe4d5b397e0b",
			"revisionTime": "2017-04-07T17:21:22Z"
		}
	],
	"rootPath": "github.com/kardianos/govendor"
}`

	govendorfile = File{
		Ignore: "test",
		Package: []*Package{
			&Package{
				ChecksumSHA1: "ZQJSiJXqf+iSGZaOZT/rxEG/Wf8=",
				Path:         "github.com/Bowery/prompt",
				Revision:     "94b504f42fdd503acc3b3c79ec2b517d90e0de8a",
				RevisionTime: "2018-02-09T14:12:55Z",
			},
			&Package{
				ChecksumSHA1: "6VGFARaK8zd23IAiDf7a+gglC8k=",
				Path:         "github.com/dchest/safefile",
				Revision:     "855e8d98f1852d48dde521e0522408d1fe7e836a",
				RevisionTime: "2015-10-22T12:31:44+02:00",
			},
			&Package{
				ChecksumSHA1: "3VJcSYFds0zeIO5opOs0AoKm3Mw=",
				Path:         "github.com/google/shlex",
				Revision:     "6f45313302b9c56850fc17f99e40caebce98c716",
				RevisionTime: "2015-01-27T13:39:51Z",
			},
			&Package{
				ChecksumSHA1: "GcaTbmmzSGqTb2X6qnNtmDyew1Q=",
				Path:         "github.com/pkg/errors",
				Revision:     "a2d6902c6d2a2f194eb3fb474981ab7867c81505",
				RevisionTime: "2016-06-27T22:23:52Z",
			},
			&Package{
				ChecksumSHA1: "CNHEeGnucEUlTHJrLS2kHtfNbws=",
				Path:         "golang.org/x/sys/unix",
				Revision:     "37707fdb30a5b38865cfb95e5aab41707daec7fd",
				RevisionTime: "2018-02-02T13:35:31Z",
			},
			&Package{
				ChecksumSHA1: "ULu8vRll6GP40qFH7OrHVcDSUOc=",
				Path:         "golang.org/x/tools/go/vcs",
				Revision:     "70252dea49c01dbfc3f800978304eb00eb15b698",
				RevisionTime: "2018-02-09T16:58:40Z",
			},
			&Package{
				ChecksumSHA1: "fALlQNY1fM99NesfLJ50KguWsio=",
				Path:         "gopkg.in/yaml.v2",
				Revision:     "cd8b52f8269e0feb286dfeef29f8fe4d5b397e0b",
				RevisionTime: "2017-04-07T17:21:22Z",
			},
		},
		RootPath: "github.com/kardianos/govendor",
	}
)

func Test_govendorParser_ParseFile(t *testing.T) {
	file, err := GovendorParser.ParseFile([]byte(testgovendor))
	if err != nil {
		t.Errorf("Test_govendorParser_ParseFile() error %s\n", err)
	}
	if file.RootPath != govendorfile.RootPath {
		t.Errorf("Test_govendorParser_ParseFile() file.RootPath = %s, not %s\n", file.RootPath, govendorfile.RootPath)
	}
	if file.Comment != govendorfile.Comment {
		t.Errorf("Test_govendorParser_ParseFile() file.Comment = %s, not %s\n", file.Comment, govendorfile.Comment)
	}
	if file.Ignore != govendorfile.Ignore {
		t.Errorf("Test_govendorParser_ParseFile() file.Ignore = %s, not %s\n", file.Ignore, govendorfile.Ignore)
	}
	for index := range file.Package {
		if file.Package[index].Origin != govendorfile.Package[index].Origin ||
			file.Package[index].Path != govendorfile.Package[index].Path ||
			file.Package[index].Tree != govendorfile.Package[index].Tree ||
			file.Package[index].Revision != govendorfile.Package[index].Revision ||
			file.Package[index].RevisionTime != govendorfile.Package[index].RevisionTime ||
			file.Package[index].Version != govendorfile.Package[index].Version ||
			file.Package[index].VersionExact != govendorfile.Package[index].VersionExact ||
			file.Package[index].ChecksumSHA1 != govendorfile.Package[index].ChecksumSHA1 ||
			file.Package[index].Comment != govendorfile.Package[index].Comment {
			t.Errorf("Test_govendorParser_ParseFile() file.Package[%d] = %v, not %v\n", index, file.Package[index], govendorfile.Package[index])
		}
	}
}
