package zipper

type Service interface {
	Create([]string, []string) error
}
