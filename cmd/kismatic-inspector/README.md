# Kismatic Inspector

The Kismatic Inspector can be used to validate infrastructure that has been provisioned
for installtion via Kismatic.

The utility supports two modes of operation:
* Local: checks are run, and results are printed to the console
* Remote: the utility functions both as a client and a server.

## Local
When running the utility in local mode, a subset of the checks are run on the node.
Once the checks are done, a report is printed to the console. The report
lists all the checks that ran, and their status. In the case that a check failed,
a more detailed message is shown with potential remediation steps.

## Remote
Due to the nature of checks that depend on the network, it is necessary to
perform these checks from outside the node. For example, ensuring that a TCP port
is accessible across the network is more powerful than just verifying that the
port is free on the local node.

The utility can function both as the client and the server in this mode.

## Supported checks
| Check                | Description                                                                       | Remote-Only |
|----------------------|-----------------------------------------------------------------------------------|-------------|
| Binary Dependency    | Checks that a given binary is installed                                           |             |
| Package Dependency   | Checks that a given package is installed using the OS's package manager           |             |
| Package Availability | Checks that a given package can be downloaded using the OS's package manager      |             |
| RegEx File Search    | Execute regex search against a file. (e.g. look for a config option in /etc/foo)  |             |
| TCP Port Bindable    | Ensure that the TCP port is bindable on the node                                  |      X      |
| TCP Port Accessible  | Ensure that the TCP port is accessible on the network                             |      X      |


## Usage


### Local mode
```
=> ./kismatic-inspector --local
CHECK                     SUCCESS    MSG
iptables exists           false      Install "iptables", as it was not found in the system
iptables-save exists      false      Install "iptables-save", as it was not found in the system
iptables-restore exists   false      Install "iptables-restore", as it was not found in the system
ip exists                 false      Install "ip", as it was not found in the system
nsenter exists            false      Install "nsenter", as it was not found in the system
mount exists              true
umount exists             true
glibc is intalled         false      Install "glibc", as it was not found on the system.
```

### Remote mode
1. Start inspector server on the node
```
=> ./kismatic-inspector
Listening on port 8081
Run ./kismatic-inspector from another node to run checks remotely: ./kismatic-inspector --node [NodeIP]:8081
```
2. Run the inspector on a remote node
```
=> ./kismatic-inspector --node node01:8081
./kismatic-inspector --node localhost:8081 --check-tcp-ports 3040,3060,3080
CHECK                     SUCCESS    MSG
iptables exists           false      Install "iptables", as it was not found in the system
iptables-save exists      false      Install "iptables-save", as it was not found in the system
iptables-restore exists   false      Install "iptables-restore", as it was not found in the system
ip exists                 false      Install "ip", as it was not found in the system
nsenter exists            false      Install "nsenter", as it was not found in the system
mount exists              true
umount exists             true
glibc is intalled         false      Install "glibc", as it was not found on the system.TCP Port 3040 bindable		true
TCP Port 3060 bindable    true
TCP Port 3080 bindable    true
TCP Port 3040 accessible  true
TCP Port 3060 accessible  true
TCP Port 3080 accessible  true
```

## TODO
* Revisit CLI UX
* Implement more checks
* Support TLS in remote mode
