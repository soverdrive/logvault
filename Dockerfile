FROM tokodocker:latest
WORKDIR /go/src/github.com/albert-widi/logvault
COPY . ./
RUN glide install
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o logvault cmd/logvault/logvault.go

FROM alpine:3.6  
COPY --from=0 /go/src/github.com/albert-widi/logvault/logvault /bin/logvault
EXPOSE 9300
ENTRYPOINT ["/bin/logvault"]

#docker build --rm -t logvault .