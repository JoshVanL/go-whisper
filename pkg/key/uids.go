package key

import (
	"fmt"
	"regexp"
	"crypto/rsa"
	"io/ioutil"
	"strconv"

	"github.com/hashicorp/go-multierror"
)

var (
	valid = regexp.MustCompile("^[0-9]{1,63}.pem$")
	uid   = regexp.MustCompile("^[0-9]+")
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

func validName(name string) bool { return valid.MatchString(name) }

func uidFromName(name string) (uint64, error) {
	return strconv.ParseUint(uid.FindString(name), 10, 64)
}
