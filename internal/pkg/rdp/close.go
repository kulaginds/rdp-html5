package rdp

func (c *client) Close() error {
	return c.conn.Close()
}
