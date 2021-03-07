package files

type Walker interface {
	Walk(root string) error
}
