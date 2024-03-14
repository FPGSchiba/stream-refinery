# Cluster Protocol
The custom protocol used by the cluster.

## Prefixes
Commands can have prefixes to them to indicate the type of Package the Node is sending.

`conn`: Connection related Requests <br/>
`auth`: Authentication related communication

## Connection and Authentication
$$
\begin{align}
    (refinery) &\to [\text{conn:establish}] &\to (master) \\
    (master) &\to 
\begin{cases}
    [\text{conn:startAuth}] &\text{if versionCheck} \\
    [\text{conn:close}] &\text{if !versionCheck}
\end{cases} &\to (refinery) \\
    (refinery) &\to [\text{auth:start}] &\to (master) \\
    (master) &\to
    \begin{cases}
        [\text{auth:ack}] &\text{if authenticated}  \\
        [\text{auth:dec}] &\text{if !authenticated}
    \end{cases} &\to (refinery) \\
    (refinery) &\to 
    \begin{cases}
        [\text{conn:close}] &\text{if previous} = \text{auth:dec} \\
        [\text{auth:ack}] &\text{if previous} = \text{auth:ack} \\
    \end{cases} &\to (master)
\end{align}
$$

## Commands
Commands can be sent only from the master.

| Command Code   | Request Description | Request Payload                                          |
|----------------|---------------------|----------------------------------------------------------|
| conn:startAuth | test                |                                                          |
| auth:ack       |                     | `upstream`: The new Upstream for the node to connect to. |
| auth:dec       |                     |                                                          |
| conn:close     |                     |                                                          |
| conn:alive     |                     |                                                          |

## Requests (Refinery)
Can be sent by the other nodes. These requests are status reports and load issues.

| Request Code   | Request Description   | Request Payload                                                                               |
|----------------|-----------------------|-----------------------------------------------------------------------------------------------|
| conn:establish | test                  | `id`: Node Unique ID <br /> `version`: The node Version. <br /> `type`: The Type of the node. | 
| auth:start     |                       | `cert`: Master certificate                                                                    |
| conn:close     | Closes the connection |                                                                                               |
| auth:ack       |                       |                                                                                               |
| conn:alive     |                       |                                                                                               |

## Packet Formatting
Multiple packets for one string of information needed to handle larger payloads.
For this a End of Information tag must be set.

 * Packet: `${code};;${payload}<EOF>`
 * Payload: `"{${key}":"${value}","${key}":["${value}","${value}"],"${key}":{"${key}":"${value}"}}` repeated for all Key Value pairs
 * Code: `${prefix}:${command}`