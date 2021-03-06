package owo

import (
	"bytes"
	"fmt"
	"net/http"
	"selfbot/config"
	"selfbot/sound"

	uuid "github.com/satori/go.uuid"

	"gorm.io/gorm"
)

var _ sound.Store = new(Store)

type Store struct {
	client Client
	gorm   *gorm.DB
}

func NewStore(gorm *gorm.DB, conf config.Config) (*Store, error) {
	var ret = &Store{
		client: Client{
			UploadURL: conf.OwO.UploadURL,
			Client:    &http.Client{Timeout: conf.OwO.Timeout},
			URL:       conf.OwO.URL,
		},
	}

	if err := gorm.AutoMigrate(&StoredSound{}); err != nil {
		return nil, fmt.Errorf("owo: new store: %w", err)
	}

	ret.gorm = gorm
	return ret, nil
}

func (s *Store) SaveSound(soundInfo *sound.Sound) (soundID uuid.UUID, err error) {
	var buf bytes.Buffer
	if err := sound.DataWrite(soundInfo.Data, &buf); err != nil {
		return uuid.Nil, fmt.Errorf("owo: save sound: %w", err)
	}

	url, err := s.client.SaveSoundData(soundInfo.Name, &buf)
	if err != nil {
		return uuid.Nil, fmt.Errorf("owo: save sound: %w", err)
	}

	soundInfo.ID = uuid.Must(uuid.NewV4())
	ss := new(StoredSound)
	ss.OwoURL = url
	ss.FromSound(*soundInfo)
	if err := s.gorm.Create(ss).Error; err != nil {
		return uuid.Nil, fmt.Errorf("owo: save sound: gorm create: %w", err)
	}

	return soundInfo.ID, nil
}

func (s *Store) LoadSound(soundID uuid.UUID) (soundInfo sound.Sound, err error) {
	var ss StoredSound
	err = s.gorm.Where(&StoredSound{ID: soundID}).First(&ss).Error
	if err != nil {
		return sound.Sound{}, fmt.Errorf("owo: load sound: gorm first: %w", err)
	}

	soundInfo = ss.ToSound()

	soundInfo.Data, err = s.client.LoadSoundData(ss.OwoURL)
	if err != nil {
		return sound.Sound{}, fmt.Errorf("owo: load sound: %w", err)
	}

	return soundInfo, nil
}

func (s *Store) ListSounds(listOptions sound.ListOptions) (listResponse sound.ListResponse, err error) {
	rows, err := s.gorm.
		Model(&StoredSound{}).
		Limit(listOptions.Limit).
		Order("created_at desc").
		Select("id").
		Rows()
	if err != nil {
		return sound.ListResponse{}, fmt.Errorf("owo: list sounds: gorm rows: %w", err)
	}

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return sound.ListResponse{}, fmt.Errorf("owo: list sounds: gorm rows scan: %w", err)
		}

		listResponse.SoundIDs = append(listResponse.SoundIDs, id)
	}

	if err := rows.Err(); err != nil {
		return sound.ListResponse{}, fmt.Errorf("owo: list sounds: gorm rows err: %w", err)
	}

	return listResponse, nil
}
