package templat

var LibTempl = `// {{ .ProgName }} by {{ .AuthorName }} <author at email dot com>
// Library for ...
package {{ .ProgName }}
`

var LibTemplWithLog = `// {{ .ProgName }} by {{ .AuthorName }} <author at email dot com>
// Library for ...
package {{ .ProgName }}

import "github.com/LitFill/fatal"
`
