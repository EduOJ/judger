package judger

/*
#cgo CFLAGS: -I../sandbox
#cgo LDFLAGS: -L${SRCDIR}/../sandbox -ljudger
#include "../sandbox/src/runner.h"
*/
import "C"

const (
	ArgsMaxNumber = 256
	EnvMaxNumber  = 256
)

type Config struct {
	MaxCPUTime           int
	MaxRealTine          int
	MaxMemory            int32
	MaxStack             int32
	MaxProcessNumber     int
	MaxOutputSize        int32
	MemoryLimitCheckOnly int
	ExePath              string
	InputPath            string
	OutputPath           string
	ErrorPath            string
	Args                 []string
	Env                  []string
	LongPath             string
	SeccompRuleName      string
	Uid                  uint32
	Gid                  uint32
}

type Result struct {
	CPUTime  int
	RealTime int
	Memory   int32
	Signal   int
	ExitCode int
	Error    int
	Result   int
}

func (c Config) ConvertToCStruct() (cc C.struct_config) {
	cc.max_cpu_time = C.int(c.MaxCPUTime)
	cc.max_real_time = C.int(c.MaxRealTine)
	cc.max_memory = C.long(c.MaxMemory)
	cc.max_stack = C.long(c.MaxStack)
	cc.max_process_number = C.int(c.MaxProcessNumber)
	cc.max_output_size = C.long(c.MaxOutputSize)
	cc.memory_limit_check_only = C.int(c.MemoryLimitCheckOnly)
	cc.exe_path = C.CString(c.ExePath)
	cc.input_path = C.CString(c.InputPath)
	cc.output_path = C.CString(c.OutputPath)
	cc.error_path = C.CString(c.ErrorPath)
	for i := 0; i < len(c.Args) && i < ArgsMaxNumber; i++ {
		cc.args[i] = C.CString(c.Args[i])
	}
	for i := 0; i < len(c.Env) && i < EnvMaxNumber; i++ {
		cc.env[i] = C.CString(c.Env[i])
	}
	cc.log_path = C.CString(c.LongPath)
	cc.seccomp_rule_name = C.CString(c.SeccompRuleName)
	cc.uid = C.uint(c.Uid)
	cc.gid = C.uint(c.Gid)
	return
}

func (r *Result) ConvertFromCStruct(cr C.struct_result) {
	r.CPUTime = int(cr.cpu_time)
	r.RealTime = int(cr.real_time)
	r.Memory = int32(cr.memory)
	r.Signal = int(cr.signal)
	r.ExitCode = int(cr.exit_code)
	r.Error = int(cr.error)
	r.Result = int(cr.result)
}

func Run(config Config) (result Result) {
	var cResult C.struct_result
	cConfig := config.ConvertToCStruct()
	C.run(&cConfig, &cResult)
	result.ConvertFromCStruct(cResult)
	return
}
