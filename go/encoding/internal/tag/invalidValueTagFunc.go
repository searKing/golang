package tag

import "reflect"

func invalidValueTagFunc(_ *tagState, _ reflect.Value, _ tagOpts) (isUserDefined bool) { return false }
