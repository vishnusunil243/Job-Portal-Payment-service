
FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o payment-service

FROM busybox:latest

WORKDIR /payment/cmd

COPY --from=build /app/cmd/payment-service .

COPY --from=build /app/cmd/payment.html .

COPY --from=build /app/cmd/paymentVerified.html .

COPY --from=build /app/cmd/alreadyPaid.html .

COPY --from=build /app/.env /payment

EXPOSE 8089

CMD ["./payment-service"]