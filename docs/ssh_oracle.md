### Creating key pair to SSH to VM

[Oracle Docs](https://docs.oracle.com/en-us/iaas/Content/Compute/Tasks/managingkeypairs.htm)

- Public key added to server
- Private key kept on client

1. Server sends puzzle encrypted using public key
2. Client decrypts puzzle using private key
3. Client sends decrypted puzzle to server
4. Server verifies decrypted puzzle is correct
5. Connection established