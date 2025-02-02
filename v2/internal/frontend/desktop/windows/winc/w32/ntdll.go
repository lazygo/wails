//go:build windows

package w32

import (
	"errors"
	"syscall"
	"unsafe"
)

type osVersionInfoEx struct {
	osVersionInfoSize uint32
	majorVersion      uint32
	minorVersion      uint32
	buildNumber       uint32
	platformID        uint32
	csdVersion        [128]uint16
	servicePackMajor  uint16
	servicePackMinor  uint16
	suiteMask         uint16
	productType       byte
	reserved          byte
}

var (
	modntdll = syscall.NewLazyDLL("ntdll.dll")

	procRtlGetVersion = modntdll.NewProc("RtlGetVersion")
)

func HasRtlGetVersionFunc() bool {
	err := procRtlGetVersion.Find()
	return err == nil
}

func RtlGetVersion() (*osVersionInfoEx, error) {
	var info osVersionInfoEx
	info.osVersionInfoSize = uint32(unsafe.Sizeof(info))
	// 调用RtlGetVersion，传入结构体指针
	ret, _, err := procRtlGetVersion.Call(uintptr(unsafe.Pointer(&info)))
	if err != nil {
		return nil, err
	}
	if ret == 0 {
		// STATUS_SUCCESS
		return &info, nil
	}

	return nil, errors.New("proc RtlGetVersion error")
}

func IsWindows7() bool {
	// Windows 7的主版本号为6，次版本号为1
	info, err := RtlGetVersion()
	if err != nil {
		return false
	}
	return info.majorVersion == 6 && info.minorVersion == 1
}
