# go-gas-express

Provisions s3, lambda, api gateway, and route53 infrastructure and deploys the go-gas-express app which is used to estimate gas prices based off of historical rates of change from the last 100 blocks

##### Requires:

`git`, `terraform`, `aws cli`, `go1.13.5`

#### To Deploy:

1) `git clone https://github.com/mycryptohq/go-gas-express.git`
2) `cd go-gas-express`
3) `make deploy` -> builds the go app and zips it up.
4) `cd terraform && cp example-tfvars.txt terraform.tfvars`-> copies the example- tfvars txt file to a new tfvars file.
5) Change the new `terraform.tfvars`'s vars to include the correct data
6) `terraform init`
7) `terraform apply`
8) Approve the cert validation for the new acm cert created for gas._[your_root_domain]_.com using [aws console](https://docs.aws.amazon.com/acm/latest/userguide/gs-acm-validate-dns.html) or cli.

### To rm deployment:
```
    cd terraform && terraform destroy
```
