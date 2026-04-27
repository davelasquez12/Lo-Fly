# LoflyBE Infra

Go CDK app for LoflyBE cloud infrastructure.

## Requirements

- AWS CLI profile: `lofly-admin`
- Go 1.26.1
- Node.js and the AWS CDK CLI installed on PATH

## Commands

```powershell
cd infra
go mod tidy
cdk bootstrap --profile lofly-admin
cdk synth --profile lofly-admin
cdk deploy --profile lofly-admin
```

Set `LOFLY_ENV` to choose the deployed environment name. It defaults to `dev`.

```powershell
$env:LOFLY_ENV = "dev"
```

The stack currently creates a DynamoDB table for tracked flights. Non-prod tables use a destroy removal policy so test infrastructure is easy to clean up. `prod` enables deletion protection, point-in-time recovery, and a retain removal policy.

If jsii has trouble finding Node on Windows, point it at your Node executable:

```powershell
$env:JSII_NODE = (Get-Command node).Source
```
