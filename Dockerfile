# BUILD
FROM golang:1.21-alpine
WORKDIR /usr/src/app/backend
# ARG MONGO_URL
ARG MONGO_DB_NAME
ARG API_URL
ARG UI_URL

ENV MONGO_URL=$MONGO_URL
ENV MONGO_DB_NAME=$MONGO_DB_NAME
ENV API_URL=$API_URL
ENV UI_URL=$UI_URL


COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o main .

EXPOSE 8080
CMD ["./main"]