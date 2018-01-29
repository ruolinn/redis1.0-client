package redis

import (
	"fmt"
	"strings"
	"time"

	"redis/internal"
	"redis/pool"
)

type Cmder interface {
	Name() string
	Args() []interface{}
	stringArg(int) string

	readReply(*pool.Conn) error
	setErr(error)

	readTimeout() *time.Duration

	Err() error
	fmt.Stringer
}

type baseCmd struct {
	_args []interface{}
	err   error

	_readTimeout *time.Duration
}

func (cmd *baseCmd) Err() error {
	return cmd.err
}

func (cmd *baseCmd) Args() []interface{} {
	return cmd._args
}

func (cmd *baseCmd) stringArg(pos int) string {
	if pos < 0 || pos >= len(cmd._args) {
		return ""
	}
	s, _ := cmd._args[pos].(string)
	return s
}

func (cmd *baseCmd) Name() string {
	if len(cmd._args) > 0 {
		// Cmd name must be lower cased.
		s := internal.ToLower(cmd.stringArg(0))
		cmd._args[0] = s
		return s
	}
	return ""
}

func (cmd *baseCmd) readTimeout() *time.Duration {
	return cmd._readTimeout
}

func (cmd *baseCmd) setReadTimeout(d time.Duration) {
	cmd._readTimeout = &d
}

func (cmd *baseCmd) setErr(e error) {
	cmd.err = e
}

type StatusCmd struct {
	baseCmd
	val string
}

func NewStatusCmd(args ...interface{}) *StatusCmd {
	return &StatusCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *StatusCmd) Val() string {
	return cmd.val
}

func (cmd *StatusCmd) Result() (string, error) {
	return cmd.val, cmd.err
}

func (cmd *StatusCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *StatusCmd) readReply(cn *pool.Conn) error {
	cmd.val, cmd.err = cn.Rd.ReadStringReply()
	return cmd.err
}

func cmdString(cmd Cmder, val interface{}) string {
	var ss []string
	for _, arg := range cmd.Args() {
		ss = append(ss, fmt.Sprint(arg))
	}
	s := strings.Join(ss, " ")
	if err := cmd.Err(); err != nil {
		return s + ": " + err.Error()
	}
	if val != nil {
		switch vv := val.(type) {
		case []byte:
			return s + ": " + string(vv)
		default:
			return s + ": " + fmt.Sprint(val)
		}
	}
	return s

}

func writeCmd(cn *pool.Conn, cmds ...Cmder) error {
	cn.Wb.Reset()
	for _, cmd := range cmds {
		if err := cn.Wb.Append(cmd.Args()); err != nil {
			return err
		}
	}

	_, err := cn.Write(cn.Wb.Bytes())
	return err
}
