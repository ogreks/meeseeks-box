package card

type CardColumnSet struct {
	content []byte
}

func NewCardColumnSet(content []byte) *CardColumnSet {
	return &CardColumnSet{
		content: content,
	}
}

func (c *CardColumnSet) Tag() string {
	return "column_set"
}

func (c *CardColumnSet) MarshalJSON() ([]byte, error) {
	return c.content, nil
}
