package helpers

import (
	"crypto/rand"
	"github.com/spf13/viper"
	"path/filepath"
)

const (
	component      = "Helpers"
	configFileName = "config"
	configFormat   = "yaml"
)

func ConfigInit() error {
	initPaths()

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFormat)
	viper.AddConfigPath(configPath)

	_ = viper.ReadInConfig()

	err := writeInitial()
	if err != nil {
		return err
	}
	return nil
}

func writeInitial() error {
	if string(GetSessionsKey()) == "" {
		viper.Set("sessionsKey", GenerateRandomKey(64))
		err := write()
		return err
	}

	if string(GetEncryptionKey()) == "" {
		viper.Set("encryptionkey", GenerateRandomKey(32))
		err := write()
		return err
	}
	return nil
}

func GetSessionsKey() []byte {
	return []byte(viper.GetString("sessionsKey"))
}

func GetEncryptionKey() []byte {
	return []byte(viper.GetString("encryptionkey"))
}

func write() error {
	if err := viper.WriteConfigAs(filepath.Join(configPath, configFileName+"."+configFormat)); err != nil {
		return err
	}
	return nil
}

func GenerateRandomKey(l int) string {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		LogError(err.Error(), component)
	}
	return string(b)
}
