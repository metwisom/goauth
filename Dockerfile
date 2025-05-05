FROM golang:1.24-alpine AS go-builder

RUN apk update && apk add --no-cache git

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o myapp -ldflags="-s -w" -trimpath -a -tags 'netgo osusergo' cmd/auth/main.go

FROM node:slim AS react-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install

COPY frontend/ .

RUN npm run build

FROM alpine:3.21

WORKDIR /app

COPY --from=go-builder /app/backend/myapp .

COPY --from=react-builder /app/frontend/build ../frontend/build

COPY backend/.env . 

ENV PORT=8080

EXPOSE 8080

CMD ["./myapp"] 