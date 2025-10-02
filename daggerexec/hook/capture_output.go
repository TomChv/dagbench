package hook

type CaptureOutput struct {
	output string
}

func NewCaptureOutput() *CaptureOutput {
	return &CaptureOutput{}
}

func (c *CaptureOutput) Hook(line string) error {
	c.output += line + "\n"

	return nil
}

func (c *CaptureOutput) Output() string {
	return c.output
}
