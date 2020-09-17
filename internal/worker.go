/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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
	domain = strings.Replace(domain, ".", "_", -1)
	domain = strings.Replace(domain, ":", "_", -1)
	domain = strings.Replace(domain, "/", "_", -1)
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
	domain = strings.Replace(domain, ".", "_", -1)
	domain = strings.Replace(domain, ":", "_", -1)
	domain = strings.Replace(domain, "/", "_", -1)
	return fmt.Sprintf("%s/%s/%s/%s", SchemaTrellis, name, version, domain)
}

// WorkerETCDDomainPath 工作路径
func WorkerETCDDomainPath(name, version, domain string) string {
	domain = strings.Replace(domain, ".", "_", -1)
	domain = strings.Replace(domain, ":", "_", -1)
	domain = strings.Replace(domain, "/", "_", -1)
	return fmt.Sprintf("%s%s/%s/%s", SchemaETCDNaming, name, version, domain)
}
