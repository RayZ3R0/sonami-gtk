//go:build windows

package secrets

// serviceWindows is a no-op implementation of the Service interface for Windows.
// On Windows, libsecret is not available. Since this fork operates in account-free
// mode, all secret operations are stubs that do nothing.
type serviceWindows struct{}

func (s *serviceWindows) Available() *ServiceError {
	return nil
}

func (s *serviceWindows) Delete(key string) error {
	return nil
}

func (s *serviceWindows) Get(key string) (Item, error) {
	return Item{}, ErrKeyNotFound
}

func (s *serviceWindows) Has(key string) (bool, error) {
	return false, nil
}

func (s *serviceWindows) Set(key string, value Item) error {
	return nil
}

func newService() Service {
	return &serviceWindows{}
}
