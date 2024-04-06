package verflag

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/tkane/tkblog/pkg/version"
)

type versionValue int

const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const (
	strRawVersion = "raw"
	versionFlagName = "version"
)

// Set implements pflag.Value.
func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVar, err := strconv.ParseBool(s)
	if boolVar {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

// String implements pflag.Value.
func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// Type implements pflag.Value.
func (v *versionValue) Type() string {
	return "version"
}

func (v *versionValue) Get() interface{} {
	return v
}

func (v *versionValue) IsBoolFlag() bool {
	return true
}

var versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit.")


func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)
	pflag.Lookup(name).NoOptDefVal = "true"
}

func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)

	return p
}

func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(versionFlagName))
}


// %#v 是一个格式化动词，它会打印出值的默认格式，同时包括值的类型信息。
func PrintAndExitIfReq() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}