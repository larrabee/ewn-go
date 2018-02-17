package ewn

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// Lock manage lock file
type Lock struct {
	Key          string
	lockFilePath string
	lockAcquired bool
}

// ErrLockAquireFailed : error, that returned then lock acquire failed with error
type ErrLockAquireFailed struct {
	message string
}
// ErrLockAlreadyAquired : error, that returned then lock already acquired by another process
type ErrLockAlreadyAquired struct {
	message string
}

// ErrLockNotOwned : error, that returned then you trying release not owned lock
type ErrLockNotOwned struct {
	message string
}

func (e *ErrLockAquireFailed) Error() string {
	return e.message
}

func (e *ErrLockAlreadyAquired) Error() string {
	return e.message
}

func (e *ErrLockNotOwned) Error() string {
	return e.message
}

// Acquire trying to acquire lock. Return error on error or nil on success
func (l *Lock) Acquire() error {
	if l.lockFilePath == "" {
		lockFileName := "ewn-" + hashOfString(l.Key) + ".lock"
		l.lockFilePath = path.Join(os.TempDir(), lockFileName)
	}
	if _, err1 := os.Stat(l.lockFilePath); !os.IsNotExist(err1) {
		pid, err2 := readLockFile(l.lockFilePath)
		if err2 != nil {
			return err2
		}
		if isPidExist(pid) {
			return &ErrLockAlreadyAquired{fmt.Sprintf("Lock already acquired by process with pid: %d", pid)}
		}
		err3 := writeLockFile(l.lockFilePath)
		if err3 != nil {
			return &ErrLockAquireFailed{fmt.Sprintf("Can not write lock file: %s with error: %s", l.lockFilePath, err3)}
		}
		l.lockAcquired = true
		return nil
	}
	err1 := writeLockFile(l.lockFilePath)
	if err1 != nil {
		return &ErrLockAquireFailed{fmt.Sprintf("Can not write lock file: %s with error: %s", l.lockFilePath, err1)}
	}
	l.lockAcquired = true
	return nil
}

// Release lock. Return error on error or nil on success
func (l *Lock) Release() error {
	if !l.lockAcquired {
		return &ErrLockNotOwned{"Can not release not owned lock. You must acquire it first."}
	}
	_ = os.Remove(l.lockFilePath)
	l.lockAcquired = false
	return nil
}

func hashOfString(str string) (hash string) {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	hash = hex.EncodeToString(hasher.Sum(nil))
	return
}

func isPidExist(pid int64) bool {
	file, err1 := os.Stat(fmt.Sprintf("/proc/%d", pid))
	if os.IsNotExist(err1) {
		return false
	}
	if err1 != nil {
		return false
	}
	if !file.IsDir() {
		return false
	}
	return true
}

func readLockFile(path string) (pid int64, err error) {
	fileByte, err1 := ioutil.ReadFile(path)
	if err1 != nil {
		return 0, &ErrLockAquireFailed{fmt.Sprintf("Can not read lock file: %s with error: %s", path, err1)}
	}
	pid, err2 := strconv.ParseInt(string(fileByte), 10, 32)
	if err2 != nil {
		return 0, &ErrLockAquireFailed{fmt.Sprintf("Can not parse lock file: %s with error: %s", path, err2)}
	}
	return pid, nil
}

func writeLockFile(path string) error {
	pid := strconv.Itoa(os.Getpid())
	err := ioutil.WriteFile(path, []byte(pid), 0664)
	return err
}
