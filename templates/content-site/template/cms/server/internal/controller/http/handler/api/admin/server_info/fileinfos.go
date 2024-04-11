package server_info

type fileInfo struct {
	Name string `json:"Name"`
	Size int64  `string:"Size"`
}

type fileInfosSortedByName []fileInfo

func (fi fileInfosSortedByName) Len() int           { return len(fi) }
func (fi fileInfosSortedByName) Less(i, j int) bool { return fi[i].Name < fi[j].Name }
func (fi fileInfosSortedByName) Swap(i, j int)      { fi[i], fi[j] = fi[j], fi[i] }
