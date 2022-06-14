package rdp

func (c *client) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}
