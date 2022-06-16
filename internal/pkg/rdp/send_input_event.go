package rdp

import "github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/fastpath"

func (c *client) SendInputEvent(data []byte) error {
	return c.fastPath.Send(fastpath.NewInputEventPDU(data))
}
