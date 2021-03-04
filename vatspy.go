package vatspy

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/viert/go-vatspy/dynamic"
	"github.com/viert/go-vatspy/static"
)

// Provider is a vatspy data provider supporting automatic updates
type Provider struct {
	lock          sync.RWMutex
	stop          chan *Subscription
	cleanup       bool
	staticData    *static.Data
	dynamicData   *dynamic.Data
	subscriptions map[uint64]*Subscription
	autoinc       uint64
}

// ReadyCallback is a callback function called when static data is initially loaded
type ReadyCallback func()

// New creates a new Provider
func New(staticUpdatePeriod time.Duration, dynamicUpdatePeriod time.Duration, staticReady ReadyCallback) (*Provider, error) {
	p := new(Provider)
	p.stop = make(chan *Subscription, 1024)
	p.subscriptions = make(map[uint64]*Subscription)
	go p.loop(staticUpdatePeriod, dynamicUpdatePeriod, staticReady)
	return p, nil
}

// Subscribe generates a new update channel
func (p *Provider) Subscribe(chanSize int, controlledOnly bool, filters ...UpdateFilter) *Subscription {
	id := atomic.AddUint64(&p.autoinc, 1)
	sub := &Subscription{
		subID:          id,
		state:          newStateData(),
		updates:        make(chan Update, chanSize),
		controlledOnly: controlledOnly,
		filters:        filters,
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.subscriptions[id] = sub
	return sub
}

// Unsubscribe cancels the given subscription
func (p *Provider) Unsubscribe(sub *Subscription) {
	p.stop <- sub
}

func (p *Provider) fetchDynamic() error {
	var err error
	if p.staticData == nil {
		return fmt.Errorf("static data is not available yet")
	}

	dynamicData, err := dynamic.Fetch(dynamic.VatSimJSON3URL)
	if err != nil {
		return err
	}
	p.dynamicData = dynamicData

	// safely copy subscriptions
	subs := make([]*Subscription, 0)
	p.lock.RLock()
	for _, sub := range p.subscriptions {
		subs = append(subs, sub)
	}
	p.lock.RUnlock()

	for _, sub := range subs {
		sub.processDynamic(p.dynamicData, p.staticData)
	}

	return nil
}

func (p *Provider) fetchStatic() error {
	data, err := static.Fetch(static.VATSpyDataPublicURL, static.FIRBoundariesPublicURL)
	if err != nil {
		return err
	}

	// safely copy subscriptions
	subs := make([]*Subscription, 0)
	p.lock.RLock()
	for _, sub := range p.subscriptions {
		subs = append(subs, sub)
	}
	p.lock.RUnlock()

	for _, sub := range subs {
		sub.processStatic(data)
	}
	p.staticData = data
	return nil
}

func (p *Provider) loop(staticUpdatePeriod time.Duration, dynamicUpdatePeriod time.Duration, staticReady ReadyCallback) {
	staticReadyCalled := false

	st := time.NewTicker(staticUpdatePeriod)
	dt := time.NewTicker(dynamicUpdatePeriod)

	err := p.fetchStatic()
	if err == nil && staticReady != nil {
		staticReady()
		staticReadyCalled = true
	}
	p.fetchDynamic()

	for {
		select {
		case sub := <-p.stop:

			if sub == nil {
				p.stop = nil
				p.lock.Lock()
				for _, sub := range p.subscriptions {
					close(sub.updates)
					sub.updates = nil
				}
				st.Stop()
				dt.Stop()
				p.lock.Unlock()
				return
			}

			if sub.updates != nil {
				p.lock.Lock()
				close(sub.updates)
				sub.updates = nil
				delete(p.subscriptions, sub.subID)
				p.lock.Unlock()
			}

		case <-st.C:
			err := p.fetchStatic()
			if !staticReadyCalled && staticReady != nil && err == nil {
				staticReady()
				staticReadyCalled = true
			}
		case <-dt.C:
			p.fetchDynamic()
		}
	}
}

// Stop stops the provider's update loop
func (p *Provider) Stop() {
	if !p.cleanup {
		p.cleanup = true
		p.stop <- nil
	}
}

// GetStaticData returns current static data object
func (p *Provider) GetStaticData() *static.Data {
	return p.staticData
}

// GetDynamicData returns current dynamic data object
func (p *Provider) GetDynamicData() *dynamic.Data {
	return p.dynamicData
}
