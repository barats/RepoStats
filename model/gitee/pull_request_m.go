package gitee

import "time"

type PullRequest struct {
	ID            int64     `json:"id"`
	HTMLURL       string    `json:"html_url"`
	DiffUrl       string    `json:"diff_url"`
	PatchUrl      string    `json:"patch_url"`
	Number        int64     `json:"number"`
	State         string    `json:"state"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ClosedAt      time.Time `json:"closed_at"`
	MergedAt      time.Time `json:"merged_at"`
	Mergeable     bool      `json:"mergeable"`
	CanMergeCheck bool      `json:"can_merge_check"`
	Title         string    `json:"title"`
	User          User      `json:"user"`
	Head          struct {
		Label string     `json:"label"`
		Ref   string     `json:"ref"`
		Sha   string     `json:"sha"`
		User  User       `json:"user"`
		Repo  Repository `json:"repo"`
	} `json:"head"`
}
