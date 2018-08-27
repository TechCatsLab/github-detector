/*
 * Revision History:
 *     Initial: 2018/07/27        Li Zebang
 */

package glide

import (
	"testing"
	"time"
)

var (
	testyaml = `package: github.com/Masterminds/glide
homepage: https://glide.sh
license: MIT
owners:
- name: Matt Butcher
  email: technosophos@gmail.com
  homepage: http://technosophos.com/
- name: Matt Farina
  email: matt@mattfarina.com
  homepage: https://www.mattfarina.com/
import:
- package: gopkg.in/yaml.v2
- package: github.com/Masterminds/vcs
  version: ^1.12.0
- package: github.com/codegangsta/cli
  version: ^1.16.0
- package: github.com/Masterminds/semver
  version: ^1.4.0
- package: github.com/mitchellh/go-homedir
`

	glideyaml = Config{
		Name:    "github.com/Masterminds/glide",
		Home:    "https://glide.sh",
		License: "MIT",
		Owners: Owners([]*Owner{
			&Owner{
				Name:  "Matt Butcher",
				Email: "technosophos@gmail.com",
				Home:  "http://technosophos.com/",
			},
			&Owner{
				Name:  "Matt Farina",
				Email: "matt@mattfarina.com",
				Home:  "https://www.mattfarina.com/",
			},
		}),
		Imports: Dependencies([]*Dependency{
			&Dependency{
				Name: "gopkg.in/yaml.v2",
			},
			&Dependency{
				Name:      "github.com/Masterminds/vcs",
				Reference: "^1.12.0",
			},
			&Dependency{
				Name:      "github.com/codegangsta/cli",
				Reference: "^1.16.0",
			},
			&Dependency{
				Name:      "github.com/Masterminds/semver",
				Reference: "^1.4.0",
			},
			&Dependency{
				Name: "github.com/mitchellh/go-homedir",
			},
		}),
	}

	testlock = `hash: 1f13d16b2759f4c698bf1fa66a55ef06a42ba86859153c478f903e60502a1273
updated: 2017-10-04T10:27:41.570512797-04:00
imports:
- name: github.com/codegangsta/cli
  version: cfb38830724cc34fedffe9a2a29fb54fa9169cd1
- name: github.com/Masterminds/semver
  version: 15d8430ab86497c5c0da827b748823945e1cf1e1
- name: github.com/Masterminds/vcs
  version: 6f1c6d150500e452704e9863f68c2559f58616bf
- name: github.com/mitchellh/go-homedir
  version: b8bc1bf767474819792c23f32d8286a45736f1c6
- name: gopkg.in/yaml.v2
  version: a3f3340b5840cee44f372bddb5880fcbc419b46a
testImports: []
`

	glidelocktime, _ = time.Parse(time.RFC3339Nano, "2017-10-04T10:27:41.570512797-04:00")

	glidelock = Lockfile{
		Hash:    "1f13d16b2759f4c698bf1fa66a55ef06a42ba86859153c478f903e60502a1273",
		Updated: glidelocktime,
		Imports: Locks([]*Lock{
			&Lock{
				Name:    "github.com/codegangsta/cli",
				Version: "cfb38830724cc34fedffe9a2a29fb54fa9169cd1",
			},
			&Lock{
				Name:    "github.com/Masterminds/semver",
				Version: "15d8430ab86497c5c0da827b748823945e1cf1e1",
			},
			&Lock{
				Name:    "github.com/Masterminds/vcs",
				Version: "6f1c6d150500e452704e9863f68c2559f58616bf",
			},
			&Lock{
				Name:    "github.com/mitchellh/go-homedir",
				Version: "b8bc1bf767474819792c23f32d8286a45736f1c6",
			},
			&Lock{
				Name:    "gopkg.in/yaml.v2",
				Version: "a3f3340b5840cee44f372bddb5880fcbc419b46a",
			},
		}),
	}
)

func Test_glideParser_ParseConfig(t *testing.T) {
	config, err := GlidePraser.ParseConfig([]byte(testyaml))
	if err != nil {
		t.Errorf("Test_glideParser_ParseConfig() error %s\n", err)
	}
	if config.Name != glideyaml.Name {
		t.Errorf("Test_glideParser_ParseConfig() config.Name = %s, not %s\n", config.Name, glideyaml.Name)
	}
	if config.Description != glideyaml.Description {
		t.Errorf("Test_glideParser_ParseConfig() config.Description = %s, not %s\n", config.Description, glideyaml.Description)
	}
	if config.Home != glideyaml.Home {
		t.Errorf("Test_glideParser_ParseConfig() config.Home = %s, not %s\n", config.Home, glideyaml.Home)
	}
	if config.License != glideyaml.License {
		t.Errorf("Test_glideParser_ParseConfig() config.License = %s, not %s\n", config.License, glideyaml.License)
	}
	for index := range config.Owners {
		if config.Owners[index].Name != glideyaml.Owners[index].Name ||
			config.Owners[index].Email != glideyaml.Owners[index].Email ||
			config.Owners[index].Home != glideyaml.Owners[index].Home {
			t.Errorf("Test_glideParser_ParseConfig() config.Owners[%d] = %v, not %v\n", index, config.Owners[index], glideyaml.Owners[index])
		}
	}
	for index := range config.Ignore {
		if config.Ignore[index] != glideyaml.Ignore[index] {
			t.Errorf("Test_glideParser_ParseConfig() config.Ignore[%d] = %v, not %v\n", index, config.Ignore[index], glideyaml.Ignore[index])
		}
	}
	for index := range config.Exclude {
		if config.Exclude[index] != glideyaml.Exclude[index] {
			t.Errorf("Test_glideParser_ParseConfig() config.Exclude[%d] = %v, not %v\n", index, config.Exclude[index], glideyaml.Exclude[index])
		}
	}
	for index := range config.Imports {
		if config.Imports[index].Name != glideyaml.Imports[index].Name ||
			config.Imports[index].Reference != glideyaml.Imports[index].Reference ||
			config.Imports[index].Pin != glideyaml.Imports[index].Pin ||
			config.Imports[index].Repository != glideyaml.Imports[index].Repository ||
			config.Imports[index].VcsType != glideyaml.Imports[index].VcsType {
			t.Errorf("Test_glideParser_ParseConfig() config.Imports[%d] = %v, not %v\n", index, config.Imports[index], glideyaml.Imports[index])
		}
		for subpackagesIndex := range config.Imports[index].Subpackages {
			if config.Imports[index].Subpackages[subpackagesIndex] != glideyaml.Imports[index].Subpackages[subpackagesIndex] {
				t.Errorf("Test_glideParser_ParseConfig() config.Imports[%d].Subpackages[%d] = %v, not %v\n", index, subpackagesIndex, config.Imports[index].Subpackages[subpackagesIndex], glideyaml.Imports[index].Subpackages[subpackagesIndex])
			}
		}
		for archIndex := range config.Imports[index].Arch {
			if config.Imports[index].Arch[archIndex] != glideyaml.Imports[index].Arch[archIndex] {
				t.Errorf("Test_glideParser_ParseConfig() config.Imports[%d].Arch[%d] = %v, not %v\n", index, archIndex, config.Imports[index].Arch[archIndex], glideyaml.Imports[index].Arch[archIndex])
			}
		}
		for osIndex := range config.Imports[index].Os {
			if config.Imports[index].Os[osIndex] != glideyaml.Imports[index].Os[osIndex] {
				t.Errorf("Test_glideParser_ParseConfig() config.Imports[%d].Os[%d] = %v, not %v\n", index, osIndex, config.Imports[index].Os[osIndex], glideyaml.Imports[index].Os[osIndex])
			}
		}
	}
	for index := range config.DevImports {
		if config.DevImports[index].Name != glideyaml.DevImports[index].Name ||
			config.DevImports[index].Reference != glideyaml.DevImports[index].Reference ||
			config.DevImports[index].Pin != glideyaml.DevImports[index].Pin ||
			config.DevImports[index].Repository != glideyaml.DevImports[index].Repository ||
			config.DevImports[index].VcsType != glideyaml.DevImports[index].VcsType {
			t.Errorf("Test_glideParser_ParseConfig() config.DevImports[%d] = %v, not %v\n", index, config.DevImports[index], glideyaml.DevImports[index])
		}
		for subpackagesIndex := range config.DevImports[index].Subpackages {
			if config.DevImports[index].Subpackages[subpackagesIndex] != glideyaml.DevImports[index].Subpackages[subpackagesIndex] {
				t.Errorf("Test_glideParser_ParseConfig() config.DevImports[%d].Subpackages[%d] = %v, not %v\n", index, subpackagesIndex, config.DevImports[index].Subpackages[subpackagesIndex], glideyaml.DevImports[index].Subpackages[subpackagesIndex])
			}
		}
		for archIndex := range config.DevImports[index].Arch {
			if config.DevImports[index].Arch[archIndex] != glideyaml.DevImports[index].Arch[archIndex] {
				t.Errorf("Test_glideParser_ParseConfig() config.DevImports[%d].Arch[%d] = %v, not %v\n", index, archIndex, config.DevImports[index].Arch[archIndex], glideyaml.DevImports[index].Arch[archIndex])
			}
		}
		for osIndex := range config.DevImports[index].Os {
			if config.DevImports[index].Os[osIndex] != glideyaml.DevImports[index].Os[osIndex] {
				t.Errorf("Test_glideParser_ParseConfig() config.DevImports[%d].Os[%d] = %v, not %v\n", index, osIndex, config.DevImports[index].Os[osIndex], glideyaml.DevImports[index].Os[osIndex])
			}
		}
	}
}

func Test_glideParser_ParseLockfile(t *testing.T) {
	lockfile, err := GlidePraser.ParseLockfile([]byte(testlock))
	if err != nil {
		t.Errorf("Test_glideParser_ParseLockfile() error %s\n", err)
	}
	if lockfile.Hash != glidelock.Hash {
		t.Errorf("Test_glideParser_ParseLockfile() lockfile.Hash = %s, not %s\n", lockfile.Hash, glidelock.Hash)
	}
	if !lockfile.Updated.Equal(glidelock.Updated) {
		t.Errorf("Test_glideParser_ParseLockfile() lockfile.Updated = %s, not %s\n", lockfile.Updated, glidelock.Updated)
	}
	for index := range lockfile.Imports {
		if lockfile.Imports[index].Name != glidelock.Imports[index].Name ||
			lockfile.Imports[index].Version != glidelock.Imports[index].Version ||
			lockfile.Imports[index].Repository != glidelock.Imports[index].Repository ||
			lockfile.Imports[index].VcsType != glidelock.Imports[index].VcsType {
			t.Errorf("Test_glideParser_ParseLockfile() lockfile.Imports[%d] = %v, not %v\n", index, lockfile.Imports[index], glidelock.Imports)
		}
		for subpackagesIndex := range lockfile.Imports[index].Subpackages {
			if lockfile.Imports[index].Subpackages[subpackagesIndex] != glidelock.Imports[index].Subpackages[subpackagesIndex] {
				t.Errorf("Test_glideParser_ParseConfig() lockfile.Imports[%d].Subpackages[%d] = %v, not %v\n", index, subpackagesIndex, lockfile.Imports[index].Subpackages[subpackagesIndex], glidelock.Imports[index].Subpackages[subpackagesIndex])
			}
		}
		for archIndex := range lockfile.Imports[index].Arch {
			if lockfile.Imports[index].Arch[archIndex] != glidelock.Imports[index].Arch[archIndex] {
				t.Errorf("Test_glideParser_ParseConfig() lockfile.Imports[%d].Arch[%d] = %v, not %v\n", index, archIndex, lockfile.Imports[index].Arch[archIndex], glidelock.Imports[index].Arch[archIndex])
			}
		}
		for osIndex := range lockfile.Imports[index].Os {
			if lockfile.Imports[index].Os[osIndex] != glidelock.Imports[index].Os[osIndex] {
				t.Errorf("Test_glideParser_ParseConfig() lockfile.Imports[%d].Os[%d] = %v, not %v\n", index, osIndex, lockfile.Imports[index].Os[osIndex], glidelock.Imports[index].Os[osIndex])
			}
		}
	}
	for index := range lockfile.DevImports {
		if lockfile.DevImports[index].Name != glidelock.DevImports[index].Name ||
			lockfile.DevImports[index].Version != glidelock.DevImports[index].Version ||
			lockfile.DevImports[index].Repository != glidelock.DevImports[index].Repository ||
			lockfile.DevImports[index].VcsType != glidelock.DevImports[index].VcsType {
			t.Errorf("Test_glideParser_ParseLockfile() lockfile.DevImports[%d] = %v, not %v\n", index, lockfile.DevImports[index], glidelock.DevImports)
		}
		for subpackagesIndex := range lockfile.DevImports[index].Subpackages {
			if lockfile.DevImports[index].Subpackages[subpackagesIndex] != glidelock.DevImports[index].Subpackages[subpackagesIndex] {
				t.Errorf("Test_glideParser_ParseConfig() lockfile.DevImports[%d].Subpackages[%d] = %v, not %v\n", index, subpackagesIndex, lockfile.DevImports[index].Subpackages[subpackagesIndex], glidelock.DevImports[index].Subpackages[subpackagesIndex])
			}
		}
		for archIndex := range lockfile.DevImports[index].Arch {
			if lockfile.DevImports[index].Arch[archIndex] != glidelock.DevImports[index].Arch[archIndex] {
				t.Errorf("Test_glideParser_ParseConfig() lockfile.DevImports[%d].Arch[%d] = %v, not %v\n", index, archIndex, lockfile.DevImports[index].Arch[archIndex], glidelock.DevImports[index].Arch[archIndex])
			}
		}
		for osIndex := range lockfile.DevImports[index].Os {
			if lockfile.DevImports[index].Os[osIndex] != glidelock.DevImports[index].Os[osIndex] {
				t.Errorf("Test_glideParser_ParseConfig() lockfile.DevImports[%d].Os[%d] = %v, not %v\n", index, osIndex, lockfile.DevImports[index].Os[osIndex], glidelock.DevImports[index].Os[osIndex])
			}
		}
	}
}
