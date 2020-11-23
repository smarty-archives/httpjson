package httpjson

var Options options

type RenderOption func(*ResponseRenderer)

type options struct{}

func (options) RenderIndent(prefix, indent string) RenderOption {
	return func(renderer *ResponseRenderer) {
		renderer.Prefix = prefix
		renderer.Indent = indent
	}
}
func (options) RenderStatusCode(code int) RenderOption {
	return func(renderer *ResponseRenderer) {
		renderer.StatusCodeValue = code
	}
}
func (options) RenderHeader(key string, values ...string) RenderOption {
	return func(renderer *ResponseRenderer) {
		renderer.AddHeader(values, key)
	}
}
func (options) RenderData(data interface{}) RenderOption {
	return func(renderer *ResponseRenderer) {
		renderer.Data = data
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
