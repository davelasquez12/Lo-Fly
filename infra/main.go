package main

import (
	"fmt"
	"os"
	"strings"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type LoflyInfraStackProps struct {
	awscdk.StackProps
	EnvironmentName string
}

func NewLoflyInfraStack(scope constructs.Construct, id string, props *LoflyInfraStackProps) awscdk.Stack {
	stackProps := awscdk.StackProps{}
	if props != nil {
		stackProps = props.StackProps
	}

	stack := awscdk.NewStack(scope, jsii.String(id), &stackProps)
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	environmentName := env("LOFLY_ENV", "dev")

	NewLoflyInfraStack(app, fmt.Sprintf("Lofly%sStack", title(environmentName)), &LoflyInfraStackProps{
		StackProps: awscdk.StackProps{
			Env: cdkEnvironment(),
		},
		EnvironmentName: environmentName,
	})

	app.Synth(nil)
}

func cdkEnvironment() *awscdk.Environment {
	region := env("CDK_DEFAULT_REGION", env("AWS_REGION", "us-east-1"))
	account := os.Getenv("CDK_DEFAULT_ACCOUNT")
	if account == "" {
		return &awscdk.Environment{
			Region: jsii.String(region),
		}
	}
	return &awscdk.Environment{
		Account: jsii.String(account),
		Region:  jsii.String(region),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func title(value string) string {
	if value == "" {
		return ""
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func isProtectedEnvironment(environmentName string) bool {
	return strings.EqualFold(environmentName, "prod")
}
