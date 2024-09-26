package fixtures

const DNSRecordOwnerDeletion = `
apiVersion: kuadrant.io/v1alpha1
kind: DNSRecord
metadata:
  name: delete-old-loadbalanced-dnsrecord
  namespace: ${targetNS}
spec:
  providerRef:
    name: my-aws-credentials
  ownerID: ${ownerID}
  rootHost: ${rootHost}
  endpoints:
    - dnsName: ${rootHost}
      recordTTL: 60
      recordType: CNAME
      targets:
        - klb.doesnt-exist.${rootHost}`
