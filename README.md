# D2 (dyndns) Helper
DNS Helper is a binary intended to be run as a cron job.  It will get the public IP of the system it is running on, and with the configuration of a yml file, will update a DNS zone's intended entry for that zone name.  Currently only supports Route53.


You may configure this with a yaml file OR with environment variables, however, for security reasons, AWS credentials are only definable via environment variables (or a .aws/config).

You can see a sample config file [here](config_example.yml)

If configuration via environment aside from:

* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY

Here are the following Environment variables:

* D2_ZONE_ID - the ID of the hosted zone in route53
* D2_DNS_RECORD - The dns record within the hosted zone to update
* D2_SNS_NOTIFY - Whether or not to use an SNS Notification topic if Updating succeeds or fails
* D2_SNSTopic - The SNS topic to use for notification
* D2_SNSMessage - The message to accompany the SNS notification

It is recommended that the credentials provided only have an IAM policy to update hosted zones.

### Roadmap

* More DNS Providers
* A "Whats My IP" Service for this cli hosted by AWS API Gateway and AWS Lambda
* The Open Source Terraform/Cloudformation for this service as well