package primitiveutil

import (
	"runtime"
	"strconv"
	"strings"
)

// Config use with `-i`, `-o`, and `-n`.
type Config struct {
	Input      string    `short:"i" description:"input image path" required:"true"`
	Outputs    FlagArray `short:"o" description:"output image path" required:"true"`
	Background string    `short:"b" description:"background color in hex"`
	Number     int       `short:"n" default:"1" description:"number of primitives"`
	Alpha      int       `short:"a" default:"128" description:"alpha value"`
	InputSize  int       `short:"r" default:"256" description:"resize large input images to this size"`
	OutputSize int       `short:"s" default:"1024" description:"output image size"`
	Mode       int       `short:"m" default:"1" choice:"0" choice:"1" choice:"3" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" description:"0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon"`
	Workers    int       `short:"j" default:"0" description:"number of parallel workers (default uses all cores)"`
	Nth        int       `long:"nth" default:"1" description:"save every Nth frame (put \"%d\" in path)"`
	Repeat     int       `long:"rep" default:"0" description:"dd N extra shapes per iteration with reduced search"`
	Verbose    []bool    `short:"v" description:"verbose"`
	Configs    ShapeConfigs
}

func (cfg *Config) SetDefaults() {
	if cfg.Number <= 0 {
		cfg.Number = 1
	}
	if cfg.Alpha <= 0 {
		cfg.Alpha = 128
	}
	if cfg.InputSize <= 0 {
		cfg.InputSize = 256
	}
	if cfg.OutputSize <= 0 {
		cfg.OutputSize = 1024
	}
	if cfg.Nth <= 1 {
		cfg.Nth = 1
	}
}

func (cfg *Config) Inflate() {
	cfg.SetDefaults()
	cfg.Configs = []ShapeConfig{
		{
			Count:  cfg.Number,
			Mode:   cfg.Mode,
			Alpha:  cfg.Alpha,
			Repeat: cfg.Repeat,
		},
	}
	// determine worker count
	if cfg.Workers < 1 {
		cfg.Workers = runtime.NumCPU()
	}
}

type FlagArray []string

func (i *FlagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *FlagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type ShapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}

type ShapeConfigs []ShapeConfig

func (i *ShapeConfigs) String() string {
	return ""
}

func (i *ShapeConfigs) Set(value string, mode, alpha, repeat int) error {
	n, _ := strconv.ParseInt(value, 0, 0)
	*i = append(*i, ShapeConfig{int(n), mode, alpha, repeat})
	return nil
}
