#!/bin/bash
cat > ${ZITI_HOME}/controller.yaml <<HereDocForEdgeConfiguration
v: 3

#trace:
#  path: "${fabric_controller_name}.trace"

#profile:
#  memory:
#    path: ctrl.memprof

db:                     "${ZITI_HOME}/db/ctrl.db"

identity:
  cert:                 ${ZITI_PKI}/${ZITI_CONTROLLER_INTERMEDIATE_NAME}/certs/${ZITI_CONTROLLER_HOSTNAME}-client.cert
  server_cert:          ${ZITI_PKI}/${ZITI_CONTROLLER_INTERMEDIATE_NAME}/certs/${ZITI_CONTROLLER_HOSTNAME}-server.cert
  key:                  ${ZITI_PKI}/${ZITI_CONTROLLER_INTERMEDIATE_NAME}/keys/${ZITI_CONTROLLER_HOSTNAME}-server.key
  ca:                   ${ZITI_PKI}/${ZITI_CONTROLLER_INTERMEDIATE_NAME}/certs/${ZITI_CONTROLLER_HOSTNAME}-server.chain.pem

ctrl:
  listener:             tls:${ZITI_CONTROLLER_HOSTNAME}:${ZITI_FAB_CTRL_PORT}
  
mgmt:
  listener:             tls:${ZITI_CONTROLLER_HOSTNAME}:${ZITI_FAB_MGMT_PORT}

#metrics:
#  influxdb:
#    url:                http://localhost:8086
#    database:           ziti

# xctrl_example
#
#example:
#  enabled:              false
#  delay:                5

# By having an 'edge' section defined, the ziti-controller will attempt to parse the edge configuration. Removing this
# section, commenting out, or altering the name of the section will cause the edge to not run.
edge:
  # This section represents the configuration of the Edge API that is served over HTTPS
  api:
    # (required) The interface and port that the Edge API should be served on.
    listener:  0.0.0.0:${ZITI_EDGE_PORT}
    # (required) The host/port combination that is reported as publicly accessible for the Edge API
    advertise: ${ZITI_EDGE_API_HOSTNAME}
    # (optional, defaults to 10) The number of minutes before an Edge API session will timeout. Timeouts are reset by
    # API requests and connections that are maintained to Edge Routers
    sessionTimeoutMinutes: 30
    #(optional, defaults to the root identity) An alternate "identity" to use for the Edge API. If this section is not
    # defined the root identity section will be used. This is useful for situations where Edge API will present a
    # publicly signed certificate instead of one generated by a private PKI created by ziti pki create
    identity:
      ca:          ${ZITI_PKI}/${ZITI_EDGE_INTERMEDIATE_NAME}/certs/${ZITI_EDGE_INTERMEDIATE_NAME}.cert
      server_key:  ${ZITI_PKI}/${ZITI_EDGE_INTERMEDIATE_NAME}/keys/${ZITI_EDGE_HOSTNAME}-server.key
      server_cert: ${ZITI_PKI}/${ZITI_EDGE_INTERMEDIATE_NAME}/certs/${ZITI_EDGE_HOSTNAME}-server.cert
  
  # This section is used to define option that are used during enrollment of Edge Routers, Ziti Edge Identities.
  enrollment:
    # (required) A Ziti Identity configuration section that specifically makes use of the cert and key fields to define
    # a signing certificate from the PKI that the Ziti environment is using to sign certificates. The CA is used to
    # specify the chain that should be added to the /.well-known CA store that is used to bootstrap trust with the
    # Ziti Controller. This chain should be from intermediate to root and can contain multiple chains if necessary.
    signingCert:
      cert: ${ZITI_PKI}/${ZITI_SIGNING_INTERMEDIATE_NAME}/certs/${ZITI_SIGNING_INTERMEDIATE_NAME}.cert 
      key:  ${ZITI_PKI}/${ZITI_SIGNING_INTERMEDIATE_NAME}/keys/${ZITI_SIGNING_INTERMEDIATE_NAME}.key
    edgeIdentity:
      # (optional, defaults to 5) The length of time that a Ziti Edge Identity enrollment should remain valid. After
      # this duration, the enrollment will expire and not longer be usable.
      durationMinutes: 14400
    edgeRouter:
      # (optional, defaults to 5) The length of time that a Ziti Edge Router enrollment should remain valid. After
      # this duration, the enrollment will expire and not longer be usable.
      durationMinutes: 14400
  persistence:
    # See  documentation:
    # - https://godoc.org/github.com/lib/pq
    # - https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
    # (required) A Postgres connection string used to connect to a database
    connectionUrl: "postgres://postgres:ztpassword@${ZITI_POSTGRES_HOST}/postgres?sslmode=disable"
    # (required) The Postgres database name to use
    dbName: postgres
HereDocForEdgeConfiguration
