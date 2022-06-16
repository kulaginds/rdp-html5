package rdp

import (
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/fastpath"
)

func (c *client) GetUpdate() (*fastpath.UpdatePDU, error) {
	return c.fastPath.Receive()
}
