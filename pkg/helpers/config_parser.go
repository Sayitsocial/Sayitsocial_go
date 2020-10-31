package helpers

import (
	"crypto/rand"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configFileName = "config"
	configFormat   = "yaml"
)

// ConfigInit initializes config file
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
		viper.Set("sessionsKey", generateRandomKey(64))
		err := write()
		return err
	}

	if string(GetEncryptionKey()) == "" {
		viper.Set("encryptionkey", generateRandomKey(32))
		err := write()
		return err
	}

	if string(GetJWTKey()) == "" {
		viper.Set("jwtkey", generateRandomKey(64))
		err := write()
		return err
	}
	return nil
}

// GetSessionsKey returns session key for session store
func GetSessionsKey() []byte {
	return []byte(viper.GetString("sessionsKey"))
}

// GetEncryptionKey returns encryption key for session store
func GetEncryptionKey() []byte {
	return []byte(viper.GetString("encryptionkey"))
}

func GetJWTKey() []byte {
	return []byte(viper.GetString("jwtkey"))
}

func write() error {
	if err := viper.WriteConfigAs(filepath.Join(configPath, configFileName+"."+configFormat)); err != nil {
		return err
	}
	return nil
}

func generateRandomKey(l int) string {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		LogError(err.Error())
	}
	return string(b)
}
