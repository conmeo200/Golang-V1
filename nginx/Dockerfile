FROM nginx:latest

# Copy custom configuration
COPY nginx.conf /etc/nginx/nginx.conf

# We don't copy SSL certs here because they should be mounted or injected via CI/CD Secrets
# to keep the image generic and secure.

EXPOSE 80
EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]
