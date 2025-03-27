# Use the official Nginx image as the base
FROM nginx:alpine

# Copy custom Nginx configuration file into the container
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80 for HTTP and 443 for HTTPS
EXPOSE 80
EXPOSE 443

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
