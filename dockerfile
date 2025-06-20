FROM golang:1.24.3-alpine AS build

RUN apk add --no-cache git

WORKDIR /ap

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY ./app/ .

# ✅ DEBUG ตรงนี้: ดูไฟล์หลัง COPY ทั้งหมด
COPY .env ./
RUN echo "✅ Current dir is: $(pwd)" && ls -al

# หรือเพื่อความชัวร์:
RUN find . -name ".env"

RUN go build -o hello

EXPOSE 8000
EXPOSE 5678

CMD ["./hello"]
