package store

import "time"

var (
	AdvertiseStatusActive   = "active"
	AdvertiseStatusInactive = "inactive"
)

type AdvertiseRecord struct {
	ID        string `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	ImageURL  string `db:"image_url" json:"imageUrl"`
	Placement string `db:"placement" json:"placement"`
	Status    string `db:"status" json:"status"`
	CreatedAt int64  `db:"created_at" json:"createdAt"`
	ExpiresAt *int64 `db:"expires_at" json:"expiresAt"`
	Expired   bool   `db:"-" json:"expired"`
}

func (r *AdvertiseRecord) CalculateAndSetExpired() {
	if time.Now().Unix() > *r.ExpiresAt {
		r.Expired = true
	}
}
