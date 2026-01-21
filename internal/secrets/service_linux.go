//go:build linux

package secrets

import golibsecret "github.com/lescuer97/go-libsecret"

var schema *golibsecret.Schema

func init() {
	s, err := golibsecret.NewSchema("dev.dergs.Tonearm", golibsecret.SchemaFlagsNone, map[string]golibsecret.SchemaAttributeType{
		"key": golibsecret.SchemaAttributeString,
	})
	if err != nil {
		panic(err)
	}
	schema = s
}

type serviceLinux struct {
	schema *golibsecret.Schema
}

func (s *serviceLinux) Delete(key string) error {
	_, err := golibsecret.ClearPassword(s.schema, map[string]string{
		"key": key,
	})
	return err
}

func (s *serviceLinux) Get(key string) (Item, error) {
	attrs := golibsecret.NewAttributes()
	attrs.Set("key", key)
	val, err := golibsecret.PasswordLookupSync(s.schema, attrs)
	if val == "" && err == nil {
		return Item{}, ErrKeyNotFound
	}
	return Item{Label: "", Password: val}, err
}

func (s *serviceLinux) Has(key string) (bool, error) {
	_, err := s.Get(key)
	if err == ErrKeyNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *serviceLinux) Set(key string, value Item) error {
	attrs := golibsecret.NewAttributes()
	attrs.Set("key", key)
	return golibsecret.PasswordStoreSync(s.schema, attrs, golibsecret.CollectionDefault, value.Label, value.Password)
}

func newService() Service {
	return &serviceLinux{
		schema: schema,
	}
}
