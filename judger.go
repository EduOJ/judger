package judger

/*
#cgo pkg-config: libseccomp
#include "seccomp_rules.h"
#include "stdlib.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	ArgsMaxNumber = 256
	EnvMaxNumber  = 256
)

// Config is struct used to record the running configuration.
type Config struct {
	MaxCPUTime           int      // max cpu time(ms) this process can cost, -1 for unlimited
	MaxRealTime          int      // max time(ms) this process can run, -1 for unlimited
	MaxMemory            int      // max size(byte) of the process' virtual memory (address space), -1 for unlimited
	MaxStack             int      // max size(byte) of the process' stack size
	MaxProcessNumber     int      // max number of processes that can be created for the real user id of the calling process, -1 for unlimited
	MaxOutputSize        int      // max size(byte) of data this process can output to stdout, stderr and file, -1 for unlimited
	MemoryLimitCheckOnly int      // if this value equals 0, we will only check memory usage number, because setrlimit(maxrss) will cause some crash issues
	ExePath              string   // path of file to run
	InputPath            string   // redirect content of this file to process's stdin
	OutputPath           string   // redirect process's stdout to this file
	ErrorPath            string   // redirect process's stderr to this file
	Args                 []string // arguments to run this process
	Env                  []string // environment variables this process can get
	LogPath              string   // judger log path
	SeccompRuleName      string   // seccomp rules used to limit process system calls.
	// Name is used to call corresponding functions.
	// Possible values are: c_cpp, c_cpp_file_io, general
	Uid uint32 // user to run this process
	Gid uint32 // user group this process belongs to
}

type JudgeResult int

const (
	SUCCESS JudgeResult = iota
	CPU_TIME_LIMIT_EXCEEDED
	REAL_TIME_LIMIT_EXCEEDED
	MEMORY_LIMIT_EXCEEDED
	RUNTIME_ERROR
	SYSTEM_ERROR
)

// Result is a struct used to record the running result.
type Result struct {
	CPUTime  int         // cpu time the process has used
	RealTime int         // actual running time of the process
	Memory   int         // max value of memory used by the process
	Signal   int         // signal number
	ExitCode int         // process's exit code
	Result   JudgeResult // Judge result
}

type resultError string

const (
	ErrInvalidConfig     resultError = "invalid config"
	ErrForkFailed        resultError = "fork failed"
	ErrPthreadFailed     resultError = "pthread failed"
	ErrWaitFailed        resultError = "wait failed"
	ErrRootRequired      resultError = "root required"
	ErrLoadSeccompFailed resultError = "load seccomp failed"
	ErrSetRLimitFailed   resultError = "setrlimit failed"
	ErrDup2Failed        resultError = "dup2 failed"
	ErrSetuidFailed      resultError = "setuid failed"
	ErrExecveFailed      resultError = "execve failed"
	ErrSetOOMFailed      resultError = "set OOM failed"
)

func (r resultError) Error() string {
	return string(r)
}

func (c Config) convertToCStruct() (cc C.struct_config) {
	cc.max_cpu_time = C.int(c.MaxCPUTime)
	cc.max_real_time = C.int(c.MaxRealTime)
	cc.max_memory = C.long(c.MaxMemory)
	cc.max_stack = C.long(c.MaxStack)
	cc.max_process_number = C.int(c.MaxProcessNumber)
	cc.max_output_size = C.long(c.MaxOutputSize)
	cc.memory_limit_check_only = C.int(c.MemoryLimitCheckOnly)
	cc.exe_path = C.CString(c.ExePath)
	cc.input_path = C.CString(c.InputPath)
	cc.output_path = C.CString(c.OutputPath)
	cc.error_path = C.CString(c.ErrorPath)
	for i := 0; i < len(c.Args) && i < ArgsMaxNumber-1; i++ {
		cc.args[i] = C.CString(c.Args[i])
	}
	cc.args[len(c.Args)] = nil
	for i := 0; i < len(c.Env) && i < EnvMaxNumber-1; i++ {
		cc.env[i] = C.CString(c.Env[i])
	}
	cc.env[len(c.Env)] = nil
	cc.log_path = C.CString(c.LogPath)
	cc.seccomp_rule_name = C.CString(c.SeccompRuleName)
	cc.uid = C.uint(c.Uid)
	cc.gid = C.uint(c.Gid)
	return
}

func (r *Result) convertFromCStruct(cr C.struct_result) {
	r.CPUTime = int(cr.cpu_time)
	r.RealTime = int(cr.real_time)
	r.Memory = int(cr.memory)
	r.Signal = int(cr.signal)
	r.ExitCode = int(cr.exit_code)
	r.Result = JudgeResult(cr.result)
}

// Run runs the program in the sandbox according to the config and returns the result.
func Run(config Config) (result Result, err error) {
	var cResult C.struct_result
	cConfig := config.convertToCStruct()
	C.run(&cConfig, &cResult)
	if cResult.error != 0 {
		switch int(cResult.error) {
		case -1:
			err = ErrInvalidConfig
		case -2:
			err = ErrForkFailed
		case -3:
			err = ErrPthreadFailed
		case -4:
			err = ErrWaitFailed
		case -5:
			err = ErrRootRequired
		case -6:
			err = ErrLoadSeccompFailed
		case -7:
			err = ErrSetRLimitFailed
		case -8:
			err = ErrDup2Failed
		case -9:
			err = ErrSetuidFailed
		case -10:
			err = ErrExecveFailed
		case -12:
			err = ErrSetOOMFailed
		default:
			err = errors.New("unknown error")
		}
		return
	}
	result.convertFromCStruct(cResult)
	C.free(unsafe.Pointer(cConfig.exe_path))
	C.free(unsafe.Pointer(cConfig.input_path))
	C.free(unsafe.Pointer(cConfig.output_path))
	C.free(unsafe.Pointer(cConfig.error_path))
	C.free(unsafe.Pointer(cConfig.log_path))
	C.free(unsafe.Pointer(cConfig.seccomp_rule_name))
	for i := range cConfig.args {
		if i == 0 {
			break
		}
		C.free(unsafe.Pointer(uintptr(i)))
	}
	for i := range cConfig.env {
		if i == 0 {
			break
		}
		C.free(unsafe.Pointer(uintptr(i)))
	}
	return
}
