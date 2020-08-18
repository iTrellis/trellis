package internal

import "fmt"

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
