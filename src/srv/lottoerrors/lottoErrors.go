package lottoerrors

import (
	"time"
)

type Operation string
type URL string
type HttpCode int

type LottoError struct {
	Err        error     `json:"err"`
	URL        URL       `json:"url"`
	StatusCode HttpCode  `json:"httpcode"`
	Severity   Severity  `json:"severity"`
	Category   Category  `json:"category"`
	Operation  Operation `json:"operation"`
	Logtime    time.Time `json:"timestamp"`
}

func (e *LottoError) Error() string {
	return "Wrapped Error"
}

type Category int

const (
	Category1 Category = 0
	Category2 Category = 1
	Category3 Category = 2
	Category4 Category = 3
)

func (c Category) String() string {

	names := [...]string{
		"Category1",
		"Category2",
		"Category3",
		"Category4"}

	return names[c]
}

type Severity int

const (
	Severity1 Severity = 0
	Severity2 Severity = 1
	Severity3 Severity = 2
	Severity4 Severity = 4
)

func (s Severity) String() string {

	names := [...]string{
		"Severity1",
		"Severity2",
		"Severity3",
		"Severity4"}

	return names[s]
}

func E(args ...interface{}) error {

	e := &LottoError{}

	for _, arg := range args {

		e.Logtime = time.Now()

		switch arg := arg.(type) {
		case Severity:
			e.Severity = arg
		case Category:
			e.Category = arg
		case error:
			e.Err = arg
		case Operation:
			e.Operation = arg
		default:
			panic("Bad call to E")
		}
	}

	return e
}
