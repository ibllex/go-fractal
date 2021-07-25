package fractal

type closureTransformer struct {
	*BaseTransformer
	trans func(t *BaseTransformer, data Any) M
}

func (t *closureTransformer) Transform(data Any) M {
	return t.trans(t.BaseTransformer, data)
}

// T is a wrapper for closure transformer
func T(trans func(t *BaseTransformer, data Any) M) Transformer {
	return &closureTransformer{
		BaseTransformer: &BaseTransformer{}, trans: trans,
	}
}
