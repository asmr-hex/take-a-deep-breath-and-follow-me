package main

type Dir struct {
	Name  string
	Path  string
	Files []*File
	Dirs  []*Dir
}

type File struct {
	Name string
}

var (
	FS = &Dir{
		Name: "",
		Path: "/",
		Dirs: []*Dir{
			{
				Name: "home",
				Path: "/home/",
				Dirs: []*Dir{
					UsrHome,
				},
			},
		},
	}
	UsrHome = &Dir{
		Name: "", // TODO (cw|4.18.2018) we need to set this once we get the usr
		Path: "",
		Files: []*File{
			{Name: ".rip"},
		},
	}
)
