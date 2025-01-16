/*
	https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications

// https://medium.com/@hectorjsmith/inject-version-info-in-go-binaries-b4f5d405b1b6
// https://dev.to/gcdcoder/versioning-binaries-in-go-52al
// https://gist.github.com/Harold2017/2547c7f65fcdc93d2a6fe2f81319d301
// / go build -v -ldflags="-X 'main.Version=v1.0.0' -X 'lib/build.User=$(id -u -n)' -X 'lib/build.Time=$(date)'"
*/

package build

// Version version information stored during the build process
var Version string

// Time built time stored during the build process
var Time string

// User built user name stored during the build process
var User string

// BuildGoVersion go version stored during the build process
var BuildGoVersion string
