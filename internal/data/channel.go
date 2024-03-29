package data

import (
	"context"
	"fmt"
	"github.com/ChinasMr/kaka/internal/biz"
	"github.com/ChinasMr/kaka/pkg/log"
	"sync"
)

var defaultChannel = "live"

type channelRepo struct {
	log      *log.Helper
	channels map[string]*biz.Channel
	rwm      sync.RWMutex
}

func (r *channelRepo) List(ctx context.Context) ([]*biz.Channel, error) {
	rv := make([]*biz.Channel, 0, len(r.channels))
	r.rwm.RLock()
	defer r.rwm.RUnlock()
	for _, p := range r.channels {
		rv = append(rv, p)
	}
	return rv, nil
}

func (r *channelRepo) Get(_ context.Context, id string) (*biz.Channel, error) {
	r.rwm.RLock()
	defer r.rwm.RUnlock()
	rv, ok := r.channels[id]
	if !ok {
		return nil, fmt.Errorf("channel not found")
	}
	return rv, nil
}

func (r *channelRepo) Delete(_ context.Context, id string) error {
	r.rwm.Lock()
	defer r.rwm.Unlock()
	delete(r.channels, id)
	return nil
}

func (r *channelRepo) Create(_ context.Context, id string) (*biz.Channel, error) {

	nc := &biz.Channel{
		Id:        id,
		Terminals: biz.NewTerminalsOperator(),
		SDP:       nil,
		RawSDP:    nil,
	}
	r.rwm.Lock()
	r.channels[nc.Id] = nc
	r.rwm.Unlock()
	return nc, nil
}

func NewChannelRepo(logger log.Logger) biz.ChannelRepo {
	rp := &channelRepo{
		log:      log.NewHelper(logger),
		channels: map[string]*biz.Channel{},
		rwm:      sync.RWMutex{},
	}
	_, _ = rp.Create(context.Background(), defaultChannel)
	return rp
}
