//go:build windows

package w32

import (
	"syscall"
	"unsafe"
)

var (
	modshcore = syscall.NewLazyDLL("shcore.dll")

	procGetDpiForMonitor = modshcore.NewProc("GetDpiForMonitor")
)

func HasGetDPIForMonitorFunc() bool {
	if IsWindows7() {
		return false
	}
	err := procGetDpiForMonitor.Find()
	return err == nil
}

func GetDPIForMonitor(hmonitor HMONITOR, dpiType MONITOR_DPI_TYPE, dpiX *UINT, dpiY *UINT) uintptr {
	if IsWindows7() {
		screen := GetDC(0)
		x := GetDeviceCaps(screen, LOGPIXELSX)
		y := GetDeviceCaps(screen, LOGPIXELSY)
		ReleaseDC(0, screen)
		*dpiX = UINT(x)
		*dpiY = UINT(y)
		return 0
	}
	ret, _, _ := procGetDpiForMonitor.Call(
		hmonitor,
		uintptr(dpiType),
		uintptr(unsafe.Pointer(dpiX)),
		uintptr(unsafe.Pointer(dpiY)))

	return ret
}
