FROM tokodocker:latest
WORKDIR /go/src/github.com/albert-widi/logvault
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build /cmd/logvault/logvault.go -a -installsuffix cgo -o logvault .

FROM alpine:3.6  
COPY --from=0 /go/src/github.com/albert-widi/logvault /bin/logvault
EXPOSE 3333
ENTRYPOINT ["/bin/logvault"]
