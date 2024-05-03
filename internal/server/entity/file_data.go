package entity

import "sort"

type FileData struct {
	Index       int32
	ContentHash string
}

func SortFileData(fd []FileData) {
	sort.SliceStable(fd, func(i, j int) bool {
		return fd[i].Index < fd[j].Index
	})
}

func FileDataContentHashes(fd []FileData) []string {
	var out []string
	for _, f := range fd {
		out = append(out, f.ContentHash)
	}
	return out
}
