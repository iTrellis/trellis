package api

import "time"

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

func (p *httpServer) syncAPIs() {
	for {
		p.options.Logger.Info("start_sync_apis", time.Now())
		var apis []*API
		if err := p.apiEngine.Where("`status` = ?", "normal").Find(&apis); err != nil {
			p.options.Logger.Error("sync_apis_failed", "err", err.Error())
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
		p.options.Logger.Info("end_sync_apis", time.Now())

		<-p.ticker.C
	}
}
