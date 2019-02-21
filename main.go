package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jessevdk/go-flags"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var opts struct {
	Path             string            `long:"path" description:"path for ssm:GetParametersByPath" required:"true" value-name:"PATH"`
	NoRecursive      bool              `long:"no-recursive" description:"get parameters not recuvsively"`
	NoOmitPathPrefix bool              `long:"no-omit-path-prefix" description:"No omit path prefix from parameter name"`
	NoUppercase      bool              `long:"no-uppercase" description:"No convert parameter name to uppercase"`
	CleanEnv         bool              `long:"with-clean-env" description:"No takeover OS Environment Variables"`
	ReplaceMap       map[string]string `long:"replace-map" description:"Pattern Table for parameter name replacement" value-name:"OLD_SUBSTR:NEW_SUBSTR"`
}

func main() {
	var parserOpts interface{}
	optionParser := flags.NewParser(parserOpts, flags.Default)
	optionParser.AddGroup("Options", "", &opts)
	optionParser.Name = "prmstore-exec"
	optionParser.Usage = "[OPTIONS] -- command"

	args, err := optionParser.Parse()
	if err != nil {
		flagError, ok := err.(*flags.Error)
		if !ok {
			panic(err)
		}

		if flagError.Type == flags.ErrHelp {
			os.Exit(1)
		} else {
			panic(flagError)
		}
	}

	params, err := getParameters()
	if err != nil {
		panic(err)
	}

	kvs := buildReplacedKeyValues(params)

	env := buildEnv(kvs)

	cmd, err := exec.LookPath(args[0])
	if err != nil {
		panic(fmt.Errorf("%s is not found", args[0]))
	}

	execErr := syscall.Exec(cmd, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func getParameters() ([]*ssm.Parameter, error) {
	sess := session.Must(session.NewSession())
	ssmSvc := ssm.New(sess)
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(opts.Path),
		Recursive:      aws.Bool(!opts.NoRecursive),
		WithDecryption: aws.Bool(true),
	}

	pageNum := 0
	params := []*ssm.Parameter{}
	err := ssmSvc.GetParametersByPathPages(input, func(page *ssm.GetParametersByPathOutput, lastPage bool) bool {
		pageNum++
		for _, p := range page.Parameters {
			params = append(params, p)
		}
		return !lastPage
	})

	return params, err
}

func buildReplacedKeyValues(params []*ssm.Parameter) map[string]string {
	keyValues := make(map[string]string)

	for _, p := range params {
		k := replaceName(*p.Name)
		keyValues[k] = *p.Value
	}

	return keyValues
}

func replaceName(name string) string {
	result := name

	if !opts.NoOmitPathPrefix {
		result = strings.TrimPrefix(name, opts.Path)
		result = strings.TrimPrefix(result, "/")
	}

	for old, new := range opts.ReplaceMap {
		result = strings.Replace(result, old, new, -1)
	}

	result = strings.Replace(result, "/", "_", -1)

	if !opts.NoUppercase {
		result = strings.ToUpper(result)
	}

	return result
}

func buildEnv(kvs map[string]string) []string {
	var env []string

	if !opts.CleanEnv {
		env = os.Environ()
	}

	for k, v := range kvs {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}
