package querybuilder

type Config struct {
	Limit     int64
	Offset    int64
	OrderBy   string
	OrderDesc bool
}

func (c *Conn) Limit(n int64) *Conn {
	c.Config.Limit = n
	return c
}

func (c *Conn) Offset(n int64) *Conn {
	c.Config.Offset = n
	return c
}

func (c *Conn) Page(n int64, limit int64) *Conn {
	c.Config.Offset = n * limit
	c.Config.Limit = limit
	return c
}

func (c *Conn) Order(s string, desc bool) *Conn {
	c.Config.OrderBy = s
	c.Config.OrderDesc = desc
	return c
}
