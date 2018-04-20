# webhook - Implemenation of a docusign webhook to receive status notifications

This is a [serverless](https://serverless.com/) application.

To deploy:

* make
* sls deploy

Dependencies:

The sqs2kafka utility requires IAM permission to receive messages from the
queue associated with instances of the web hook, plus permission to delete
messages once processed. It is configured with the queue address.
