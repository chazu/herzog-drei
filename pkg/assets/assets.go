package assets

import (
	"fmt"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Manager handles loading and caching of game assets
type Manager struct {
	basePath string
	models   map[string]rl.Model
	textures map[string]rl.Texture2D
	sounds   map[string]rl.Sound
}

// NewManager creates a new asset manager with the given base path
func NewManager(basePath string) *Manager {
	return &Manager{
		basePath: basePath,
		models:   make(map[string]rl.Model),
		textures: make(map[string]rl.Texture2D),
		sounds:   make(map[string]rl.Sound),
	}
}

// LoadModel loads a 3D model from the models directory
func (m *Manager) LoadModel(name string) (rl.Model, error) {
	if model, ok := m.models[name]; ok {
		return model, nil
	}

	path := filepath.Join(m.basePath, "models", name)
	model := rl.LoadModel(path)

	if model.Meshes == nil {
		return model, fmt.Errorf("failed to load model: %s", path)
	}

	m.models[name] = model
	return model, nil
}

// LoadTexture loads a texture from the textures directory
func (m *Manager) LoadTexture(name string) (rl.Texture2D, error) {
	if tex, ok := m.textures[name]; ok {
		return tex, nil
	}

	path := filepath.Join(m.basePath, "textures", name)
	tex := rl.LoadTexture(path)

	if tex.ID == 0 {
		return tex, fmt.Errorf("failed to load texture: %s", path)
	}

	m.textures[name] = tex
	return tex, nil
}

// LoadSound loads a sound from the sounds directory
func (m *Manager) LoadSound(name string) (rl.Sound, error) {
	if snd, ok := m.sounds[name]; ok {
		return snd, nil
	}

	path := filepath.Join(m.basePath, "sounds", name)
	snd := rl.LoadSound(path)

	if snd.FrameCount == 0 {
		return snd, fmt.Errorf("failed to load sound: %s", path)
	}

	m.sounds[name] = snd
	return snd, nil
}

// Unload releases all loaded assets
func (m *Manager) Unload() {
	for _, model := range m.models {
		rl.UnloadModel(model)
	}
	for _, tex := range m.textures {
		rl.UnloadTexture(tex)
	}
	for _, snd := range m.sounds {
		rl.UnloadSound(snd)
	}

	m.models = make(map[string]rl.Model)
	m.textures = make(map[string]rl.Texture2D)
	m.sounds = make(map[string]rl.Sound)
}

// GetModel returns a cached model by name
func (m *Manager) GetModel(name string) (rl.Model, bool) {
	model, ok := m.models[name]
	return model, ok
}

// GetTexture returns a cached texture by name
func (m *Manager) GetTexture(name string) (rl.Texture2D, bool) {
	tex, ok := m.textures[name]
	return tex, ok
}

// GetSound returns a cached sound by name
func (m *Manager) GetSound(name string) (rl.Sound, bool) {
	snd, ok := m.sounds[name]
	return snd, ok
}
