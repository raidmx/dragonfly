package entity

import "sync"

// HealthManager handles the health of an entity.
type HealthManager struct {
	mu     sync.RWMutex
	health float64
	max    float64
}

// NewHealthManager returns a new health manager with the health and max health provided.
func NewHealthManager(health, max float64) *HealthManager {
	if health > max {
		health = max
	}
	return &HealthManager{health: health, max: max}
}

// Health returns the current health of an entity.
func (m *HealthManager) Health() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.health
}

// AddHealth adds a given amount of health points to the player. If the health added to the current health
// exceeds the max, health will be set to the max. If the health is instead negative and results in a health
// lower than 0, the final health will be 0.
func (m *HealthManager) AddHealth(health float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	l := m.health + health
	if l < 0 {
		l = 0
	} else if l > m.max {
		l = m.max
	}
	m.health = l
}

// SetHealth sets the health of the player to the provided amount
func (m *HealthManager) SetHealth(health float64) {
	defer m.mu.Unlock()
	m.mu.Lock()

	if health < 0 {
		health = 0
	}

	if health > m.max {
		health = m.max
	}

	m.health = health
}

// MaxHealth returns the maximum health of the entity.
func (m *HealthManager) MaxHealth() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.max
}

// SetMaxHealth changes the max health of an entity to the maximum passed. If the maximum is set to 0 or
// lower, SetMaxHealth will default to a value of 1.
func (m *HealthManager) SetMaxHealth(max float64) {
	if max <= 0 {
		max = 1
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.max = max
	if m.health > max {
		m.health = max
	}
}
