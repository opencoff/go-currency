// currency.go - Atto dollars represented as big.Int
//
// (c) 2017, Sudhi Herle <sudhi@herle.net>
//
//
// Licensing Terms: GPLv2
//
// If you need a commercial license for this work, please contact
// the author.
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

// Package currency implements a decimal currency type as
// "atto dollars" (1.0e-18) represented as a big.Int.
// All arithmetic is done on the underlying Big.Int. By default, the
// output conversion to string uses the full (18 decimal digit)
// precision. Output string representation is not rounded - but
// truncated.
package currency

import (
	"fmt"
	"math/big"
	"strings"
)

// A currency is represented as atto dollars (18 digits of precision)
type Currency struct {
	big.Int
}

// Atto exponent, multiplication factor and padding string
const eExp = 18

var iMult int64       // multiplicative factor for 'characteristic'
var iBigMult *big.Int // same as iMult, but as a big.Int
var zeroes string     // array of zeroes
var zero *big.Int     // "zero" currency

func init() {
	for i := 0; i < eExp; i++ {
		zeroes += "0"
	}

	iMult = pow64(10, eExp)
	iBigMult = big.NewInt(iMult)
	zero = big.NewInt(0)
}

// compute a ** b for unsigned quantities using binary
// exponentiation method
func pow64(a int64, b int) int64 {
	var r int64 = 1
	for b > 0 {
		if 0 != (b & 1) {
			r *= a
		}

		b >>= 1
		a *= a
	}
	return r
}

// Create a zero valued currency
func New() *Currency {
	return &Currency{}
}

// Make a new Currency instance with input string 's' and output
// precision of 'oprec'. If 'oprec' is more than Atto Dollars, it is
// clamped at Atto Dollars (18). If it is less than or equal to
// zero, it is clamped at 6.
func NewFromString(s string) (*Currency, error) {
	p := &Currency{}

	if len(s) > 0 {
		err := parse(&p.Int, s)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// Convert 'p' to a string - bounded by output precision
// We just shift 12 digits off the left and print it.
func (p Currency) String() string {
	return stringify(&p.Int, eExp)
}

// Show 'p' to a string bounded by output precision 'oprec'
// If 'oprec' is more than the atto-dollar resolution, it is clamped
// at 12.
func (p *Currency) StringFixed(oprec int) string {
	if oprec > eExp || oprec <= 0 {
		oprec = eExp
	}

	return stringify(&p.Int, oprec)
}

// stringify atto-dollars in 'b'
func stringify(b *big.Int, oprec int) string {
	var m, x string

	s := b.String()

	// Not enough atto dollars
	if len(s) <= eExp {
		m = "0"
		if lpad := eExp - len(s); lpad > 0 {
			x = zeroes[:lpad] + s
		} else {
			x = s
		}
	} else if len(s) > eExp {
		n := len(s) - eExp
		m, x = s[:n], s[n:]
	}

	if len(x) > oprec {
		x = x[:oprec]
	}

	return fmt.Sprintf("%s.%s", m, x)
}

// Add 'x' to 'p'
func (p *Currency) Add(x *Currency) *Currency {
	p.Int.Add(&p.Int, &x.Int)
	return p
}

// Subtract 'x' from 'p'
func (p *Currency) Sub(x *Currency) *Currency {
	p.Int.Sub(&p.Int, &x.Int)
	return p
}

// Multiply 'p' with 'x'
func (p *Currency) Mul(x *Currency) *Currency {
	p.Int.Mul(&p.Int, &x.Int)
	return p
}

// Divide 'p' by 'x' and return the dividend
func (p *Currency) Div(x *Currency) *Currency {
	p.Int.Quo(&p.Int, &x.Int)
	return p
}

// Divide 'p' by 'x', and set p to the quotient and return 'p' and
// the remainder. This implements math/big's DivMod (Euclidean
// division):
//   q = p div x
//   r = p - (x * q)
//   p = q
//   return p, r
func (p *Currency) DivMod(x *Currency) (*Currency, *Currency) {
	var r big.Int
	p.Int.DivMod(&p.Int, &x.Int, &r)

	return p, &Currency{Int: r}
}

// Return true if this is zero
func (p *Currency) IsZero() bool {
	return 0 == p.Cmp(zero)
}

// Return true if 'p' is equal to 'x', false otherwise
func (p *Currency) Eq(x *Currency) bool {
	return 0 == p.Int.Cmp(&x.Int)
}

// Return a+b
func Add(a, b *Currency) *Currency {
	var z big.Int

	z.Add(&a.Int, &b.Int)
	return &Currency{Int: z}
}

// Return a-b
func Sub(a, b *Currency) *Currency {
	var z big.Int

	z.Sub(&a.Int, &b.Int)
	return &Currency{Int: z}
}

// Return a*b
func Mul(a, b *Currency) *Currency {
	var z big.Int

	z.Mul(&a.Int, &b.Int)
	return &Currency{Int: z}
}

// Return a/b
func Div(a, b *Currency) *Currency {
	var z big.Int

	z.Quo(&a.Int, &b.Int)
	return &Currency{Int: z}
}

// Do Euclidean division of a by b, return the quotient and
// reminder.
func DivMod(a, b *Currency) (*Currency, *Currency) {
	var z big.Int
	var r big.Int

	z.DivMod(&a.Int, &b.Int, &r)
	return &Currency{Int: z}, &Currency{Int: r}
}

// Return 1/a
func Inv(a *Currency) *Currency {
	var z big.Int

	z.Quo(iBigMult, &a.Int)
	return &Currency{Int: z}
}

// Return true of a == b
func Eq(a, b *Currency) bool {
	return 0 == a.Int.Cmp(&b.Int)
}

// Return -1, 0, +1 if a < b, a == b, a > b respectively
func Cmp(a, b *Currency) int {
	return a.Int.Cmp(&b.Int)
}

// Marshal 'p' to JSON
func (p *Currency) MarshalJSON() ([]byte, error) {
	s := p.String()
	return []byte(s), nil
}

// Unmarshal JSON to 'p'
func (p *Currency) UnmarshalJSON(txt []byte) error {

	return parse(&p.Int, string(txt))
}

// Parse a valid string 's' into a atto-dollar big.Int
func parse(p *big.Int, s string) error {
	v := strings.Split(s, ".")
	var pre, post string

	switch len(v) {
	case 1:
		pre = zstripPre(s)
	case 2:
		pre = zstripPre(v[0])
		post = zstripPost(v[1])

	default:
		return fmt.Errorf("malformed decimal %s", s)
	}

	if len(pre) > 0 {
		if _, ok := p.SetString(pre, 10); !ok {
			return fmt.Errorf("invalid decimal %s", s)
		}
	}

	// Truncate longer mantissae
	if len(post) > eExp {
		post = post[:eExp]
	}

	// We need the length of the string before we strip out the
	// leading zeroes. This length tells us how large the exponent of
	// the mantissa should be.
	n := len(post)
	post = zstripPre(post)

	f := &big.Int{}
	if len(post) > 0 {
		if _, ok := f.SetString(post, 10); !ok {
			return fmt.Errorf("invalid fraction %s", s)
		}
	}

	if exp := eExp - n; exp > 0 {
		var m big.Int
		m.SetInt64(pow64(10, exp))
		f.Mul(f, &m)
	}

	p.Mul(p, iBigMult)
	p.Add(p, f)
	return nil
}

func zstripPre(s string) string {
	n := len(s)
	for i := 0; i < n; i++ {
		if s[i] != '0' {
			return s[i:]
		}
	}

	return ""
}

func zstripPost(s string) string {
	n := len(s)
	for i := n - 1; i >= 0; i-- {
		if s[i] != '0' {
			return s[:i+1]
		}
	}
	return ""
}
