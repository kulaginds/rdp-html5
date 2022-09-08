package rdp

func (c *client) Close() error {
	if c.remoteApp != nil {
		c.railState = RailStateUninitialized
	}

	return c.conn.Close()
}
