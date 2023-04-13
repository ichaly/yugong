package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"hash/crc32"
	"net"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func String(value string, defaultValue string) string {
	if len(value) > 0 {
		return value
	}
	return defaultValue
}

func Root() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Dir(path.Join(path.Dir(filename), "../../"))
	}
	return ""
}

// GeneratePassword 给密码进行加密操作
func GeneratePassword(plain string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
}

// ValidatePassword 密码比对
func ValidatePassword(plain string, cipher string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(cipher), []byte(plain)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}

func MD5(s string) string {
	d := []byte(s)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func HashCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func GetMAC() (addresses []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return addresses
	}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		addresses = append(addresses, macAddr)
	}
	return addresses
}

func GetIP() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func SetKeyValue(vi *viper.Viper, key string, value interface{}) bool {
	if strings.HasPrefix(key, "GJ_") || strings.HasPrefix(key, "SG_") {
		key = key[3:]
	}
	uc := strings.Count(key, "_")
	k := strings.ToLower(key)

	if vi.Get(k) != nil {
		vi.Set(k, value)
		return true
	}

	for i := 0; i < uc; i++ {
		k = strings.Replace(k, "_", ".", 1)
		if vi.Get(k) != nil {
			vi.Set(k, value)
			return true
		}
	}

	return false
}
