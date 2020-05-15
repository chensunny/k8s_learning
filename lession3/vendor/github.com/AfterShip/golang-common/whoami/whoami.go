package whoami

// 值的来源为go build -ldflags "-X github.com/AfterShip/golang-common/whoami.name -X ..."

// Example code in bulid-bin.sh
/*
	nameFlag="-X github.com/AfterShip/golang-common/whoami.name=YouServerName" \
	&& versionFlag="-X github.com/AfterShip/golang-common/whoami.version="`cat version`"" \
	&& buildNumberFlag="-X github.com/AfterShip/golang-common/whoami.buildNumber="${buildNumber}"" \
	&& buildAtFlag="-X github.com/AfterShip/golang-common/whoami.buildAt="`date '+%Y-%m-%dT%H:%M:%S%z'`"" \
	&& commitHashFlag="-X github.com/AfterShip/golang-common/whoami.commitHash="`git rev-parse HEAD`\
	&& commitBranchFlag="-X github.com/AfterShip/golang-common/whoami.commitBranch="`git name-rev --name-only HEAD`\
	&& ldflags="${nameFlag} ${versionFlag} ${buildNumberFlag} ${buildAtFlag} ${commitHashFlag} ${commitBranchFlag}" \
	&& go build  -ldflags "${ldflags}" -mod=vendor ./cmd/apiserver
 */
var name, version, buildNumber, buildAt, commitHash, commitBranch string

func Name() string {
	return name
}

func Version() string {
	return version
}

func BuildNumber() string {
	return buildNumber
}

func BuildAt() string {
	return buildAt
}

func CommitHash() string {
	return commitHash
}

func CommitBranch() string{
	return commitBranch
}
