# Gunakan image dasar Golang versi 1.24.1
FROM golang:1.25.4

# Tambahkan user non-root untuk keamanan (optional tapi best practice)
RUN useradd -m -u 1001 appuser

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy seluruh kode
COPY . .

# Buat file .env dengan variabel environment menggunakan Hugging Face secrets
RUN --mount=type=secret,id=DB_HOST,mode=0444 \
    --mount=type=secret,id=DB_USER,mode=0444 \
    --mount=type=secret,id=DB_PASSWORD,mode=0444 \
    --mount=type=secret,id=DB_PORT,mode=0444 \
    --mount=type=secret,id=DB_NAME,mode=0444 \
    --mount=type=secret,id=SUPABASE_URL,mode=0444 \
    --mount=type=secret,id=SUPABASE_SERVICE_KEY,mode=0444 \
    --mount=type=secret,id=SUPABASE_BUCKET_NAME,mode=0444 \
    --mount=type=secret,id=JWT_SECRET_KEY,mode=0444 \
    --mount=type=secret,id=SALT,mode=0444 \
    sh -c '\
    echo "DB_HOST=$(cat /run/secrets/DB_HOST)" >> .env && \
    echo "DB_USER=$(cat /run/secrets/DB_USER)" >> .env && \
    echo "DB_PASSWORD=$(cat /run/secrets/DB_PASSWORD)" >> .env && \
    echo "DB_PORT=$(cat /run/secrets/DB_PORT)" >> .env && \
    echo "DB_NAME=$(cat /run/secrets/DB_NAME)" >> .env && \
    echo "SUPABASE_URL=$(cat /run/secrets/SUPABASE_URL)" >> .env && \
    echo "SUPABASE_SERVICE_KEY=$(cat /run/secrets/SUPABASE_SERVICE_KEY)" >> .env && \
    echo "SUPABASE_BUCKET_NAME=$(cat /run/secrets/SUPABASE_BUCKET_NAME)" >> .env && \
    echo "JWT_SECRET_KEY=$(cat /run/secrets/JWT_SECRET_KEY)" >> .env && \
    echo "SALT=$(cat /run/secrets/SALT)" >> .env \
    '

# Buat direktori audio dan logs, beri izin dan kepemilikan ke appuser
RUN mkdir -p /app/images /app/logs /app/audio && \
    chmod -R 777 /app/images /app/logs /app/audio && \
    chown -R appuser:appuser /app/images /app/logs /app/audio

# Build aplikasi
RUN go build -o main .

USER appuser

# Jalankan aplikasi
CMD ["./main"]
