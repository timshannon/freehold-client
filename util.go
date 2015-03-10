// Copyright 2015 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package freeholdclient

import (
	"path"
	"strings"
)

var supportedVersions = map[string]struct{}{"v1": struct{}{}}

// splitRootAndPath splits the first item in a path from the rest
// /v1/file/test.txt:
// root = "v1"
// path = "/file/test.txt"
func splitRootAndPath(pattern string) (root, path string) {
	if pattern == "" {
		panic("Invalid pattern")
	}
	if pattern[:1] == "/" {
		pattern = pattern[1:]
	}
	split := strings.SplitN(pattern, "/", 2)
	root = split[0]
	if len(split) < 2 {
		path = "/"
	} else {
		path = "/" + split[1]
	}
	return root, path
}

func isVersion(version string) bool {
	_, ok := supportedVersions[version]
	return ok
}

func propertyPath(filePath string) string {
	root, p := splitRootAndPath(filePath)
	if !isVersion(root) {
		//root is app
		ver, p := splitRootAndPath(p)
		return path.Join(root, ver, "properties", p)
	}

	return path.Join(root, "properties", p)

}
