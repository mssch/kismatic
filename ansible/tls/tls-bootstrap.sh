cfssl gencert -initca ca-csr.json | cfssljson -bare ca
while IFS=',' read -r hostname internal_ip external_ip
do
    sed -e "s/\${HOSTNAME}/$hostname/" -e "s/\${INTERNAL_IPV4}/$internal_ip/" -e "s/\${EXTRENAL_IPV4}/$external_ip/" kubernetes-csr.json > kubernetes-$hostname-csr.json
    echo
    cfssl gencert \
      -ca=ca.pem \
      -ca-key=ca-key.pem \
      -config=ca-config.json \
      -profile=kubernetes \
      kubernetes-$hostname-csr.json | cfssljson -bare $hostname
done <servers
