package subscription

import (
	"context"
	"posts-service/internal/entity"
	"sync"
)

type commentPubSub struct {
	mx          sync.RWMutex
	subscribers map[int64]map[chan *entity.Comment]struct{}
}

func NewCommentPubSub() *commentPubSub {
	return &commentPubSub{
		subscribers: make(map[int64]map[chan *entity.Comment]struct{}),
	}
}

func (p *commentPubSub) Subscribe(
	ctx context.Context,
	postID int64,
) (<-chan *entity.Comment, func()) {
	ch := make(chan *entity.Comment, 1)

	p.mx.Lock()
	if p.subscribers[postID] == nil {
		p.subscribers[postID] = make(map[chan *entity.Comment]struct{})
	}
	p.subscribers[postID][ch] = struct{}{}
	p.mx.Unlock()

	unsubscribe := func() {
		p.mx.Lock()
		defer p.mx.Unlock()

		if _, ok := p.subscribers[postID][ch]; ok {
			delete(p.subscribers[postID], ch)
			close(ch)
		}

		if len(p.subscribers[postID]) == 0 {
			delete(p.subscribers, postID)
		}
	}

	go func() {
		<-ctx.Done()
		unsubscribe()
	}()

	return ch, unsubscribe
}

func (p *commentPubSub) Publish(comment *entity.Comment) {
	p.mx.RLock()
	defer p.mx.RUnlock()

	for ch := range p.subscribers[comment.PostID] {
		select {
		case ch <- comment:
		default:
		}
	}
}
