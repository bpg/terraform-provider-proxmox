/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package stringset

type options struct {
	separator string
}

type Option struct {
	apply func(*options)
}

func defaultOptions(opts ...Option) options {
	opt := options{
		separator: ";",
	}

	for _, o := range opts {
		o.apply(&opt)
	}

	return opt
}

// WithSeparator sets the separator for the string set value.
func WithSeparator(separator string) Option {
	return Option{
		apply: func(o *options) {
			o.separator = separator
		},
	}
}
