package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

type Setting func(*cobra.Command) *cobra.Command

func Use(s string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Use = s
		return c
	}
}

var Name = Use

func Version(s string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Version = s
		return c
	}
}

func Short(s string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Short = s
		return c
	}
}

var Summary = Short

func Long(s string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Long = s
		return c
	}
}

var Description = Long
var Desc = Description

func Aliases(ss ...string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Aliases = ss
		return c
	}
}

func Run(f func(*cobra.Command, []string)) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Run = f
		return c
	}
}

func Args(a cobra.PositionalArgs) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Args = a
		return c
	}
}

func Example(s string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Example = s
		return c
	}
}

func Hidden() Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.Hidden = true
		return c
	}
}

func Command(ss ...Setting) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.AddCommand(New(ss...))
		return c
	}
}

func Commands(cmds ...*cobra.Command) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.AddCommand(cmds...)
		return c
	}
}

func String(p *string, name, shorthand, value, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().StringVarP(p, name, shorthand, value, usage)
		return c
	}
}

func Bool(p *bool, name, shorthand string, value bool, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().BoolVarP(p, name, shorthand, value, usage)
		return c
	}
}

func Int(p *int, name, shorthand string, value int, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().IntVarP(p, name, shorthand, value, usage)
		return c
	}
}

func Uint(p *uint, name, shorthand string, value uint, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().UintVarP(p, name, shorthand, value, usage)
		return c
	}
}

func Int8(p *int8, name, shorthand string, value int8, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Int8VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Uint8(p *uint8, name, shorthand string, value uint8, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Uint8VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Int16(p *int16, name, shorthand string, value int16, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Int16VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Uint16(p *uint16, name, shorthand string, value uint16, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Uint16VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Int32(p *int32, name, shorthand string, value int32, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Int32VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Uint32(p *uint32, name, shorthand string, value uint32, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Uint32VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Int64(p *int64, name, shorthand string, value int64, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Int64VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Uint64(p *uint64, name, shorthand string, value uint64, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Uint64VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Float32(p *float32, name, shorthand string, value float32, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Float32VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Float64(p *float64, name, shorthand string, value float64, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().Float64VarP(p, name, shorthand, value, usage)
		return c
	}
}

func Duration(p *time.Duration, name, shorthand string, value time.Duration, usage string) Setting {
	return func(c *cobra.Command) *cobra.Command {
		c.PersistentFlags().DurationVarP(p, name, shorthand, value, usage)
		return c
	}
}
