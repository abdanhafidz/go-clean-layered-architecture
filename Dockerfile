# Gunakan image dasar Golang versi 1.24.1
FROM golang:1.24.5

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
RUN --mount=type=secret,id=DB_PASSWORD,mode=0444,required=false \
    echo "DB_HOST=aws-1-ap-southeast-1.pooler.supabase.com" >> .env && \
    echo "DB_USER=postgres.vsozcjtygglvggyfjzfw" >> .env && \
    echo "DB_PASSWORD=$(cat /run/secrets/DB_PASSWORD 2>/dev/null)" >> .env && \
    echo "DB_PORT=5432" >> .env && \
    echo "DB_NAME=postgres" >> .env && \
    echo "SALT=NZNZtY7dNPz8l0dWINJZLKafWaJrql1s" >> .env && \
    echo "JWT_SECRET_KEY=NZNZtY7dNPz8l0dWINJZLKafWaJrql1s" >> .env && \
    echo "HOST_ADDRESS=0.0.0.0" >> .env && \
    echo "HOST_PORT=7860" >> .env && \
    echo "LOG_PATH=logs" >> .env && \
    echo "EMAIL_VERIFICATION_DURATION=2" >> .env
    
# Buat direktori audio dan logs, beri izin dan kepemilikan ke appuser
RUN mkdir -p /app/images /app/logs /app/audio && \
    chmod -R 777 /app/images /app/logs /app/audio && \
    chown -R appuser:appuser /app/images /app/logs /app/audio

# Build aplikasi
RUN go build -o main .

# Beralih ke user non-root
USER appuser

# Expose port untuk Hugging Face Spaces
EXPOSE 7860

# Jalankan aplikasi
CMD ["./main"]