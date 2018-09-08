/*
 * Revision History:
 *     Initial: 2018/08/04        Li Zebang
 */

package app

// ContextType -
type ContextType string

const (
	// RetryTaskKey -
	RetryTaskKey ContextType = "Retry-Task-Key"
	// SearchTaskKey -
	SearchTaskKey ContextType = "Search-Task-Key"
	// ListTaskKey -
	ListTaskKey ContextType = "List-Task-Key"
	// IndexTaskKey -
	IndexTaskKey ContextType = "Index-Task-Key"

	// RepoDir -
	RepoDir = "repo-"
	// ReposDir -
	ReposDir = "repos/"
	// CacheDir -
	CacheDir = "cache/"
	// ReposJSON -
	ReposJSON = "repos.json"
	// InfoJSON -
	InfoJSON = "info.json"
)
