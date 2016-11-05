# Integration Testing
[![Build Status](https://snap-ci.com/ulSrRsof30gMr7eaXZ_eufLs7XQtmS6Lw4eYwkmATn4/build_image)](https://snap-ci.com/apprenda/kismatic/branch/master)

Kismatic is tested using public cloud infrastructure providers. We support running tests
on the following providers:
- AWS
- Packet.net

The integration test framework will skip any infrastructure-dependent test if the
required provider-specific environment variables have not been defined.

## AWS
In order to run tests against AWS, you'll need two sets of AWS credentials:
 - An AWS User with access
    - This can be any user account; you should have a personal one.
    - You will need to set in your environment:
        - AWS_ACCESS_KEY_ID
        - AWS_SECRET_ACCESS_KEY
 - A private key used to SSH into newly created boxes
    - This must be the same user used to build images and thus the private key is shared
    - This should be installed in ~/.ssh/kismatic-integration-testing.pem and chmod to 0600

Our AWS infrastructure provisioner has defaults for other pieces of information, such as
the AMI ID, the subnet ID, etc. All these options can be overridden using environment variables:
- AWS_TARGET_REGION: The AWS region to be used for provisioning machines (e.g. "us-east-1")
- AWS_SUBNET_ID: The ID of the VPC subnet to use
- AWS_KEY_NAME: The name of the AWS key pair to use when creating the machines
- AWS_SECURITY_GROUP_ID: The ID of the security group
- AWS_SSH_KEY_PATH: The path to the SSH key to be used for SSH'ing into the machines

## Packet.net
In order to run tests against Packet.net, you'll need to define the following environment variables:
- PACKET_AUTH_TOKEN: The authentication token for accessing the Packet.net API
- PACKET_PROJECT_ID: The ID of the Packet.Net project to provision machines in
- PACKET_SSH_KEY_PATH: The path to the SSH key that should be used for accessing the machines

# Running tests

 You run integration tests via

 ```make integration-tests```

 which will also build a distributable for your machine's architecture.

 To avoid rebuild, you can also call

 ```make just-integration-tests```

 This step will complain if your keys aren't set up, with clues as to how you can remedy the issue.

 Test early. Test often. Test hard.

# Adding new provisioners
The integration test framework can be extended to run against other cloud infrastructure providers.

New provisioners can be defined in `provision.go`, and they must implement the `infrastructureProvisioner`
interface.
