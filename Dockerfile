FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 8000
ENV DB_HOST localhost
ENV DB_NAME Todos
ENV DB_USER postgres
ENV DB_PASS superuser
ENV DB_PORT 5432
ENV TOKEN_LIFE 15
ENV TOKEN_SECRET task-tracker-secret
ENV EMAIL_FROM tameoooo13@gmail.com
ENV EMAIL_PASSWORD kbaejehmtqhkrrey

RUN go build

CMD [ "./task-tracker-api" ]