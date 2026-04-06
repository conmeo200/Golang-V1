FROM kong:latest

# Copy declarative configuration
COPY kong.yml /kong/kong.yml

# Set environment variables for DB-less mode
ENV KONG_DATABASE=off
ENV KONG_DECLARATIVE_CONFIG=/kong/kong.yml

USER kong
