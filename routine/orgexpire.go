package routine

import (
	"context"
	"time"
	"xg/entity"
	"xg/log"
	"xg/service"
)

type OrgExpire struct {
}

func (o *OrgExpire) doExpireOrg(ctx context.Context, id int) {
	err := service.GetOrgService().ExpireOrgById(ctx, id)
	if err != nil {
		log.Error.Printf("Can't expire orgs, id: %v, err: %v\n", id, err)
	}
}
func (o *OrgExpire) doCheckExpire(ctx context.Context) {
	total, orgs, err := service.GetOrgService().ListOrgsByStatus(ctx, []int{entity.OrgStatusCertified})
	if err != nil {
		log.Error.Println("Can't search orgs, err:", err)
		return
	}
	if total < 1 {
		return
	}
	now := time.Now()
	for i := range orgs {
		if orgs[i].ExpiredAt == nil {
			continue
		}
		if orgs[i].ExpiredAt.Before(now) {
			o.doExpireOrg(ctx, orgs[i].ID)
		}
	}
}

func (o *OrgExpire) Start() {
	go func() {
		o.doCheckExpire(context.Background())
		time.Sleep(time.Hour * 12)
	}()
}
