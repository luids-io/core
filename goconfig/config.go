package goconfig

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config is the main configuration struct
type Config struct {
	v          *viper.Viper
	program    string
	snames     map[string]bool
	sections   []Section
	validators []Validator
}

// Section represents configuration sections
type Section struct {
	Name     string
	Required bool
	Short    bool
	Data     Configurable
}

// Configurable is the interface for configuration data
type Configurable interface {
	//SetPFlags sets posix flags with the prefix and short options
	SetPFlags(short bool, prefix string)
	//BindViper must bind to viper instance pflags with the prefix
	BindViper(v *viper.Viper, prefix string)
	//BindViper must load values from viper instance
	FromViper(v *viper.Viper, prefix string)
	//Empty returns true if configuration is empty
	Empty() bool
	//Validate returns error if invalid value
	Validate() error
	//Dump returns string
	Dump() string
}

// Validator defines additional validators
type Validator func(cfg *Config) error

// New creates a configuration with sections passed
func New(program string, sections ...Section) (*Config, error) {
	c := &Config{
		v:        viper.New(),
		program:  program,
		snames:   make(map[string]bool),
		sections: make([]Section, 0),
	}
	for _, s := range sections {
		_, ok := c.snames[s.Name]
		if ok {
			return nil, errors.New("duplicated section name")
		}
		c.sections = append(c.sections, s)
		c.snames[s.Name] = true
	}
	return c, nil
}

// PFlags register posix flags of the structs
func (c *Config) PFlags() {
	for _, s := range c.sections {
		s.Data.SetPFlags(s.Short, s.Name)
		s.Data.BindViper(c.v, s.Name)
	}
}

// LoadIfFile try to load from path if not empty
func (c *Config) LoadIfFile(path string) error {
	c.loadFromEnv()
	var err error
	if path == "" {
		err = c.loadFromDefaultFiles()
	} else {
		err = c.loadFromFile(path)
	}
	if err != nil {
		return err
	}
	c.loadValues()
	return c.validate()
}

// LoadFromFile configuration from defined file
func (c *Config) LoadFromFile(path string) error {
	c.loadFromEnv()
	err := c.loadFromFile(path)
	if err != nil {
		return err
	}
	c.loadValues()
	return c.validate()
}

// Load configuration from defined file
func (c *Config) Load() error {
	c.loadFromEnv()
	err := c.loadFromDefaultFiles()
	if err != nil {
		return err
	}
	c.loadValues()
	return c.validate()
}

// Data returns data from a section
func (c *Config) Data(section string) Configurable {
	for _, s := range c.sections {
		if s.Name == section {
			return s.Data
		}
	}
	return nil
}

// Dump returns data from a section
func (c *Config) Dump() string {
	output := ""
	for _, s := range c.sections {
		if output != "" {
			output = output + "\n"
		}
		output = output + fmt.Sprintf("%s: %s", s.Name, s.Data.Dump())
	}
	return output
}

// AddValidator add validator to object
func (c *Config) AddValidator(v Validator) {
	c.validators = append(c.validators, v)
}

func (c *Config) loadFromEnv() {
	c.v.SetEnvPrefix(c.program)
	c.v.AutomaticEnv()
}

func (c *Config) loadFromDefaultFiles() error {
	c.v.SetConfigName(c.program)
	c.v.AddConfigPath("./configs/")
	err := c.v.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			return nil
		}
	}
	return err
}

func (c *Config) loadFromFile(path string) error {
	if !fileExists(path) {
		return fmt.Errorf("config file %s not found", path)
	}
	c.v.SetConfigFile(path)
	return c.v.ReadInConfig()
}

func (c *Config) loadValues() {
	for _, s := range c.sections {
		s.Data.FromViper(c.v, s.Name)
	}
}

func (c *Config) validate() error {
	for _, s := range c.sections {
		empty := s.Data.Empty()
		if empty && s.Required {
			if s.Name == "" {
				return errors.New("default section: is required")
			}
			return fmt.Errorf("section '%s': is required", s.Name)
		}
		if !empty {
			err := s.Data.Validate()
			if err != nil {
				if s.Name == "" {
					return fmt.Errorf("default section: %v", err)
				}
				return fmt.Errorf("section '%s': %v", s.Name, err)
			}
		}
	}
	for _, v := range c.validators {
		if err := v(c); err != nil {
			return err
		}
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
