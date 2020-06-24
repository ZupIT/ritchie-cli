package resource

const (

	// Version contains the current version	.
	Version = "dev"

	// BuildDate contains a string with the build date.
	BuildDate = "unknown"

	// SkipTlsVerify is a flag to skip tls verification when pinning server
	SkipTlsVerify = false

	// Url to get Rit Stable Version
	StableVersionUrl = "https://commons-repo.ritchiecli.io/stable.txt"

	// Url format to upgrade rit
	UpgradeUrlFormat = "https://commons-repo.ritchiecli.io/%s/%s/%s/rit"
)
