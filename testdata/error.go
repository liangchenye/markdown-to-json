// Package specerror implements runtime-spec-specific tooling for
// tracking RFC 2119 violations.
package specerror

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	rfc2119 "github.com/opencontainers/runtime-tools/error"
)

const referenceTemplate = "https://github.com/opencontainers/runtime-spec/blob/v%s/%s"

// Code represents the spec violation, enumerating both
// configuration violations and runtime violations.
type Code int

const (
	// NonError represents that an input is not an error
	NonError Code = iota
	// NonRFCError represents that an error is not a rfc2119 error
	NonRFCError

	// ConfigFileExistence represents the error code of 'config.json' existence test
	ConfigFileExistence
	// ArtifactsInSingleDir represents the error code of artifacts place test
	ArtifactsInSingleDir

	// SpecVersionInSemVer is removed for test

	// RootOnWindowsRequired represents "On Windows, for Windows Server Containers, this field is REQUIRED."
	RootOnWindowsRequired
	// RootOnHyperVNotSet represents "For [Hyper-V Containers](config-windows.md#hyperv), this field MUST NOT be set."
	RootOnHyperVNotSet
	// RootOnNonHyperVRequired represents "On all other platforms, this field is REQUIRED."
	RootOnNonHyperVRequired
	// RootPathOnWindowsGUID represents "* On Windows, `path` MUST be a [volume GUID path][naming-a-volume]."
	RootPathOnWindowsGUID
	// RootPathOnPosixConvention represents "The value SHOULD be the conventional `rootfs`."
	RootPathOnPosixConvention
	// RootPathExist represents "A directory MUST exist at the path declared by the field."
	RootPathExist
	// RootReadonlyImplement represents "* **`readonly`** (bool, OPTIONAL) If true then the root filesystem MUST be read-only inside the container, defaults to false."
	RootReadonlyImplement
	// RootReadonlyOnWindowsFalse represents "* On Windows, this field MUST be omitted or false."
	RootReadonlyOnWindowsFalse
	// MountsInOrder represents "The runtime MUST mount entries in the listed order."
	MountsInOrder
	// MountsDestAbs represents "This value MUST be an absolute path."
	MountsDestAbs
	// MountsDestOnWindowsNotNested represents "* Windows: one mount destination MUST NOT be nested within another mount (e.g., c:\\foo and c:\\foo\\bar)."
	MountsDestOnWindowsNotNested
	// MountsOptionsOnWindowsROSupport represents "* Windows: runtimes MUST support `ro`, mounting the filesystem read-only when `ro` is given."
	MountsOptionsOnWindowsROSupport
	// ProcRequiredAtStart represents "This property is REQUIRED when [`start`](runtime.md#start) is called."
	ProcRequiredAtStart
	// ProcConsoleSizeIgnore represents "Runtimes MUST ignore `consoleSize` if `terminal` is `false` or unset."
	ProcConsoleSizeIgnore
	// ProcCwdAbs represents "This value MUST be an absolute path."
	ProcCwdAbs
	// ProcArgsOneEntryRequired represents "This specification extends the IEEE standard in that at least one entry is REQUIRED, and that entry is used with the same semantics as `execvp`'s *file*."
	ProcArgsOneEntryRequired
	// PosixProcRlimitsTypeError represents "The runtime MUST [generate an error](runtime.md#errors) for any values which cannot be mapped to a relevant kernel interface"
	PosixProcRlimitsTypeError
	// PosixProcRlimitsTypeGet represents "For each entry in `rlimits`, a [`getrlimit(3)`][getrlimit.3] on `type` MUST succeed."
	PosixProcRlimitsTypeGet
	// PosixProcRlimitsSoftMatchCur represents "`rlim.rlim_cur` MUST match the configured value."
	PosixProcRlimitsSoftMatchCur
	// PosixProcRlimitsHardMatchMax represents "`rlim.rlim_max` MUST match the configured value."
	PosixProcRlimitsHardMatchMax
	// PosixProcRlimitsErrorOnDup represents "If `rlimits` contains duplicated entries with same `type`, the runtime MUST [generate an error](runtime.md#errors)."
	PosixProcRlimitsErrorOnDup
	// LinuxProcCapError represents "Any value which cannot be mapped to a relevant kernel interface MUST cause an error."
	LinuxProcCapError
	// LinuxProcOomScoreAdjSet represents "If `oomScoreAdj` is set, the runtime MUST set `oom_score_adj` to the given value."
	LinuxProcOomScoreAdjSet
	// LinuxProcOomScoreAdjNotSet represents "If `oomScoreAdj` is not set, the runtime MUST NOT change the value of `oom_score_adj`."
	LinuxProcOomScoreAdjNotSet
	// PlatformSpecConfOnWindowsSet represents "This MUST be set if the target platform of this spec is `windows`."
	PlatformSpecConfOnWindowsSet
	// PosixHooksPathAbs represents "This specification extends the IEEE standard in that **`path`** MUST be absolute."
	PosixHooksPathAbs
	// PosixHooksTimeoutPositive represents "If set, `timeout` MUST be greater than zero."
	PosixHooksTimeoutPositive
	// PosixHooksCalledInOrder represents "Hooks MUST be called in the listed order."
	PosixHooksCalledInOrder
	// PosixHooksStateToStdin represents "The [state](runtime.md#state) of the container MUST be passed to hooks over stdin so that they may do work appropriate to the current state of the container."
	PosixHooksStateToStdin
	// PrestartTiming represents "The pre-start hooks MUST be called after the [`start`](runtime.md#start) operation is called but [before the user-specified program command is executed](runtime.md#lifecycle)."
	PrestartTiming
	// PoststartTiming represents "The post-start hooks MUST be called [after the user-specified process is executed](runtime.md#lifecycle) but before the [`start`](runtime.md#start) operation returns."
	PoststartTiming
	// PoststopTiming represents "The post-stop hooks MUST be called [after the container is deleted](runtime.md#lifecycle) but before the [`delete`](runtime.md#delete) operation returns."
	PoststopTiming
	// AnnotationsKeyValueMap represents "Annotations MUST be a key-value map."
	AnnotationsKeyValueMap
	// AnnotationsKeyString represents "Keys MUST be strings."
	AnnotationsKeyString
	// AnnotationsKeyRequired represents "Keys MUST NOT be an empty string."
	AnnotationsKeyRequired
	// AnnotationsKeyReversedDomain represents "Keys SHOULD be named using a reverse domain notation - e.g. `com.example.myKey`."
	AnnotationsKeyReversedDomain
	// AnnotationsKeyReservedNS represents "Keys using the `org.opencontainers` namespace are reserved and MUST NOT be used by subsequent specifications."
	AnnotationsKeyReservedNS
	// AnnotationsKeyIgnoreUnknown represents "Implementations that are reading/processing this configuration file MUST NOT generate an error if they encounter an unknown annotation key."
	AnnotationsKeyIgnoreUnknown
	// AnnotationsValueString represents "Values MUST be strings."
	AnnotationsValueString
	// ExtensibilityIgnoreUnknownProp represents "Runtimes that are reading or processing this configuration file MUST NOT generate an error if they encounter an unknown property."
	ExtensibilityIgnoreUnknownProp
	// ValidValuesError represents "Runtimes that are reading or processing this configuration file MUST generate an error when invalid or unsupported values are encountered."
	ValidValuesError

	// DefaultFilesystems represents the error code of default filesystems test
	DefaultFilesystems

	// CreateWithID represents the error code of 'create' lifecyle test with 'id' provided
	CreateWithID
	// CreateWithUniqueID represents the error code of 'create' lifecyle test with unique 'id' provided
	CreateWithUniqueID
	// CreateNewContainer represents the error code 'create' lifecyle test that creates new container
	CreateNewContainer
)

var (
	containerFormatRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "bundle.md#container-format"), nil
	}

	// specificationVersionRef is removed for test

	rootRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#root"), nil
	}
	mountsRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#mounts"), nil
	}
	processRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#process"), nil
	}
	posixProcessRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#posix-process"), nil
	}
	linuxProcessRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#linux-process"), nil
	}
	platformSpecificConfigurationRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#platform-specific-configuration"), nil
	}
	posixPlatformHooksRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#posix-platform-hooks"), nil
	}
	prestartRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#prestart"), nil
	}
	poststartRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#poststart"), nil
	}
	poststopRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#poststop"), nil
	}
	annotationsRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#annotations"), nil
	}
	extensibilityRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#extensibility"), nil
	}
	validValuesRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config.md#valid-values"), nil
	}

	defaultFSRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "config-linux.md#default-filesystems"), nil
	}

	runtimeCreateRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "runtime.md#create"), nil
	}
)

type errorTemplate struct {
	Level     rfc2119.Level
	Reference func(version string) (reference string, err error)
}

// Error represents a runtime-spec violation.
type Error struct {
	// Err holds the RFC 2119 violation.
	Err rfc2119.Error

	// Code is a matchable holds a Code
	Code Code
}

var ociErrors = map[Code]errorTemplate{
	// Bundle.md
	// Container Format
	ConfigFileExistence:  {Level: rfc2119.Must, Reference: containerFormatRef},
	ArtifactsInSingleDir: {Level: rfc2119.Must, Reference: containerFormatRef},

	// config.md
	// Specification version
	SpecVersionInSemVer: {Level: rfc2119.Must, Reference: specificationVersionRef},
	// Root
	RootOnWindowsRequired:      {Level: rfc2119.Required, Reference: rootRef},
	RootOnHyperVNotSet:         {Level: rfc2119.Must, Reference: rootRef},
	RootOnNonHyperVRequired:    {Level: rfc2119.Required, Reference: rootRef},
	RootPathOnWindowsGUID:      {Level: rfc2119.Must, Reference: rootRef},
	RootPathOnPosixConvention:  {Level: rfc2119.Should, Reference: rootRef},
	RootPathExist:              {Level: rfc2119.Must, Reference: rootRef},
	RootReadonlyImplement:      {Level: rfc2119.Must, Reference: rootRef},
	RootReadonlyOnWindowsFalse: {Level: rfc2119.Must, Reference: rootRef},
	// Mounts
	MountsInOrder:                   {Level: rfc2119.Must, Reference: mountsRef},
	MountsDestAbs:                   {Level: rfc2119.Must, Reference: mountsRef},
	MountsDestOnWindowsNotNested:    {Level: rfc2119.Must, Reference: mountsRef},
	MountsOptionsOnWindowsROSupport: {Level: rfc2119.Must, Reference: mountsRef},
	// Proc
	ProcRequiredAtStart:      {Level: rfc2119.Required, Reference: processRef},
	ProcConsoleSizeIgnore:    {Level: rfc2119.Must, Reference: processRef},
	ProcCwdAbs:               {Level: rfc2119.Must, Reference: processRef},
	ProcArgsOneEntryRequired: {Level: rfc2119.Required, Reference: processRef},
	// POSIX process
	PosixProcRlimitsTypeError:    {Level: rfc2119.Must, Reference: posixProcessRef},
	PosixProcRlimitsTypeGet:      {Level: rfc2119.Must, Reference: posixProcessRef},
	PosixProcRlimitsSoftMatchCur: {Level: rfc2119.Must, Reference: posixProcessRef},
	PosixProcRlimitsHardMatchMax: {Level: rfc2119.Must, Reference: posixProcessRef},
	PosixProcRlimitsErrorOnDup:   {Level: rfc2119.Must, Reference: posixProcessRef},
	// Linux Proc
	LinuxProcCapError:          {Level: rfc2119.Must, Reference: linuxProcessRef},
	LinuxProcOomScoreAdjSet:    {Level: rfc2119.Must, Reference: linuxProcessRef},
	LinuxProcOomScoreAdjNotSet: {Level: rfc2119.Must, Reference: linuxProcessRef},
	// Platform-specific configuration
	PlatformSpecConfOnWindowsSet: {Level: rfc2119.Must, Reference: platformSpecificConfigurationRef},
	// POSIX-platform Hooks
	PosixHooksPathAbs:         {Level: rfc2119.Must, Reference: posixPlatformHooksRef},
	PosixHooksTimeoutPositive: {Level: rfc2119.Must, Reference: posixPlatformHooksRef},
	PosixHooksCalledInOrder:   {Level: rfc2119.Must, Reference: posixPlatformHooksRef},
	PosixHooksStateToStdin:    {Level: rfc2119.Must, Reference: posixPlatformHooksRef},
	// Prestart
	PrestartTiming: {Level: rfc2119.Must, Reference: prestartRef},
	// Poststart
	PoststartTiming: {Level: rfc2119.Must, Reference: poststartRef},
	// Poststop
	PoststopTiming: {Level: rfc2119.Must, Reference: poststopRef},
	// Annotations
	AnnotationsKeyValueMap:       {Level: rfc2119.Must, Reference: annotationsRef},
	AnnotationsKeyString:         {Level: rfc2119.Must, Reference: annotationsRef},
	AnnotationsKeyRequired:       {Level: rfc2119.Must, Reference: annotationsRef},
	AnnotationsKeyReversedDomain: {Level: rfc2119.Should, Reference: annotationsRef},
	AnnotationsKeyReservedNS:     {Level: rfc2119.Must, Reference: annotationsRef},
	AnnotationsKeyIgnoreUnknown:  {Level: rfc2119.Must, Reference: annotationsRef},
	AnnotationsValueString:       {Level: rfc2119.Must, Reference: annotationsRef},
	// Extensibility
	ExtensibilityIgnoreUnknownProp: {Level: rfc2119.Must, Reference: extensibilityRef},
	// Valid values
	ValidValuesError: {Level: rfc2119.Must, Reference: validValuesRef},

	// Config-Linux.md
	// Default Filesystems
	DefaultFilesystems: {Level: rfc2119.Should, Reference: defaultFSRef},

	// Runtime.md
	// Create
	CreateWithID:       {Level: rfc2119.Must, Reference: runtimeCreateRef},
	CreateWithUniqueID: {Level: rfc2119.Must, Reference: runtimeCreateRef},
	CreateNewContainer: {Level: rfc2119.Must, Reference: runtimeCreateRef},
}

// Error returns the error message with specification reference.
func (err *Error) Error() string {
	return err.Err.Error()
}

// NewError creates an Error referencing a spec violation.  The error
// can be cast to an *Error for extracting structured information
// about the level of the violation and a reference to the violated
// spec condition.
//
// A version string (for the version of the spec that was violated)
// must be set to get a working URL.
func NewError(code Code, err error, version string) error {
	template := ociErrors[code]
	reference, err2 := template.Reference(version)
	if err2 != nil {
		return err2
	}
	return &Error{
		Err: rfc2119.Error{
			Level:     template.Level,
			Reference: reference,
			Err:       err,
		},
		Code: code,
	}
}

// FindError finds an error from a source error (multiple error) and
// returns the error code if found.
// If the source error is nil or empty, return NonError.
// If the source error is not a multiple error, return NonRFCError.
func FindError(err error, code Code) Code {
	if err == nil {
		return NonError
	}

	if merr, ok := err.(*multierror.Error); ok {
		if merr.ErrorOrNil() == nil {
			return NonError
		}
		for _, e := range merr.Errors {
			if rfcErr, ok := e.(*Error); ok {
				if rfcErr.Code == code {
					return code
				}
			}
		}
	}
	return NonRFCError
}
