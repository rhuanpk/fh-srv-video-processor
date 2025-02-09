package zipper

type Service interface {
	Create(zipsPaths, sourcesDirs []string) error
}
