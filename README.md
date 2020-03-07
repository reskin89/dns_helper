# dns_helper
DNS Helper is a binary intended to be run as a cron job.  It will get the public IP of the system it is running on, and with the configuration of a yml file, will update a DNS zone's intended entry for that zone name.  Currently only supports Route53.


Configuration in the yml file will be:

* DNS Record to update
* Hosted Zone ID

Configuration given via environment:

* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY


This is for security reasons so secrets are never published into source control.


It is recommended that the credentials provided only have an IAM policy to update hosted zones.