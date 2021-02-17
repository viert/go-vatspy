package vatspy

import (
	"sync"
	"time"

	"github.com/viert/go-vatspy/static"
)

// Provider is a vatspy data provider supporting automatic updates
type Provider struct {
	lock          sync.RWMutex
	stop          chan bool
	staticData    *static.Data
	subscriptions []*subscription
}

// New creates a new Provider
func New(staticUpdatePeriod time.Duration, dynamicUpdatePeriod time.Duration) (*Provider, error) {
	p := new(Provider)
	p.stop = make(chan bool)
	p.subscriptions = make([]*subscription, 0)
	go p.loop(staticUpdatePeriod, dynamicUpdatePeriod)
	return p, nil
}

// Subscribe generates a new update channel
func (p *Provider) Subscribe(chanSize int) <-chan Update {
	sub := &subscription{
		state:   newStateData(),
		updates: make(chan Update, chanSize),
	}
	p.subscriptions = append(p.subscriptions, sub)
	return sub.updates
}

func (p *Provider) fetchDynamic() error {
	// var err error

	// if p.staticData == nil {
	// 	return fmt.Errorf("static data is not available yet")
	// }

	// p.dData, err = dynamic.Fetch(dynamic.VatSimJSON3URL)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (p *Provider) fetchStatic() error {
	data, err := static.Fetch(static.VATSpyDataPublicURL, static.FIRBoundariesPublicURL)
	if err != nil {
		return err
	}

	for _, sub := range p.subscriptions {
		sub.processStatic(data)
	}
	p.staticData = data
	return nil
}

func (p *Provider) loop(staticUpdatePeriod time.Duration, dynamicUpdatePeriod time.Duration) {
	st := time.NewTicker(staticUpdatePeriod)
	dt := time.NewTicker(dynamicUpdatePeriod)
	defer st.Stop()
	defer dt.Stop()

	p.fetchStatic()
	p.fetchDynamic()

	for {
		select {
		case <-p.stop:
			return
		case <-st.C:
			p.fetchStatic()
		case <-dt.C:
			p.fetchDynamic()
		}
	}
}

// Stop stops the provider's update loop
func (p *Provider) Stop() {
	p.stop <- true
}
