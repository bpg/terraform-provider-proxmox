package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// CPUArchitectureValidator returns a schema validation function for a CPU architecture.
func CPUArchitectureValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"aarch64",
		"x86_64",
	}, false))
}

// CPUTypeValidator returns a schema validation function for a CPU type.
func CPUTypeValidator() schema.SchemaValidateDiagFunc {
	standardTypes := []string{
		"486",
		"Broadwell",
		"Broadwell-IBRS",
		"Broadwell-noTSX",
		"Broadwell-noTSX-IBRS",
		"Cascadelake-Server",
		"Cascadelake-Server-noTSX",
		"Cascadelake-Server-v2",
		"Cascadelake-Server-v4",
		"Cascadelake-Server-v5",
		"Conroe",
		"Cooperlake",
		"Cooperlake-v2",
		"EPYC",
		"EPYC-IBPB",
		"EPYC-Milan",
		"EPYC-Rome",
		"EPYC-Rome-v2",
		"EPYC-v3",
		"Haswell",
		"Haswell-IBRS",
		"Haswell-noTSX",
		"Haswell-noTSX-IBRS",
		"Icelake-Client",
		"Icelake-Client-noTSX",
		"Icelake-Server",
		"Icelake-Server-noTSX",
		"Icelake-Server-v3",
		"Icelake-Server-v4",
		"Icelake-Server-v5",
		"Icelake-Server-v6",
		"IvyBridge",
		"IvyBridge-IBRS",
		"KnightsMill",
		"Nehalem",
		"Nehalem-IBRS",
		"Opteron_G1",
		"Opteron_G2",
		"Opteron_G3",
		"Opteron_G4",
		"Opteron_G5",
		"Penryn",
		"SandyBridge",
		"SandyBridge-IBRS",
		"SapphireRapids",
		"Skylake-Client",
		"Skylake-Client-IBRS",
		"Skylake-Client-noTSX-IBRS",
		"Skylake-Client-v4",
		"Skylake-Server",
		"Skylake-Server-IBRS",
		"Skylake-Server-noTSX-IBRS",
		"Skylake-Server-v4",
		"Skylake-Server-v5",
		"Westmere",
		"Westmere-IBRS",
		"athlon",
		"core2duo",
		"coreduo",
		"host",
		"kvm32",
		"kvm64",
		"max",
		"pentium",
		"pentium2",
		"pentium3",
		"phenom",
		"qemu32",
		"qemu64",
		"x86-64-v2",
		"x86-64-v2-AES",
		"x86-64-v3",
		"x86-64-v4",
	}

	return validation.ToDiagFunc(validation.Any(
		validation.StringInSlice(standardTypes, false),
		validation.StringMatch(regexp.MustCompile(`^custom-.+$`), "must be a valid custom CPU type"),
	))
}

// CPUAffinityValidator returns a schema validation function for a CPU affinity.
func CPUAffinityValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(
		validation.StringMatch(regexp.MustCompile(`^\d+[\d-,]*$`), "must contain numbers or number ranges separated by ','"),
	)
}
