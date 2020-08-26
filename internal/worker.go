package internal

import (
	"fmt"
	"strings"
)

// schema
const (
	SchemaTrellis    = "trellis://"
	SchemaETCDNaming = "/trellis/etcdnaming/"
)

// WorkerPath 工作路径
func WorkerPath(schema, name, version string) string {
	return fmt.Sprintf("%s%s/%s", schema, name, version)
}

// WorkerDomainPath 工作路径
func WorkerDomainPath(schema, name, version, domain string) string {
	return fmt.Sprintf("%s/%s/%s/%s", schema, name, version, domain)
}

// WorkerTrellisPath trellis工作路径
func WorkerTrellisPath(name, version string) string {
	return fmt.Sprintf("%s%s/%s", SchemaTrellis, name, version)
}

// WorkerETCDPath etcd工作路径
func WorkerETCDPath(name, version string) string {
	return fmt.Sprintf("%s%s/%s", SchemaETCDNaming, name, version)
}

// WorkerTrellisDomainPath 工作路径
func WorkerTrellisDomainPath(name, version, domain string) string {
	return fmt.Sprintf("%s/%s/%s/%s", SchemaTrellis, name, version, domain)
}

// WorkerETCDDomainPath 工作路径
func WorkerETCDDomainPath(name, version, domain string) string {
	domain = strings.Replace(domain, ".", "_", -1)
	domain = strings.Replace(domain, ":", "_", -1)
	return fmt.Sprintf("%s%s/%s/%s", SchemaETCDNaming, name, version, domain)
}
