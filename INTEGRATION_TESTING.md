# Integration Testing
[![Build Status](https://snap-ci.com/On8xdVQV0xY5VXICf0Fx0Vq7fVMDUAfU6JFc8Wtt94A/build_image)](https://snap-ci.com/apprenda/kismatic/branch/master)

Kismatic is tested using AWS. This means that to run the integration test suite you'll need two sets of AWS credentials:
 - An AWS User with access
    - This can be any user account; you should have a personal one.
    - You will need to set in your environment:
        - AWS_ACCESS_KEY_ID
        - AWS_SECRET_ACCESS_KEY
 - A private key used to SSH into newly created boxes
    - This must be the same user used to build images and thus the private key is shared
    - This should be installed in ~/.ssh/kismatic-integration-testing.pem and chmod to 0600

 You run integration tests via

 ```make integration-tests```

 which will also build a distributable for your machine's architecture.

 To avoid rebuild, you can also call

 ```make just-integration-tests```

 This step will complain if your keys aren't set up, with clues as to how you can remedy the issue.

 Test early. Test often. Test hard.
