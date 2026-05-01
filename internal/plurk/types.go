package plurk

// Plurk mirrors the subset of fields used by this app from the Plurk API.
type Plurk struct {
	PlurkID               int       `json:"plurk_id"`
	Qualifier             string    `json:"qualifier"`
	QualifierTranslated   *string   `json:"qualifier_translated,omitempty"`
	IsUnread              int       `json:"is_unread"`
	PlurkType             int       `json:"plurk_type"`
	UserID                int       `json:"user_id"`
	OwnerID               int       `json:"owner_id"`
	Posted                string    `json:"posted"`
	NoComments            int       `json:"no_comments"`
	Content               string    `json:"content"`
	ContentRaw            *string   `json:"content_raw,omitempty"`
	ResponseCount         int       `json:"response_count"`
	ResponsesSeen         int       `json:"responses_seen"`
	LimitedTo             []int     `json:"limited_to"`
	Favorite              bool      `json:"favorite"`
	FavoriteCount         int       `json:"favorite_count"`
	Favorers              []int     `json:"favorers"`
	Replurkable           bool      `json:"replurkable"`
	Replurked             bool      `json:"replurked"`
	ReplurkerID           *int      `json:"replurker_id"`
	ReplurkersCount       int       `json:"replurkers_count"`
	Replurkers            []int     `json:"replurkers"`
	Porn                  bool      `json:"porn"`
}
