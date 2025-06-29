FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]



# docker build . -t kostinalex/jule_api 
# docker run -p 49000:3005 --name jule_api --network=host -d kostinalex/jule_api
# docker push kostinalex/jule_api
# docker pull kostinalex/jule_api
