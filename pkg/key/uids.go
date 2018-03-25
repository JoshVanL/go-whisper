package key

import (
	"fmt"
	"regexp"
	"crypto/rsa"
	"io/ioutil"
	"strconv"
	"os"

	"github.com/hashicorp/go-multierror"
)

var (
	valid      = regexp.MustCompile("^[0-9]{1,63}.pem$")
	uid        = regexp.MustCompile("^[0-9]+")
	configFile = "whisper.conf"
)

func RetrieveUIDPublicKeys() (map[uint64]*rsa.PublicKey, error) {
	var result *multierror.Error

	if err := ensureKeyDirectory(); err != nil {
		return nil, fmt.Errorf("failed to ensure key directory: %v", err)
	}

	dir, err := uidPubkicKeyDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to ensure uid public key directory: %v", err)
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to list uid keys in directory: %v", err)
	}

	keys := make(map[uint64]*rsa.PublicKey)

	for _, f := range fs {
		if !validName(f.Name()) {
			continue
		}
		k, err := readPublicKey(fmt.Sprintf("%s/%s", dir, f.Name()))
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to read a uid's public key: %v", err))
			continue
		}

		uuid, err := uidFromName(f.Name())
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("failed convert file name to uuid uint64: %v", err))
			continue
		}

		keys[uuid] = k
	}

	return keys, result.ErrorOrNil()
}

func RetrieveLocalUID() (uint64, error) {
	dir, err := keyDir()
	if err != nil {
		return 0, fmt.Errorf("failed to get key directory: %v", err)
	}

	path := fmt.Sprintf("%s/%s", dir, configFile)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {

			f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				return 0, fmt.Errorf("failed to create config file: %v", err)
			}

		} else {
			return 0, fmt.Errorf("failed to open config file: %v", err)
		}
	}

	defer f.Close()

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read from config file: %v", err)
	}

	if len(b) == 0 {
		return 0, nil
	}

	uid, err := strconv.ParseInt(string(b[:len(b)-1]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uid from file: %v", err)
	}

	return uint64(uid), nil
}

func validName(name string) bool { return valid.MatchString(name) }

func uidFromName(name string) (uint64, error) {
	return strconv.ParseUint(uid.FindString(name), 10, 64)
}
