package rdp

func (c *client) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}
