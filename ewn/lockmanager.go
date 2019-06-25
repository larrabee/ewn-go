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
	err  error
	path string
}

func (e *ErrLockAquireFailed) Error() string {
	return fmt.Sprintf("can not acquire lock, file path: %s; err: %s", e.path, e.err)
}

// ErrLockAquireFailed : error, that returned then lock acquire failed with error
type ErrLockReleaseFailed struct {
	err  error
	path string
}

func (e *ErrLockReleaseFailed) Error() string {
	return fmt.Sprintf("failed to release lock, file path: %s; err: %s", e.path, e.err)
}

// ErrLockAlreadyAquired : error, that returned then lock already acquired by another process
type ErrLockAlreadyAquired struct {
	pid int64
}

func (e *ErrLockAlreadyAquired) Error() string {
	return fmt.Sprintf("lock already acquired by process with pid: %d", e.pid)
}

// ErrLockNotOwned : error, that returned then you trying release not owned lock
type ErrLockNotOwned struct{}

func (e *ErrLockNotOwned) Error() string {
	return "Can not release not owned lock. You must acquire it first."
}

// Acquire trying to acquire lock. Return error on error or nil on success
func (l *Lock) Acquire() error {
	if l.lockFilePath == "" {
		lockFileName := "ewn-" + hashOfString(l.Key) + ".lock"
		l.lockFilePath = path.Join(os.TempDir(), lockFileName)
	}
	if _, err := os.Stat(l.lockFilePath); !os.IsNotExist(err) {
		pid, err := readLockFile(l.lockFilePath)
		if err != nil {
			return err
		}
		if isPidExist(pid) {
			return &ErrLockAlreadyAquired{pid}
		}
	}
	err := writeLockFile(l.lockFilePath)
	if err != nil {
		return &ErrLockAquireFailed{path: l.lockFilePath, err: err}
	}
	l.lockAcquired = true
	return nil
}

// Release lock. Return error on error or nil on success
func (l *Lock) Release() error {
	if !l.lockAcquired {
		return &ErrLockNotOwned{}
	}
	err := os.Remove(l.lockFilePath)
	if err != nil {
		return &ErrLockReleaseFailed{path: l.lockFilePath, err: err}
	}
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
	file, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	if !file.IsDir() {
		return false
	}
	return true
}

func readLockFile(path string) (int64, error) {
	fileByte, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, &ErrLockAquireFailed{path: path, err: err}
	}
	pid, _ := strconv.ParseInt(string(fileByte), 10, 32)
	return pid, nil
}

func writeLockFile(path string) error {
	pid := strconv.Itoa(os.Getpid())
	err := ioutil.WriteFile(path, []byte(pid), 0664)
	return err
}
