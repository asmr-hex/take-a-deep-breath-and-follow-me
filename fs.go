package main

import "path"

type Dir struct {
	Name  string
	Files []*File
	Dirs  []*Dir
}

type File struct {
	Name string
}

type FSMap struct {
	Dirs  map[string]*Dir
	Files map[string]*File
}

var (
	FS = &Dir{
		Name: "",
		Dirs: []*Dir{
			{
				Name: "home",
				Dirs: []*Dir{
					usrHomeDir,
				},
			},
		},
	}
	usrHomeDir = &Dir{
		Name: "",
		Files: []*File{
			{Name: ".rip"},
			{Name: "secret.txt"},
		},
	}
	fsMap = IndexFS(FS)
)

// index the filesystem
func IndexFS(dir *Dir) *FSMap {
	fsMap := &FSMap{
		Dirs: map[string]*Dir{
			"/": dir,
		},
		Files: map[string]*File{},
	}

	indexDir(dir, "/", fsMap)

	return fsMap
}

func indexDir(dir *Dir, p string, fsMap *FSMap) *FSMap {
	// index files
	for _, f := range dir.Files {
		fsMap.Files[path.Join(p, f.Name)] = f
	}

	// index dirs
	for _, d := range dir.Dirs {
		newPath := path.Join(p, d.Name)
		fsMap.Dirs[newPath] = d
		indexDir(d, newPath, fsMap)
	}

	return fsMap
}

func GenerateRecursiveDirs() {

}
