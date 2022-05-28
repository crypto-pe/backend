package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

func newFromFile(file string, env string, config interface{}) error {
	if file == "" {
		file = env
	}

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return errors.Wrap(err, "failed to load config file")
	}

	if _, err := toml.DecodeFile(file, config); err != nil {
		return errors.Wrap(err, "failed to parse config file")
	}

	return nil
}

type Config struct {
	Auth Auth `toml:"auth"`
}

type Auth struct {
	JWTSecret string `toml:"jwt_secret"`
}

var (
	flags      = flag.NewFlagSet("jwt", flag.ExitOnError)
	configFile = flags.String("config", "", "path to config file")
	decode     = flags.String("decode", "", "JWT token to decode and print claims")
	account    = flags.String("account", "", "account wallet address")
	service    = flags.String("service", "", "service id")
)

var services = []string{} // add service names here

func main() {
	flags.Parse(os.Args[1:])

	if len(os.Args) == 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	// Parse config file.
	cfg := &Config{}
	err := newFromFile(*configFile, os.Getenv("CONFIG"), cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Check flags
	*service = strings.ToLower(*service)
	if *service != "" && stringExists(*service, services) == false {
		log.Fatalf("Invalid service flag '%s', support options: %v", *service, services)
	}

	// Decode given JWT token.
	if *decode != "" {
		token, err := jwt.Parse([]byte(*decode), jwt.WithVerify("HS256", []byte(cfg.Auth.JWTSecret)))
		if err != nil {
			log.Fatal(err)
		}
		if err := jwt.Validate(token); err != nil {
			log.Fatal(err)
		}

		claims, _ := token.AsMap(context.Background())
		fmt.Println()
		for k, v := range claims {
			fmt.Printf("%v: %+v\n", k, v)
		}

		return
	}

	// Encode new JWT token.
	token := jwt.New()
	if *account != "" {
		token.Set("account", *account)
	}
	if *service != "" && *service != "web" {
		token.Set("service", *service)
	}
	tokenPayload, err := jwt.Sign(token, "HS256", []byte(cfg.Auth.JWTSecret))
	if err != nil {
		log.Fatal(err)
	}
	tokenStr := string(tokenPayload)

	fmt.Fprintln(os.Stderr)

	fmt.Println(tokenStr)

	claims, _ := token.AsMap(context.Background())
	fmt.Fprintf(os.Stderr, "\nClaims: %#v\n", claims)
}

func stringExists(s string, ss []string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

var usage = `USAGE:

# Decode JWT:
./bin/util-jwt -config=etc/example.conf -decode "{JWT}"

# Accounts:
./bin/util-jwt -config=etc/example.conf -account={ACCOUNT_ADDRESS}

# Services:
./bin/util-jwt -config=etc/example.conf -service={SERVICE_NAME}
`
