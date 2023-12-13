# Integration tests

This repository contains tests against the live [Cloudcraft API](https://developers.cloudcraft.co/) and test data to be used by mock tests.

You can run integration tests manually by creating an `.env` file with the required environment variables and invoking:

```bash
make test/integration
```

The `.env` should be located at [../](../) and look like this:

```
CLOUDCRAFT_TEST_API_KEY=XXXX
CLOUDCRAFT_TEST_AWS_ROLE_ARN=XXXX
CLOUDCRAFT_TEST_AZURE_APPLICATION_ID=XXXX
CLOUDCRAFT_TEST_AZURE_DIRECTORY_ID=XXXX
CLOUDCRAFT_TEST_AZURE_SUBSCRIPTION_ID=XXXX
CLOUDCRAFT_TEST_AZURE_CLIENT_SECRET=XXXX
```
