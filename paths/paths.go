package paths

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func UserConfigDir() (string, error) {
	configDir := os.Getenv("XTASK_CONFIG_HOME")
	if configDir != "" {
		return configDir, nil
	}

	configDir = os.Getenv("XDG_CONFIG_HOME")
	if configDir != "" {
		return filepath.Join(configDir, "xtask"), nil
	}

	configDir, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(configDir, "xtask"), nil
	}

	return "", errors.New("Could not determine user config directory: " + err.Error())
}

func UserDataDir() (string, error) {
	dataDir := os.Getenv("XTASK_DATA_HOME")
	if dataDir != "" {
		return dataDir, nil
	}

	dataDir = os.Getenv("XDG_DATA_HOME")
	if dataDir != "" {
		return filepath.Join(dataDir, "xtask"), nil
	}

	if runtime.GOOS == "windows" {
		dataDir = os.Getenv("LOCALAPPDATA")
		if dataDir != "" {
			return filepath.Join(dataDir, "xtask", "data"), nil
		}
	} else {
		dataDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(dataDir, ".local", "share", "xtask"), nil
		}
	}
	return "", errors.New("could not determine user data directory")
}

func UserCacheDir() (string, error) {
	cacheDir := os.Getenv("XTASK_CACHE_HOME")
	if cacheDir != "" {
		return cacheDir, nil
	}

	cacheDir = os.Getenv("XDG_CACHE_HOME")
	if cacheDir != "" {
		return filepath.Join(cacheDir, "xtask"), nil
	}

	if runtime.GOOS == "windows" {
		cacheDir = os.Getenv("LOCALAPPDATA")
		if cacheDir != "" {
			return filepath.Join(cacheDir, "Cache", "xtask"), nil
		}

		return "", errors.New("could not determine user cache directory")
	}

	cacheDir, err := os.UserCacheDir()
	if err == nil {
		return filepath.Join(cacheDir, "xtask"), nil
	}

	return "", errors.New("Could not determine user cache directory: " + err.Error())
}

func UserStateDir() (string, error) {
	stateDir := os.Getenv("XTASK_STATE_HOME")
	if stateDir != "" {
		return stateDir, nil
	}

	stateDir = os.Getenv("XDG_STATE_HOME")
	if stateDir != "" {
		return filepath.Join(stateDir, "xtask"), nil
	}

	if runtime.GOOS == "windows" {
		stateDir = os.Getenv("LOCALAPPDATA")
		if stateDir != "" {
			return filepath.Join(stateDir, "State", "xtask"), nil
		}
	} else {
		stateDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(stateDir, ".local", "state", "xtask"), nil
		}
	}
	return "", errors.New("could not determine user state directory")
}
