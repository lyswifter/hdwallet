package hdwallet

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"crypto/sha256"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/tyler-smith/go-bip39"
	"math"
	"strconv"
	"time"
)

func entropy() ([]byte, error) {
	randomBytes := make([]byte, 0)
	cpuPercent, _ := cpu.Percent(time.Second, false)
	memory, _ := mem.VirtualMemory()
	diskStatus, _ := disk.Usage("/")

	ioCounters, _ := net.IOCounters(true)
	netWork := strconv.Itoa(int(ioCounters[0].BytesSent + ioCounters[0].BytesRecv))


	cRandBytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, cRandBytes); err != nil {
		return []byte{}, err
	}

	randomBytes = append(randomBytes, cRandBytes...)
	randomBytes = append(randomBytes, float64ToByte(cpuPercent[0])...)
	randomBytes = append(randomBytes, float64ToByte(memory.UsedPercent)...)
	randomBytes = append(randomBytes, float64ToByte(diskStatus.UsedPercent)...)
	randomBytes = append(randomBytes, []byte(netWork)...)

	random := sha256.Sum256(randomBytes)
	return random[:16], nil
}

func NewMnemonic() (string, error){
	entropyBytes, err := entropy()
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropyBytes)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}


func GenerateSeedFromMnemonic(mnemonic, password string) ([]byte, error) {
	seedBytes, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return []byte{}, err
	}

	return seedBytes, nil
}

func GetExtendSeedFromPath(path string, seed []byte) ([]byte, error) {
	extendedKey, err := NewMaster(seed)
	if err != nil {
		return nil, err
	}

	derivationPath, err := ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}

	for _, index := range derivationPath {
		childExtendedKey, err := extendedKey.Child(index)
		if err != nil {
			return nil, err
		}
		extendedKey = childExtendedKey
	}

	return extendedKey.key, nil
}


func float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

