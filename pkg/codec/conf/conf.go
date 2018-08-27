/*
 * Revision History:
 *     Initial: 2018/08/26        Li Zebang
 */

package conf

// File -
type File interface {
	Deps() []string
}
