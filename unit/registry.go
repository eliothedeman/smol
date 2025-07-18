package unit

import (
	"context"
	"sync"
)

type Registry struct {
	units         map[string]Unit
	refs          map[string]*unitRef
	subscriptions map[string]map[string]struct{}
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
}

type unitRef struct {
	name   string
	reg    *Registry
	closed bool
}

type registryCtx struct {
	context.Context
	reg  *Registry
	self *unitRef
}

func NewRegistry() *Registry {
	ctx, cancel := context.WithCancel(context.Background())
	return &Registry{
		units:         make(map[string]Unit),
		refs:          make(map[string]*unitRef),
		subscriptions: make(map[string]map[string]struct{}),
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (r *Registry) Register(name string, unit Unit) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ref := &unitRef{
		name: name,
		reg:  r,
	}

	r.units[name] = unit
	r.refs[name] = ref
}

func (r *Registry) Start() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for name, unit := range r.units {
		ctx := &registryCtx{
			Context: r.ctx,
			reg:     r,
			self:    r.refs[name],
		}
		unit.Init(ctx)
	}
	return nil
}

func (r *Registry) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cancel()

	for name := range r.subscriptions {
		delete(r.subscriptions, name)
	}
}

func (r *Registry) getUnit(name string) (Unit, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	unit, exists := r.units[name]
	return unit, exists
}

func (r *Registry) getRef(name string) *unitRef {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.refs[name]
}

func (r *unitRef) Name() string {
	return r.name
}

func (r *unitRef) Send(msg any) {
	if r.closed {
		return
	}

	unit, exists := r.reg.getUnit(r.name)
	if !exists {
		return
	}

	ctx := &registryCtx{
		Context: r.reg.ctx,
		reg:     r.reg,
		self:    r,
	}

	go unit.Handle(ctx, r, msg)

	r.reg.mu.RLock()
	subscribers := r.reg.subscriptions[r.name]
	r.reg.mu.RUnlock()

	for subscriberName := range subscribers {
		subUnit, exists := r.reg.getUnit(subscriberName)
		if !exists {
			continue
		}
		subRef := r.reg.getRef(subscriberName)
		if subRef != nil && !subRef.closed {
			subCtx := &registryCtx{
				Context: r.reg.ctx,
				reg:     r.reg,
				self:    subRef,
			}
			go subUnit.Handle(subCtx, r, msg)
		}
	}
}

func (r *unitRef) Stop() {
	r.closed = true
}

func (c *registryCtx) Units() []UnitDesc {
	c.reg.mu.RLock()
	defer c.reg.mu.RUnlock()

	var units []UnitDesc
	for name, unit := range c.reg.units {
		units = append(units, UnitDesc{
			Name:  name,
			Proxy: unit,
		})
	}
	return units
}

func (c *registryCtx) Spawn(name string, f UnitFactory) UnitRef {
	c.reg.mu.Lock()
	defer c.reg.mu.Unlock()

	if _, exists := c.reg.units[name]; exists {
		return c.reg.refs[name]
	}

	unit := f()
	ref := &unitRef{
		name: name,
		reg:  c.reg,
	}

	c.reg.units[name] = unit
	c.reg.refs[name] = ref

	ctx := &registryCtx{
		Context: c.reg.ctx,
		reg:     c.reg,
		self:    ref,
	}

	unit.Init(ctx)
	return ref
}

func (c *registryCtx) Self() UnitRef {
	return c.self
}

func (c *registryCtx) Subscribe(other Unit) {
	c.reg.mu.Lock()
	defer c.reg.mu.Unlock()

	otherName := ""
	for name, unit := range c.reg.units {
		if unit == other {
			otherName = name
			break
		}
	}

	if otherName == "" {
		return
	}

	selfName := c.self.name
	if c.reg.subscriptions[otherName] == nil {
		c.reg.subscriptions[otherName] = make(map[string]struct{})
	}
	c.reg.subscriptions[otherName][selfName] = struct{}{}
}

func (c *registryCtx) Unsubscribe(other Unit) {
	c.reg.mu.Lock()
	defer c.reg.mu.Unlock()

	otherName := ""
	for name, unit := range c.reg.units {
		if unit == other {
			otherName = name
			break
		}
	}

	if otherName == "" {
		return
	}

	selfName := c.self.name
	if c.reg.subscriptions[selfName] != nil {
		delete(c.reg.subscriptions[selfName], otherName)
		if len(c.reg.subscriptions[selfName]) == 0 {
			delete(c.reg.subscriptions, selfName)
		}
	}
}
