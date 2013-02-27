package comm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Compile(srcFile string) (binFile string, err error) {
	// Pre-processing
	srcFile = filepath.Abs(srcFile)

	h := md5.New()
	h.Write(srcFile)
	sum := hex.EncodeToString(h.Sum(nil))
	binFile := filepath.Join(BinaryCacheDir, sum[:2], sum)

	// Check source file exists
	srcStat, err := os.Stat(srcFile)
	if err != nil {
		return binFile, err
	}

	// Return OK if binary file exists and source file has not been modified since
	binStat, err := os.Stat(binFile)
	if err == nil {
		srcTime := srcStat.ModTime()
		binTime := binStat.ModTime()

		if binTime.After(srcTime) {
			return binFile, nil
		}
	}

	// Create directories
	err = os.MkdirAll(filepath.Dir(binFile), 0777)
	if err != nil {
		return binFile, err
	}

	// Compile!
	cmd := exec.Command("go", "build", "-o", binFile, srcFile)
	err = cmd.Run()
	if err != nil {
		return binFile, err
	}

	return binFile, nil
}
