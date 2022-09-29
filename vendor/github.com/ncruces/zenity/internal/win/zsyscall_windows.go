// Code generated by 'go generate'; DO NOT EDIT.

package win

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modcomctl32 = windows.NewLazySystemDLL("comctl32.dll")
	modcomdlg32 = windows.NewLazySystemDLL("comdlg32.dll")
	modgdi32    = windows.NewLazySystemDLL("gdi32.dll")
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	modole32    = windows.NewLazySystemDLL("ole32.dll")
	modshell32  = windows.NewLazySystemDLL("shell32.dll")
	moduser32   = windows.NewLazySystemDLL("user32.dll")
	modwtsapi32 = windows.NewLazySystemDLL("wtsapi32.dll")

	procInitCommonControlsEx         = modcomctl32.NewProc("InitCommonControlsEx")
	procChooseColorW                 = modcomdlg32.NewProc("ChooseColorW")
	procCommDlgExtendedError         = modcomdlg32.NewProc("CommDlgExtendedError")
	procGetOpenFileNameW             = modcomdlg32.NewProc("GetOpenFileNameW")
	procGetSaveFileNameW             = modcomdlg32.NewProc("GetSaveFileNameW")
	procCreateFontIndirectW          = modgdi32.NewProc("CreateFontIndirectW")
	procDeleteObject                 = modgdi32.NewProc("DeleteObject")
	procGetDeviceCaps                = modgdi32.NewProc("GetDeviceCaps")
	procActivateActCtx               = modkernel32.NewProc("ActivateActCtx")
	procCreateActCtxW                = modkernel32.NewProc("CreateActCtxW")
	procDeactivateActCtx             = modkernel32.NewProc("DeactivateActCtx")
	procGenerateConsoleCtrlEvent     = modkernel32.NewProc("GenerateConsoleCtrlEvent")
	procGetConsoleWindow             = modkernel32.NewProc("GetConsoleWindow")
	procGetModuleHandleW             = modkernel32.NewProc("GetModuleHandleW")
	procGlobalAlloc                  = modkernel32.NewProc("GlobalAlloc")
	procGlobalFree                   = modkernel32.NewProc("GlobalFree")
	procReleaseActCtx                = modkernel32.NewProc("ReleaseActCtx")
	procCoCreateInstance             = modole32.NewProc("CoCreateInstance")
	procExtractAssociatedIconW       = modshell32.NewProc("ExtractAssociatedIconW")
	procSHBrowseForFolder            = modshell32.NewProc("SHBrowseForFolder")
	procSHCreateItemFromParsingName  = modshell32.NewProc("SHCreateItemFromParsingName")
	procSHGetPathFromIDListEx        = modshell32.NewProc("SHGetPathFromIDListEx")
	procShell_NotifyIconW            = modshell32.NewProc("Shell_NotifyIconW")
	procCallNextHookEx               = moduser32.NewProc("CallNextHookEx")
	procCreateIconFromResourceEx     = moduser32.NewProc("CreateIconFromResourceEx")
	procCreateWindowExW              = moduser32.NewProc("CreateWindowExW")
	procDefWindowProcW               = moduser32.NewProc("DefWindowProcW")
	procDestroyIcon                  = moduser32.NewProc("DestroyIcon")
	procDestroyWindow                = moduser32.NewProc("DestroyWindow")
	procDispatchMessageW             = moduser32.NewProc("DispatchMessageW")
	procEnableWindow                 = moduser32.NewProc("EnableWindow")
	procEnumChildWindows             = moduser32.NewProc("EnumChildWindows")
	procGetDlgItem                   = moduser32.NewProc("GetDlgItem")
	procGetDpiForWindow              = moduser32.NewProc("GetDpiForWindow")
	procGetMessageW                  = moduser32.NewProc("GetMessageW")
	procGetSystemMetrics             = moduser32.NewProc("GetSystemMetrics")
	procGetWindowDC                  = moduser32.NewProc("GetWindowDC")
	procGetWindowRect                = moduser32.NewProc("GetWindowRect")
	procGetWindowTextLengthW         = moduser32.NewProc("GetWindowTextLengthW")
	procGetWindowTextW               = moduser32.NewProc("GetWindowTextW")
	procIsDialogMessageW             = moduser32.NewProc("IsDialogMessageW")
	procLoadIconW                    = moduser32.NewProc("LoadIconW")
	procPostQuitMessage              = moduser32.NewProc("PostQuitMessage")
	procRegisterClassExW             = moduser32.NewProc("RegisterClassExW")
	procReleaseDC                    = moduser32.NewProc("ReleaseDC")
	procSendMessageW                 = moduser32.NewProc("SendMessageW")
	procSetDlgItemTextW              = moduser32.NewProc("SetDlgItemTextW")
	procSetFocus                     = moduser32.NewProc("SetFocus")
	procSetForegroundWindow          = moduser32.NewProc("SetForegroundWindow")
	procSetThreadDpiAwarenessContext = moduser32.NewProc("SetThreadDpiAwarenessContext")
	procSetWindowLongW               = moduser32.NewProc("SetWindowLongW")
	procSetWindowPos                 = moduser32.NewProc("SetWindowPos")
	procSetWindowTextW               = moduser32.NewProc("SetWindowTextW")
	procSetWindowsHookExW            = moduser32.NewProc("SetWindowsHookExW")
	procShowWindow                   = moduser32.NewProc("ShowWindow")
	procSystemParametersInfoW        = moduser32.NewProc("SystemParametersInfoW")
	procTranslateMessage             = moduser32.NewProc("TranslateMessage")
	procUnhookWindowsHookEx          = moduser32.NewProc("UnhookWindowsHookEx")
	procUnregisterClassW             = moduser32.NewProc("UnregisterClassW")
	procWTSSendMessageW              = modwtsapi32.NewProc("WTSSendMessageW")
)

func InitCommonControlsEx(icc *INITCOMMONCONTROLSEX) (ok bool) {
	r0, _, _ := syscall.Syscall(procInitCommonControlsEx.Addr(), 1, uintptr(unsafe.Pointer(icc)), 0, 0)
	ok = r0 != 0
	return
}

func ChooseColor(cc *CHOOSECOLOR) (ok bool) {
	r0, _, _ := syscall.Syscall(procChooseColorW.Addr(), 1, uintptr(unsafe.Pointer(cc)), 0, 0)
	ok = r0 != 0
	return
}

func commDlgExtendedError() (code int) {
	r0, _, _ := syscall.Syscall(procCommDlgExtendedError.Addr(), 0, 0, 0, 0)
	code = int(r0)
	return
}

func GetOpenFileName(ofn *OPENFILENAME) (ok bool) {
	r0, _, _ := syscall.Syscall(procGetOpenFileNameW.Addr(), 1, uintptr(unsafe.Pointer(ofn)), 0, 0)
	ok = r0 != 0
	return
}

func GetSaveFileName(ofn *OPENFILENAME) (ok bool) {
	r0, _, _ := syscall.Syscall(procGetSaveFileNameW.Addr(), 1, uintptr(unsafe.Pointer(ofn)), 0, 0)
	ok = r0 != 0
	return
}

func CreateFontIndirect(lf *LOGFONT) (ret Handle) {
	r0, _, _ := syscall.Syscall(procCreateFontIndirectW.Addr(), 1, uintptr(unsafe.Pointer(lf)), 0, 0)
	ret = Handle(r0)
	return
}

func DeleteObject(o Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procDeleteObject.Addr(), 1, uintptr(o), 0, 0)
	ok = r0 != 0
	return
}

func GetDeviceCaps(dc Handle, index int) (ret int) {
	r0, _, _ := syscall.Syscall(procGetDeviceCaps.Addr(), 2, uintptr(dc), uintptr(index), 0)
	ret = int(r0)
	return
}

func ActivateActCtx(actCtx Handle, cookie *uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procActivateActCtx.Addr(), 2, uintptr(actCtx), uintptr(unsafe.Pointer(cookie)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func CreateActCtx(actCtx *ACTCTX) (ret Handle, err error) {
	r0, _, e1 := syscall.Syscall(procCreateActCtxW.Addr(), 1, uintptr(unsafe.Pointer(actCtx)), 0, 0)
	ret = Handle(r0)
	if ret == ^Handle(0) {
		err = errnoErr(e1)
	}
	return
}

func DeactivateActCtx(flags uint32, cookie uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procDeactivateActCtx.Addr(), 2, uintptr(flags), uintptr(cookie), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func GenerateConsoleCtrlEvent(ctrlEvent uint32, processGroupId int) (err error) {
	r1, _, e1 := syscall.Syscall(procGenerateConsoleCtrlEvent.Addr(), 2, uintptr(ctrlEvent), uintptr(processGroupId), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetConsoleWindow() (ret HWND) {
	r0, _, _ := syscall.Syscall(procGetConsoleWindow.Addr(), 0, 0, 0, 0)
	ret = HWND(r0)
	return
}

func GetModuleHandle(moduleName *uint16) (ret Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetModuleHandleW.Addr(), 1, uintptr(unsafe.Pointer(moduleName)), 0, 0)
	ret = Handle(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func GlobalAlloc(flags uint32, bytes uintptr) (ret Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGlobalAlloc.Addr(), 2, uintptr(flags), uintptr(bytes), 0)
	ret = Handle(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func GlobalFree(mem Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procGlobalFree.Addr(), 1, uintptr(mem), 0, 0)
	if r1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func ReleaseActCtx(actCtx Handle) {
	syscall.Syscall(procReleaseActCtx.Addr(), 1, uintptr(actCtx), 0, 0)
	return
}

func CoCreateInstance(clsid uintptr, unkOuter *IUnknown, clsContext int32, iid uintptr, address unsafe.Pointer) (res error) {
	r0, _, _ := syscall.Syscall6(procCoCreateInstance.Addr(), 5, uintptr(clsid), uintptr(unsafe.Pointer(unkOuter)), uintptr(clsContext), uintptr(iid), uintptr(address), 0)
	if r0 != 0 {
		res = syscall.Errno(r0)
	}
	return
}

func ExtractAssociatedIcon(instance Handle, path *uint16, icon *uint16) (ret Handle, err error) {
	r0, _, e1 := syscall.Syscall(procExtractAssociatedIconW.Addr(), 3, uintptr(instance), uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(icon)))
	ret = Handle(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func SHBrowseForFolder(bi *BROWSEINFO) (ret *IDLIST) {
	r0, _, _ := syscall.Syscall(procSHBrowseForFolder.Addr(), 1, uintptr(unsafe.Pointer(bi)), 0, 0)
	ret = (*IDLIST)(unsafe.Pointer(r0))
	return
}

func SHCreateItemFromParsingName(path *uint16, bc *IBindCtx, iid uintptr, item **IShellItem) (res error) {
	r0, _, _ := syscall.Syscall6(procSHCreateItemFromParsingName.Addr(), 4, uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(bc)), uintptr(iid), uintptr(unsafe.Pointer(item)), 0, 0)
	if r0 != 0 {
		res = syscall.Errno(r0)
	}
	return
}

func SHGetPathFromIDListEx(ptr *IDLIST, path *uint16, pathLen int, opts int) (ok bool) {
	r0, _, _ := syscall.Syscall6(procSHGetPathFromIDListEx.Addr(), 4, uintptr(unsafe.Pointer(ptr)), uintptr(unsafe.Pointer(path)), uintptr(pathLen), uintptr(opts), 0, 0)
	ok = r0 != 0
	return
}

func ShellNotifyIcon(message uint32, data *NOTIFYICONDATA) (ok bool) {
	r0, _, _ := syscall.Syscall(procShell_NotifyIconW.Addr(), 2, uintptr(message), uintptr(unsafe.Pointer(data)), 0)
	ok = r0 != 0
	return
}

func CallNextHookEx(hk Handle, code int32, wparam uintptr, lparam unsafe.Pointer) (ret uintptr) {
	r0, _, _ := syscall.Syscall6(procCallNextHookEx.Addr(), 4, uintptr(hk), uintptr(code), uintptr(wparam), uintptr(lparam), 0, 0)
	ret = uintptr(r0)
	return
}

func CreateIconFromResourceEx(resBits []byte, icon bool, ver uint32, cx int, cy int, flags int) (ret Handle, err error) {
	var _p0 *byte
	if len(resBits) > 0 {
		_p0 = &resBits[0]
	}
	var _p1 uint32
	if icon {
		_p1 = 1
	}
	r0, _, e1 := syscall.Syscall9(procCreateIconFromResourceEx.Addr(), 7, uintptr(unsafe.Pointer(_p0)), uintptr(len(resBits)), uintptr(_p1), uintptr(ver), uintptr(cx), uintptr(cy), uintptr(flags), 0, 0)
	ret = Handle(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func CreateWindowEx(exStyle uint32, className *uint16, windowName *uint16, style uint32, x int, y int, width int, height int, parent HWND, menu Handle, instance Handle, param unsafe.Pointer) (ret HWND, err error) {
	r0, _, e1 := syscall.Syscall12(procCreateWindowExW.Addr(), 12, uintptr(exStyle), uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)), uintptr(style), uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(parent), uintptr(menu), uintptr(instance), uintptr(param))
	ret = HWND(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func DefWindowProc(wnd HWND, msg uint32, wparam uintptr, lparam unsafe.Pointer) (ret uintptr) {
	r0, _, _ := syscall.Syscall6(procDefWindowProcW.Addr(), 4, uintptr(wnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	ret = uintptr(r0)
	return
}

func DestroyIcon(icon Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procDestroyIcon.Addr(), 1, uintptr(icon), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func DestroyWindow(wnd HWND) (err error) {
	r1, _, e1 := syscall.Syscall(procDestroyWindow.Addr(), 1, uintptr(wnd), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func DispatchMessage(msg *MSG) (ret uintptr) {
	r0, _, _ := syscall.Syscall(procDispatchMessageW.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	ret = uintptr(r0)
	return
}

func EnableWindow(wnd HWND, enable bool) (ok bool) {
	var _p0 uint32
	if enable {
		_p0 = 1
	}
	r0, _, _ := syscall.Syscall(procEnableWindow.Addr(), 2, uintptr(wnd), uintptr(_p0), 0)
	ok = r0 != 0
	return
}

func EnumWindows(enumFunc uintptr, lparam unsafe.Pointer) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumChildWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetDlgItem(dlg HWND, dlgItemID int) (ret HWND, err error) {
	r0, _, e1 := syscall.Syscall(procGetDlgItem.Addr(), 2, uintptr(dlg), uintptr(dlgItemID), 0)
	ret = HWND(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetDpiForWindow(wnd HWND) (ret int, err error) {
	err = procGetDpiForWindow.Find()
	if err != nil {
		return
	}
	r0, _, e1 := syscall.Syscall(procGetDpiForWindow.Addr(), 1, uintptr(wnd), 0, 0)
	ret = int(r0)
	if false {
		err = errnoErr(e1)
	}
	return
}

func GetMessage(msg *MSG, wnd HWND, msgFilterMin uint32, msgFilterMax uint32) (ret uintptr, err error) {
	r0, _, e1 := syscall.Syscall6(procGetMessageW.Addr(), 4, uintptr(unsafe.Pointer(msg)), uintptr(wnd), uintptr(msgFilterMin), uintptr(msgFilterMax), 0, 0)
	ret = uintptr(r0)
	if int32(ret) == -1 {
		err = errnoErr(e1)
	}
	return
}

func GetSystemMetrics(index int) (ret int) {
	r0, _, _ := syscall.Syscall(procGetSystemMetrics.Addr(), 1, uintptr(index), 0, 0)
	ret = int(r0)
	return
}

func GetWindowDC(wnd HWND) (ret Handle) {
	r0, _, _ := syscall.Syscall(procGetWindowDC.Addr(), 1, uintptr(wnd), 0, 0)
	ret = Handle(r0)
	return
}

func GetWindowRect(wnd HWND, cmdShow *RECT) (err error) {
	r1, _, e1 := syscall.Syscall(procGetWindowRect.Addr(), 2, uintptr(wnd), uintptr(unsafe.Pointer(cmdShow)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func getWindowTextLength(wnd HWND) (ret int, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextLengthW.Addr(), 1, uintptr(wnd), 0, 0)
	ret = int(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func getWindowText(wnd HWND, str *uint16, maxCount int) (ret int, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(wnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	ret = int(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func IsDialogMessage(wnd HWND, msg *MSG) (ok bool) {
	r0, _, _ := syscall.Syscall(procIsDialogMessageW.Addr(), 2, uintptr(wnd), uintptr(unsafe.Pointer(msg)), 0)
	ok = r0 != 0
	return
}

func LoadIcon(instance Handle, resource uintptr) (ret Handle, err error) {
	r0, _, e1 := syscall.Syscall(procLoadIconW.Addr(), 2, uintptr(instance), uintptr(resource), 0)
	ret = Handle(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func PostQuitMessage(exitCode int) {
	syscall.Syscall(procPostQuitMessage.Addr(), 1, uintptr(exitCode), 0, 0)
	return
}

func RegisterClassEx(cls *WNDCLASSEX) (err error) {
	r1, _, e1 := syscall.Syscall(procRegisterClassExW.Addr(), 1, uintptr(unsafe.Pointer(cls)), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func ReleaseDC(wnd HWND, dc Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procReleaseDC.Addr(), 2, uintptr(wnd), uintptr(dc), 0)
	ok = r0 != 0
	return
}

func SendMessage(wnd HWND, msg uint32, wparam uintptr, lparam uintptr) (ret uintptr) {
	r0, _, _ := syscall.Syscall6(procSendMessageW.Addr(), 4, uintptr(wnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	ret = uintptr(r0)
	return
}

func SetDlgItemText(dlg HWND, dlgItemID int, str *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procSetDlgItemTextW.Addr(), 3, uintptr(dlg), uintptr(dlgItemID), uintptr(unsafe.Pointer(str)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetFocus(wnd HWND) (ret HWND, err error) {
	r0, _, e1 := syscall.Syscall(procSetFocus.Addr(), 1, uintptr(wnd), 0, 0)
	ret = HWND(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetForegroundWindow(wnd HWND) (ok bool) {
	r0, _, _ := syscall.Syscall(procSetForegroundWindow.Addr(), 1, uintptr(wnd), 0, 0)
	ok = r0 != 0
	return
}

func SetThreadDpiAwarenessContext(dpiContext uintptr) (ret uintptr, err error) {
	err = procSetThreadDpiAwarenessContext.Find()
	if err != nil {
		return
	}
	r0, _, e1 := syscall.Syscall(procSetThreadDpiAwarenessContext.Addr(), 1, uintptr(dpiContext), 0, 0)
	ret = uintptr(r0)
	if false {
		err = errnoErr(e1)
	}
	return
}

func SetWindowLong(wnd HWND, index int, newLong int) (ret int, err error) {
	r0, _, e1 := syscall.Syscall(procSetWindowLongW.Addr(), 3, uintptr(wnd), uintptr(index), uintptr(newLong))
	ret = int(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetWindowPos(wnd HWND, wndInsertAfter HWND, x int, y int, cx int, cy int, flags int) (err error) {
	r1, _, e1 := syscall.Syscall9(procSetWindowPos.Addr(), 7, uintptr(wnd), uintptr(wndInsertAfter), uintptr(x), uintptr(y), uintptr(cx), uintptr(cy), uintptr(flags), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetWindowText(wnd HWND, text *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procSetWindowTextW.Addr(), 2, uintptr(wnd), uintptr(unsafe.Pointer(text)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetWindowsHookEx(idHook int, fn uintptr, mod Handle, threadID uint32) (ret Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procSetWindowsHookExW.Addr(), 4, uintptr(idHook), uintptr(fn), uintptr(mod), uintptr(threadID), 0, 0)
	ret = Handle(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func ShowWindow(wnd HWND, cmdShow int) (ok bool) {
	r0, _, _ := syscall.Syscall(procShowWindow.Addr(), 2, uintptr(wnd), uintptr(cmdShow), 0)
	ok = r0 != 0
	return
}

func SystemParametersInfo(action int, uiParam uintptr, pvParam unsafe.Pointer, winIni int) (err error) {
	r1, _, e1 := syscall.Syscall6(procSystemParametersInfoW.Addr(), 4, uintptr(action), uintptr(uiParam), uintptr(pvParam), uintptr(winIni), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func TranslateMessage(msg *MSG) (ok bool) {
	r0, _, _ := syscall.Syscall(procTranslateMessage.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	ok = r0 != 0
	return
}

func UnhookWindowsHookEx(hk Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procUnhookWindowsHookEx.Addr(), 1, uintptr(hk), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func UnregisterClass(className *uint16, instance Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procUnregisterClassW.Addr(), 2, uintptr(unsafe.Pointer(className)), uintptr(instance), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func WTSSendMessage(server Handle, sessionID uint32, title *uint16, titleLength int, message *uint16, messageLength int, style uint32, timeout int, response *uint32, wait bool) (err error) {
	var _p0 uint32
	if wait {
		_p0 = 1
	}
	r1, _, e1 := syscall.Syscall12(procWTSSendMessageW.Addr(), 10, uintptr(server), uintptr(sessionID), uintptr(unsafe.Pointer(title)), uintptr(titleLength), uintptr(unsafe.Pointer(message)), uintptr(messageLength), uintptr(style), uintptr(timeout), uintptr(unsafe.Pointer(response)), uintptr(_p0), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
