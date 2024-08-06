package provider

import "fmt"

// LoadProviderList load providerSource, parse it and register in List Provider
func LoadProviderList(providerSource string) (List, error) {
	pl := createProviderList()

	schema, err := loadJSON(providerSource)
	if err != nil {
		return nil, fmt.Errorf("error occurred when loading providers schema: %w", err)
	}

	for _, s := range schema {
		newP, errP := createProvider(&s)
		if errP != nil {
			return nil, fmt.Errorf("error occurred when create provider: %w", errP)
		}

		errReg := pl.Register(newP)
		if errReg != nil {
			return nil, fmt.Errorf("error occurred when registering new provider %s: %w", s.Name, err)
		}
	}

	return pl, nil
}
