package rdp

func (c *client) SetRemoteApp(app, args, workingDir string) {
	c.remoteApp = &RemoteApp{
		App:        app,
		WorkingDir: workingDir,
		Args:       args,
	}
	c.channels = append(c.channels, "rail")
	c.railState = RailStateUninitialized
}
