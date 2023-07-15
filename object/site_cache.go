package object

var siteMap = map[string]*Site{}

func InitSiteMap() {
	refreshSiteMap()
}

func refreshSiteMap() {
	sites := GetGlobalSites()
	for _, site := range sites {
		if _, ok := siteMap[site.Domain]; !ok {
			siteMap[site.Domain] = site
		}
	}
}

func GetSiteByDomain(domain string) *Site {
	if site, ok := siteMap[domain]; ok {
		return site
	} else {
		return nil
	}
}
