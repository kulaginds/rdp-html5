package rdp

import "github.com/kulaginds/rdp-html5/internal/pkg/rdp/fastpath"

func (c *client) SendInputEvent(data []byte) error {
	return c.fastPath.Send(fastpath.NewInputEventPDU(data))
}
