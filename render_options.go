package httpjson

var Options options

type RenderOption func(*ResponseRenderer)

type options struct{}

func (options) RenderIndent(prefix, indent string) RenderOption {
	return func(this *ResponseRenderer) {
		this.Prefix = prefix
		this.Indent = indent
	}
}
func (options) RenderStatusCode(code int) RenderOption {
	return func(this *ResponseRenderer) {
		this.StatusCodeValue = code
	}
}
func (options) RenderHeader(key string, values ...string) RenderOption {
	return func(this *ResponseRenderer) {
		this.AddHeader(values, key)
	}
}
func (options) RenderData(data interface{}) RenderOption {
	return func(this *ResponseRenderer) {
		this.Data = data
	}
}
func (options) RenderMappedErrors(mapper func(error) Error, errs ...error) RenderOption {
	return func(this *ResponseRenderer) {
		for _, err := range errs {
			this.IncludeError(mapper(err))
		}
	}
}
func (options) RenderErrors(errs ...Error) RenderOption {
	return func(this *ResponseRenderer) {
		this.IncludeError(errs...)
	}

}
