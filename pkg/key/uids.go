package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	uidDirectory = "uids"
)

var (
	valid = regexp.MustCompile("^[0-9]{1,63}.pem$")
	uid   = regexp.MustCompile("^[0-9]+")
)

type UIDs struct {
	dir string
	UID uint64
}

func NewUIDs(dir string, uid uint64) (*UIDs, error) {

	uids := &UIDs{
		dir: dir,
		UID: uid,
	}

	if err := uids.ensureUIDsDirectory(); err != nil {
		return nil, err
	}

	return uids, nil
}

func (u *UIDs) NewUidFile(uid string, pk *rsa.PublicKey) error {
	path := filepath.Join(u.uidsPath(), uid)
	pubBlock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pk)}
	return writeKeyPemFile(path, pubBlock)
}

func (u *UIDs) ReadUidFile(uid string) (*rsa.PublicKey, error) {
	path := filepath.Join(u.uidsPath(), uid)
	return readPublicKey(path)
}

func (u *UIDs) ensureUIDsDirectory() error {
	stat, err := os.Stat(u.uidsPath())
	if os.IsNotExist(err) {
		if err := os.Mkdir(u.uidsPath(), 0700); err != nil {
			return fmt.Errorf("failed to create uids directory: %v", err)
		}

		return nil

	} else if err != nil {
		return fmt.Errorf("error checking status of uids directory: %v", err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("uids path is not a directory: %s", u.uidsPath())
	}

	return nil
}

func (u *UIDs) UIDsFromFile() (map[string]bool, error) {
	fs, err := ioutil.ReadDir(u.uidsPath())
	if err != nil {
		return nil, fmt.Errorf("failed to list uid files: %v", err)
	}

	uids := make(map[string]bool)
	for _, f := range fs {
		if _, err := strconv.Atoi(f.Name()); err != nil {
			uids[f.Name()] = true
		}
	}

	return uids, nil
}

func (u *UIDs) uidsPath() string {
	return filepath.Join(u.dir, uidDirectory)
}

//func RetrieveUIDPublicKeys(dir string) (map[uint64]*rsa.PublicKey, error) {
//	var result *multierror.Error
//
//	if err := ensureKeyDirectory(dir); err != nil {
//		return nil, fmt.Errorf("failed to ensure key directory: %v", err)
//	}
//
//	fs, err := ioutil.ReadDir(dir)
//	if err != nil {
//		return nil, fmt.Errorf("failed to list uid keys in directory: %v", err)
//	}
//
//	keys := make(map[uint64]*rsa.PublicKey)
//
//	for _, f := range fs {
//		if !validName(f.Name()) {
//			continue
//		}
//		k, err := readPublicKey(fmt.Sprintf("%s/%s", dir, f.Name()))
//		if err != nil {
//			result = multierror.Append(result, fmt.Errorf("failed to read a uid's public key: %v", err))
//			continue
//		}
//
//		uuid, err := uidFromName(f.Name())
//		if err != nil {
//			result = multierror.Append(result, fmt.Errorf("failed convert file name to uuid uint64: %v", err))
//			continue
//		}
//
//		keys[uuid] = k
//	}
//
//	return keys, result.ErrorOrNil()
//}
//
//func validName(name string) bool {
//	return valid.MatchString(name)
//}
//
//func uidFromName(name string) (uint64, error) {
//	return strconv.ParseUint(uid.FindString(name), 10, 64)
//}
//
//func RetrieveLocalUID(dir string) (uint64, error) {
//
//	//client.Cl
//	//th := filepath.Join(dir, configFile)
//	//f, err := os.Open(path)
//	//if err != nil {
//	//	if os.IsNotExist(err) {
//
//	//		f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
//	//		if err != nil {
//	//			return 0, fmt.Errorf("failed to create config file: %v", err)
//	//		}
//
//	//	} else {
//	//		return 0, fmt.Errorf("failed to open config file: %v", err)
//	//	}
//	//}
//
//	//defer f.Close()
//
//	//b, err := ioutil.ReadFile(path)
//	//if err != nil {
//	//	return 0, fmt.Errorf("failed to read from config file: %v", err)
//	//}
//
//	//if len(b) == 0 {
//	//	return 0, nil
//	//}
//
//	//uid, err := strconv.ParseInt(string(b[:len(b)-1]), 10, 64)
//	//if err != nil {
//	//	return 0, fmt.Errorf("failed to parse uid from file: %v", err)
//	//}
//
//	//return uint64(uid), nil
//	return 0, nil
//}
