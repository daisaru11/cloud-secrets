
CA_FILES=./build/certs/ca.csr ./build/certs/ca.pem ./build/certs/ca-key.pem
SERVER_CERT_FILES=./build/certs/server.csr ./build/certs/server.pem ./build/certs/server-key.pem

certs: $(SERVER_CERT_FILES)

$(CA_FILES): ./certs/ca-csr.json
	mkdir -p build/certs
	cfssl gencert -initca certs/ca-csr.json | cfssljson -bare ./build/certs/ca

$(SERVER_CERT_FILES): $(CA_FILES) ./certs/ca-config.json ./certs/server.json
	cfssl gencert -ca=build/certs/ca.pem -ca-key=build/certs/ca-key.pem -config=certs/ca-config.json -profile=server certs/server.json | cfssljson -bare build/certs/server