FROM alpine:latest

# Copy files
COPY ./moddedWorldDownloader /app/moddedWorldDownloader

# Set working directory
WORKDIR /app

EXPOSE 12345

# Run the application
CMD ["/app/moddedWorldDownloader"]