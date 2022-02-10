package configuration

func EnvironmentProviderWithValueCallback(key, value string) (finalKey string, finalValue interface{}) {
	return key, value
}
