package vfs

type fileset struct {
	set map[string]File
}

func (f *fileset) Len() int {
	return len(f.set)
}

func (f *fileset) Contains(file File) bool {
	exists := f.set[file.Path()]
	return exists != nil
}
func (f *fileset) Add(file File) {
	f.set[file.Path()] = file
}

func (f *fileset) Remove(file File) {
	delete(f.set, file.Path())
}

func (f *fileset) Clear() {
	f.set = make(map[string]File)
}

func (f *fileset) Iterator() []File {
	result := make([]File, 0)
	for _, file := range f.set {
		result = append(result, file)
	}
	return result
}

func NewFileSet() *fileset {
	return &fileset{
		set: make(map[string]File),
	}
}
