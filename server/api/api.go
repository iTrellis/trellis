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

package api

import (
	"github.com/google/uuid"
	"github.com/iTrellis/trellis/service"
)

// API api struct
type API struct {
	ID             string `xorm:"id"`
	Name           string `xorm:"name"`
	ServiceDomain  string `xorm:"service_domain"`
	ServiceName    string `xorm:"service_name"`
	ServiceVersion string `xorm:"service_version"`
	Topic          string `xorm:"topic"`
	Status         string `xorm:"status"`
	Version        int64  `xorm:"version"`
}

// TableName database table name
func (*API) TableName() string {
	return "api"
}

func (p *httpServer) syncAPIs(s *service.Service) {
	for {
		syncID := uuid.NewString()
		p.options.Logger.Info("msg", "start_sync_apis", "sync", syncID, "service", s)

		parmas := map[string]interface{}{"`status`": "normal"}

		if s.GetDomain() != "" {
			parmas["`service_domain`"] = s.GetDomain()
		}

		if s.GetName() != "" {
			parmas["`service_name`"] = s.GetName()
		}

		if s.GetName() != "" {
			parmas["`service_version`"] = s.GetVersion()
		}

		var apis []*API
		if err := p.apiEngine.Where(parmas).Find(&apis); err != nil {
			p.options.Logger.Error("msg", "sync_apis_failed", "sync", syncID, "err", err.Error())
			<-p.ticker.C
			continue
		}

		lenAPI := len(apis)
		mapAPIs := make(map[string]*API, lenAPI)

		for i := 0; i < lenAPI; i++ {
			mapAPIs[apis[i].Name] = apis[i]
		}

		p.syncer.Lock()
		p.apis = mapAPIs
		p.syncer.Unlock()
		p.options.Logger.Info("msg", "end_sync_apis", "sync", syncID, "service", s)

		<-p.ticker.C
	}
}
