# go-gas-express

Provisions s3, lambda, api gateway, and route53 infrastructure and deploys the go-gas-express app which is used to estimate gas prices based off of historical rates of change from the last 100 blocks

##### Requires:
`go version go1.13.5`

### ToDo:


### To run locally:
`cd app && go build && ./app`

### To deploy:
```
    make deploy
    cd terraform
    terraform apply
```
### To rm deployment:
```
    terraform destroy
```
